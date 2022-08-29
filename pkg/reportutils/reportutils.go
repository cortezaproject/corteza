package reportutils

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
)

type (
	ReportWorkload struct {
		Pipeline  dal.Pipeline
		FrameDefs types.ReportFrameDefinitionSet
	}

	PipelineRunner interface {
		FindModel(dal.ModelRef) *dal.Model
		Run(context.Context, dal.Pipeline) (dal.Iterator, error)
		Dryrun(context.Context, dal.Pipeline) error
	}
)

// Workloads creates a set of workloads for the given pipeline and frame definitions
func Workloads(pr PipelineRunner, steps types.ReportStepSet, defs types.ReportFrameDefinitionSet) (out []ReportWorkload, err error) {
	// Construct a pipeline from the steps; we'll slice it later
	base, err := Pipeline(pr, steps)
	if err != nil {
		return
	}

	// Prepare workloads based on the provided definitions
	//
	// - If consecutive definitions point to the same source with the same name
	//   consider them to fall under the same workload (the link step)
	// - else, one def per workload
	auxDefs := make(types.ReportFrameDefinitionSet, 0)
	var auxOut ReportWorkload
	for i, def := range defs {
		if i == 0 {
			auxDefs = append(auxDefs, def)
			continue
		}

		// This is for the link step
		if def.Name == defs[i-1].Name && def.Source == defs[i-1].Source {
			auxDefs = append(auxDefs, def)
			continue
		}

		// This is for the rest
		auxOut, err = makeWorkload(base, auxDefs)
		if err != nil {
			return
		}
		out = append(out, auxOut)
		auxDefs = make(types.ReportFrameDefinitionSet, 0)
		auxDefs = append(auxDefs, def)
	}

	// Handle the ones (potentially) not covered by the above loop
	if len(auxDefs) > 0 {
		auxOut, err = makeWorkload(base, auxDefs)
		if err != nil {
			return
		}
		out = append(out, auxOut)
	}

	return
}

// Frames returns a set of frames for the given workload & iterator combo
func Frames(ctx context.Context, iter dal.Iterator, workload ReportWorkload) (ff []*types.ReportFrame, err error) {
	// Preprocessing on the workload's frame definitions; assure all columns/metdata are there
	// to avoid nonesense later down the line
	updateDefAttrs(workload)

	// @todo perhaps need to change; for now only this scenario introduces multiple
	//       frame defs per workload
	if len(workload.FrameDefs) > 1 {
		return framifyLinkIter(ctx, iter, workload)
	}

	return framifyIter(ctx, iter, workload)
}

// framifyLinkIter is a handler dedicated for the link step due to it's unique output
func framifyLinkIter(ctx context.Context, iter dal.Iterator, workload ReportWorkload) (ff []*types.ReportFrame, err error) {
	defs := workload.FrameDefs
	// @note this will only be called for the link step so it can freely panic if violated
	r := workload.Pipeline[0].(*dal.Link)

	// Unpack frame definitions for the link
	defLeft, defRight := unpackLinkDefs(defs, workload.Pipeline)

	// Init vars to keep track of the progress
	// @note true is left, false is right
	counters := make(map[bool]uint)

	builders := make(map[bool]*reportFrameBuilder)
	builders[true] = newReportFrameBuilder(defLeft)
	builders[false] = newReportFrameBuilder(defRight)
	builders[false].withRefs(r.On.Right)

	limits := make(map[bool]uint)
	if defLeft.Paging != nil {
		limits[true] = defLeft.Paging.Limit
	}
	if defRight.Paging != nil {
		limits[false] = defRight.Paging.Limit
	}

	// Helper to determine if we need a next cursor
	nextCursor := false

	// Helpers for reading iterators
	var (
		ref    string
		row    = &dal.Row{}
		doingF = false
	)

	for iter.Next(ctx) {
		if limits[true] > 0 && counters[true] >= limits[true] {
			nextCursor = true
			break
		}
		row.Reset()

		_ = iter.Scan(row)

		// Determine ref and which vars to use
		aux, _ := row.GetValue(dal.LinkRefIdent, 0)
		ref = cast.ToString(aux)
		if ref == "" {
			ref = defLeft.Ref
		}
		left := ref == defLeft.Ref

		// When needed, flush the finished frames to the output
		if left && doingF {
			ff = append(ff, builders[false].done())
			doingF = false
		} else if !left {
			doingF = true
		}

		builders[left].addRow(row)
		counters[left]++
	}
	if err = iter.Err(); err != nil {
		return
	}

	// If the loop ended before the limit cut it off, we need to finish the
	// last right frame as it wasn't yet in the above loop
	if !nextCursor {
		if doingF {
			ff = append(ff, builders[false].done())
		}
	}

	// Apply paging cursor to the left frame
	// @todo consider applying them to the right as well, for now, no
	if nextCursor {
		if builders[true].frame.Paging == nil {
			builders[true].frame.Paging = &filter.Paging{}
		}
		builders[true].frame.Paging.NextPage, err = iter.ForwardCursor(row)
		if err != nil {
			return
		}
	}

	// Complete the output with the left frame
	// @note the left frame goes to the start and the right frames are in the same order
	//       as the related rows from the left frame.
	return append([]*types.ReportFrame{builders[true].done()}, ff...), nil
}

// unpackLinkDefs returns the left and the right frame definition disregarding
// the order of the definitions in the input
func unpackLinkDefs(defs types.ReportFrameDefinitionSet, pp dal.Pipeline) (left, right *types.ReportFrameDefinition) {
	r := pp[0]

	l := r.(*dal.Link)

	find := func(defs types.ReportFrameDefinitionSet, ref string) *types.ReportFrameDefinition {
		for _, def := range defs {
			if def.Ref == ref {
				return def
			}
		}
		return nil
	}

	return find(defs, l.RelLeft), find(defs, l.RelRight)
}

// framifyIter is a generic iter to frame handler
func framifyIter(ctx context.Context, iter dal.Iterator, workload ReportWorkload) (ff []*types.ReportFrame, err error) {
	defs := workload.FrameDefs

	// @note only the link step takes multiple defs and that one is not covered
	//       by this function
	if len(defs) != 1 {
		panic(fmt.Sprintf("impossible state: expecting one frame definition, got %d", len(defs)))
	}
	def := defs[0]

	// Init vars to keep track of the progress
	limit := uint(0)
	counter := uint(0)
	builder := newReportFrameBuilder(def)
	if def.Paging != nil {
		limit = def.Paging.Limit
	}

	// Helper to determine if we need a next cursor
	nextCursor := false

	// Helpers for reading iterators
	row := &dal.Row{}

	for iter.Next(ctx) {
		if limit > 0 && counter >= limit {
			nextCursor = true
			break
		}
		row.Reset()

		_ = iter.Scan(row)
		builder.addRow(row)
		counter++
	}
	if err = iter.Err(); err != nil {
		return
	}

	// Apply paging cursor to the frame
	if nextCursor {
		if builder.frame.Paging == nil {
			builder.frame.Paging = &filter.Paging{}
		}
		builder.frame.Paging.NextPage, err = iter.ForwardCursor(row)
		if err != nil {
			return
		}
	}

	return append(ff, builder.done()), nil
}

// Pipeline creates a pipeline from the given steps
func Pipeline(pr PipelineRunner, steps types.ReportStepSet) (out dal.Pipeline, err error) {
	out = make(dal.Pipeline, 0, len(steps))

	for _, step := range steps {
		switch {
		case step.Load != nil:
			aux, err := makeStepLoad(pr, *step.Load)
			if err != nil {
				return nil, err
			}
			out = append(out, aux)

		case step.Join != nil:
			aux, err := makeStepJoin(*step.Join)
			if err != nil {
				return nil, err
			}
			out = append(out, aux)

		case step.Link != nil:
			aux, err := makeStepLink(*step.Link)
			if err != nil {
				return nil, err
			}
			out = append(out, aux)

		case step.Aggregate != nil:
			aux, err := makeStepAggregate(*step.Aggregate)
			if err != nil {
				return nil, err
			}
			out = append(out, aux)
		}
	}

	return out, out.LinkSteps()
}

// DescribePipeline returns a set of frame descriptions based on the given pipeline
func DescribePipeline(pp dal.Pipeline, sources []string) (out types.FrameDescriptionSet, err error) {
	var aux types.FrameDescriptionSet

	for _, src := range sources {
		sub := pp.Slice(src)
		s := sub[0]

		aux, err = describePipeline(s, src)
		if err != nil {
			return
		}

		out = append(out, aux...)
	}

	return
}

func describePipeline(s dal.PipelineStep, src string) (out types.FrameDescriptionSet, err error) {
	aa := s.Attributes()

	out = make(types.FrameDescriptionSet, len(aa))
	for i, a := range aa {
		out[i].Source = src
		out[i].Columns = mappingToFrameCols(a)
	}

	// @note this case is only possible for the link step; expand when/if needed
	if len(out) == 2 {
		l := s.(*dal.Link)
		out[0].Ref = l.RelLeft
		out[1].Ref = l.RelRight
	}

	return
}

func makeWorkload(pp dal.Pipeline, defs types.ReportFrameDefinitionSet) (out ReportWorkload, err error) {
	// We'll trust all of the defs point to the same source; this should be validated
	// way sooner
	def := defs[0]

	out.FrameDefs = defs
	out.Pipeline = pp.Slice(def.Source)

	return
}

// @todo address with col/attr rework/rethink
func mappingToFrameCols(mm []dal.AttributeMapping) types.ReportFrameColumnSet {
	out := make(types.ReportFrameColumnSet, 0, len(mm))

	for _, m := range mm {
		out = append(out, types.ReportFrameColumn{
			Name: m.Identifier(),
			Kind: "String",
		})
	}

	return out
}

// Report step -> DAL step conversion

func makeStepJoin(step types.ReportStepJoin) (out *dal.Join, err error) {
	out = &dal.Join{
		Ident:    step.Name,
		RelLeft:  step.LocalSource,
		RelRight: step.ForeignSource,

		On: dal.JoinPredicate{
			Left:  step.LocalColumn,
			Right: step.ForeignColumn,
		},
		Filter: dal.FilterForExpr(step.Filter.Node()),
	}
	return
}

func makeStepLink(step types.ReportStepLink) (out *dal.Link, err error) {
	out = &dal.Link{
		Ident:    step.Name,
		RelLeft:  step.LocalSource,
		RelRight: step.ForeignSource,

		On: dal.LinkPredicate{
			Left:  step.LocalColumn,
			Right: step.ForeignColumn,
		},
		Filter: dal.FilterForExpr(step.Filter.Node()),
	}
	return
}

func makeStepAggregate(step types.ReportStepAggregate) (out *dal.Aggregate, err error) {
	out = &dal.Aggregate{
		Ident:     step.Name,
		RelSource: step.Source,
		Filter:    dal.FilterForExpr(step.Filter.Node()),

		Group:         step.Keys.DalMapping(),
		OutAttributes: step.Columns.DalMapping(),
	}
	return
}

func makeStepLoad(pr PipelineRunner, step types.ReportStepLoad) (out *dal.Datasource, err error) {
	mfr, c, err := getModelRef(step)
	if err != nil {
		return
	}

	return &dal.Datasource{
		Ident: step.Name,
		Filter: dal.FilterForExpr(step.Filter.Node()).
			WithConstraints(c),
		ModelRef:      mfr,
		OutAttributes: filteredModelAttributes(pr, step, mfr),
	}, nil

}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities...

// filteredModelAttributes returns the requested attributes based on the step
// definition or all attributes if none are specified
//
// The function collects the attributes from the DAL model to omit the attribute
// construction step.
func filteredModelAttributes(pr PipelineRunner, step types.ReportStepLoad, mfr dal.ModelRef) (out []dal.AttributeMapping) {
	out = make([]dal.AttributeMapping, 0, 100)

	// All of the attributes
	fullAttrs, err := getModelAttrs(pr, mfr)
	if err != nil {
		return
	}

	// No filtering, return all
	if len(step.Columns) == 0 {
		return attrToMapping(fullAttrs...)
	}

	// Filter out the ones we don't want
	reqIndex := make(map[string]bool)
	for _, col := range step.Columns {
		reqIndex[col.Name] = true
	}

	for _, attr := range fullAttrs {
		if reqIndex[attr.Ident] {
			out = append(out, attrToMapping(attr)...)
		}
	}

	return
}

func getModelAttrs(pr PipelineRunner, mfr dal.ModelRef) (attrs dal.AttributeSet, err error) {
	m := pr.FindModel(mfr)
	if m == nil {
		return nil, fmt.Errorf("model not found: %v", mfr)
	}

	return m.Attributes, nil
}

func attrToMapping(aa ...*dal.Attribute) (out []dal.AttributeMapping) {
	for _, a := range aa {
		out = append(out, dal.SimpleAttr{
			Ident: a.Ident,
			Src:   a.Ident,
			Props: dal.MapProperties{
				Type:      a.Type,
				Nullable:  a.Type.IsNullable(),
				IsPrimary: a.PrimaryKey,
			},
		})
	}
	return
}

// @todo support for ns/mod by handle
func getModelRef(step types.ReportStepLoad) (out dal.ModelRef, constraints map[string][]any, err error) {
	var (
		moduleID, connectionID, namespaceID uint64
		aux                                 any
		ok                                  bool
	)

	constraints = make(map[string][]any)

	if aux, ok = step.Definition["moduleID"]; ok {
		moduleID = cast.ToUint64(aux)
	} else if aux, ok = step.Definition["module"]; ok {
		moduleID = cast.ToUint64(aux)
	} else {
		err = fmt.Errorf("step definition is missing moduleID")
		return
	}

	if aux, ok = step.Definition["namespaceID"]; ok {
		namespaceID = cast.ToUint64(aux)
	} else if aux, ok = step.Definition["module"]; ok {
		namespaceID = cast.ToUint64(aux)
	} else {
		err = fmt.Errorf("step definition is missing namespaceID")
		return
	}

	if aux, ok = step.Definition["connectionID"]; ok {
		connectionID = cast.ToUint64(aux)
	} else if aux, ok = step.Definition["connection"]; ok {
		connectionID = cast.ToUint64(aux)
	}

	constraints["moduleID"] = []any{moduleID}
	constraints["namespaceID"] = []any{namespaceID}

	return dal.ModelRef{ConnectionID: connectionID, ResourceID: moduleID}, constraints, nil
}

func updateDefAttrs(workload ReportWorkload) {
	ppAttrs := workload.Pipeline[0].Attributes()
	for i, def := range workload.FrameDefs {
		if len(def.Columns) > 0 {
			continue
		}

		def.Columns = mappingToFrameCols(ppAttrs[i])
	}
}

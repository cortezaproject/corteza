// Package reporting implements utilities for working with system reports and
// DAL pipelines
//
// The package exists primarily for making the system service codebase smaller
// and prepare for reusability with compose record reports.
package reporting

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
)

type (
	// run is a wrapper around the pipeline and corresponding
	// output frame definitions
	run struct {
		Pipeline dal.Pipeline
		Defs     FrameDefinitionSet
	}

	modelFinder interface {
		FindModel(dal.ModelRef) *dal.Model
	}

	// @todo not a fan of this anymore; consider reworking pkg/dal/service to
	//       provide a more sequential approach we can use.
	dryRunner interface {
		modelFinder
		Dryrun(context.Context, dal.Pipeline) error
	}
)

// Runs create a set of report runs based on step and frame definitions
func Runs(pr modelFinder, steps types.ReportStepSet, defs FrameDefinitionSet) (out []run, err error) {
	// Prepare runs based on the provided definitions
	//
	// - If consecutive definitions point to the same source with the same name
	//   consider them to fall under the same workload (the link step)
	// - else, one def per workload
	auxDefs := make(FrameDefinitionSet, 0)
	var aux run
	for i, def := range defs {
		if i == 0 {
			auxDefs = append(auxDefs, def)
			continue
		}

		// Definitions fall together
		if def.Name == defs[i-1].Name && def.Source == defs[i-1].Source {
			auxDefs = append(auxDefs, def)
			continue
		}

		// Make run for the previous definition (exclude current!!)
		aux, err = makeRun(pr, steps, auxDefs)
		if err != nil {
			return
		}
		out = append(out, aux)

		// Prepare next definition batch including the current one
		auxDefs = make(FrameDefinitionSet, 0)
		auxDefs = append(auxDefs, def)
	}

	// Handle the ones (potentially) not covered by the above loop
	if len(auxDefs) > 0 {
		aux, err = makeRun(pr, steps, auxDefs)
		if err != nil {
			return
		}
		out = append(out, aux)
	}

	return
}

// Frames returns a set of Frame for the given workload & iterator combo
func Frames(ctx context.Context, iter dal.Iterator, r run) (ff []*Frame, err error) {
	// Preprocessing on the workload's frame definitions; assure all
	// columns/metdata are there to avoid nonesense later down the line
	updateDefinitionColumns(r)

	// @todo perhaps need to change; for now only this scenario introduces multiple
	//       frame defs per workload
	if len(r.Defs) > 1 {
		return stepLinkFrames(ctx, iter, r)
	}

	return stepFrames(ctx, iter, r)
}

// Describe returns a set of frame descriptions based on the given pipeline
func Describe(ctx context.Context, rr dryRunner, ss types.ReportStepSet, sources []string) (out []FrameDescription, err error) {
	// Make a run for the whole thing
	pp, err := makePipeline(rr, ss, nil)
	if err != nil {
		return
	}

	var aux []FrameDescription
	for _, src := range sources {
		// Use the requested source as root
		sub := pp.Slice(src)
		err = rr.Dryrun(ctx, sub)
		if err != nil {
			return
		}

		s := sub[0]
		aux, err = describePipeline(s, src)
		if err != nil {
			return
		}
		out = append(out, aux...)
	}

	return
}

// stepLinkFrames is dedicated for the link step due to it's unique output
func stepLinkFrames(ctx context.Context, iter dal.Iterator, r run) (ff []*Frame, err error) {
	defs := r.Defs
	// @note this will only be called for the link step so it can freely panic if violated
	defLink := r.Pipeline[0].(*dal.Link)

	// Unpack frame definitions for the link
	defLeft, defRight := unpackLinkDefs(defs, r.Pipeline)

	// Init vars to keep track of the progress
	// @note true is left, false is right
	counters := make(map[bool]uint)

	builders := make(map[bool]*reportFrameBuilder)
	builders[true] = newReportFrameBuilder(defLeft)
	builders[false] = newReportFrameBuilder(defRight)
	builders[false].linked(defLink.On.Right, defLink.On.Left, defLink.RelLeft)

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
	return append([]*Frame{builders[true].done()}, ff...), nil
}

// unpackLinkDefs returns the left and the right frame definition disregarding
// the order of the definitions in the input
func unpackLinkDefs(defs FrameDefinitionSet, pp dal.Pipeline) (left, right *FrameDefinition) {
	r := pp[0]

	l := r.(*dal.Link)

	find := func(defs FrameDefinitionSet, ref string) *FrameDefinition {
		for _, def := range defs {
			if def.Ref == ref {
				return def
			}
		}
		return nil
	}

	return find(defs, l.RelLeft), find(defs, l.RelRight)
}

// stepFrames is a generic iter to Frame converter
func stepFrames(ctx context.Context, iter dal.Iterator, r run) (ff []*Frame, err error) {
	defs := r.Defs

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

func describePipeline(s dal.PipelineStep, src string) (out []FrameDescription, err error) {
	aa := s.Attributes()

	out = make([]FrameDescription, len(aa))
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

func makeRun(pr modelFinder, ss types.ReportStepSet, defs FrameDefinitionSet) (out run, err error) {
	var pp dal.Pipeline
	pp, err = makePipeline(pr, ss, defs)
	if err != nil {
		return
	}

	out.Defs = defs
	out.Pipeline = pp.Slice(defs[0].Source)
	return
}

// makePipeline converts the given report steps into the DAL pipeline
func makePipeline(mf modelFinder, ss types.ReportStepSet, defs FrameDefinitionSet) (pp dal.Pipeline, err error) {
	for _, step := range ss {
		switch {
		case step.Load != nil:
			aux, err := convStepLoad(mf, *step.Load, defs.FilterBySource(step.Load.Name))
			if err != nil {
				return nil, err
			}
			pp = append(pp, aux)

		case step.Aggregate != nil:
			aux, err := convStepAggregate(*step.Aggregate, defs.FilterBySource(step.Aggregate.Name))
			if err != nil {
				return nil, err
			}
			pp = append(pp, aux)

		case step.Join != nil:
			aux, err := convStepJoin(*step.Join, defs.FilterBySource(step.Join.Name))
			if err != nil {
				return nil, err
			}
			pp = append(pp, aux)

		case step.Link != nil:
			aux, err := convStepLink(*step.Link, defs.FilterBySource(step.Link.Name))
			if err != nil {
				return nil, err
			}
			pp = append(pp, aux)

		default:
			// this should never happen
			panic(fmt.Errorf("unknown step type: %v", step.Kind))
		}
	}

	return pp, pp.LinkSteps()
}

// mappingToFrameCols converts pipeline AttributeMapping to FrameColumnSet
func mappingToFrameCols(mm []dal.AttributeMapping) FrameColumnSet {
	out := make(FrameColumnSet, 0, len(mm))

	for _, m := range mm {
		out = append(out, mappingToFrameCol(m))
	}

	return out
}

// @note current implementation a bit _rushed_ since I'll probably rethink
//       how the pipeline handles attributes -- will revisit then.
func mappingToFrameCol(m dal.AttributeMapping) FrameColumn {
	p := m.Properties()

	const (
		// Coppied around to reduce imports
		emailLength = 254
		urlLength   = 2048

		attachmentResType = "corteza::system:attachment"
		userResType       = "corteza::system:user"
		moduleResType     = "corteza::compose:module"
	)

	l := m.Properties().Label
	if l == "" {
		l = m.Identifier()
	}
	out := FrameColumn{
		Name:  m.Identifier(),
		Label: l,
		Kind:  "String",

		Primary: p.IsPrimary,
	}

	switch t := p.Type.(type) {
	case *dal.TypeBoolean:
		out.Kind = "Boolean"
	case *dal.TypeDate, *dal.TypeTime, *dal.TypeTimestamp:
		out.Kind = "DateTime"

	case *dal.TypeNumber:
		out.Kind = "Number"

	case *dal.TypeEnum:
		out.Kind = "Select"

	case *dal.TypeText:
		// @note temporary solution; we should push some meta along with it
		if t.Length == emailLength {
			out.Kind = "Email"
		} else if t.Length == urlLength {
			out.Kind = "URL"
		} else {
			out.Kind = "String"
		}

	case *dal.TypeRef:
		switch t.RefModel.ResourceType {
		case moduleResType:
			out.Kind = "Record"
		case userResType:
			out.Kind = "User"
		case attachmentResType:
			out.Kind = "File"
		}
	}

	return out
}

// Report step -> DAL step conversion

// convStepLoad converts ReportStepLoad to dal.Datasource
func convStepLoad(pr modelFinder, step types.ReportStepLoad, defs FrameDefinitionSet) (out *dal.Datasource, err error) {
	// Validation
	if len(defs) > 1 {
		err = fmt.Errorf("cannot convert load step: expecting at most one definition, got %d", len(defs))
		return
	}

	// Get additional filtering
	var extf filter.Filter
	if len(defs) == 1 {
		extf = filterFromDef(defs[0])
	}

	// Prepare model ref
	mfr, err := makeModelRef(step)
	if err != nil {
		return
	}

	f, err := dal.FilterFromExpr(step.Filter.Node()).
		MergeFilters(extf)
	if err != nil {
		return
	}

	// @todo refactor this out after we support other resources with potentially missing soft delete fields
	f, err = f.MergeFilters(filter.Generic(filter.WithStateConstraint("deletedAt", filter.StateExcluded)))
	if err != nil {
		return
	}

	// Make pipeline step
	return &dal.Datasource{
		Ident:         step.Name,
		Filter:        f,
		ModelRef:      mfr,
		OutAttributes: modelAttributes(pr, step, mfr),
	}, nil
}

// convStepAggregate converts ReportStepAggregate to dal.Aggregate
func convStepAggregate(step types.ReportStepAggregate, defs FrameDefinitionSet) (out *dal.Aggregate, err error) {
	// Validation
	if len(defs) > 1 {
		err = fmt.Errorf("cannot convert aggregate step: expecting at most one definition, got %d", len(defs))
		return
	}

	// Get additional filtering
	var extf filter.Filter
	if len(defs) == 1 {
		extf = filterFromDef(defs[0])
	}

	f, err := dal.FilterFromExpr(step.Filter.Node()).
		MergeFilters(extf)
	if err != nil {
		return
	}

	// Make pipeline step
	out = &dal.Aggregate{
		Ident:     step.Name,
		RelSource: step.Source,
		Filter:    f,

		Group:         step.Keys.DalMapping(),
		OutAttributes: step.Columns.DalMapping(),
	}
	return
}

// convStepJoin converts ReportStepJoin to dal.Join
func convStepJoin(step types.ReportStepJoin, defs FrameDefinitionSet) (out *dal.Join, err error) {
	// Validation
	if len(defs) > 1 {
		err = fmt.Errorf("cannot convert join step: expecting at most one definition, got %d", len(defs))
		return
	}

	// Get additional filtering
	var extf filter.Filter
	if len(defs) == 1 {
		extf = filterFromDef(defs[0])
	}

	f, err := dal.FilterFromExpr(step.Filter.Node()).
		MergeFilters(extf)
	if err != nil {
		return
	}

	// Make pipeline step
	out = &dal.Join{
		Ident:    step.Name,
		RelLeft:  step.LocalSource,
		RelRight: step.ForeignSource,

		Filter: f,

		On: dal.JoinPredicate{
			Left:  step.LocalColumn,
			Right: step.ForeignColumn,
		},
	}
	return
}

// convStepLink converts ReportStepLink to dal.Link
func convStepLink(step types.ReportStepLink, defs FrameDefinitionSet) (out *dal.Link, err error) {
	// Validation
	if len(defs) > 2 {
		err = fmt.Errorf("cannot convert join step: expecting at most two definitions, got %d", len(defs))
		return
	}

	// @todo additional filtering; will need to split the dal.Link filter into
	//       left and right for more control and clarity

	// Make pipeline step
	out = &dal.Link{
		Ident:    step.Name,
		RelLeft:  step.LocalSource,
		RelRight: step.ForeignSource,

		On: dal.LinkPredicate{
			Left:  step.LocalColumn,
			Right: step.ForeignColumn,
		},
		Filter: dal.FilterFromExpr(step.Filter.Node()),
	}
	return
}

// updateDefinitionColumns assures run's frame column completeness
func updateDefinitionColumns(r run) {
	ppAttrs := r.Pipeline[0].Attributes()
	for i, def := range r.Defs {
		if len(def.Columns) > 0 {
			continue
		}

		def.Columns = mappingToFrameCols(ppAttrs[i])
	}
}

// makeModelRef returns the model ref based on the step load definition
// @todo should be expanded when we support models that are not compose modules
func makeModelRef(step types.ReportStepLoad) (out dal.ModelRef, err error) {
	var (
		connectionID          uint64
		moduleID, namespaceID uint64
		module, namespace     string

		aux any
		ok  bool
	)

	if aux, ok = step.Definition["moduleID"]; ok {
		moduleID = cast.ToUint64(aux)
	} else if aux, ok = step.Definition["module"]; ok {
		module = cast.ToString(aux)
	} else {
		err = fmt.Errorf("step definition is missing moduleID or module")
		return
	}

	if aux, ok = step.Definition["namespaceID"]; ok {
		namespaceID = cast.ToUint64(aux)
	} else if aux, ok = step.Definition["namespace"]; ok {
		namespace = cast.ToString(aux)
	} else {
		err = fmt.Errorf("step definition is missing namespaceID or namespace")
		return
	}

	// Connection is optional, default is primary connection
	if aux, ok = step.Definition["connectionID"]; ok {
		connectionID = cast.ToUint64(aux)
	}

	out.ConnectionID = connectionID
	out.Refs = make(map[string]any)

	// Use only one of the two identifier variations with priority to ID
	if moduleID > 0 {
		out.Refs["moduleID"] = moduleID
	} else {
		out.Refs["module"] = module
	}

	if namespaceID > 0 {
		out.Refs["namespaceID"] = namespaceID
	} else {
		out.Refs["namespace"] = namespace
	}

	return
}

// modelAttributes returns the attributes defined on the referenced model
func modelAttributes(pr modelFinder, step types.ReportStepLoad, mfr dal.ModelRef) []dal.AttributeMapping {
	// All of the attributes
	attrs, err := getModelAttrs(pr, mfr)
	if err != nil {
		return nil
	}

	return attrToMapping(attrs...)
}

func getModelAttrs(pr modelFinder, mfr dal.ModelRef) (attrs dal.AttributeSet, err error) {
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
				Label:     a.Label,
				Type:      a.Type,
				Nullable:  a.Type.IsNullable(),
				IsPrimary: a.PrimaryKey,
			},
		})
	}
	return
}

func filterFromDef(def *FrameDefinition) (out filter.Filter) {
	aux := filter.Generic(
		filter.WithExpressionParsed(def.Filter.Node()),
		filter.WithOrderBy(def.Sort),
	)

	if def.Paging != nil {
		aux = aux.With(
			filter.WithCursor(def.Paging.PageCursor),
			filter.WithLimit(def.Paging.Limit),
		)
	}

	return aux
}

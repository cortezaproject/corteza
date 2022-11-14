package dal

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	// Join produces a series of joined rows from the provided sources based
	// on the JoinPredicate.
	//
	// The join step produces an SQL left join-like output where all right rows
	// have a corresponding left row.
	Join struct {
		Ident    string
		RelLeft  string
		RelRight string
		// @todo allow multiple join predicates; for now (for easier indexing)
		// only allow one (this is the same as we had before)
		On     JoinPredicate
		Filter filter.Filter
		filter internalFilter

		OutAttributes   []AttributeMapping
		LeftAttributes  []AttributeMapping
		RightAttributes []AttributeMapping

		relLeft  PipelineStep
		relRight PipelineStep
		plan     joinPlan
		analysis map[string]OpAnalysis
	}

	// JoinPredicate determines the attributes the two datasets should get joined on
	JoinPredicate struct {
		Left  string
		Right string
	}

	// joinPlan outlines how the optimizer determined the two datasets should be
	// joined on.
	joinPlan struct {
		// @todo add strategy when we have different strategies implemented
		// strategy string

		// partialScan indicates we can partially pull data from the two sources
		// as the data is provided in the correct order.
		partialScan bool
	}
)

func (def *Join) Identifier() string {
	return def.Ident
}

func (def *Join) Sources() []string {
	return []string{def.RelLeft, def.RelRight}
}

func (def *Join) Attributes() [][]AttributeMapping {
	return [][]AttributeMapping{def.OutAttributes}
}

func (def *Join) Analyze(ctx context.Context) (err error) {
	// @todo proper analysis; for now we'll leave this as defaults
	def.analysis = map[string]OpAnalysis{
		OpAnalysisIterate: {
			ScanCost:   CostUnknown,
			SearchCost: CostUnknown,
			FilterCost: CostUnknown,
			SortCost:   CostUnknown,
			OutputSize: SizeUnknown,
		},
	}

	return
}

func (def *Join) Analysis() map[string]OpAnalysis {
	return def.analysis
}

func (def *Join) Optimize(req internalFilter) (res internalFilter, err error) {
	err = fmt.Errorf("not implemented")
	return
}

// iterator initializes an iterator based on the provided pipeline step definition
func (def *Join) iterator(ctx context.Context, left, right Iterator) (out Iterator, err error) {
	exec, err := def.init(ctx, left, right)
	if err != nil {
		return
	}

	// Get pred. type
	// @todo should validate that both pred. types are the same/compatible
	var pt Type
	for _, a := range def.LeftAttributes {
		if a.Identifier() == def.On.Left {
			pt = a.Properties().Type
			break
		}
	}

	return exec, exec.init(ctx, pt)
}

// dryrun performs step execution without interacting with the data
// @todo consider rewording this
func (def *Join) dryrun(ctx context.Context) (err error) {
	_, err = def.init(ctx, nil, nil)
	return
}

func (def *Join) init(ctx context.Context, left, right Iterator) (exec *joinLeft, err error) {
	exec = &joinLeft{
		leftSource:  left,
		rightSource: right,
	}

	// Convert the provided filter into an internal filter
	if def.Filter != nil {
		def.filter, err = toInternalFilter(def.Filter)
		if err != nil {
			return
		}
	}

	// Collect attributes from the underlaying step in case own are not provided
	if len(def.LeftAttributes) == 0 {
		def.LeftAttributes = collectAttributes(def.relLeft)
	}
	if len(def.RightAttributes) == 0 {
		def.RightAttributes = collectAttributes(def.relRight)
	}

	// When the output attributes are not provided, we determine them based on the input
	// attributes.
	//
	// The join step prefixes each attribute with the identifier of the source they came from.
	// This avoids name collisions.
	// @todo consider improving this to only prefix when there is a name collision
	//       and allow the nested attr. to be more permissive; for example, instead of
	//       requiring a.b.c only b.c or c would be enough.
	if len(def.OutAttributes) == 0 {
		ins := func(srcIdent string, aa []AttributeMapping) {
			for _, a := range aa {
				p := a.Properties()
				def.OutAttributes = append(def.OutAttributes, SimpleAttr{
					Ident: fmt.Sprintf("%s%s%s", srcIdent, attributeNestingSeparator, a.Identifier()),
					Expr:  a.Expression(),
					Src:   a.Identifier(),
					Props: p,
				})
			}
		}

		def.OutAttributes = make([]AttributeMapping, 0, len(def.LeftAttributes)+len(def.RightAttributes))
		ins(def.RelLeft, def.LeftAttributes)
		ins(def.RelRight, def.RightAttributes)
	}

	// Assure and attempt to correct the provided sort to conform with the data set and the
	// paging cursor (if any)
	def.filter, err = assureSort(def.filter, exec.collectPrimaryAttributes(def.OutAttributes))
	if err != nil {
		return
	}

	// Index attrs for validations
	leftSrcAttrs := indexAttrs(def.LeftAttributes...)
	rightSrcAttrs := indexAttrs(def.RightAttributes...)
	outAttrs := indexAttrs(def.OutAttributes...)

	for _, a := range def.OutAttributes {
		if !leftSrcAttrs[a.Source()] && !rightSrcAttrs[a.Source()] {
			return nil, fmt.Errorf("unknown attribute %s", a.Source())
		}
	}

	// Generic validation
	if !leftSrcAttrs[def.On.Left] {
		return nil, fmt.Errorf("unknown join predicate attribute %s", def.On.Left)
	}
	if !rightSrcAttrs[def.On.Right] {
		return nil, fmt.Errorf("unknown join predicate attribute %s", def.On.Right)
	}

	if len(def.OutAttributes) == 0 {
		return nil, fmt.Errorf("no attributes specified")
	}
	if len(def.LeftAttributes) == 0 {
		return nil, fmt.Errorf("no left attributes specified")
	}
	if len(def.RightAttributes) == 0 {
		return nil, fmt.Errorf("no right attributes specified")
	}

	if def.On.Left == "" {
		return nil, fmt.Errorf("no left attribute in the join predicate specified")
	}
	if def.On.Right == "" {
		return nil, fmt.Errorf("no right attribute in the join predicate specified")
	}

	// order
	for _, s := range def.filter.OrderBy() {
		if _, ok := outAttrs[s.Column]; !ok {
			return nil, fmt.Errorf("order attribute %s does not exist", s.Column)
		}
	}

	exec.filter = def.filter
	exec.def = *def
	return
}

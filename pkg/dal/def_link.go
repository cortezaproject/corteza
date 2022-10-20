package dal

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	// Link produces a series of left and corresponding right rows based on the
	// provided sources and the LinkPredicate.
	//
	// The Link step produces an SQL left join-like output where left and right
	// rows are served separately and the left rows are not duplicated.
	Link struct {
		Ident    string
		RelLeft  string
		RelRight string
		// @todo allow multiple link predicates; for now (for easier indexing)
		// only allow one (this is the same as we had before)
		On LinkPredicate
		// @todo consider splitting filter into left and right filter
		filter internalFilter
		Filter filter.Filter

		OutLeftAttributes  []AttributeMapping
		OutRightAttributes []AttributeMapping
		LeftAttributes     []AttributeMapping
		RightAttributes    []AttributeMapping

		relLeft  PipelineStep
		relRight PipelineStep

		plan     linkPlan
		analysis map[string]OpAnalysis
	}

	// LinkPredicate determines the attributes the two datasets should get joined on
	LinkPredicate struct {
		Left  string
		Right string
	}

	// linkPlan outlines how the optimizer determined the two datasets should be
	// joined on.
	linkPlan struct {
		// @todo add strategy when we have different strategies implemented
		// strategy string

		// partialScan indicates we can partially pull data from the two sources
		// as the data is provided in the correct order.
		partialScan bool
	}

	rowLink struct {
		a, b ValueGetter
	}
)

const (
	LinkRefIdent = "$sys.ref"
)

func (def *Link) Identifier() string {
	return def.Ident
}

func (def *Link) Sources() []string {
	return []string{def.RelLeft, def.RelRight}
}

func (def *Link) Attributes() [][]AttributeMapping {
	return [][]AttributeMapping{def.OutLeftAttributes, def.OutRightAttributes}
}

func (def *Link) Analyze(ctx context.Context) (err error) {
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

func (def *Link) Analysis() map[string]OpAnalysis {
	return def.analysis
}

func (def *Link) Optimize(req internalFilter) (res internalFilter, err error) {
	err = fmt.Errorf("not implemented")
	return
}

// iterator initializes an iterator based on the provided pipeline step definition
func (def *Link) iterator(ctx context.Context, left, right Iterator) (out Iterator, err error) {
	x, err := def.init(ctx, left, right)
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

	// Collect order by attributes which appear in the right source
	rightAttrs := indexAttrs(def.RightAttributes...)
	rightOrderAttrs := make([]string, 0, len(def.filter.OrderBy()))
	for _, o := range def.filter.OrderBy() {
		if rightAttrs[o.Column] {
			rightOrderAttrs = append(rightOrderAttrs, o.Column)
		}
	}

	return x, x.init(ctx, pt, rightOrderAttrs)
}

// dryrun performs step execution without interacting with the data
// @todo consider rewording this
func (def *Link) dryrun(ctx context.Context) (err error) {
	_, err = def.init(ctx, nil, nil)
	return
}

func (def *Link) init(ctx context.Context, left, right Iterator) (exec *linkLeft, err error) {
	exec = &linkLeft{
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

	// @todo this isn't quite ok -- the ident of left/right must become the src of the out.
	//       An edge-case but it should be covered.
	if len(def.OutLeftAttributes) == 0 {
		def.OutLeftAttributes = def.LeftAttributes
	}
	if len(def.OutRightAttributes) == 0 {
		def.OutRightAttributes = def.RightAttributes
	}

	// Index source attributes for group/aggregate definition validation
	leftSrcAttrs := indexAttrs(def.LeftAttributes...)
	rightSrcAttrs := indexAttrs(def.RightAttributes...)
	outAttrs := indexAttrs(def.OutLeftAttributes...)
	indexAttrsInto(outAttrs, def.OutRightAttributes...)

	// Check if source attributes are present
	for _, a := range append(def.OutLeftAttributes, def.OutRightAttributes...) {
		if !leftSrcAttrs[a.Source()] && !rightSrcAttrs[a.Source()] {
			return nil, fmt.Errorf("unknown attribute %s", a.Source())
		}
	}

	// Check link predicates
	if !leftSrcAttrs[def.On.Left] {
		return nil, fmt.Errorf("left link predicate %s does not exist", def.On.Left)
	}
	if !rightSrcAttrs[def.On.Right] {
		return nil, fmt.Errorf("right link predicate %s does not exist", def.On.Right)
	}

	// General validation
	if len(def.OutLeftAttributes) == 0 {
		return nil, fmt.Errorf("no left output attributes specified")
	}
	if len(def.OutRightAttributes) == 0 {
		return nil, fmt.Errorf("no right output attributes specified")
	}
	if len(def.LeftAttributes) == 0 {
		return nil, fmt.Errorf("no left attributes specified")
	}
	if len(def.RightAttributes) == 0 {
		return nil, fmt.Errorf("no right attributes specified")
	}

	if def.On.Left == "" {
		return nil, fmt.Errorf("no left attribute in the link predicate specified")
	}
	if def.On.Right == "" {
		return nil, fmt.Errorf("no right attribute in the link predicate specified")
	}

	for _, s := range def.filter.OrderBy() {
		if _, ok := outAttrs[s.Column]; !ok {
			return nil, fmt.Errorf("order attribute %s does not exist", s.Column)
		}
	}

	exec.filter = def.filter
	exec.def = *def
	return
}

func (r *rowLink) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	return r.GetValue(k, 0)
}

func (r *rowLink) GetValue(name string, pos uint) (v any, err error) {
	a := r.a.CountValues()
	if cc, ok := a[name]; ok {
		if pos >= cc {
			return nil, nil
		}
		return r.a.GetValue(name, pos)
	}

	b := r.b.CountValues()
	if cc, ok := b[name]; ok {
		if pos >= cc {
			return nil, nil
		}
		return r.b.GetValue(name, pos)
	}

	return
}

func (r *rowLink) CountValues() (out map[string]uint) {
	out = make(map[string]uint)

	for k, c := range r.a.CountValues() {
		out[k] = c
	}
	for k, c := range r.b.CountValues() {
		out[k] = c
	}

	return
}

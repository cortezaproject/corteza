package dal

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/ql"
)

type (
	// Aggregate produces a series of aggregated rows from the provided sources based
	// on the specified group.
	Aggregate struct {
		Ident     string
		RelSource string
		Filter    filter.Filter
		filter    internalFilter

		Group         []AttributeMapping
		OutAttributes []AttributeMapping

		SourceAttributes []AttributeMapping

		rel      PipelineStep
		plan     aggregatePlan
		analysis stepAnalysis
	}

	// aggregatePlan outlines how the optimizer determined the dataset should be
	// aggregated
	aggregatePlan struct {
		// partialScan indicates we can partially pull data from the two sources
		// as the data is provided in the correct order.
		partialScan bool
	}

	groupKey       []any
	aggregateGroup struct {
		key groupKey
		agg *aggregator
	}

	// aggregateAttr is a simple wrapper to outline aggregated attribute definitions
	aggregateAttr struct {
		ident    string
		expr     *ql.ASTNode
		attrType Type
	}
)

func (def *Aggregate) Identifier() string {
	return def.Ident
}

func (def *Aggregate) Sources() []string {
	return []string{def.RelSource}
}

func (def *Aggregate) Attributes() [][]AttributeMapping {
	return [][]AttributeMapping{append(def.Group, def.OutAttributes...)}
}

func (def *Aggregate) Analyze(ctx context.Context) (err error) {
	// @todo proper analysis; for now we'll leave this as defaults
	def.analysis = stepAnalysis{
		scanCost:   costUnknown,
		searchCost: costUnknown,
		filterCost: costUnknown,
		sortCost:   costUnknown,
		outputSize: sizeUnknown,
	}
	return
}

func (def *Aggregate) Analysis() stepAnalysis {
	return def.analysis
}

func (def *Aggregate) Optimize(reqFilter internalFilter) (rspFilter internalFilter, err error) {
	err = fmt.Errorf("not implemented")
	return
}

// iterator initializes an iterator based on the provided pipeline step definition
func (def *Aggregate) iterator(ctx context.Context, src Iterator) (out Iterator, err error) {
	exec, err := def.init(ctx, src)
	if err != nil {
		return
	}

	return exec, exec.init(ctx)
}

// dryrun performs step execution without interacting with the data
// @todo consider rewording this
func (def *Aggregate) dryrun(ctx context.Context) (err error) {
	_, err = def.init(ctx, nil)
	return
}

func (def *Aggregate) init(ctx context.Context, src Iterator) (exec *aggregate, err error) {
	exec = &aggregate{
		source: src,
	}

	// Convert the provided filter into an internal filter
	if def.Filter != nil {
		def.filter, err = toInternalFilter(def.Filter)
		if err != nil {
			return
		}
	}

	// Collect attributes from the underlaying step in case own are not provided
	if len(def.SourceAttributes) == 0 {
		def.SourceAttributes = collectAttributes(def.rel)
	}

	// Index source attributes for group/aggregate definition validation
	srcAttrs := indexAttrs(def.SourceAttributes...)
	pp := newQlParser(func(ident ql.Ident) (_ ql.Ident, err error) {
		if _, ok := srcAttrs[ident.Value]; !ok {
			return ident, fmt.Errorf("unknown attribute %s", ident.Value)
		}
		return ident, nil
	})

	// Convert & validate group definitions
	// - groups
	var gd aggregateAttr
	outAttrs := make(map[string]bool, len(def.Group)+len(def.OutAttributes))
	for _, attr := range def.Group {
		idtf := attr.Identifier()
		gd, err = aggregateAttrFromExpr(pp, attr.Properties().Type, idtf, attr.Expression())
		if err != nil {
			return
		}
		exec.groupDefs = append(exec.groupDefs, gd)
		outAttrs[idtf] = true
	}
	// - aggregates
	for _, attr := range def.OutAttributes {
		idtf := attr.Identifier()
		gd, err = aggregateAttrFromExpr(pp, attr.Properties().Type, idtf, attr.Expression())
		if err != nil {
			return
		}
		exec.aggregateDefs = append(exec.aggregateDefs, gd)
		outAttrs[idtf] = true
	}

	// Generic validation
	if len(def.Group) == 0 {
		err = fmt.Errorf("no group attributes specified")
		return
	}

	if len(def.OutAttributes) == 0 {
		err = fmt.Errorf("no output attributes specified")
		return
	}

	if len(def.SourceAttributes) == 0 {
		err = fmt.Errorf("no source attributes specified")
		return
	}

	// order
	for _, s := range def.filter.OrderBy() {
		if _, ok := outAttrs[s.Column]; !ok {
			err = fmt.Errorf("order by attribute %s does not exist", s.Column)
			return
		}
	}

	// Finishup
	exec.filter = def.filter
	exec.def = *def
	return
}

func aggregateAttrFromExpr(pp *ql.Parser, t Type, ident, expr string) (out aggregateAttr, err error) {
	n, err := pp.Parse(expr)
	if err != nil {
		return
	}

	return aggregateAttr{
		ident:    ident,
		expr:     n,
		attrType: t,
	}, nil
}

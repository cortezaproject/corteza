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

		Group         []AggregateAttr
		OutAttributes []AggregateAttr

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

	// AggregateAttr is a simple wrapper to outline aggregated attribute definitions
	AggregateAttr struct {
		Key bool

		// @todo change; temporary for compose service
		RawExpr string

		Identifier string
		Label      string
		Expression *ql.ASTNode
		Type       Type
		Store      Codec
	}
)

func (def *Aggregate) Identifier() string {
	return def.Ident
}

func (def *Aggregate) Sources() []string {
	return []string{def.RelSource}
}

func (def *Aggregate) Attributes() [][]AttributeMapping {
	aa := append(def.Group, def.OutAttributes...)
	out := make([]AttributeMapping, 0, len(aa))

	for _, a := range aa {
		out = append(out, a.toSimpleAttr())
	}

	return [][]AttributeMapping{out}
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

	prepAttr := func(attr AggregateAttr) (_ AggregateAttr, err error) {
		if attr.RawExpr != "" {
			// Parse (it already validates) raw expressions
			attr.Expression, err = pp.Parse(attr.RawExpr)
			if err != nil {
				return
			}
		} else {
			// Manually validate already parsed expressions
			err = attr.Expression.Traverse(func(a *ql.ASTNode) (bool, *ql.ASTNode, error) {
				if a.Symbol == "" {
					return true, a, nil
				}
				if _, ok := srcAttrs[a.Symbol]; !ok {
					return false, nil, fmt.Errorf("unknown attribute %s", a.Symbol)
				}

				return true, a, nil
			})
		}
		if err != nil {
			return
		}
		return def.determineAttrType(attr, def.SourceAttributes)
	}

	// Convert & validate group definitions
	// - groups
	outAttrs := make(map[string]bool, len(def.Group)+len(def.OutAttributes))
	for i, attr := range def.Group {
		attr, err = prepAttr(attr)
		if err != nil {
			return
		}

		def.Group[i] = attr
		idtf := attr.Identifier

		exec.groupDefs = append(exec.groupDefs, attr)
		outAttrs[idtf] = true
	}
	// - aggregates
	for i, attr := range def.OutAttributes {
		attr, err = prepAttr(attr)
		if err != nil {
			return
		}

		def.OutAttributes[i] = attr
		idtf := attr.Identifier

		exec.aggregateDefs = append(exec.aggregateDefs, attr)
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

// determineAttrType determines the type of the AggregateAttr based on it's definition
// and source attributes
func (def *Aggregate) determineAttrType(base AggregateAttr, ss []AttributeMapping) (out AggregateAttr, err error) {
	out = base
	if out.Type != nil {
		return
	}

	var root *ql.ASTNode
	var t Type

	// If we have a symbol, then we'll use it to determine the type.
	// All current operations should return the same output type as the input one.
	//
	// In case of a function, use the output type of the root most function which has
	// a known type.
	//
	// Note, some refs (group and add for example) may not know their types so we need
	// to dig deeper.
	base.Expression.Traverse(func(a *ql.ASTNode) (bool, *ql.ASTNode, error) {
		if a.Symbol != "" {
			root = a
			return false, a, nil
		}

		if a.Ref != "" {
			tmp := refToGvalExp[a.Ref]
			if tmp == nil || tmp.OutType == nil || tmp.OutTypeUnknown {
				return true, a, nil
			}

			if tmp.OutType != nil {
				t = tmp.OutType
				return false, a, nil
			}
		}
		return true, a, nil
	})

	if root != nil {
		for _, s := range ss {
			if s.Identifier() == root.Symbol {
				t = s.Properties().Type
				break
			}
		}
	}

	out.Type = t
	return
}

func (a AggregateAttr) toSimpleAttr() SimpleAttr {
	return SimpleAttr{
		Ident: a.Identifier,
		Props: MapProperties{
			Label:     a.Label,
			IsPrimary: a.Key,
			Type:      a.Type,
		},

		// @todo won't matter for now
		Expr: "",
		Src:  "",
	}
}

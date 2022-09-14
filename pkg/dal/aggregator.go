package dal

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/modern-go/reflect2"
	"github.com/spf13/cast"
)

type (
	// aggregator performs value aggregation primarily used for the aggregate step
	// @todo consider moving into a separate package
	//
	// The aggregator performs the requested value expressions as well as
	// aggregation operations over the evaluated values.
	//
	// The aggregator computes the values on the fly (for the ones it can).
	aggregator struct {
		// aggregates holds output aggregation values
		// @note we'll use float64 for all values, but it might make more sense to split it up
		//       in the future. For now, it'll be ok.
		aggregates []float64

		// counts holds the number of values for each aggregate including multi value fields.
		// Counts are currently only used for average.
		counts []int

		// def provides a list of aggregates.
		//
		// Eac index corresponds to the aggregates and counts slices
		def []aggregateDef
		// scanned indicates whether the aggregator has been scanned since
		// we need to block writes after the first scan.
		scanned bool
	}

	// aggregateDef is a wrapper to outline some aggregate value
	aggregateDef struct {
		outIdent string

		aggOp string

		inIdent string
		eval    evaluator
	}
)

var (
	// aggregateFunctionIndex specifies all of the registered aggregate functions
	//
	// @todo consider making this expandable via some registry/plugin/...
	aggregateFunctionIndex = map[string]bool{
		"count": true,
		"sum":   true,
		"min":   true,
		"max":   true,
		"avg":   true,
	}
)

// Aggregator initializes a new aggregator for the given set of mappings
//
// The aggregator is not routine safe; consider defining multiple aggregators
// and then combining them together.
func Aggregator(attributes ...AttributeMapping) (_ *aggregator, err error) {
	// Convert attribute mappings into an internal structure
	// @todo consider dropping this internal structure; not sure how much I like it
	//       but I think having something in between is best.
	attrs := make([]aggregateDef, len(attributes))

	for i, a := range attributes {
		attrs[i], err = mappingToAggregateDef(a)
		if err != nil {
			return
		}
	}

	return &aggregator{
		def: attrs,

		// init appropriate slices
		aggregates: make([]float64, len(attributes)),
		counts:     make([]int, len(attributes)),
	}, nil
}

// Aggregate aggregates the given value
func (a *aggregator) Aggregate(ctx context.Context, v ValueGetter) (err error) {
	if a.scanned {
		// If we attempt to add data to the aggregator, the previous data might
		// no longer be valid.
		return fmt.Errorf("cannot call Aggregate on an already scanned aggregator")
	}

	// @todo consider throwing all of the aggregates into a routine.
	//       I think (for now) performance gains will be negligible.
	for i, attr := range a.def {
		err = a.aggregate(ctx, attr, i, v)
		if err != nil {
			return
		}
	}

	return nil
}

// Scan scans the aggregated values into the setter
func (a *aggregator) Scan(s ValueSetter) (err error) {
	// On first scan, complete partial aggregates
	if !a.scanned {
		a.completePartials()
	}

	a.scanned = true

	// Set the values
	for i, attr := range a.def {
		// @note each aggregated value can be at most one so no need for multi-value
		//       suport here.
		err = s.SetValue(attr.outIdent, 0, a.aggregates[i])
		if err != nil {
			return
		}
	}

	return
}

// aggregate applies the provided value into the requested aggregate
func (a *aggregator) aggregate(ctx context.Context, attr aggregateDef, i int, v ValueGetter) (err error) {
	switch attr.aggOp {
	case "count":
		return a.count(ctx, attr, i, v)

	case "sum":
		return a.sum(ctx, attr, i, v)

	case "min":
		return a.min(ctx, attr, i, v)

	case "max":
		return a.max(ctx, attr, i, v)

	case "avg":
		return a.avg(ctx, attr, i, v)
	}

	return fmt.Errorf("unsupported aggregate function: %s", attr.aggOp)
}

// walkValues traverses the available values for the specified attribute
func (a *aggregator) walkValues(ctx context.Context, r ValueGetter, cc map[string]uint, attr aggregateDef, run func(v any, isNil bool)) (err error) {
	var out any
	if attr.inIdent == "" {
		out, err = attr.eval.Eval(ctx, r)
		if err != nil {
			return
		}

		run(out, reflect2.IsNil(out))
		return nil
	}

	for i := uint(0); i < cc[attr.inIdent]; i++ {
		v, err := r.GetValue(attr.inIdent, i)
		if err != nil {
			return err
		}
		run(v, reflect2.IsNil(v))
	}

	return nil
}

// Aggregate methods

func (a *aggregator) count(ctx context.Context, attr aggregateDef, i int, v ValueGetter) (err error) {
	err = a.walkValues(ctx, v, v.CountValues(), attr, func(v any, isNil bool) {
		if isNil {
			return
		}

		a.aggregates[i]++
		a.counts[i]++
	})
	if err != nil {
		return
	}

	return
}

func (a *aggregator) sum(ctx context.Context, attr aggregateDef, i int, v ValueGetter) (err error) {
	err = a.walkValues(ctx, v, v.CountValues(), attr, func(v any, isNil bool) {
		if isNil {
			return
		}
		a.aggregates[i] += cast.ToFloat64(v)
		a.counts[i]++
	})
	if err != nil {
		return
	}

	return
}

func (a *aggregator) min(ctx context.Context, attr aggregateDef, i int, v ValueGetter) (err error) {
	err = a.walkValues(ctx, v, v.CountValues(), attr, func(v any, isNil bool) {
		if isNil {
			return
		}

		if a.counts[i] == 0 {
			a.aggregates[i] = cast.ToFloat64(v)
		} else {
			a.aggregates[i] = math.Min(a.aggregates[i], cast.ToFloat64(v))
		}
		a.counts[i]++
	})
	if err != nil {
		return
	}

	return
}

func (a *aggregator) max(ctx context.Context, attr aggregateDef, i int, v ValueGetter) (err error) {
	err = a.walkValues(ctx, v, v.CountValues(), attr, func(v any, isNil bool) {
		if isNil {
			return
		}

		if a.counts[i] == 0 {
			a.aggregates[i] = cast.ToFloat64(v)
		} else {
			a.aggregates[i] = math.Max(a.aggregates[i], cast.ToFloat64(v))
		}
		a.counts[i]++
	})
	if err != nil {
		return
	}

	return
}

func (a *aggregator) avg(ctx context.Context, attr aggregateDef, i int, v ValueGetter) (err error) {
	err = a.walkValues(ctx, v, v.CountValues(), attr, func(v any, isNil bool) {
		if isNil {
			return
		}

		a.aggregates[i] += cast.ToFloat64(v)

		a.counts[i]++
	})
	if err != nil {
		return
	}

	return
}

func (a *aggregator) completePartials() {
	a.completeAverage()
}

func (a *aggregator) completeAverage() {
	for i, attr := range a.def {
		if attr.aggOp == "avg" {
			if a.counts[i] == 0 {
				return
			}
			a.aggregates[i] = a.aggregates[i] / float64(a.counts[i])
		}
	}
}

// Utilities

// mappingToAggregateDef constructs an aggregate definition from the provided mapping
// @note this will probably change when I change this AttributeMapping thing
func mappingToAggregateDef(a AttributeMapping) (def aggregateDef, err error) {
	def = aggregateDef{
		outIdent: a.Identifier(),
	}

	inIdent, expr, err := unpackMappingSource(a)
	if err != nil {
		return
	}

	// Take it from the source
	if inIdent != "" {
		def.inIdent = inIdent
		return
	}

	// Take it from the expression
	// - agg. op.
	def.aggOp, expr, err = unpackExpressionNode(expr)
	if err != nil {
		return
	}
	// - make evaluator
	def.eval, err = newRunnerGvalParsed(expr)
	if err != nil {
		return
	}

	return
}

func unpackMappingSource(a AttributeMapping) (ident string, expr *ql.ASTNode, err error) {
	base := a.Expression()
	if base == "" {
		ident = a.Identifier()
		return
	}

	// @note the converter is being reused so this is ok
	expr, err = newConverterGval().Parse(base)
	return
}

func unpackExpressionNode(n *ql.ASTNode) (aggOp string, expr *ql.ASTNode, err error) {
	// @todo check for supported aggregators
	if n.Ref != "" {
		aggOp = strings.ToLower(n.Ref)
	}
	if !aggregateFunctionIndex[aggOp] {
		err = fmt.Errorf("root expression must be an aggregate function")
		return
	}

	if len(n.Args) != 1 {
		err = fmt.Errorf("impossible state: aggregate function must have exactly one argument")
		return
	}

	expr = n.Args[0]
	return
}

// reset is a benchmarking utility to reset the aggregator
//
// Don't use it in production code.
func (a *aggregator) reset() {
	for i := 0; i < len(a.aggregates); i++ {
		a.aggregates[i] = 0
		a.counts[i] = 0
	}
	a.scanned = false
}

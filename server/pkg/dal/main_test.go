package dal

import (
	"context"
	"testing"
)

type (
	simpleAttribute struct {
		ident      string
		source     string
		expr       string
		primary    bool
		multivalue bool
		t          Type
	}

	simpleRow map[string]any
)

func saToMapping(sa ...simpleAttribute) []AttributeMapping {
	out := make([]AttributeMapping, 0, len(sa))
	for _, a := range sa {
		out = append(out, a)
	}

	return out
}

func saToAggAttr(sa ...simpleAttribute) []AggregateAttr {
	out := make([]AggregateAttr, 0, len(sa))
	for _, a := range sa {
		aux := AggregateAttr{
			RawExpr:    a.expr,
			Identifier: a.ident,
			Type:       a.t,
		}
		if aux.RawExpr == "" {
			aux.RawExpr = a.source
		}
		if aux.RawExpr == "" {
			aux.RawExpr = a.ident
		}
		out = append(out, aux)
	}

	return out
}

func (sa simpleAttribute) Properties() MapProperties {
	return MapProperties{
		IsPrimary:    sa.primary,
		Type:         sa.t,
		IsMultivalue: sa.multivalue,
	}
}

func (sa simpleAttribute) Identifier() (ident string) {
	return sa.ident
}

func (sa simpleAttribute) Source() (expr string) {
	if sa.source != "" {
		return sa.source
	}

	return sa.ident
}

func (sa simpleAttribute) Expression() (expr string) {
	if sa.expr != "" {
		return sa.expr
	}

	if sa.source != "" {
		return sa.source
	}

	return sa.ident
}

func bootstrapAggregate(rootTest *testing.T, run func(context.Context, *testing.T, *Aggregate, Buffer)) {
	bootstrapAggregateNoOpt(rootTest, run)
}

func bootstrapAggregateNoOpt(rootTest *testing.T, run func(context.Context, *testing.T, *Aggregate, Buffer)) {
	rootTest.Run("no optimization", func(t *testing.T) {
		ctx := context.Background()
		buff := InMemoryBuffer()
		agg := &Aggregate{}

		run(ctx, t, agg, buff)
	})
}

func (r simpleRow) CountValues() map[string]uint {
	out := make(map[string]uint)

	for k := range r {
		out[k]++
	}

	return out
}

func (r simpleRow) GetValue(k string, place uint) (any, error) {
	return r[k], nil
}

func (r simpleRow) SetValue(k string, place uint, v any) error {
	r[k] = v
	return nil
}

package dal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAggregator(t *testing.T) {
	tcc := []struct {
		name string

		rows       []simpleRow
		attrubutes []simpleAttribute
		out        simpleRow
	}{
		// Plain operations
		{
			name: "simple count",
			rows: []simpleRow{
				{"v": 1},
				{"v": 5},
				{"v": 35},
				{"v": 11.5},
			},
			out: simpleRow{"count": float64(4)},
			attrubutes: []simpleAttribute{
				{ident: "count", expr: "count(v)"},
			},
		},
		{
			name: "simple count nulls",
			rows: []simpleRow{
				{"v": nil},
				{"v": 5},
				{"v": 35},
				{"v": 11.5},
			},
			out: simpleRow{"count": float64(3)},
			attrubutes: []simpleAttribute{
				{ident: "count", expr: "count(v)"},
			},
		},
		{
			name: "simple sum",
			rows: []simpleRow{
				{"v": 1},
				{"v": 5},
				{"v": 35},
				{"v": 11.5},
			},
			out: simpleRow{"sum": float64(52.5)},
			attrubutes: []simpleAttribute{
				{ident: "sum", expr: "sum(v)"},
			},
		},
		{
			name: "simple sum nulls",
			rows: []simpleRow{
				{"v": nil},
				{"v": 5},
				{"v": 35},
				{"v": 11.5},
			},
			out: simpleRow{"sum": float64(51.5)},
			attrubutes: []simpleAttribute{
				{ident: "sum", expr: "sum(v)"},
			},
		},
		{
			name: "simple min",
			rows: []simpleRow{
				{"v": 1},
				{"v": 5},
				{"v": 35},
				{"v": 11.5},
			},
			out: simpleRow{"min": float64(1)},
			attrubutes: []simpleAttribute{
				{ident: "min", expr: "min(v)"},
			},
		},
		{
			name: "simple min nulls",
			rows: []simpleRow{
				{"v": nil},
				{"v": 5},
				{"v": 35},
				{"v": 11.5},
			},
			out: simpleRow{"min": float64(5)},
			attrubutes: []simpleAttribute{
				{ident: "min", expr: "min(v)"},
			},
		},
		{
			name: "simple max",
			rows: []simpleRow{
				{"v": 1},
				{"v": 5},
				{"v": 35},
				{"v": 11.5},
			},
			out: simpleRow{"max": float64(35)},
			attrubutes: []simpleAttribute{
				{ident: "max", expr: "max(v)"},
			},
		},
		{
			name: "simple max nulls",
			rows: []simpleRow{
				{"v": 1},
				{"v": 5},
				{"v": nil},
				{"v": 11.5},
			},
			out: simpleRow{"max": float64(11.5)},
			attrubutes: []simpleAttribute{
				{ident: "max", expr: "max(v)"},
			},
		},
		{
			name: "simple avg",
			rows: []simpleRow{
				{"v": 1},
				{"v": 5},
				{"v": 35},
				{"v": 11.5},
			},
			out: simpleRow{"avg": float64(13.125)},
			attrubutes: []simpleAttribute{
				{ident: "avg", expr: "avg(v)"},
			},
		},
		{
			name: "simple avg nulls",
			rows: []simpleRow{
				{"v": nil},
				{"v": 5},
				{"v": 35},
				{"v": nil},
			},
			out: simpleRow{"avg": float64(20)},
			attrubutes: []simpleAttribute{
				{ident: "avg", expr: "avg(v)"},
			},
		},

		// With a nested expression
		// @todo tests to assure nil values; omitting due to the gval issue
		{
			name: "nested expression",
			rows: []simpleRow{
				{"v": 1},
				{"v": 5},
				{"v": 35},
				{"v": 11.5},
			},
			out: simpleRow{"sum": float64(60.5)},
			attrubutes: []simpleAttribute{
				{ident: "sum", expr: "sum(v + 2)"},
			},
		},
	}

	ctx := context.Background()
	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			agg := Aggregator()
			for _, a := range tc.attrubutes {
				require.NoError(t, agg.AddAggregateE(a.ident, a.expr))
			}

			for _, r := range tc.rows {
				require.NoError(t, agg.Aggregate(ctx, r))
			}

			out := make(simpleRow)
			require.NoError(t, agg.Scan(out))
			require.Equal(t, tc.out, out)
		})
	}
}

func TestAggregatorInit(t *testing.T) {
	t.Run("non supported agg. op.", func(t *testing.T) {
		agg := Aggregator()
		require.Error(t, agg.AddAggregateE("count", "div(v)"))
	})

	t.Run("invalid expression", func(t *testing.T) {
		agg := Aggregator()
		require.Error(t, agg.AddAggregateE("x", "sum(q we)"))
	})
}

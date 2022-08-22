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

		// With a nested expression
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
			agg, err := Aggregator(saToMapping(tc.attrubutes...)...)
			require.NoError(t, err)

			for _, r := range tc.rows {
				require.NoError(t, agg.Aggregate(ctx, r))
			}

			out := make(simpleRow)
			err = agg.Scan(out)
			require.NoError(t, err)
			require.Equal(t, tc.out, out)
		})
	}
}

func TestAggregatorInit(t *testing.T) {
	t.Run("non supported agg. op.", func(t *testing.T) {
		_, err := Aggregator(
			simpleAttribute{
				expr: "div(v)",
			},
		)
		require.Error(t, err)
	})

	t.Run("invalid expression", func(t *testing.T) {
		_, err := Aggregator(
			simpleAttribute{
				expr: "sum(q we)",
			},
		)
		require.Error(t, err)
	})
}

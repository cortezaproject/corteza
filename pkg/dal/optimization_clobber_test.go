package dal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func makeAnalysisDsAggregate() map[string]OpAnalysis {
	return map[string]OpAnalysis{
		OpAnalysisAggregate: {},
	}
}

func TestClobberStep(t *testing.T) {
	t.Run("no steps", func(t *testing.T) {
		out, err := pipelineClobberSteps(nil)

		require.NoError(t, err)
		require.Nil(t, out)
	})

	t.Run("one step no optimize", func(t *testing.T) {
		out, err := pipelineClobberSteps(Pipeline{&Datasource{}})

		require.NoError(t, err)
		require.Len(t, out, 1)
	})

	t.Run("agg ds", func(t *testing.T) {
		ds := &Datasource{
			Ident:    "ds_1",
			analysis: makeAnalysisDsAggregate(),
		}
		agg := &Aggregate{
			Ident:     "agg_1",
			RelSource: "ds_1",
			rel:       ds,
		}

		out, err := pipelineClobberSteps(Pipeline{agg, ds})

		require.NoError(t, err)
		require.Len(t, out, 1)
		c := out[0].(*Datasource)
		require.Len(t, c.clobbered, 1)
	})

	t.Run("agg agg ds", func(t *testing.T) {
		// @note for now we're only offloading one aggregation
		ds := &Datasource{
			Ident:    "ds_1",
			analysis: makeAnalysisDsAggregate(),
		}
		agg1 := &Aggregate{
			Ident:     "agg_1",
			RelSource: "ds_1",
			rel:       ds,
		}
		agg2 := &Aggregate{
			Ident:     "agg_2",
			RelSource: "agg_1",
			rel:       agg1,
		}

		out, err := pipelineClobberSteps(Pipeline{agg2, agg1, ds})

		require.NoError(t, err)
		require.Len(t, out, 2)
		c := out[1].(*Datasource)
		require.Len(t, c.clobbered, 1)
	})

	t.Run("join ds ds", func(t *testing.T) {
		ds1 := &Datasource{
			Ident:    "ds_1",
			analysis: makeAnalysisDsAggregate(),
		}
		ds2 := &Datasource{
			Ident:    "ds_2",
			analysis: makeAnalysisDsAggregate(),
		}

		join := &Join{
			Ident:    "join_1",
			RelLeft:  "ds_1",
			RelRight: "ds_2",

			relLeft:  ds1,
			relRight: ds2,
		}

		out, err := pipelineClobberSteps(Pipeline{join, ds1, ds2})

		require.NoError(t, err)
		require.Len(t, out, 3)
	})

	t.Run("join join ds ds ds", func(t *testing.T) {
		ds1 := &Datasource{
			Ident:    "ds_1",
			analysis: makeAnalysisDsAggregate(),
		}
		ds2 := &Datasource{
			Ident:    "ds_2",
			analysis: makeAnalysisDsAggregate(),
		}
		ds3 := &Datasource{
			Ident:    "ds_3",
			analysis: makeAnalysisDsAggregate(),
		}

		join1 := &Join{
			Ident:    "join_1",
			RelLeft:  "ds_1",
			RelRight: "ds_2",

			relLeft:  ds1,
			relRight: ds2,
		}
		join2 := &Join{
			Ident:    "join_2",
			RelLeft:  "join_1",
			RelRight: "ds_3",

			relLeft:  join1,
			relRight: ds3,
		}

		out, err := pipelineClobberSteps(Pipeline{join2, join1, ds1, ds2, ds3})

		require.NoError(t, err)
		require.Len(t, out, 5)
	})

	t.Run("join join agg ds ds ds", func(t *testing.T) {
		ds1 := &Datasource{
			Ident:    "ds_1",
			analysis: makeAnalysisDsAggregate(),
		}
		agg := &Aggregate{
			Ident:     "agg_1",
			RelSource: "ds_1",
			rel:       ds1,
		}
		ds2 := &Datasource{
			Ident:    "ds_2",
			analysis: makeAnalysisDsAggregate(),
		}
		ds3 := &Datasource{
			Ident:    "ds_3",
			analysis: makeAnalysisDsAggregate(),
		}

		join1 := &Join{
			Ident:    "join_1",
			RelLeft:  "agg_1",
			RelRight: "ds_2",

			relLeft:  agg,
			relRight: ds2,
		}
		join2 := &Join{
			Ident:    "join_2",
			RelLeft:  "join_1",
			RelRight: "ds_3",

			relLeft:  join1,
			relRight: ds3,
		}

		out, err := pipelineClobberSteps(Pipeline{join2, join1, agg, ds1, ds2, ds3})

		require.NoError(t, err)
		require.Len(t, out, 5)

		require.Len(t, ds1.clobbered, 1)
	})
}

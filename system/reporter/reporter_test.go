package reporter

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

type (
	mockModelFinder struct{}
)

func (mockModelFinder) FindModel(dal.ModelRef) *dal.Model {
	return &dal.Model{}
}

func TestRuns(t *testing.T) {
	simpleSS := types.ReportStepSet{
		{Load: &types.ReportStepLoad{
			Name:   "l1",
			Source: "src",
			Definition: map[string]interface{}{
				"module": "mod",
				"ns":     "ns",
			},
		}},
	}

	t.Run("no definitions", func(t *testing.T) {
		rr, err := Runs(mockModelFinder{}, simpleSS, nil)
		require.NoError(t, err)
		require.Empty(t, rr)
	})

	t.Run("one definition", func(t *testing.T) {
		defs := types.ReportFrameDefinitionSet{
			{Name: "d1", Source: "l1"},
		}
		rr, err := Runs(mockModelFinder{}, simpleSS, defs)
		require.NoError(t, err)
		require.Len(t, rr, 1)
	})

	t.Run("two different definitions for same source", func(t *testing.T) {
		defs := types.ReportFrameDefinitionSet{
			{Name: "d1", Source: "l1"},
			{Name: "d2", Source: "l1"},
		}
		rr, err := Runs(mockModelFinder{}, simpleSS, defs)
		require.NoError(t, err)
		require.Len(t, rr, 2)
		require.Equal(t, "d1", rr[0].Defs[0].Name)
		require.Equal(t, "d2", rr[1].Defs[0].Name)
	})

	t.Run("two defs for link", func(t *testing.T) {
		defs := types.ReportFrameDefinitionSet{
			{Name: "d1", Source: "lnk", Ref: "l1"},
			{Name: "d1", Source: "lnk", Ref: "l2"},
		}

		linkSS := types.ReportStepSet{
			{Load: &types.ReportStepLoad{
				Name:   "l1",
				Source: "src",
				Definition: map[string]interface{}{
					"module": "mod",
					"ns":     "ns",
				},
			}},
			{Load: &types.ReportStepLoad{
				Name:   "l2",
				Source: "src",
				Definition: map[string]interface{}{
					"module": "mod",
					"ns":     "ns",
				},
			}},
			{Link: &types.ReportStepLink{
				Name:          "lnk",
				LocalSource:   "l1",
				ForeignSource: "l2",
				LocalColumn:   "lc",
				ForeignColumn: "fc",
			}},
		}

		rr, err := Runs(mockModelFinder{}, linkSS, defs)
		require.NoError(t, err)
		require.Len(t, rr, 1)
		require.Equal(t, "d1", rr[0].Defs[0].Name)
		require.Equal(t, "l1", rr[0].Defs[0].Ref)
		require.Equal(t, "d1", rr[0].Defs[1].Name)
		require.Equal(t, "l2", rr[0].Defs[1].Ref)
	})

	t.Run("link and regular", func(t *testing.T) {
		defs := types.ReportFrameDefinitionSet{
			{Name: "d1", Source: "lnk", Ref: "l1"},
			{Name: "d1", Source: "lnk", Ref: "l2"},

			{Name: "d2", Source: "l1"},
		}

		linkSS := types.ReportStepSet{
			{Load: &types.ReportStepLoad{
				Name:   "l1",
				Source: "src",
				Definition: map[string]interface{}{
					"module": "mod",
					"ns":     "ns",
				},
			}},
			{Load: &types.ReportStepLoad{
				Name:   "l2",
				Source: "src",
				Definition: map[string]interface{}{
					"module": "mod",
					"ns":     "ns",
				},
			}},
			{Link: &types.ReportStepLink{
				Name:          "lnk",
				LocalSource:   "l1",
				ForeignSource: "l2",
				LocalColumn:   "lc",
				ForeignColumn: "fc",
			}},
		}

		rr, err := Runs(mockModelFinder{}, linkSS, defs)
		require.NoError(t, err)
		require.Len(t, rr, 2)
		require.Equal(t, "d1", rr[0].Defs[0].Name)
		require.Equal(t, "l1", rr[0].Defs[0].Ref)
		require.Equal(t, "d1", rr[0].Defs[1].Name)
		require.Equal(t, "l2", rr[0].Defs[1].Ref)

		require.Equal(t, "d2", rr[1].Defs[0].Name)
	})
}

func TestFrames_regular(t *testing.T) {
	ctx := context.Background()

	buff := dal.InMemoryBuffer()
	buff.Add(ctx, (&dal.Row{}).WithValue("id", 0, 1).WithValue("f1", 0, "f1 v1").WithValue("f2", 0, "f2 v1"))
	buff.Add(ctx, (&dal.Row{}).WithValue("id", 0, 2).WithValue("f1", 0, "f1 v2").WithValue("f2", 0, "f2 v2"))
	buff.Add(ctx, (&dal.Row{}).WithValue("id", 0, 3).WithValue("f1", 0, "f1 v3").WithValue("f2", 0, "f2 v3"))

	r := run{
		Pipeline: dal.Pipeline{&dal.Datasource{Ident: "l1"}},

		Defs: types.ReportFrameDefinitionSet{
			{
				Name:   "d1",
				Source: "l1",
				Columns: types.ReportFrameColumnSet{
					{Name: "id", Label: "id", Kind: "Record"},
					{Name: "f1", Label: "f1", Kind: "String"},
				},
			},
		},
	}

	ff, err := Frames(ctx, buff, r)
	require.NoError(t, err)

	require.Len(t, ff, 1)
	f := ff[0]
	require.Len(t, f.Rows, 3)
	require.Equal(t, types.ReportFrameRow{"1", "f1 v1"}, f.Rows[0])
	require.Equal(t, types.ReportFrameRow{"2", "f1 v2"}, f.Rows[1])
	require.Equal(t, types.ReportFrameRow{"3", "f1 v3"}, f.Rows[2])
}

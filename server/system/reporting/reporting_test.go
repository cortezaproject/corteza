package reporting

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/ql"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

type (
	mockModelFinder struct{}
)

func (mockModelFinder) FindModel(dal.ModelRef) *dal.Model {
	return &dal.Model{
		Attributes: dal.AttributeSet{{
			Ident: "attr1",
			Label: "Attr 1",
			Type:  dal.TypeText{Length: 100},
		}, {
			Ident: "attr2",
			Label: "Attr 2",
			Type:  dal.TypeNumber{},
		}},
	}
}

func TestRuns(t *testing.T) {
	simpleSS := types.ReportStepSet{
		{Load: &types.ReportStepLoad{
			Name: "l1",
			Definition: map[string]interface{}{
				"module":    "mod",
				"namespace": "ns",
			},
		}},
	}

	t.Run("no definitions", func(t *testing.T) {
		rr, err := Runs(mockModelFinder{}, simpleSS, nil)
		require.NoError(t, err)
		require.Empty(t, rr)
	})

	t.Run("one definition", func(t *testing.T) {
		defs := FrameDefinitionSet{
			{Name: "d1", Source: "l1"},
		}
		rr, err := Runs(mockModelFinder{}, simpleSS, defs)
		require.NoError(t, err)
		require.Len(t, rr, 1)
	})

	t.Run("two different definitions for same source", func(t *testing.T) {
		defs := FrameDefinitionSet{
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
		defs := FrameDefinitionSet{
			{Name: "d1", Source: "lnk", Ref: "l1"},
			{Name: "d1", Source: "lnk", Ref: "l2"},
		}

		linkSS := types.ReportStepSet{
			{Load: &types.ReportStepLoad{
				Name: "l1",
				Definition: map[string]interface{}{
					"module":    "mod",
					"namespace": "ns",
				},
			}},
			{Load: &types.ReportStepLoad{
				Name: "l2",
				Definition: map[string]interface{}{
					"module":    "mod",
					"namespace": "ns",
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
		defs := FrameDefinitionSet{
			{Name: "d1", Source: "lnk", Ref: "l1"},
			{Name: "d1", Source: "lnk", Ref: "l2"},

			{Name: "d2", Source: "l1"},
		}

		linkSS := types.ReportStepSet{
			{Load: &types.ReportStepLoad{
				Name: "l1",
				Definition: map[string]interface{}{
					"module":    "mod",
					"namespace": "ns",
				},
			}},
			{Load: &types.ReportStepLoad{
				Name: "l2",
				Definition: map[string]interface{}{
					"module":    "mod",
					"namespace": "ns",
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

	// Uses all possible step definitions in one go
	t.Run("full definitions", func(t *testing.T) {
		defs := FrameDefinitionSet{
			// Load
			{Name: "d1", Source: "l1"},

			// Aggregate
			{Name: "d2", Source: "agg1"},

			// Link
			{Name: "d3", Source: "lnk", Ref: "l1"},
			{Name: "d3", Source: "lnk", Ref: "l2"},

			// join
			{Name: "d4", Source: "jn"},
		}

		rr, err := Runs(
			mockModelFinder{},
			types.ReportStepSet{
				// Loads
				{Load: &types.ReportStepLoad{
					Name: "l1",
					Definition: map[string]interface{}{
						"module":    "mod",
						"namespace": "ns",
					},
				}},
				{Load: &types.ReportStepLoad{
					Name: "l2",
					Definition: map[string]interface{}{
						"module":    "mod",
						"namespace": "ns",
					},
				}},

				// Aggregate
				{Aggregate: &types.ReportStepAggregate{
					Name:   "agg1",
					Source: "l1",
					Keys: types.ReportAggregateColumnSet{{
						Name: "k1",
						Def:  &types.ReportFilterExpr{ASTNode: &ql.ASTNode{Raw: "f1"}},
					}},
					Columns: types.ReportAggregateColumnSet{{
						Name: "c1",
						Def:  &types.ReportFilterExpr{ASTNode: &ql.ASTNode{Raw: "count(f1)"}},
					}},
				}},

				// Link
				{Link: &types.ReportStepLink{
					Name:          "lnk",
					LocalSource:   "l1",
					ForeignSource: "l2",
					LocalColumn:   "lc",
					ForeignColumn: "fc",
				}},

				// Join
				{Join: &types.ReportStepJoin{
					Name:          "jn",
					LocalSource:   "l1",
					ForeignSource: "l2",
					LocalColumn:   "lc",
					ForeignColumn: "fc",
				}},
			},
			defs,
		)
		require.NoError(t, err)
		require.Len(t, rr, 4)
	})
}

func TestStepLoadConv(t *testing.T) {
	t.Run("ref from ID", func(t *testing.T) {
		c, err := convStepLoad(
			mockModelFinder{},
			types.ReportStepLoad{
				Name: "l1",
				Definition: map[string]interface{}{
					"moduleID":     42,
					"namespaceID":  43,
					"connectionID": 44,
				},
			},
			nil,
		)

		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("ref from handles", func(t *testing.T) {
		c, err := convStepLoad(
			mockModelFinder{},
			types.ReportStepLoad{
				Name: "l1",
				Definition: map[string]interface{}{
					"module":       "mod1",
					"namespace":    "ns1",
					"connectionID": 44,
				},
			},
			nil,
		)

		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("missing module ref", func(t *testing.T) {
		_, err := convStepLoad(
			mockModelFinder{},
			types.ReportStepLoad{
				Name: "l1",
				Definition: map[string]interface{}{
					"namespace":    "ns1",
					"connectionID": 44,
				},
			},
			nil,
		)

		require.Error(t, err)
	})

	t.Run("missing namespace ref", func(t *testing.T) {
		_, err := convStepLoad(
			mockModelFinder{},
			types.ReportStepLoad{
				Name: "l1",
				Definition: map[string]interface{}{
					"module":       "mod1",
					"connectionID": 44,
				},
			},
			nil,
		)

		require.Error(t, err)
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

		Defs: FrameDefinitionSet{
			{
				Name:   "d1",
				Source: "l1",
				Columns: FrameColumnSet{
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
	require.Equal(t, FrameRow{"1", "f1 v1"}, f.Rows[0])
	require.Equal(t, FrameRow{"2", "f1 v2"}, f.Rows[1])
	require.Equal(t, FrameRow{"3", "f1 v3"}, f.Rows[2])
}

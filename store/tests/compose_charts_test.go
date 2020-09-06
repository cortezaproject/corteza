package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testComposeCharts(t *testing.T, s store.Storable) {
	var (
		ctx = context.Background()
		req = require.New(t)

		namespaceID = id.Next()

		makeNew = func(name, handle string) *types.Chart {
			// minimum data set for new composeChart
			return &types.Chart{
				ID:          id.Next(),
				NamespaceID: namespaceID,
				CreatedAt:   time.Now(),
				Name:        name,
				Handle:      handle,
			}
		}
	)

	t.Run("create", func(t *testing.T) {
		composeChart := makeNew("ComposeChartCRUD", "compose-chart-crud")
		req.NoError(s.CreateComposeChart(ctx, composeChart))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup by ID", func(t *testing.T) {
		composeChart := makeNew("look up by id", "look-up-by-id")
		req.NoError(s.CreateComposeChart(ctx, composeChart))
		fetched, err := s.LookupComposeChartByID(ctx, composeChart.ID)
		req.NoError(err)
		req.Equal(composeChart.Name, fetched.Name)
		req.Equal(composeChart.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("Delete", func(t *testing.T) {
		composeChart := makeNew("Delete", "Delete")
		req.NoError(s.CreateComposeChart(ctx, composeChart))
		req.NoError(s.DeleteComposeChart(ctx))
	})

	t.Run("Delete by ID", func(t *testing.T) {
		composeChart := makeNew("Delete by id", "Delete-by-id")
		req.NoError(s.CreateComposeChart(ctx, composeChart))
		req.NoError(s.DeleteComposeChart(ctx))
	})

	t.Run("update", func(t *testing.T) {
		composeChart := makeNew("update me", "update-me")
		req.NoError(s.CreateComposeChart(ctx, composeChart))

		composeChart = &types.Chart{
			ID:        composeChart.ID,
			CreatedAt: composeChart.CreatedAt,
			Name:      "ComposeChartCRUD+2",
		}
		req.NoError(s.UpdateComposeChart(ctx, composeChart))

		updated, err := s.LookupComposeChartByID(ctx, composeChart.ID)
		req.NoError(err)
		req.Equal(composeChart.Name, updated.Name)
	})

	t.Run("update with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Chart{
			makeNew("/one-one", "chart-1-1"),
			makeNew("/one-two", "chart-1-2"),
			makeNew("/two-one", "chart-2-1"),
			makeNew("/two-two", "chart-2-2"),
			makeNew("/two-deleted", "chart-2-d"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateComposeCharts(ctx))
		req.NoError(s.CreateComposeChart(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchComposeCharts(ctx, types.ChartFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchComposeCharts(ctx, types.ChartFilter{Deleted: rh.FilterStateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchComposeCharts(ctx, types.ChartFilter{Deleted: rh.FilterStateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		set, f, err = s.SearchComposeCharts(ctx, types.ChartFilter{Handle: "chart-2-1"})
		req.NoError(err)
		req.Len(set, 1)

		// find all prefixed
		set, f, err = s.SearchComposeCharts(ctx, types.ChartFilter{Query: "/two-"})
		req.NoError(err)
		req.Len(set, 2)

		_ = f // dummy
	})
}

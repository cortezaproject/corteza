package tests

import (
	"context"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testComposeCharts(t *testing.T, s store.Storer) {
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

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Chart) {
			req := require.New(t)
			req.NoError(s.TruncateComposeCharts(ctx))
			res := makeNew(string(rand.Bytes(10)), string(rand.Bytes(10)))
			req.NoError(s.CreateComposeChart(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		composeChart := makeNew("ComposeChartCRUD", "compose-chart-crud")
		req.NoError(s.CreateComposeChart(ctx, composeChart))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup", func(t *testing.T) {
		t.Run("by ID", func(t *testing.T) {
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

		t.Run("by Namespace ID, Handle", func(t *testing.T) {
			composeChart := makeNew("look up by handle", "look-up-by-handle")
			req.NoError(s.CreateComposeChart(ctx, composeChart))
			fetched, err := s.LookupComposeChartByNamespaceIDHandle(ctx, composeChart.NamespaceID, composeChart.Handle)
			req.NoError(err)
			req.Equal(composeChart.Name, fetched.Name)
			req.Equal(composeChart.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})
	})

	t.Run("update", func(t *testing.T) {
		req, composeChart := truncAndCreate(t)
		composeChart.Name = "ComposeChartCRUD+2"

		req.NoError(s.UpdateComposeChart(ctx, composeChart))

		updated, err := s.LookupComposeChartByID(ctx, composeChart.ID)
		req.NoError(err)
		req.Equal(composeChart.Name, updated.Name)
	})

	t.Run("update with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, composeChart := truncAndCreate(t)
			composeChart.Name = "ComposeChartCRUD+2"

			req.NoError(s.UpsertComposeChart(ctx, composeChart))

			upserted, err := s.LookupComposeChartByID(ctx, composeChart.ID)
			req.NoError(err)
			req.Equal(composeChart.Name, upserted.Name)
		})

		t.Run("new", func(t *testing.T) {
			composeChart := makeNew("upsert me", "upsert-me")
			composeChart.Name = "ComposeChartCRUD+2"

			req.NoError(s.UpsertComposeChart(ctx, composeChart))

			upserted, err := s.LookupComposeChartByID(ctx, composeChart.ID)
			req.NoError(err)
			req.Equal(composeChart.Name, upserted.Name)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Chart", func(t *testing.T) {
			req, composeChart := truncAndCreate(t)
			req.NoError(s.DeleteComposeChart(ctx, composeChart))
			_, err := s.LookupComposeChartByID(ctx, composeChart.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, composeChart := truncAndCreate(t)
			req.NoError(s.DeleteComposeChartByID(ctx, composeChart.ID))
			_, err := s.LookupComposeChartByID(ctx, composeChart.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
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
		set, f, err = s.SearchComposeCharts(ctx, types.ChartFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchComposeCharts(ctx, types.ChartFilter{Deleted: filter.StateExclusive})
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

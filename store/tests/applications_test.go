package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testApplications(t *testing.T, s store.Applications) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(name string) *types.Application {
			// minimum data set for new application
			return &types.Application{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Name:      name,
				Unify:     &types.ApplicationUnify{},
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Application) {
			req := require.New(t)
			req.NoError(s.TruncateApplications(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateApplication(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		application := makeNew("ApplicationCRUD")
		req.NoError(s.CreateApplication(ctx, application))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, application := truncAndCreate(t)
		fetched, err := s.LookupApplicationByID(ctx, application.ID)
		req.NoError(err)
		req.Equal(application.Name, fetched.Name)
		req.Equal(application.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, application := truncAndCreate(t)
		application.Name = "ApplicationCRUD+2"

		req.NoError(s.UpdateApplication(ctx, application))

		updated, err := s.LookupApplicationByID(ctx, application.ID)
		req.NoError(err)
		req.Equal(application.Name, updated.Name)
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, application := truncAndCreate(t)
			application.Name = "ApplicationCRUD+2"

			req.NoError(s.UpsertApplication(ctx, application))

			updated, err := s.LookupApplicationByID(ctx, application.ID)
			req.NoError(err)
			req.Equal(application.Name, updated.Name)
		})

		t.Run("new", func(t *testing.T) {
			application := makeNew("upsert me")
			application.Name = "ComposeChartCRUD+2"

			req.NoError(s.UpsertApplication(ctx, application))

			upserted, err := s.LookupApplicationByID(ctx, application.ID)
			req.NoError(err)
			req.Equal(application.Name, upserted.Name)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Application", func(t *testing.T) {
			req, application := truncAndCreate(t)
			req.NoError(s.DeleteApplication(ctx, application))
			_, err := s.LookupApplicationByID(ctx, application.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, application := truncAndCreate(t)
			req.NoError(s.DeleteApplicationByID(ctx, application.ID))
			_, err := s.LookupApplicationByID(ctx, application.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Application{
			makeNew("/one-one"),
			makeNew("/one-two"),
			makeNew("/two-one"),
			makeNew("/two-two"),
			makeNew("/two-deleted"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateApplications(ctx))
		req.NoError(s.CreateApplication(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchApplications(ctx, types.ApplicationFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchApplications(ctx, types.ApplicationFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchApplications(ctx, types.ApplicationFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		set, f, err = s.SearchApplications(ctx, types.ApplicationFilter{Name: "/two-one"})
		req.NoError(err)
		req.Len(set, 1)

		// find all prefixed
		set, f, err = s.SearchApplications(ctx, types.ApplicationFilter{Query: "/two"})
		req.NoError(err)
		req.Len(set, 2)

		_ = f // dummy
	})

	t.Run("metrics", func(t *testing.T) {
		var (
			req = require.New(t)

			e = &types.ApplicationMetrics{
				Total:   5,
				Valid:   2,
				Deleted: 3,
			}
		)

		req.NoError(s.TruncateApplications(ctx))
		req.NoError(s.CreateApplication(ctx, &types.Application{ID: id.Next(), CreatedAt: *now(), UpdatedAt: now(), Unify: &types.ApplicationUnify{}}))
		req.NoError(s.CreateApplication(ctx, &types.Application{ID: id.Next(), CreatedAt: *now(), UpdatedAt: now(), Unify: &types.ApplicationUnify{}}))
		req.NoError(s.CreateApplication(ctx, &types.Application{ID: id.Next(), CreatedAt: *now(), DeletedAt: now(), Unify: &types.ApplicationUnify{}}))
		req.NoError(s.CreateApplication(ctx, &types.Application{ID: id.Next(), CreatedAt: *now(), DeletedAt: now(), Unify: &types.ApplicationUnify{}}))
		req.NoError(s.CreateApplication(ctx, &types.Application{ID: id.Next(), CreatedAt: *now(), DeletedAt: now(), Unify: &types.ApplicationUnify{}}))

		m, err := store.ApplicationMetrics(ctx, s)
		req.NoError(err)
		req.Equal(e, m)
	})
}

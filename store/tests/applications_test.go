package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
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
	)

	t.Run("create", func(t *testing.T) {
		application := makeNew("ApplicationCRUD")
		req.NoError(s.CreateApplication(ctx, application))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		application := makeNew("look up by id")
		req.NoError(s.CreateApplication(ctx, application))
		fetched, err := s.LookupApplicationByID(ctx, application.ID)
		req.NoError(err)
		req.Equal(application.Name, fetched.Name)
		req.Equal(application.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("Delete", func(t *testing.T) {
		application := makeNew("Delete")
		req.NoError(s.CreateApplication(ctx, application))
		req.NoError(s.DeleteApplication(ctx))
	})

	t.Run("Delete by ID", func(t *testing.T) {
		application := makeNew("Delete by id")
		req.NoError(s.CreateApplication(ctx, application))
		req.NoError(s.DeleteApplication(ctx))
	})

	t.Run("update", func(t *testing.T) {
		application := makeNew("update me")
		req.NoError(s.CreateApplication(ctx, application))

		application = &types.Application{
			ID:        application.ID,
			CreatedAt: application.CreatedAt,
			Name:      "ApplicationCRUD+2",
			Unify:     application.Unify,
		}
		req.NoError(s.UpdateApplication(ctx, application))

		updated, err := s.LookupApplicationByID(ctx, application.ID)
		req.NoError(err)
		req.Equal(application.Name, updated.Name)
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
}

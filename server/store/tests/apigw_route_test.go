package tests

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testApigwRoutes(t *testing.T, s store.ApigwRoutes) {
	var (
		ctx = context.Background()
		new = &types.ApigwRoute{
			CreatedAt: *now(),
			ID:        42,
			Endpoint:  "/foo",
			Enabled:   true,
			CreatedBy: 1}

		disabled = &types.ApigwRoute{
			CreatedAt: *now(),
			ID:        4242,
			Endpoint:  "/foo_disabled",
			Enabled:   false,
			CreatedBy: 1}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwRoutes(ctx))
		req.NoError(s.CreateApigwRoute(ctx, new))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwRoutes(ctx))
		req.NoError(s.CreateApigwRoute(ctx, new))
		fetched, err := s.LookupApigwRouteByID(ctx, new.ID)
		req.NoError(err)
		req.Equal(new.ID, fetched.ID)
		req.Equal(new.Endpoint, fetched.Endpoint)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("lookup by endpoint", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwRoutes(ctx))
		req.NoError(s.CreateApigwRoute(ctx, new))
		fetched, err := s.LookupApigwRouteByEndpoint(ctx, new.Endpoint)
		req.NoError(err)
		req.Equal(new.ID, fetched.ID)
		req.Equal(new.Endpoint, fetched.Endpoint)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwRoutes(ctx))
		req.NoError(s.UpdateApigwRoute(ctx, new))
	})

	t.Run("search", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwRoutes(ctx))
		req.NoError(s.CreateApigwRoute(ctx,
			new,
		))

		set, _, err := s.SearchApigwRoutes(ctx, types.ApigwRouteFilter{Endpoint: "fo"})
		req.NoError(err)
		req.Len(set, 1)
	})

	t.Run("search only disabled", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwRoutes(ctx))
		req.NoError(s.CreateApigwRoute(ctx,
			disabled,
		))

		set, _, err := s.SearchApigwRoutes(ctx, types.ApigwRouteFilter{
			Disabled: filter.StateExclusive,
		})

		req.NoError(err)
		req.Len(set, 1)
		req.Equal(set[0].ID, disabled.ID)
	})

	t.Run("search enabled and disabled", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwRoutes(ctx))
		req.NoError(s.CreateApigwRoute(ctx,
			new,
			disabled,
		))

		set, _, err := s.SearchApigwRoutes(ctx, types.ApigwRouteFilter{
			Disabled: filter.StateInclusive,
		})

		req.NoError(err)
		req.Len(set, 2)
	})
}

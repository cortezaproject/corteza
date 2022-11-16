package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testApigwFilter(t *testing.T, s store.ApigwFilters) {
	var (
		ctx = context.Background()
		new = &types.ApigwFilter{
			ID:        42,
			Enabled:   true,
			CreatedAt: time.Now(),
			CreatedBy: 1}
		disabled = &types.ApigwFilter{
			ID:        4242,
			Enabled:   false,
			CreatedAt: time.Now(),
			CreatedBy: 1}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwFilters(ctx))
		req.NoError(s.CreateApigwFilter(ctx, new))
	})

	t.Run("update", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwFilters(ctx))
		req.NoError(s.UpdateApigwFilter(ctx, new))
	})

	t.Run("search", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwFilters(ctx))
		req.NoError(s.CreateApigwFilter(ctx,
			new,
		))

		set, _, err := s.SearchApigwFilters(ctx, types.ApigwFilterFilter{})
		req.NoError(err)
		req.Len(set, 1)
	})

	t.Run("search only disabled", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwFilters(ctx))
		req.NoError(s.CreateApigwFilter(ctx,
			disabled,
		))

		set, _, err := s.SearchApigwFilters(ctx, types.ApigwFilterFilter{
			Disabled: filter.StateExclusive,
		})

		req.NoError(err)
		req.Len(set, 1)
		req.Equal(set[0].ID, disabled.ID)
	})

	t.Run("search enabled and disabled", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateApigwFilters(ctx))
		req.NoError(s.CreateApigwFilter(ctx,
			new,
			disabled,
		))

		set, _, err := s.SearchApigwFilters(ctx, types.ApigwFilterFilter{
			Disabled: filter.StateInclusive,
		})

		req.NoError(err)
		req.Len(set, 2)
	})
}

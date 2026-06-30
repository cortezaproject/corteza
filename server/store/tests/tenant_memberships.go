package tests

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func testTenantMemberships(t *testing.T, s store.TenantMemberships) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(tenantID, userID uint64) *types.TenantMembership {
			return &types.TenantMembership{
				ID:        id.Next(),
				TenantID:  tenantID,
				UserID:    userID,
				Role:      types.TenantRoleMember,
				Status:    types.TenantMemberStatusActive,
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.TenantMembership) {
			req := require.New(t)
			req.NoError(s.TruncateTenantMemberships(ctx))
			res := makeNew(id.Next(), id.Next())
			req.NoError(s.CreateTenantMembership(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateTenantMemberships(ctx))
		tm := makeNew(id.Next(), id.Next())
		req.NoError(s.CreateTenantMembership(ctx, tm))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, tm := truncAndCreate(t)
		fetched, err := s.LookupTenantMembershipByID(ctx, tm.ID)
		req.NoError(err)
		req.Equal(tm.ID, fetched.ID)
		req.Equal(tm.TenantID, fetched.TenantID)
	})

	t.Run("lookup by user ID", func(t *testing.T) {
		req, tm := truncAndCreate(t)
		fetched, err := s.LookupTenantMembershipByUserID(ctx, tm.UserID)
		req.NoError(err)
		req.Equal(tm.ID, fetched.ID)
	})

	t.Run("lookup by tenant and user", func(t *testing.T) {
		req, tm := truncAndCreate(t)
		fetched, err := s.LookupTenantMembershipByTenantIDUserID(ctx, tm.TenantID, tm.UserID)
		req.NoError(err)
		req.Equal(tm.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, tm := truncAndCreate(t)
		tm.Role = types.TenantRoleAdmin
		req.NoError(s.UpdateTenantMembership(ctx, tm))
		fetched, err := s.LookupTenantMembershipByID(ctx, tm.ID)
		req.NoError(err)
		req.Equal(types.TenantRoleAdmin, fetched.Role)
	})

	t.Run("delete", func(t *testing.T) {
		req, tm := truncAndCreate(t)
		req.NoError(s.DeleteTenantMembershipByID(ctx, tm.ID))
	})
}

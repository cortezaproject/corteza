package tests

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/rand"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func testTenants(t *testing.T, s store.Tenants) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.Tenant {
			return &types.Tenant{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
				Status:    types.TenantStatusActive,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Tenant) {
			req := require.New(t)
			req.NoError(s.TruncateTenants(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateTenant(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateTenants(ctx))
		tenant := makeNew("TenantCRUD")
		req.NoError(s.CreateTenant(ctx, tenant))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, tenant := truncAndCreate(t)
		fetched, err := s.LookupTenantByID(ctx, tenant.ID)
		req.NoError(err)
		req.Equal(tenant.Handle, fetched.Handle)
		req.Equal(tenant.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, tenant := truncAndCreate(t)
		tenant.Status = types.TenantStatusSuspended
		req.NoError(s.UpdateTenant(ctx, tenant))
		fetched, err := s.LookupTenantByID(ctx, tenant.ID)
		req.NoError(err)
		req.Equal(types.TenantStatusSuspended, fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		req, tenant := truncAndCreate(t)
		req.NoError(s.DeleteTenantByID(ctx, tenant.ID))
	})
}

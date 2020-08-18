package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testRoles(t *testing.T, s rolesStore) {
	var (
		ctx = context.Background()
		req = require.New(t)

		//err  error
		role *types.Role
	)

	t.Run("create", func(t *testing.T) {
		role = &types.Role{
			ID:        42,
			CreatedAt: time.Now(),
			Name:      "RoleCRUD",
			Handle:    "rolecrud",
		}
		role.ArchivedAt = &role.CreatedAt
		req.NoError(s.CreateRole(ctx, role))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		fetched, err := s.LookupRoleByID(ctx, role.ID)
		req.NoError(err)
		req.Equal(role.Name, fetched.Name)
		req.Equal(role.Handle, fetched.Handle)
		req.Equal(role.ID, fetched.ID)
		req.NotNil(fetched.ArchivedAt)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		role = &types.Role{
			ID:        42,
			CreatedAt: time.Now(),
			Name:      "RoleCRUD+2",
			Handle:    "rolecrud+2",
		}
		req.NoError(s.UpdateRole(ctx, role))
	})

	t.Run("lookup by handle", func(t *testing.T) {
		fetched, err := s.LookupRoleByHandle(ctx, role.Handle)
		req.NoError(err)
		req.Equal(role.ID, fetched.ID)
	})

	t.Run("search", func(t *testing.T) {
		set, f, err := s.SearchRoles(ctx, types.RoleFilter{RoleID: []uint64{role.ID}})
		req.NoError(err)
		req.Equal([]uint64{role.ID}, f.RoleID)
		req.Len(set, 1)
		//req.Equal(set[0].ID, role.ID)
	})

	t.Run("search", func(t *testing.T) {
		set, f, err := s.SearchRoles(ctx, types.RoleFilter{Name: role.Name})
		req.NoError(err)
		req.Len(set, 1)

		_ = f // dummy
	})

	t.Run("search by *", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("ordered search", func(t *testing.T) {
		t.Skip("not implemented")
	})
}

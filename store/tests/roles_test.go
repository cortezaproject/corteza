package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func testRoles(t *testing.T, s store.Roles) {
	var (
		ctx = context.Background()

		makeNew = func(nn ...string) *types.Role {
			name := strings.Join(nn, "")
			return &types.Role{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Name:      "RoleCRUD+" + name,
				Handle:    "rolecrud+" + name,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Role) {
			req := require.New(t)
			req.NoError(s.TruncateRoles(ctx))
			role := makeNew()
			req.NoError(s.CreateRole(ctx, role))
			return req, role
		}

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, types.RoleSet) {
			req := require.New(t)
			req.NoError(s.TruncateRoles(ctx))

			set := make([]*types.Role, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(string(rand.Bytes(10)))
			}

			req.NoError(s.CreateRole(ctx, set...))
			return req, set
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.CreateRole(ctx, makeNew()))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, role := truncAndCreate(t)

		fetched, err := s.LookupRoleByID(ctx, role.ID)
		req.NoError(err)
		req.Equal(role.Name, fetched.Name)
		req.Equal(role.Handle, fetched.Handle)
		req.Equal(role.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, role := truncAndCreate(t)
		req.NoError(s.UpdateRole(ctx, role))
	})

	t.Run("lookup by handle", func(t *testing.T) {
		req, role := truncAndCreate(t)

		fetched, err := s.LookupRoleByHandle(ctx, role.Handle)
		req.NoError(err)
		req.Equal(role.ID, fetched.ID)
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by ID", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, f, err := s.SearchRoles(ctx, types.RoleFilter{RoleID: []uint64{prefill[0].ID}})
			req.NoError(err)
			req.Equal([]uint64{prefill[0].ID}, f.RoleID)
			req.Len(set, 1)
		})

		t.Run("by query", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, _, err := s.SearchRoles(ctx, types.RoleFilter{Query: prefill[0].Handle})
			req.NoError(err)
			req.Len(set, 1)
		})

		t.Run("by state", func(t *testing.T) {
			t.Run("deleted", func(t *testing.T) {
				req, prefill := truncAndFill(t, 5)

				prefill[0].DeletedAt = &(prefill[0].CreatedAt)
				s.UpdateRole(ctx, prefill[0])

				set, _, err := s.SearchRoles(ctx, types.RoleFilter{Deleted: rh.FilterStateExcluded})
				req.NoError(err)
				req.Len(set, 4)

				set, _, err = s.SearchRoles(ctx, types.RoleFilter{Deleted: rh.FilterStateInclusive})
				req.NoError(err)
				req.Len(set, 5)

				set, _, err = s.SearchRoles(ctx, types.RoleFilter{Deleted: rh.FilterStateExclusive})
				req.NoError(err)
				req.Len(set, 1)
			})

			t.Run("archived", func(t *testing.T) {
				req, prefill := truncAndFill(t, 5)

				prefill[0].ArchivedAt = &(prefill[0].CreatedAt)
				s.UpdateRole(ctx, prefill[0])

				set, _, err := s.SearchRoles(ctx, types.RoleFilter{Archived: rh.FilterStateExcluded})
				req.NoError(err)
				req.Len(set, 4)

				set, _, err = s.SearchRoles(ctx, types.RoleFilter{Archived: rh.FilterStateInclusive})
				req.NoError(err)
				req.Len(set, 5)

				set, _, err = s.SearchRoles(ctx, types.RoleFilter{Archived: rh.FilterStateExclusive})
				req.NoError(err)
				req.Len(set, 1)
			})
		})

		t.Run("with check", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, _, err := s.SearchRoles(ctx, types.RoleFilter{
				Check: func(role *types.Role) (bool, error) {
					return role.ID == prefill[0].ID, nil
				},
			})
			req.NoError(err)
			req.Len(set, 1)
			req.Equal(prefill[0].ID, set[0].ID)
		})
	})

	t.Run("ordered search", func(t *testing.T) {
		t.Skip("not implemented")
	})
}

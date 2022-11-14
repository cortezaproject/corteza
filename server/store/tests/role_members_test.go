package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"testing"
)

func testRoleMembers(t *testing.T, s store.RoleMembers) {
	var (
		ctx = context.Background()

		makeNew = func(nn ...string) *types.RoleMember {
			// minimum data set for new RoleMember
			return &types.RoleMember{
				RoleID: id.Next(),
				UserID: id.Next(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.RoleMember) {
			req := require.New(t)
			req.NoError(s.TruncateRoleMembers(ctx))
			roleMember := makeNew()
			req.NoError(s.CreateRoleMember(ctx, roleMember))
			return req, roleMember
		}

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, types.RoleMemberSet) {
			req := require.New(t)
			req.NoError(s.TruncateRoleMembers(ctx))

			set := make([]*types.RoleMember, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew()
			}

			req.NoError(s.CreateRoleMember(ctx, set...))
			return req, set
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.CreateRoleMember(ctx, makeNew()))
	})

	t.Run("create multiple", func(t *testing.T) {
		truncAndFill(t, 5)
	})

	// t.Run("update", func(t *testing.T) {
	// 	req, roleMember := truncAndCreate(t)

	// 	newID := id.Next()
	// 	roleMember = &types.RoleMember{
	// 		OwnerID:        roleMember.OwnerID,
	// 		RoleID:        newID,
	// 	}
	// 	req.NoError(s.UpdateRoleMember(ctx, roleMember))

	// 	updated, f, err := s.SearchRoleMembers(ctx, types.RoleMemberFilter{OwnerID: roleMember.OwnerID})
	// 	req.NoError(err)
	// 	req.Equal(updated[0].OwnerID, f.OwnerID)
	// })

	// t.Run("search", func(t *testing.T) {
	// 	t.Run("by RoleID", func(t *testing.T) {
	// 		req, prefill := truncAndFill(t, 5)

	// 		set, f, err := s.SearchRoleMembers(ctx, types.RoleMemberFilter{RoleID: prefill[0].RoleID})
	// 		req.NoError(err)
	// 		req.Equal(prefill[0].RoleID, f.RoleID)
	// 		req.Len(set, 1)
	// 	})
	// })

	t.Run("delete", func(t *testing.T) {
		t.Run("by role member", func(t *testing.T) {
			req, roleMember := truncAndCreate(t)
			req.NoError(s.DeleteRoleMember(ctx, roleMember))
			roleMembers, _, _ := s.SearchRoleMembers(ctx, types.RoleMemberFilter{UserID: roleMember.UserID, RoleID: roleMember.RoleID})
			req.Len(roleMembers, 0)
		})

		t.Run("by RoleID and UserID", func(t *testing.T) {
			req, roleMember := truncAndCreate(t)
			req.NoError(s.DeleteRoleMemberByUserIDRoleID(ctx, roleMember.UserID, roleMember.RoleID))
			roleMembers, _, _ := s.SearchRoleMembers(ctx, types.RoleMemberFilter{UserID: roleMember.UserID, RoleID: roleMember.RoleID})
			req.Len(roleMembers, 0)
		})
	})
}

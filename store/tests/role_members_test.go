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
		req = require.New(t)

		role *types.RoleMember
	)

	t.Run("create", func(t *testing.T) {
		role = &types.RoleMember{
			RoleID: id.Next(),
			UserID: id.Next(),
		}
		req.NoError(s.CreateRoleMember(ctx, role))
	})

	t.Run("search", func(t *testing.T) {
		_, _, err := s.SearchRoleMembers(ctx, types.RoleMemberFilter{})
		req.NoError(err)
	})

}

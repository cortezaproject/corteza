package service

import (
	"context"
	"testing"

	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/sqlite3"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestUser_ProtectedSearch(t *testing.T) {
	const testRoleID = 123

	var (
		req = require.New(t)
		ctx = context.Background()

		masked   = &types.User{Email: "email.masked@us.er", Name: "Name MSKD", ID: nextID(), CreatedAt: *now()}
		unmasked = &types.User{Email: "email.unmasked@us.er", Name: "Name UNMSKD", ID: nextID(), CreatedAt: *now()}

		set types.UserSet
		err error

		testUser = &types.User{ID: 42}

		acRBAC = rbac.NewService(zap.NewNop(), nil)

		s store.Storer
	)

	if s, err = sqlite3.ConnectInMemory(ctx); err != nil {
		req.NoError(err)
	} else if err = store.Upgrade(ctx, zap.NewNop(), s); err != nil {
		req.NoError(err)
	}

	acRBAC.UpdateRoles(rbac.CommonRole.Make(testRoleID, "test-role"))
	req.NoError(acRBAC.Grant(ctx,
		rbac.AllowRule(testRoleID, types.ComponentRbacResource(), "users.search"),
		rbac.AllowRule(testRoleID, types.UserRbacResource(0), "read"),
		rbac.DenyRule(testRoleID, masked.RbacResource(), "email.unmask"),
		rbac.AllowRule(testRoleID, unmasked.RbacResource(), "email.unmask"),
		rbac.DenyRule(testRoleID, masked.RbacResource(), "name.unmask"),
		rbac.AllowRule(testRoleID, unmasked.RbacResource(), "name.unmask"),
	))

	testUser.SetRoles(testRoleID)
	ctx = a.SetIdentityToContext(ctx, testUser)

	svc := &user{
		settings: &types.AppSettings{},
		ac:       &accessControl{rbac: acRBAC},
		eventbus: eventbus.New(),
		store:    s,
	}

	req.NoError(store.CreateUser(ctx, svc.store, masked, unmasked))

	t.Run("with disabled masking", func(t *testing.T) {
		// Masking disabled, expecting to fetch both users
		svc.settings.Privacy.Mask.Email = false
		svc.settings.Privacy.Mask.Name = false
		set, _, err = svc.Find(ctx, types.UserFilter{Query: "email"})
		req.NoError(err)
		req.Len(set, 2)
		req.NotEqual(set[0].Email, maskPrivateDataEmail)
		req.NotEqual(set[1].Email, maskPrivateDataEmail)
	})

	t.Run("with enabled privacy", func(t *testing.T) {
		// Masking enabled, expecting to fetch only unmasked
		svc.settings.Privacy.Mask.Email = true
		svc.settings.Privacy.Mask.Name = true

		set, _, err = svc.Find(ctx, types.UserFilter{Query: "email"})
		req.NoError(err)
		req.Len(set, 1)
	})

	t.Run("email search with enabled privacy", func(t *testing.T) {
		// Masking enabled, expecting to fetch only unmasked
		svc.settings.Privacy.Mask.Email = true
		svc.settings.Privacy.Mask.Name = true

		set, _, err = svc.Find(ctx, types.UserFilter{Email: "email.masked@us.er"})
		req.NoError(err)
		req.Len(set, 0)
	})
}

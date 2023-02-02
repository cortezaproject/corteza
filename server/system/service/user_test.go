package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
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

	if s, err = sqlite.ConnectInMemory(ctx); err != nil {
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

func Test_processAvatarInitials(t *testing.T) {
	// Define test cases
	tests := []struct {
		name            string
		user            *types.User
		expectedInitial string
	}{
		{
			name: "Test with valid name",
			user: &types.User{
				Name: "John Doe",
			},
			expectedInitial: "JD",
		},
		{
			name: "Test with valid handle",
			user: &types.User{
				Handle: "johndoe",
			},
			expectedInitial: "J",
		},
		{
			name: "Test handle with a delimiter",
			user: &types.User{
				Handle: "john_doe",
			},
			expectedInitial: "JD",
		},
		{
			name: "Test with valid handle",
			user: &types.User{
				Email: "johndoe@example.com",
			},
			expectedInitial: "J",
		},
		{
			name: "Test email with a delimiter",
			user: &types.User{
				Email: "john-doe@example.com",
			},
			expectedInitial: "JD",
		},
		{
			name: "Test with one letter name",
			user: &types.User{
				Name: "K",
			},
			expectedInitial: "K",
		},
		{
			name: "Test with one name",
			user: &types.User{
				Name: "Doe",
			},
			expectedInitial: "DO",
		},
		{
			name: "Test with three names",
			user: &types.User{
				Name: "John David Doe",
			},
			expectedInitial: "JDD",
		},
		{
			name: "Test with a name that has a valid and invalid letter",
			user: &types.User{
				Name: "k-",
			},
			expectedInitial: "K",
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			initial := processAvatarInitials(tc.user)
			assert.Equal(t, tc.expectedInitial, initial)
		})
	}
}

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	u struct {
		store.Users
		L func(ctx context.Context, handle string) (*types.User, error)
	}
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

		s store.Storer
	)

	if s, err = sqlite.ConnectInMemory(ctx); err != nil {
		req.NoError(err)
	} else if err = store.Upgrade(ctx, zap.NewNop(), s); err != nil {
		req.NoError(err)
	}

	var (
		acRBAC = rbac.NewServiceMust(ctx, zap.NewNop(), s, rbac.Config{
			Synchronous:        true,
			DecayInterval:      time.Hour * 2,
			CleanupInterval:    time.Hour * 2,
			ReindexInterval:    time.Hour * 2,
			IndexFlushInterval: time.Hour * 2,
			RuleStorage:        s,
			RoleStorage:        s,
		})
	)

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

func Test_createUserHandle(t *testing.T) {
	ctx := context.Background()

	// Define test cases
	tests := []struct {
		name     string
		user     *types.User
		expected string
		s        u
	}{
		{
			name: "Test with valid name",
			user: &types.User{
				Name: "John Doe",
			},
			expected: "JohnDoe_1",
			s: u{
				L: func(ctx context.Context, handle string) (*types.User, error) {
					return nil, store.ErrNotFound
				},
			},
		},
		{
			name: "Test with valid name and username",
			user: &types.User{
				Name:     "John",
				Username: "username_provided_here123",
			},
			expected: "John_username_provided_here123",
			s: u{
				L: func(ctx context.Context, handle string) (*types.User, error) {
					if handle == "John1" {
						return nil, store.ErrNotUnique
					}
					return nil, store.ErrNotFound
				},
			},
		},
		{
			name: "Test with valid name and email",
			user: &types.User{
				Name:  "John John",
				Email: "jj+test@example.ltd",
			},
			expected: "jjTest",
			s: u{
				L: func(ctx context.Context, handle string) (*types.User, error) {
					return nil, store.ErrNotFound
				},
			},
		},
		{
			name: "Test with valid name, not unique handle",
			user: &types.User{
				Name:  "John Gardener",
				Email: "johng@example.ltd",
			},
			expected: "JohnGardener_1",
			s: u{
				L: func(ctx context.Context, handle string) (*types.User, error) {
					if handle == "johng" {
						return nil, store.ErrNotUnique
					}
					return nil, store.ErrNotFound
				},
			},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			createUserHandle(ctx, tc.s, tc.user)
			assert.Equal(t, tc.expected, tc.user.Handle)
		})
	}
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

func (uu u) LookupUserByHandle(ctx context.Context, handle string) (*types.User, error) {
	return uu.L(ctx, handle)
}

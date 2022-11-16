package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Mock auth service with nil for current time, dummy provider validator and mock db
func makeMockAuthService() *auth {
	var (
		ctx = context.Background()

		mem, err = sqlite.ConnectInMemory(ctx)

		svc = &auth{
			providerValidator: func(s string) error {
				// All providers are valid.
				return nil
			},

			settings: &types.AppSettings{},
			eventbus: eventbus.New(),
		}
	)

	if err != nil {
		panic(err)
	}

	if err = store.Upgrade(ctx, zap.NewNop(), mem); err != nil {
		panic(err)
	}

	svc.store = mem

	return svc
}

func TestAuth_External(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()

		// Create some virtual user and credentials
		validUser     = &types.User{Email: "valid@test.cortezaproject.org", ID: nextID(), CreatedAt: *now()}
		suspendedUser = &types.User{Email: "suspended@test.cortezaproject.org", ID: nextID(), CreatedAt: *now(), SuspendedAt: now()}

		freshProfileID = func() string {
			return fmt.Sprintf("fresh-profile-id-%d", nextID())
		}

		fooCredentials = &types.Credential{
			ID:          nextID(),
			OwnerID:     validUser.ID,
			Label:       "credentials for foo provider",
			Kind:        "foo",
			Credentials: freshProfileID(),
			CreatedAt:   time.Time{},
		}

		barCredentials = &types.Credential{
			ID:          nextID(),
			OwnerID:     validUser.ID,
			Label:       "credentials for bar provider",
			Kind:        "bar",
			Credentials: freshProfileID(),
			CreatedAt:   time.Time{},
		}

		cases = []struct {
			name    string
			profile types.ExternalAuthUser
			user    *types.User
			err     error
		}{
			{
				"matching by user email",
				types.ExternalAuthUser{goth.User{UserID: freshProfileID(), Provider: "-", Email: validUser.Email}},
				validUser,
				nil},
			{
				"unknown profile",
				types.ExternalAuthUser{goth.User{UserID: freshProfileID(), Provider: "-", Email: "fresh-from-foo@test.cortezaproject.org"}},
				&types.User{Email: "fresh-from-foo@test.cortezaproject.org"},
				nil},
			{
				"profile match by provider ID",
				types.ExternalAuthUser{goth.User{UserID: fooCredentials.Credentials, Provider: fooCredentials.Kind, Email: "valid+2nd+email@test.cortezaproject.org"}},
				validUser,
				nil},
		}
	)

	svc := makeMockAuthService()
	svc.settings.Auth.External.Enabled = true
	req.NoError(store.TruncateUsers(ctx, svc.store))
	req.NoError(store.TruncateCredentials(ctx, svc.store))
	req.NoError(store.CreateUser(ctx, svc.store, validUser, suspendedUser))
	req.NoError(store.CreateCredential(ctx, svc.store, fooCredentials, barCredentials))

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req = require.New(t)

			var (
				ru, rerr = svc.External(ctx, c.profile)
			)

			if c.err != nil {
				req.EqualError(unwrapGeneric(rerr), c.err.Error())
				return
			}

			req.NoError(unwrapGeneric(rerr))
			req.NotNil(ru)

			if c.user == nil {
				panic("invalid test case, user should not be nil")
			}

			req.Equal(c.user.Email, ru.Email)
		})
	}
}

func TestAuth_InternalSignUp(t *testing.T) {
	var (
		req = require.New(t)

		ctx = context.Background()
		svc = makeMockAuthService()

		existingUserID = id.Next()
	)

	svc.settings.Auth.Internal.Enabled = true
	svc.settings.Auth.Internal.Signup.Enabled = true

	req.NoError(store.CreateUser(ctx, svc.store, &types.User{Email: "existing@internal-signup-test.tld", ID: existingUserID, CreatedAt: *now()}))
	req.NoError(svc.SetPassword(ctx, existingUserID, "secure password"))

	t.Run("invalid email", func(t *testing.T) {
		var (
			req = require.New(t)
		)

		u, err := svc.InternalSignUp(ctx, &types.User{}, "")
		req.Nil(u)
		req.EqualError(err, AuthErrInvalidEmailFormat().Error())
	})

	t.Run("invalid handle", func(t *testing.T) {
		var (
			req = require.New(t)
		)

		u, err := svc.InternalSignUp(ctx, &types.User{Email: "new@internal-signup-test.tld", Handle: "123"}, "")
		req.Nil(u)
		req.EqualError(err, AuthErrInvalidHandle().Error())
	})

	t.Run("invalid password", func(t *testing.T) {
		var (
			req = require.New(t)
		)

		u, err := svc.InternalSignUp(ctx, &types.User{Email: "new@internal-signup-test.tld"}, "")
		req.Nil(u)
		req.EqualError(err, AuthErrPasswordNotSecure().Error())
	})

	t.Run("valid input", func(t *testing.T) {
		var (
			req = require.New(t)
		)

		u, err := svc.InternalSignUp(ctx, &types.User{Email: "new@internal-signup-test.tld"}, "secure password")
		req.NoError(err)
		req.NotNil(u)
	})

	t.Run("existing user", func(t *testing.T) {
		var (
			req = require.New(t)
		)

		u, err := svc.InternalSignUp(ctx, &types.User{Email: "existing@internal-signup-test.tld"}, "secure password")
		req.NoError(err)
		req.NotNil(u)
		req.Equal(existingUserID, u.ID)
	})

	t.Run("invalid password for existing user", func(t *testing.T) {
		var (
			req = require.New(t)
		)

		u, err := svc.InternalSignUp(ctx, &types.User{Email: "existing@internal-signup-test.tld"}, "invalid password")
		req.EqualError(err, AuthErrInvalidCredentials().Error())
		req.Nil(u)
	})
}

func TestAuth_InternalLogin(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()

		validPass     = "this is a valid password !! 42"
		validUser     = &types.User{Email: "valid@test.cortezaproject.org", ID: nextID(), CreatedAt: *now(), EmailConfirmed: true}
		suspendedUser = &types.User{Email: "suspended@test.cortezaproject.org", ID: nextID(), CreatedAt: *now(), SuspendedAt: now()}

		tests = []struct {
			name     string
			email    string
			password string
			err      error
		}{
			{
				"with no email",
				"",
				"",
				fmt.Errorf("invalid email")},
			{
				"with bad email",
				"test",
				"",
				fmt.Errorf("invalid email")},
			{
				"with empty password",
				"test@domain.tld",
				"",
				fmt.Errorf("invalid username and password combination")},
			{
				"with valid credentials",
				validUser.Email,
				validPass,
				nil},
			{
				"with invalid password",
				validUser.Email,
				"invalid password",
				fmt.Errorf("invalid username and password combination")},
			{
				"with suspended user",
				suspendedUser.Email,
				validPass,
				fmt.Errorf("invalid username and password combination")},
		}
	)

	svc := makeMockAuthService()
	svc.settings.Auth.Internal.Enabled = true
	req.NoError(store.TruncateUsers(ctx, svc.store))
	req.NoError(store.TruncateCredentials(ctx, svc.store))
	req.NoError(store.CreateUser(ctx, svc.store, validUser, suspendedUser))
	req.NoError(svc.SetPasswordCredentials(ctx, validUser.ID, validPass))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req = require.New(t)

			var (
				usr, err = svc.InternalLogin(ctx, tt.email, tt.password)
			)

			if tt.err == nil {
				req.NoError(err)
				req.NotNil(usr)
			} else {
				req.EqualError(err, tt.err.Error())
			}

		})
	}
}

func TestAuth_createUserToken(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()

		validUser = &types.User{Email: "valid@test.cortezaproject.org", ID: nextID(), CreatedAt: *now(), EmailConfirmed: true}

		tests = []struct {
			name string
			user *types.User
			kind string
			err  error
		}{
			{
				"no user",
				nil,
				"",
				AuthErrGeneric(),
			},
			{
				"zero ID",
				&types.User{},
				"",
				AuthErrGeneric(),
			},
			{
				"valid user",
				validUser,
				credentialsTypeResetPasswordToken,
				nil,
			},
		}
	)

	svc := makeMockAuthService()
	req.NoError(store.TruncateUsers(ctx, svc.store))
	req.NoError(store.TruncateCredentials(ctx, svc.store))
	req.NoError(store.CreateUser(ctx, svc.store, validUser))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req = require.New(t)

			var (
				token, err = svc.createUserToken(ctx, tt.user, tt.kind)
			)

			if tt.err == nil {
				req.NoError(err)
				req.NotEmpty(token)
			} else {
				req.EqualError(err, tt.err.Error())
			}

		})
	}
}

// ensure that existing password reset tokens are invalidated AND rate limiting kicks in
func TestAuth_multiCreateUserTokenForPasswordReset(t *testing.T) {
	var (
		err           error
		pToken, token string

		svc = makeMockAuthService()

		req = require.New(t)
		ctx = context.Background()

		validUser = &types.User{Email: "valid@test.cortezaproject.org", ID: nextID(), CreatedAt: *now(), EmailConfirmed: true}

		// load credentials from token
		t2c = func(token string) *types.Credential {
			id, _ := validateToken(token)
			req.NotZero(id)
			c, err := store.LookupCredentialByID(ctx, svc.store, id)
			req.NoError(err)
			return c
		}
	)

	req.NoError(store.TruncateUsers(ctx, svc.store))
	req.NoError(store.TruncateCredentials(ctx, svc.store))
	req.NoError(store.CreateUser(ctx, svc.store, validUser))

	for try := 0; try <= tokenReqMaxCount+1; try++ {
		token, err = svc.createUserToken(ctx, validUser, credentialsTypeResetPasswordToken)
		t.Log("got token", token)

		if try == tokenReqMaxCount+1 {
			t.Log("rate limiting should kicked in")
			req.EqualError(err, AuthErrRateLimitExceeded().Error())
		} else {
			if try > 0 {
				t.Log("checking if previous token", pToken, "is deleted")
				req.NotNil(t2c(pToken).DeletedAt)
			}

			req.NoError(err)
			req.Nil(t2c(token).DeletedAt)
			pToken = token
		}
	}

}

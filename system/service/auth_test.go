package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/sqlite"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
	req.NoError(svc.store.TruncateUsers(ctx))
	req.NoError(svc.store.TruncateCredentials(ctx))
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

	req.NoError(svc.store.CreateUser(ctx, &types.User{Email: "existing@internal-signup-test.tld", ID: existingUserID, CreatedAt: *now()}))
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
	req.NoError(svc.store.TruncateUsers(ctx))
	req.NoError(svc.store.TruncateCredentials(ctx))
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
	req.NoError(svc.store.TruncateUsers(ctx))
	req.NoError(svc.store.TruncateCredentials(ctx))
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

	req.NoError(svc.store.TruncateUsers(ctx))
	req.NoError(svc.store.TruncateCredentials(ctx))
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

func Test_auth_checkPassword(t *testing.T) {
	plainPassword := " ... plain password ... "
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	type args struct {
		password string
		cc       types.CredentialSet
	}
	tests := []struct {
		name string
		args args
		rval bool
	}{
		{
			name: "empty set",
			rval: false,
			args: args{}},
		{
			name: "bad pwd",
			rval: false,
			args: args{
				password: " foo ",
				cc:       types.CredentialSet{&types.Credential{ID: 1, Credentials: string(hashedPassword)}}}},
		{
			name: "invalid credentials",
			rval: false,
			args: args{
				password: " foo ",
				cc:       types.CredentialSet{&types.Credential{ID: 0, Credentials: string(hashedPassword)}}}},
		{
			name: "ok",
			rval: true,
			args: args{
				password: plainPassword,
				cc:       types.CredentialSet{&types.Credential{ID: 1, Credentials: string(hashedPassword)}}}},
		{
			name: "multipass",
			rval: true,
			args: args{
				password: plainPassword,
				cc: types.CredentialSet{
					&types.Credential{ID: 0, Credentials: string(hashedPassword)},
					&types.Credential{ID: 1, Credentials: "$2a$10$8sOZxfZinxnu3bAtpkqEx.wBBwOfci6aG1szgUyxm5.BL2WiLu.ni"},
					&types.Credential{ID: 2, Credentials: string(hashedPassword)},
					&types.Credential{ID: 3, Credentials: ""},
				}}},
	}

	svc := auth{
		settings: &types.AppSettings{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.rval != svc.checkPassword(tt.args.password, tt.args.cc) {
				t.Errorf("auth.checkPassword() expecting rval to be %v", tt.rval)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name            string
		args            args
		wantID          uint64
		wantCredentials string
	}{
		{
			name:            "empty",
			wantID:          0,
			wantCredentials: "",
			args:            args{token: ""}},
		{
			name:            "foo",
			wantID:          0,
			wantCredentials: "",
			args:            args{token: "foo1"}},
		{
			name:            "semivalid",
			wantID:          0,
			wantCredentials: "",
			args:            args{token: "foofoofoofoofoofoofoofoofoofoofo0"}},
		{
			name:            "valid",
			wantID:          1,
			wantCredentials: "foofoofoofoofoofoofoofoofoofoofo",
			args:            args{token: "foofoofoofoofoofoofoofoofoofoofo1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotCredentials := validateToken(tt.args.token)

			if gotID != tt.wantID {
				t.Errorf("auth.validateToken() gotID = %v, want %v", gotID, tt.wantID)
			}
			if gotCredentials != tt.wantCredentials {
				t.Errorf("auth.validateToken() gotCredentials = %v, want %v", gotCredentials, tt.wantCredentials)
			}
		})
	}
}

package handlers

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/quasoft/memstore"
	"github.com/stretchr/testify/require"
)

type (
	expectPayload struct {
		kv    map[string]string
		alert []request.Alert
	}
)

func Test_loginForm_setValues(t *testing.T) {
	var (
		ctx      = context.Background()
		memStore = memstore.NewMemStore()
		user     = makeMockUser(ctx)

		req = &http.Request{}

		authService  *mockAuthService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		rq = require.New(t)
	)

	service.CurrentSettings = &types.AppSettings{}
	service.CurrentSettings.Auth.Internal.Enabled = true

	authSettings := &settings.Settings{}

	authService = prepareClientAuthService(ctx, user, memStore)
	authReq = prepareClientAuthReq(ctx, req, user, memStore)
	authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

	payload := map[string]string{"key": "value"}
	authReq.SetKV(payload)

	authHandlers.Settings = &settings.Settings{
		EmailConfirmationRequired: true,
	}

	err := authHandlers.loginForm(authReq)

	rq.NoError(err)
	rq.Equal(TmplLogin, authReq.Template)
	rq.Equal(payload, authReq.Data["form"])
}

func Test_loginProc(t *testing.T) {
	var (
		ctx      = context.Background()
		memStore = memstore.NewMemStore()
		user     = makeMockUser(ctx)

		req = &http.Request{}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}

		rq = require.New(t)
	)

	service.CurrentSettings = &types.AppSettings{}

	tcc := []testingExpect{
		{
			name:    "successful login",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "primary", Text: "You are now logged-in", Html: ""}},
			link:    GetLinks().Profile,
			fn: func() {
				authService = &authServiceMocked{
					internalLogin: func(ctx context.Context, email, password string) (u *types.User, err error) {
						u = &types.User{Meta: &types.UserMeta{}}
						u.Meta.SecurityPolicy.MFA.EnforcedEmailOTP = true
						u.Meta.SecurityPolicy.MFA.EnforcedTOTP = false

						err = nil
						return
					},
				}
			},
		},
		{
			name:    "internal login is not enabled",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "danger", Text: "Local accounts disabled", Html: ""}},
			link:    GetLinks().Profile,
			fn: func() {
				authService = &authServiceMocked{
					internalLogin: func(ctx context.Context, email, password string) (u *types.User, err error) {
						err = service.AuthErrInternalLoginDisabledByConfig()
						return
					},
				}
			},
		},
		{
			name:    "invalid email format",
			payload: map[string]string{"email": "email@", "error": "invalid email"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Login,
			fn: func() {
				req.PostForm.Add("email", "email@")

				authService = &authServiceMocked{
					internalLogin: func(ctx context.Context, email, password string) (u *types.User, err error) {
						err = service.AuthErrInvalidEmailFormat()
						return
					},
				}
			},
		},
		{
			name:    "invalid credentials",
			payload: map[string]string{"email": "mockuser@example.tld", "error": "invalid username and password combination"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Login,
			fn: func() {
				req.PostForm.Add("email", "mockuser@example.tld")

				authService = &authServiceMocked{
					internalLogin: func(ctx context.Context, email, password string) (u *types.User, err error) {
						err = service.AuthErrInvalidCredentials()
						return
					},
				}
			},
		},
		{
			name:    "credentials linked to invalid user",
			payload: map[string]string{"email": "mockuser@example.tld", "error": "credentials {credentials.kind} linked to disabled or deleted user {user}"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Login,
			fn: func() {
				req.PostForm.Add("email", "mockuser@example.tld")

				authService = &authServiceMocked{
					internalLogin: func(ctx context.Context, email, password string) (u *types.User, err error) {
						err = service.AuthErrCredentialsLinkedToInvalidUser()
						return
					},
				}
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			// reset from previous
			req.Form = url.Values{}
			req.PostForm = url.Values{}
			user.Meta = &types.UserMeta{}

			tc.fn()

			authReq = prepareClientAuthReq(ctx, req, user, memStore)
			authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

			err := authHandlers.loginProc(authReq)

			rq.NoError(err)
			rq.Equal(tc.payload, authReq.GetKV())
			rq.Equal(tc.alerts, authReq.NewAlerts)
			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}

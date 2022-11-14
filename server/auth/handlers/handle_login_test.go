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
		user = makeMockUser()

		req = &http.Request{}

		authService  *mockAuthService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		rq = require.New(t)
	)

	service.CurrentSettings = &types.AppSettings{}
	service.CurrentSettings.Auth.Internal.Enabled = true

	authSettings := &settings.Settings{}

	authService = prepareClientAuthService()
	authHandlers = prepareClientAuthHandlers(authService, authSettings)
	authReq = prepareClientAuthReq(authHandlers, req, user)

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
		user = makeMockUser()

		req = &http.Request{}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq
	)

	service.CurrentSettings = &types.AppSettings{}

	tcc := []testingExpect{
		{
			name:    "successful login",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "primary", Text: "login.alerts.logged-in", Html: ""}},
			link:    GetLinks().Profile,
			fn: func(_ *settings.Settings) {
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
			alerts:  []request.Alert{{Type: "danger", Text: "login.alert.local-disabled", Html: ""}},
			link:    GetLinks().Profile,
			fn: func(_ *settings.Settings) {
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
			fn: func(_ *settings.Settings) {
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
			fn: func(_ *settings.Settings) {
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
			payload: map[string]string{"email": "mockuser@example.tld", "error": "credentials {{credentials.kind}} linked to disabled or deleted user {{user}}"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Login,
			fn: func(*settings.Settings) {
				req.PostForm.Add("email", "mockuser@example.tld")

				authService = &authServiceMocked{
					internalLogin: func(ctx context.Context, email, password string) (u *types.User, err error) {
						err = service.AuthErrCredentialsLinkedToInvalidUser()
						return
					},
				}
			},
		},
		{
			name:    "split credentials check",
			payload: map[string]string{"email": "mockuser@example.tld"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Login,
			fn: func(authSettings *settings.Settings) {
				req.PostForm.Add("email", "mockuser@example.tld")
				authSettings.SplitCredentialsCheck = true

				authService = &authServiceMocked{
					passwordSet: func(ctx context.Context, email string) bool {
						return false
					},
				}
			},
		},
		{
			name:    "split credentials check with providers",
			payload: map[string]string{"email": "mockuser@example.tld"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().External + "/test-idp",
			fn: func(authSettings *settings.Settings) {
				req.PostForm.Add("email", "mockuser@example.tld")
				authSettings.SplitCredentialsCheck = true
				authSettings.Providers = []settings.Provider{
					{Handle: "test-idp"},
				}

				authService = &authServiceMocked{
					passwordSet: func(ctx context.Context, email string) bool {
						return false
					},
				}
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			// reset from previous
			req.Form = url.Values{}
			req.PostForm = url.Values{}
			user.Meta = &types.UserMeta{}

			authSettings := &settings.Settings{}

			tc.fn(authSettings)

			authHandlers = prepareClientAuthHandlers(authService, authSettings)
			authReq = prepareClientAuthReq(authHandlers, req, user)

			err := authHandlers.loginProc(authReq)

			rq.NoError(err)
			rq.Equal(tc.payload, authReq.GetKV())
			rq.Equal(tc.alerts, authReq.NewAlerts)
			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}

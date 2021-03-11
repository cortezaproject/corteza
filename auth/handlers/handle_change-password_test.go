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

func Test_changePasswordForm_setValues(t *testing.T) {
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

	err := authHandlers.changePasswordForm(authReq)

	rq.NoError(err)
	rq.Equal(TmplChangePassword, authReq.Template)
	rq.Equal(payload, authReq.Data["form"])
}

func Test_changePasswordProc(t *testing.T) {
	var (
		ctx      = context.Background()
		memStore = memstore.NewMemStore()
		user     = makeMockUser(ctx)

		req = &http.Request{}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		rq = require.New(t)

		authSettings = &settings.Settings{}
	)
	service.CurrentSettings = &types.AppSettings{}

	tcc := []testingExpect{
		{
			name:    "successful password change",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "primary", Text: "Password successfully changed.", Html: ""}},
			link:    GetLinks().Profile,
			fn: func() {
				authService = &authServiceMocked{
					changePassword: func(ctx context.Context, userID uint64, oldPassword, newPassword string) (err error) {
						return nil
					},
				}
			},
		},
		{
			name:    "provided password is not secure",
			payload: map[string]string{"error": "provided password is not secure; use longer password with more non-alphanumeric character"},
			link:    GetLinks().ChangePassword,
			fn: func() {
				authService = &authServiceMocked{
					changePassword: func(ctx context.Context, userID uint64, oldPassword, newPassword string) (err error) {
						return service.AuthErrPasswordNotSecure()
					},
				}
			},
		},
		{
			name:    "internal login is not enabled",
			payload: map[string]string{"error": "internal login (username/password) is disabled"},
			link:    GetLinks().ChangePassword,
			fn: func() {
				authService = &authServiceMocked{
					changePassword: func(ctx context.Context, userID uint64, oldPassword, newPassword string) (err error) {
						return service.AuthErrInternalLoginDisabledByConfig()
					},
				}
			},
		},
		{
			name:    "password change failed old password does not match",
			payload: map[string]string{"error": "failed to change password, old password does not match"},
			link:    GetLinks().ChangePassword,
			fn: func() {
				authService = &authServiceMocked{
					changePassword: func(ctx context.Context, userID uint64, oldPassword, newPassword string) (err error) {
						return service.AuthErrPasswodResetFailedOldPasswordCheckFailed()
					},
				}
			},
		},
		{
			name:    "password change failed for unknown user",
			payload: map[string]string{"error": "failed to change password for the unknown user"},
			link:    GetLinks().ChangePassword,
			fn: func() {
				authService = &authServiceMocked{
					changePassword: func(ctx context.Context, userID uint64, oldPassword, newPassword string) (err error) {
						return service.AuthErrPasswordChangeFailedForUnknownUser()
					},
				}
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			req.PostForm = url.Values{}

			tc.fn()

			authReq = prepareClientAuthReq(ctx, req, user, memStore)
			authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

			err := authHandlers.changePasswordProc(authReq)

			rq.NoError(err)
			rq.Equal(tc.payload, authReq.GetKV())

			if tc.alerts != nil {
				rq.Equal(tc.alerts, authReq.NewAlerts)
			}

			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}

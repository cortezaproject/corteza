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

func Test_mfaProc(t *testing.T) {
	var (
		ctx  = context.Background()
		user = makeMockUser(ctx)

		req = &http.Request{}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}
	)

	service.CurrentSettings = &types.AppSettings{}

	tcc := []testingExpect{
		{
			name:    "Email: successful login",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "primary", Text: "Email OTP valid"}},
			link:    GetLinks().Profile,
			fn: func() {
				req.Form.Set("action", "verifyEmailOtp")
				req.PostForm.Add("code", "123456")

				authService = &authServiceMocked{
					validateEmailOTP: func(ctx context.Context, code string) (err error) {
						return nil
					},
				}
			},
		},
		{
			name:    "TOTP: successful login",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "primary", Text: "TOTP valid"}},
			link:    GetLinks().Mfa,
			fn: func() {
				req.Form.Set("action", "verifyTotp")
				req.PostForm.Add("code", "123456")

				authService = &authServiceMocked{
					validateTOTP: func(ctx context.Context, code string) (err error) {
						return nil
					},
				}
			},
		},
		{
			name:    "Email: disabled",
			payload: map[string]string{"emailOtpError": "multi factor authentication with email OTP is disabled"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Mfa,
			fn: func() {
				req.Form.Set("action", "verifyEmailOtp")

				authService = &authServiceMocked{
					validateEmailOTP: func(ctx context.Context, code string) (err error) {
						return service.AuthErrDisabledMFAWithEmailOTP()
					},
				}
			},
		},
		{
			name:    "Email: auth failed for disabled user",
			payload: map[string]string{"emailOtpError": "invalid username and password combination"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Mfa,
			fn: func() {
				req.Form.Set("action", "verifyEmailOtp")

				authService = &authServiceMocked{
					validateEmailOTP: func(ctx context.Context, code string) (err error) {
						return service.AuthErrFailedForUnknownUser()
					},
				}
			},
		},
		{
			name:    "Email: invalid token",
			payload: map[string]string{"emailOtpError": "invalid code"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Mfa,
			fn: func() {
				req.Form.Set("action", "verifyEmailOtp")
				req.PostForm.Add("code", "token_TOO_LONG")

				authService = &authServiceMocked{
					validateEmailOTP: func(ctx context.Context, code string) (err error) {
						return service.AuthErrInvalidEmailOTP()
					},
				}
			},
		},
		{
			name:    "Email: no token in credentials db",
			payload: map[string]string{"emailOtpError": "invalid code"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Mfa,
			fn: func() {
				req.Form.Set("action", "verifyEmailOtp")
				req.PostForm.Add("code", "123456")

				authService = &authServiceMocked{
					validateEmailOTP: func(ctx context.Context, code string) (err error) {
						return service.AuthErrInvalidEmailOTP()
					},
				}
			},
		},
		{
			name:    "TOTP: disabled",
			payload: map[string]string{"totpError": "multi factor authentication with TOTP is disabled"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Mfa,
			fn: func() {
				req.Form.Set("action", "verifyTotp")

				authService = &authServiceMocked{
					validateTOTP: func(ctx context.Context, code string) (err error) {
						return service.AuthErrDisabledMFAWithTOTP()
					},
				}
			},
		},
		{
			name:    "TOTP: auth failed for disabled user",
			payload: map[string]string{"totpError": "invalid username and password combination"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Mfa,
			fn: func() {
				req.Form.Set("action", "verifyTotp")

				authService = &authServiceMocked{
					validateTOTP: func(ctx context.Context, code string) (err error) {
						return service.AuthErrFailedForUnknownUser()
					},
				}
			},
		},
		{
			name:    "TOTP: invalid token",
			payload: map[string]string{"totpError": "invalid code"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Mfa,
			fn: func() {
				req.Form.Set("action", "verifyTotp")
				req.PostForm.Add("code", "token_TOO_LONG")

				authService = &authServiceMocked{
					validateTOTP: func(ctx context.Context, code string) (err error) {
						return service.AuthErrInvalidTOTP()
					},
				}
			},
		},
		{
			name:    "TOTP: no token in credentials db",
			payload: map[string]string{"totpError": "invalid code"},
			alerts:  []request.Alert(nil),
			link:    GetLinks().Mfa,
			fn: func() {
				req.Form.Set("action", "verifyTotp")
				req.PostForm.Add("code", "123456")

				authService = &authServiceMocked{
					validateTOTP: func(ctx context.Context, code string) (err error) {
						return service.AuthErrInvalidTOTP()
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

			tc.fn()

			authReq = prepareClientAuthReq(ctx, req, user)
			authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

			err := authHandlers.mfaProc(authReq)

			rq.NoError(err)
			rq.Equal(tc.payload, authReq.GetKV())
			rq.Equal(tc.alerts, authReq.NewAlerts)
			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}

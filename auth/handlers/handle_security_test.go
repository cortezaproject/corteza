package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func Test_securityForm(t *testing.T) {
	var (
		ctx  = context.Background()
		user = makeMockUser(ctx)

		req = &http.Request{
			URL: &url.URL{},
		}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}

		rq = require.New(t)
	)

	authReq = prepareClientAuthReq(ctx, req, user)
	authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

	err := authHandlers.securityForm(authReq)

	rq.NoError(err)
	rq.Equal(TmplSecurity, authReq.Template)
	rq.Equal(true, authReq.Data["emailOtpEnforced"])
	rq.Equal(false, authReq.Data["totpEnforced"])
}

func Test_securityProc(t *testing.T) {
	var (
		ctx  = context.Background()
		user = makeMockUser(ctx)

		req = &http.Request{
			PostForm: url.Values{},
		}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}
	)

	tcc := []testingExpect{
		{
			name:    "reconfigureTOTP",
			link:    GetLinks().MfaTotpNewSecret,
			payload: map[interface{}]interface{}{},
			fn: func() {
				req.Form.Set("action", "reconfigureTOTP")
			},
		},
		{
			name:    "disableTOTP",
			link:    GetLinks().MfaTotpDisable,
			payload: map[interface{}]interface{}{"totpSecret": "SECRET_VALUE"},
			fn: func() {
				req.Form.Set("action", "disableTOTP")
			},
		},
		{
			name:    "disableEmailOTP failure",
			err:     "custom error",
			link:    GetLinks().Security,
			payload: map[interface{}]interface{}{"totpSecret": "SECRET_VALUE"},
			fn: func() {
				req.Form.Set("action", "disableEmailOTP")

				authService = &authServiceMocked{
					configureEmailOTP: func(c context.Context, u uint64, b bool) (user *types.User, err error) {
						return nil, errors.New("custom error")
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

			tc.fn()

			authReq = prepareClientAuthReq(ctx, req, user)
			authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

			authReq.Session.Values = map[interface{}]interface{}{"totpSecret": "SECRET_VALUE"}

			err := authHandlers.securityProc(authReq)

			rq.Equal(tc.payload, authReq.Session.Values)
			rq.Equal(tc.link, authReq.RedirectTo)

			if tc.alerts != nil {
				rq.Equal(tc.alerts, authReq.NewAlerts)
			}

			if tc.err != "" {
				rq.Equal(tc.err, err.Error())
			}
		})
	}
}

func Test_securityProcDisableEmailOTPSuccess(t *testing.T) {
	var (
		ctx  = context.Background()
		user = makeMockUser(ctx)

		req = &http.Request{
			Form:     url.Values{},
			PostForm: url.Values{},
		}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}

		rq = require.New(t)
	)

	req.Form.Set("action", "disableEmailOTP")

	authService = &authServiceMocked{
		configureEmailOTP: func(c context.Context, u uint64, b bool) (*types.User, error) {
			return user, nil
		},
	}

	authReq = prepareClientAuthReq(ctx, req, user)
	authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

	authReq.Session.Values = map[interface{}]interface{}{"totpSecret": "SECRET_VALUE"}

	err := authHandlers.securityProc(authReq)

	rq.NoError(err)
	rq.Equal([]request.Alert{{Type: "primary", Text: "Two factor authentication with TOTP disabled", Html: ""}}, authReq.NewAlerts)
}

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

func Test_profileForm(t *testing.T) {
	var (
		user = makeMockUser()

		req = &http.Request{
			URL: &url.URL{},
		}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{
			EmailConfirmationRequired: false,
		}

		rq = require.New(t)
	)

	authHandlers = prepareClientAuthHandlers(authService, authSettings)
	authReq = prepareClientAuthReq(authHandlers, req, user)

	userForm := map[string]string{
		"email":             user.Email,
		"handle":            user.Handle,
		"name":              user.Name,
		"preferredLanguage": user.Meta.PreferredLanguage,
	}

	authReq.SetKV(map[string]string{})

	err := authHandlers.profileForm(authReq)

	rq.NoError(err)
	rq.Equal(TmplProfile, authReq.Template)
	rq.Equal(authReq.Data["form"], userForm)
	rq.Equal(authReq.Data["emailConfirmationRequired"], false)
}

func Test_profileFormProc(t *testing.T) {
	var (
		user = makeMockUser()

		req = &http.Request{
			PostForm: url.Values{},
		}

		authService  authService
		userService  userService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq
	)

	tcc := []testingExpect{
		{
			name:    "success",
			err:     "",
			alerts:  []request.Alert{{Type: "primary", Text: "profile.alerts.profile-updated", Html: ""}},
			link:    GetLinks().Profile,
			payload: map[string]string(nil),
			fn: func(_ *settings.Settings) {
				req.PostForm.Add("handle", "handle")
				req.PostForm.Add("name", "name")

				authService = &authServiceMocked{}

				userService = &userServiceMocked{
					update: func(c context.Context, u *types.User) (*types.User, error) {
						u = makeMockUser()
						u.SetRoles()

						return u, nil
					},
				}
			},
		},
		{
			name:    "proc invalid ID",
			err:     "",
			alerts:  []request.Alert{{Type: "danger", Text: "profile.alerts.profile-update-fail", Html: ""}},
			link:    GetLinks().Profile,
			payload: map[string]string{"email": "mockuser@example.tld", "error": "invalid ID", "handle": "handle", "name": "name"},
			fn: func(_ *settings.Settings) {
				req.PostForm.Add("handle", "handle")
				req.PostForm.Add("name", "name")

				userService = &userServiceMocked{
					update: func(c context.Context, u *types.User) (*types.User, error) {
						return nil, service.UserErrInvalidID()
					},
				}
			},
		},
		{
			name:    "proc invalid handle",
			err:     "",
			alerts:  []request.Alert{{Type: "danger", Text: "profile.alerts.profile-update-fail", Html: ""}},
			link:    GetLinks().Profile,
			payload: map[string]string{"email": "mockuser@example.tld", "error": "invalid handle", "handle": "handle", "name": "name"},
			fn: func(_ *settings.Settings) {
				req.PostForm.Add("handle", "handle")
				req.PostForm.Add("name", "name")

				userService = &userServiceMocked{
					update: func(c context.Context, u *types.User) (*types.User, error) {
						return nil, service.UserErrInvalidHandle()
					},
				}
			},
		},
		{
			name:    "proc invalid email",
			err:     "",
			alerts:  []request.Alert{{Type: "danger", Text: "profile.alerts.profile-update-fail", Html: ""}},
			link:    GetLinks().Profile,
			payload: map[string]string{"email": "mockuser@example.tld", "error": "invalid email", "handle": "handle", "name": "name"},
			fn: func(_ *settings.Settings) {
				req.PostForm.Add("handle", "handle")
				req.PostForm.Add("name", "name")

				userService = &userServiceMocked{
					update: func(c context.Context, u *types.User) (*types.User, error) {
						return nil, service.UserErrInvalidEmail()
					},
				}
			},
		},
		{
			name:    "proc handle not unique",
			err:     "",
			alerts:  []request.Alert{{Type: "danger", Text: "profile.alerts.profile-update-fail", Html: ""}},
			link:    GetLinks().Profile,
			payload: map[string]string{"email": "mockuser@example.tld", "error": "handle not unique", "handle": "handle", "name": "name"},
			fn: func(_ *settings.Settings) {
				req.PostForm.Add("handle", "handle")
				req.PostForm.Add("name", "name")

				userService = &userServiceMocked{
					update: func(c context.Context, u *types.User) (*types.User, error) {
						return nil, service.UserErrHandleNotUnique()
					},
				}
			},
		},
		{
			name:    "user.errors.notAllowedToUpdate",
			err:     "",
			alerts:  []request.Alert{{Type: "danger", Text: "profile.alerts.profile-update-fail", Html: ""}},
			link:    GetLinks().Profile,
			payload: map[string]string{"email": "mockuser@example.tld", "error": "not allowed to update this user", "handle": "handle", "name": "name"},
			fn: func(_ *settings.Settings) {
				req.PostForm.Add("handle", "handle")
				req.PostForm.Add("name", "name")

				userService = &userServiceMocked{
					update: func(c context.Context, u *types.User) (*types.User, error) {
						return nil, service.UserErrNotAllowedToUpdate()
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

			authSettings := &settings.Settings{}

			tc.fn(authSettings)

			authHandlers = prepareClientAuthHandlers(authService, authSettings)
			authHandlers.UserService = userService
			authReq = prepareClientAuthReq(authHandlers, req, user)

			err := authHandlers.profileProc(authReq)

			rq.NoError(err)
			rq.Equal(tc.template, authReq.Template)
			rq.Equal(tc.payload, authReq.GetKV())
			rq.Equal(tc.alerts, authReq.NewAlerts)
			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}

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

func Test_logoutProc(t *testing.T) {
	var (
		ctx  = context.Background()
		user = makeMockUser(ctx)

		req = &http.Request{}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}

		rq = require.New(t)
	)

	service.CurrentSettings = &types.AppSettings{}
	service.CurrentSettings.Auth.Internal.Enabled = true

	authService = &authServiceMocked{}
	authReq = prepareClientAuthReq(ctx, req, user)
	authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

	req.PostForm = url.Values{}
	req.PostForm.Add("back", "/back")
	authReq.Session.Values = map[interface{}]interface{}{"key": url.Values{"key": []string{"value"}}}

	err := authHandlers.logoutProc(authReq)
	rq.NoError(err)
	rq.Empty(authReq.Session.Values)
	rq.Empty(authReq.AuthUser)
	rq.Empty(authReq.Client)
	rq.Equal("/back", authReq.Data["backlink"])
	rq.Equal(TmplLogout, authReq.Template)
}

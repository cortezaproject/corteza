package helpers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/steinfletcher/apitest"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/system/types"
)

func BindAuthMiddleware(r chi.Router) {
	r.Use(
		auth.DefaultJwtHandler.HttpVerifier(),
		auth.DefaultJwtHandler.HttpAuthenticator(),
	)
}

func ReqHeaderAuthBearer(user *types.User) apitest.Intercept {
	return func(req *http.Request) {
		if user == nil {
			req.Header.Del("Authorization")
		} else {
			req.Header.Set("Authorization", "Bearer "+auth.DefaultJwtHandler.Encode(user))
		}
	}
}

func ReqHeaderRawAuthBearer(token string) apitest.Intercept {
	return func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

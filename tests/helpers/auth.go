package helpers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/steinfletcher/apitest"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

func BindAuthMiddleware(r chi.Router) {
	r.Use(
		auth.JWT().HttpVerifier(),
		auth.JWT().HttpValidator("api"),
	)
}

func ReqHeaderRawAuthBearer(token []byte) apitest.Intercept {
	return func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+string(token))
	}
}

package auth

import (
	"errors"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"net/http"
)

func MiddlewareValidOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		if !CheckScope(ctx.Value(scopeCtxKey{}), "api") {
			api.Send(w, r, errors.New("Unauthorized scope"))
			return
		}

		if !GetIdentityFromContext(ctx).Valid() {
			api.Send(w, r, errors.New("Unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

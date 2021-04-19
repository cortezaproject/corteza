package auth

import (
	"errors"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"net/http"
)

func MiddlewareValidOnly(next http.Handler) http.Handler {
	return AccessTokenCheck("api")(next)
}

func AccessTokenCheck(scope ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			for _, s := range scope {
				if !CheckScope(ctx.Value(scopeCtxKey{}), s) {
					w.WriteHeader(http.StatusUnauthorized)
					api.Send(w, r, errors.New("unauthorized scope"))
					return
				}
			}

			if !GetIdentityFromContext(ctx).Valid() {
				w.WriteHeader(http.StatusUnauthorized)
				api.Send(w, r, errors.New("unauthorized"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

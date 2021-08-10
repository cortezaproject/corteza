package auth

import (
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/errors"
)

func MiddlewareValidOnly(next http.Handler) http.Handler {
	return AccessTokenCheck("api")(next)
}

func AccessTokenCheck(scope ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			if !GetIdentityFromContext(ctx).Valid() {
				errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
				return
			}

			for _, s := range scope {
				if !CheckScope(ctx.Value(scopeCtxKey{}), s) {
					errors.ProperlyServeHTTP(w, r, ErrUnauthorizedScope(), false)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

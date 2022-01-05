package auth

import (
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/go-chi/jwtauth"
)

func MiddlewareValidOnly(next http.Handler) http.Handler {
	return AccessTokenCheck("api")(next)
}

func AccessTokenCheck(scope ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			jwtauth.Authenticator()

			// retrieve token and claims from context
			tkn, _, err := jwtauth.FromContext(ctx)
			if err != nil || !tkn.Valid {
				errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
				return
			}

			// check valid scope
			for _, s := range scope {
				if !CheckScope(ctx.Value(scopeCtxKey{}), s) {
					errors.ProperlyServeHTTP(w, r, ErrUnauthorizedScope(), false)
					return
				}
			}

			// verify JWT from store
			_, err = DefaultJwtStore.LookupAuthOa2tokenByAccess(ctx, tkn.Raw)
			if err != nil {
				errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

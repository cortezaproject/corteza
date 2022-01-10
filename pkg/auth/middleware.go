package auth

import (
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/jwtauth"
)

func MiddlewareValidOnly(next http.Handler) http.Handler {
	return AccessTokenCheck("api")(next)
}

func AccessTokenCheck(scope ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			token, _, err := jwtauth.FromContext(ctx)
			spew.Dump(token, err)

			if err != nil {
				errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
				return
			}

			if !CheckJwtScope(token, scope...) {
				errors.ProperlyServeHTTP(w, r, ErrUnauthorizedScope(), false)
			}

			// @todo we need to check if token is in store!!
			// @todo we need to check if token is in store!!
			// @todo we need to check if token is in store!!
			// @todo we need to check if token is in store!!
			// @todo we need to check if token is in store!!
			// @todo we need to check if token is in store!!
			// @todo we need to check if token is in store!!
			// @todo we need to check if token is in store!!
			// @todo we need to check if token is in store!!
			// @todo we need to check if token is in store!!
			//
			//// verify JWT from store
			//_, err = DefaultJwtStore.LookupAuthOa2tokenByAccess(ctx, tkn.Raw)
			//if err != nil {
			//	errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
			//	return
			//}

			next.ServeHTTP(w, r)
		})
	}
}

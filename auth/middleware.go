package auth

import (
	"errors"
	"github.com/titpetric/factory/resputil"
	"net/http"
)

func AuthenticationMiddlewareValidOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		if !GetIdentityFromContext(ctx).Valid() {
			resputil.JSON(w, errors.New("Unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

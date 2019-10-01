package auth

import (
	"errors"
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func MiddlewareValidOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		if !GetIdentityFromContext(ctx).Valid() {
			resputil.JSON(w, errors.New("Unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

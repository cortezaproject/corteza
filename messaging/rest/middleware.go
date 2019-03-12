package rest

import (
	"net/http"

	"github.com/crusttech/crust/messaging/service"
)

func middlewareAllowedAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !service.DefaultPermissions.With(r.Context()).CanAccess() {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

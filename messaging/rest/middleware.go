package rest

import (
	"net/http"

	"github.com/cortezaproject/corteza-server/messaging/internal/service"
)

func middlewareAllowedAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !service.DefaultAccessControl.CanAccess(r.Context()) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

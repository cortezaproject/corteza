package crm

import (
	"net/http"
)

var pass = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (*ModuleHandlers) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*TypesHandlers) Authenticator() func(http.Handler) http.Handler {
	return pass
}

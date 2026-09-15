package server

import (
	"net/http"
	"strings"
)

// skipCORS bypasses the CORS middleware for request paths with the given prefix
func skipCORS(prefix string, cors func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		withCORS := cors(next)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, prefix) {
				next.ServeHTTP(w, r)
				return
			}

			withCORS.ServeHTTP(w, r)
		})
	}
}

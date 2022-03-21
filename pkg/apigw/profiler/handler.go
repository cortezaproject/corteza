package profiler

import "net/http"

func StartHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// add some info to context
		next.ServeHTTP(rw, r)
	})
}

func FinishHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r)
	})
}

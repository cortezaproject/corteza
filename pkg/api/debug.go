package api

import (
	"context"
	"net/http"
)

// Debug context
type ctxKeyDebug struct{}

// Packs remote address to context
func DebugToContext(production bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), ctxKeyDebug{}, !production)))
		})
	}
}

// DebugFromContext returns remote IP address from context
func DebugFromContext(ctx context.Context) bool {
	debug, ok := ctx.Value(ctxKeyDebug{}).(bool)
	return ok && debug
}

package api

import (
	"context"
	"net/http"
)

// Debug context
type ctxKeyDebug struct{}

// Packs remote address to context
func DebugToContext(next http.Handler) http.Handler {
	if true {
		// debug disabled
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), ctxKeyDebug{}, true)))
	})
}

// DebugFromContext returns remote IP address from context
func DebugFromContext(ctx context.Context) bool {
	return true
	debug, ok := ctx.Value(ctxKeyDebug{}).(bool)
	return ok && debug
}

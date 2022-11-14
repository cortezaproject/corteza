package api

import (
	"context"
	"net/http"
)

// Key to use when setting the request ID.
type ctxKeyRemoteAddr int

// RemoteAddrKey is the key that holds th unique request ID in a request context.
const remoteAddrKey ctxKeyRemoteAddr = 0

// Packs remote address to context
func RemoteAddrToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), remoteAddrKey, req.RemoteAddr)))
	})
}

// RemoteAddrFromContext returns remote IP address from context
func RemoteAddrFromContext(ctx context.Context) string {
	v := ctx.Value(remoteAddrKey)
	if str, ok := v.(string); ok {
		return str
	}

	return ""
}

package actionlog

import (
	"context"
)

// Key to use when setting the request ID.
type ctxKey int

const (
	RequestOrigin_APP_Init      = "app/init"
	RequestOrigin_APP_Upgrade   = "app/upgrade"
	RequestOrigin_APP_Activate  = "app/activate"
	RequestOrigin_APP_Provision = "app/provision"
	RequestOrigin_APP_Run       = "app/run"
	RequestOrigin_API_REST      = "api/rest"
	RequestOrigin_API_GRPC      = "api/grpc"
	RequestOrigin_Auth          = "auth"
)

// RequestOriginKey is the key that holds th unique request ID in a request context.
const requestOriginKey ctxKey = 0

// RequestOriginToContext stores request origin to context
func RequestOriginToContext(ctx context.Context, origin string) context.Context {
	return context.WithValue(ctx, requestOriginKey, origin)
}

// RequestOriginFromContext returns remote IP address from context
func RequestOriginFromContext(ctx context.Context) string {
	v := ctx.Value(requestOriginKey)
	if str, ok := v.(string); ok {
		return str
	}

	return ""
}

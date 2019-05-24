package logger

import (
	"context"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type (
	ctxLogKey struct{}
)

func ContextWithValue(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogKey{}, log)
}

func ContextValue(ctx context.Context) *zap.Logger {
	return ctx.Value(ctxLogKey{}).(*zap.Logger)
}

// NamedDefault returns default logger with requestID (from context) and extended name
func AddRequestID(ctx context.Context, log *zap.Logger) *zap.Logger {
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		log = log.With(zap.String("requestID", reqID))
	}

	return log
}

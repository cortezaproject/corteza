package logger

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type (
	ctxLogKey struct{}
)

// ContextWithValue allows us to pack custom logger to context and pass that to the
func ContextWithValue(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogKey{}, logger)
}

// ContextValue retrieves logger from given context or falls back to
// any of the logger passed to it. If no loggers are found it uses default logger from pkg/logger
func ContextValue(ctx context.Context, fallbacks ...*zap.Logger) *zap.Logger {
	if ctx != nil {
		if ctxLogger := ctx.Value(ctxLogKey{}); ctxLogger != nil {
			// This will panic if we somehow manage to set
			return ctxLogger.(*zap.Logger)
		}
	}

	for _, l := range fallbacks {
		if l != nil {
			return l
		}
	}

	return Default()
}

// AddRequestID sets requestID field from context to logger and returns it
func AddRequestID(ctx context.Context, logger *zap.Logger) *zap.Logger {
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		logger = logger.With(zap.String("requestID", reqID))
	}

	return logger
}

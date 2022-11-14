package server

import (
	"context"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func BaseMiddleware(isProduction bool, log *zap.Logger) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		handleCORS,
		locale.DetectLanguage(locale.Global()),
		middleware.RealIP,
		api.RemoteAddrToContext,
		middleware.RequestID,
		api.DebugToContext(isProduction),
		contextLogger(log),
	}
}

func panicRecovery(ctx context.Context, w http.ResponseWriter) {
	if err := recover(); err != nil {
		log := logger.ContextValue(ctx)
		if err, ok := err.(error); ok {
			log = log.With(zap.Error(err))
		} else {
			log = log.With(zap.Any("recover-value", err))
		}

		log.Debug("crashed on http request", zap.ByteString("stack", debug.Stack()))

		w.WriteHeader(500)

		if _, has := os.LookupEnv("DEBUG_DUMP_STACK_IN_RESPONSE"); has {
			// Provide nice call stack on endpoint when
			// we crash
			_, _ = w.Write(debug.Stack())
		}

		return
	}
}

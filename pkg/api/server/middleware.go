package server

import (
	"net/http"
	"os"
	"runtime/debug"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/getsentry/sentry-go/http"
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

func sentryMiddleware() func(http.Handler) http.Handler {
	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}).Handle
}

// HandlePanic sends 500 error when panic occurs inside the request call
func handlePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log := logger.Default()
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
					w.Write(debug.Stack())
				}

				return
			}
		}()

		next.ServeHTTP(w, req)
	})
}

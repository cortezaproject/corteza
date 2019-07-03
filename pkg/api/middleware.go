package api

import (
	"net/http"
	"os"
	"runtime/debug"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	sentryhttp "github.com/getsentry/sentry-go/http"
)

func Base(log *zap.Logger) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		handleCORS,
		middleware.RealIP,
		middleware.RequestID,
		contextLogger(log),
	}
}

func Sentry() func(http.Handler) http.Handler {
	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}).Handle
}

// HandlePanic sends 500 error when panic occurs inside the request call
func HandlePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
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

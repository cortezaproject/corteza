package api

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func Base() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		handleCORS,
		middleware.RealIP,
		middleware.RequestID,
	}
}

func Logging(log *zap.Logger) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		contextLogger(log),
		LogRequest,
		LogResponse,
	}
}

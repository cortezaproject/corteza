package server

import (
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// contextLogger middleware binds logger to request's context.
//
// This allows us to use logger from context (with requestID)
// inside our (generated) handlers and controllers
func contextLogger(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var requestID = middleware.GetReqID(req.Context())

			w.Header().Add("X-Request-Id", requestID)

			req = req.WithContext(logger.ContextWithValue(
				req.Context(),
				log.With(zap.String("requestID", requestID)).Named("rest"),
			))

			next.ServeHTTP(w, req)
		})
	}
}

// LogRequest middleware logs request details
//
// It uses logger from context, see contextLogger()
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var remote = req.RemoteAddr
		if l := strings.LastIndex(remote, ":"); l > -1 {
			remote = remote[:l]
		}

		logger.ContextValue(req.Context()).Info(
			"HTTP request "+req.Method+" "+req.URL.Path,
			zap.String("method", req.Method),
			zap.String("path", req.URL.Path),
			zap.Int64("size", req.ContentLength),
			zap.String("remote", remote),
		)
		next.ServeHTTP(w, req)
	})
}

// LogResponse middleware logs response details
//
// It uses logger from context, see contextLogger()
func LogResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		wrapped := middleware.NewWrapResponseWriter(w, req.ProtoMajor)
		t := time.Now()

		defer func() {
			logger.ContextValue(req.Context()).Info(
				"HTTP response "+req.Method+" "+req.URL.Path,
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", wrapped.Status()),
				zap.Int("size", wrapped.BytesWritten()),
				zap.Float64("duration", time.Since(t).Seconds()),
			)
		}()

		next.ServeHTTP(wrapped, req)
	})

}

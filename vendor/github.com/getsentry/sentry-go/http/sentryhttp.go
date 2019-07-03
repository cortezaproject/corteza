package sentryhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

type Handler struct {
	repanic         bool
	waitForDelivery bool
	timeout         time.Duration
}

type Options struct {
	// Repanic configures whether Sentry should repanic after recovery, in most cases it should be set to true,
	// as iris.Default includes it's own Recovery middleware what handles http responses.
	Repanic bool
	// WaitForDelivery configures whether you want to block the request before moving forward with the response.
	// Because Iris's default `Recovery` handler doesn't restart the application,
	// it's safe to either skip this option or set it to `false`.
	WaitForDelivery bool
	// Timeout for the event delivery requests.
	Timeout time.Duration
}

// New returns a struct that provides Handle and HandleFunc methods
// that satisfy http.Handler and http.HandlerFunc interfaces.
func New(options Options) *Handler {
	handler := Handler{
		repanic:         false,
		timeout:         time.Second * 2,
		waitForDelivery: false,
	}

	if options.Repanic {
		handler.repanic = true
	}

	if options.WaitForDelivery {
		handler.waitForDelivery = true
	}

	return &handler
}

// Handle wraps http.Handler and recovers from caught panics.
func (h *Handler) Handle(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hub := sentry.CurrentHub().Clone()
		hub.Scope().SetRequest(sentry.Request{}.FromHTTPRequest(r))
		ctx := sentry.SetHubOnContext(
			r.Context(),
			hub,
		)
		defer h.recoverWithSentry(hub, r)
		handler.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// HandleFunc wraps http.HandleFunc and recovers from caught panics.
func (h *Handler) HandleFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		hub := sentry.CurrentHub().Clone()
		hub.Scope().SetRequest(sentry.Request{}.FromHTTPRequest(r))
		ctx := sentry.SetHubOnContext(
			r.Context(),
			hub,
		)
		defer h.recoverWithSentry(hub, r)
		handler(rw, r.WithContext(ctx))
	}
}

func (h *Handler) recoverWithSentry(hub *sentry.Hub, r *http.Request) {
	if err := recover(); err != nil {
		eventID := hub.RecoverWithContext(
			context.WithValue(r.Context(), sentry.RequestContextKey, r),
			err,
		)
		if eventID != nil && h.waitForDelivery {
			hub.Flush(h.timeout)
		}
		if h.repanic {
			panic(err)
		}
	}
}

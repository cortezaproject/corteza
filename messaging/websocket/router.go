package websocket

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

func middlewareAllowedAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !service.DefaultAccessControl.CanAccess(r.Context()) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Watch(ctx context.Context) {
	events := repository.Events()

	go func() {
		defer sentry.Recover()

		for {
			if err := eq.feedSessions(ctx, events, store); err != nil {
				if err == context.Canceled {
					return
				}

				logger.Default().Error("session event feed error", zap.Error(err))
			}
		}
	}()
	eq.store(ctx, events)
}

func (ws Websocket) ApiServerRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/websocket", func(r chi.Router) {
			r.Use(middlewareAllowedAccess)
			r.Get("/", ws.Open)
		})
	})
}

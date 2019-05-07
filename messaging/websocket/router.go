package websocket

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/messaging/internal/repository"
	"github.com/crusttech/crust/messaging/internal/service"
)

func MountRoutes(ctx context.Context, config *repository.Flags) func(chi.Router) {
	return func(r chi.Router) {
		events := repository.Events()

		go func() {
			for {
				if err := eq.feedSessions(ctx, config, events, store); err != nil {
					logger.Default().Error("session event feed error", zap.Error(err))
				}
			}
		}()
		eq.store(ctx, events)

		websocket := Websocket{}.New(config)
		r.Group(func(r chi.Router) {
			r.Route("/websocket", func(r chi.Router) {
				r.Use(middlewareAllowedAccess)
				r.Get("/", websocket.Open)
			})
		})
	}
}

func middlewareAllowedAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !service.DefaultPermissions.With(r.Context()).CanAccess() {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

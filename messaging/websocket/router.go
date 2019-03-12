package websocket

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/crusttech/crust/messaging/repository"
	"github.com/crusttech/crust/messaging/service"
)

func MountRoutes(ctx context.Context, config *repository.Flags) func(chi.Router) {
	return func(r chi.Router) {
		events := repository.Events()

		go func() {
			for {
				if err := eq.feedSessions(ctx, config, events, store); err != nil {
					log.Printf("Error in sessions event feed: %+v", err)
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

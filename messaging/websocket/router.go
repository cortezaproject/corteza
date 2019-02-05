package websocket

import (
	"context"
	"log"

	"github.com/go-chi/chi"

	"github.com/crusttech/crust/messaging/repository"
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
				r.Get("/", websocket.Open)
			})
		})
	}
}

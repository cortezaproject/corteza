package websocket

import (
	"context"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/service"
	"github.com/go-chi/chi"
)

func MountRoutes(ctx context.Context) func(chi.Router) {
	return func(r chi.Router) {
		var (
			// @todo move this 1 level up & join with rest init functions
			svcUser = service.User()
		)

		repo := repository.New()

		eq.feedSessions(ctx, repo, store)
		eq.store(ctx, repo)

		websocket := Websocket{}.New(svcUser)
		r.Group(func(r chi.Router) {
			r.Route("/websocket", func(r chi.Router) {
				r.Get("/", websocket.Open)
			})
		})
	}
}

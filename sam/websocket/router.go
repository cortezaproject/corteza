package websocket

import (
	"context"
	"fmt"

	"github.com/go-chi/chi"

	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/service"
)

func MountRoutes(ctx context.Context, config *repository.Flags) func(chi.Router) {
	return func(r chi.Router) {
		var (
			// @todo move this 1 level up & join with rest init functions
			svcUser = service.User()
		)

		repo := repository.New()

		go func() {
			if err := eq.feedSessions(ctx, config, repo, store); err != nil {
				panic(fmt.Sprintf("Error when starting sessions event feed: %+v", err))
			}
		}()
		eq.store(ctx, repo)

		websocket := Websocket{}.New(svcUser, config)
		r.Group(func(r chi.Router) {
			r.Route("/websocket", func(r chi.Router) {
				r.Get("/", websocket.Open)
			})
		})
	}
}

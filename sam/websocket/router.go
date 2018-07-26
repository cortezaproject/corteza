package websocket

import (
	"github.com/crusttech/crust/sam/service"
	"github.com/go-chi/chi"
)

func MountRoutes() func(chi.Router) {
	return func(r chi.Router) {
		var (
			// @todo move this 1 level up & join with rest init functions
			svcUser = service.User()
		)

		websocket := Websocket{}.New(svcUser)
		r.Group(func(r chi.Router) {
			r.Route("/websocket", func(r chi.Router) {
				r.Get("/", websocket.Open)
			})
		})
	}
}

package websocket

import (
	"github.com/go-chi/chi"
)

func MountRoutes(r chi.Router) {
	websocket := Websocket{}.New()
	r.Group(func(r chi.Router) {
		r.Use(websocket.Authenticator())
		r.Route("/websocket", func(r chi.Router) {
			r.Get("/", websocket.Open)
		})
	})
}

package rest

import (
	"github.com/go-chi/chi"

	"github.com/cortezaproject/corteza-server/federation/rest/handlers"
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

func MountRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		handlers.NewNodeHandshake(NodeHandshake{}.New()).MountRoutes(r)
	})
	var (
		module = Module{}.New()
	)

	// Protect all _private_ routes
	r.Group(func(r chi.Router) {
		r.Use(auth.MiddlewareValidOnly)
		handlers.NewNode(Node{}.New()).MountRoutes(r)
		handlers.NewModule(module).MountRoutes(r)
	})
}

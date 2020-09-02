package rest

import (
	"github.com/go-chi/chi"

	"github.com/cortezaproject/corteza-server/federation/rest/handlers"
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

func MountRoutes(r chi.Router) {
	var (
		foobar = Foobar{}.New()
	)

	// Protect all _private_ routes
	r.Group(func(r chi.Router) {
		r.Use(auth.MiddlewareValidOnly)
		r.Use(middlewareAllowedAccess)

		handlers.NewFoobar(foobar).MountRoutes(r)
	})
}

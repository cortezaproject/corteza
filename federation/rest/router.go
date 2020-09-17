package rest

import (
	"github.com/go-chi/chi"

	"github.com/cortezaproject/corteza-server/federation/rest/handlers"
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

func MountRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		handlers.NewPairRequest(NodePairRequest{}.New()).MountRoutes(r)

		// temporary because of acl
		handlers.NewModule((Module{}.New())).MountRoutes(r)
		handlers.NewSyncStructure((SyncStructure{}.New())).MountRoutes(r)
	})

	// Protect all _private_ routes
	r.Group(func(r chi.Router) {
		r.Use(auth.MiddlewareValidOnly)
		r.Use(middlewareAllowedAccess)

		handlers.NewIdentity(NodeIdentity{}.New()).MountRoutes(r)
		handlers.NewPair(NodePair{}.New()).MountRoutes(r)
	})
}

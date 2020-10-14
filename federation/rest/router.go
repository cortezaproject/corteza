package rest

import (
	"github.com/go-chi/chi"

	"github.com/cortezaproject/corteza-server/federation/rest/handlers"
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

func MountRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		handlers.NewNodeHandshake(NodeHandshake{}.New()).MountRoutes(r)

		// temporary because of acl
		handlers.NewManageStructure((ManageStructure{}.New())).MountRoutes(r)
		handlers.NewSyncStructure((SyncStructure{}.New())).MountRoutes(r)
		handlers.NewSyncData((SyncData{}.New())).MountRoutes(r)
	})

	// Protect all _private_ routes
	r.Group(func(r chi.Router) {
		r.Use(auth.MiddlewareValidOnly)
		r.Use(middlewareAllowedAccess)

		handlers.NewNode(Node{}.New()).MountRoutes(r)
	})
}

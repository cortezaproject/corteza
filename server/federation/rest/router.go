package rest

import (
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/go-chi/chi/v5"

	"github.com/cortezaproject/corteza/server/federation/rest/handlers"
	"github.com/cortezaproject/corteza/server/pkg/auth"
)

func MountRoutes(opts options.LimitOpt) func(r chi.Router) {
	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			handlers.NewNodeHandshake(NodeHandshake{}.New()).MountRoutes(r)
		})

		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.HttpTokenValidator("api"))

			handlers.NewPermissions(Permissions{}.New()).MountRoutes(r)

			handlers.NewNode(Node{}.New()).MountRoutes(r)
			handlers.NewManageStructure((ManageStructure{}.New())).MountRoutes(r)

			handlers.NewSyncData((SyncData{}.New(opts))).MountRoutes(r)
			handlers.NewSyncStructure((SyncStructure{}.New())).MountRoutes(r)
		})
	}
}

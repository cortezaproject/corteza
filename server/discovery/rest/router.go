package rest

import (
	"github.com/go-chi/chi/v5"

	"github.com/cortezaproject/corteza/server/discovery/rest/handlers"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/options"
)

func MountRoutes(discoveryOpts options.DiscoveryOpt) func(r chi.Router) {
	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(auth.HttpTokenValidator("discovery"))

			handlers.NewResources(Resources()).MountRoutes(r)
			handlers.NewFeed(Feed()).MountRoutes(r)
			handlers.NewMappings(Mappings(discoveryOpts)).MountRoutes(r)
		})
	}
}

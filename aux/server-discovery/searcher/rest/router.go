package rest

import (
	"github.com/cortezaproject/corteza-server-discovery/pkg/auth"
	"github.com/cortezaproject/corteza-server-discovery/searcher/rest/handlers"
	"github.com/go-chi/chi/v5"
)

func MountRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(auth.HttpTokenValidator("discovery"))
			handlers.NewSearch(Search()).MountRoutes(r)
		})
	}
}

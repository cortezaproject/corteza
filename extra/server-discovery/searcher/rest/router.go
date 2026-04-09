package rest

import (
	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/auth"
	"github.com/cortezaproject/corteza/extra/server-discovery/searcher/rest/handlers"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func MountRoutes(log *zap.Logger) func(r chi.Router) {
	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			search := Search()

			r.Use(auth.HttpTokenValidator("discovery"))

			// traditional search endpoint for Traditional and vectorsearch
			handlers.NewSearch(search).MountRoutes(r)
		})
	}
}

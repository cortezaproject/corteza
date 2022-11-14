package handlers

import (
	"context"
	"github.com/cortezaproject/corteza-server-discovery/searcher/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// SearchAPI Internal API interface
	SearchAPI interface {
		SearchResources(context.Context, *request.SearchResources) (interface{}, error)
	}

	// Search HTTP API interface
	Search struct {
		SearchResources func(http.ResponseWriter, *http.Request)
	}
)

func NewSearch(h SearchAPI) *Search {
	return &Search{
		SearchResources: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSearchListResources()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.SearchResources(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Search) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/", h.SearchResources)
	})
}

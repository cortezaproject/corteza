package handlers

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"context"
	"github.com/crusttech/human/server/automation/rest/request"
	"github.com/crusttech/human/server/pkg/api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	ConstructLibraryAPI interface {
		Functions(context.Context, *request.ConstructLibraryFunctions) (interface{}, error)
		Triggers(context.Context, *request.ConstructLibraryTriggers) (interface{}, error)
	}

	// HTTP API interface
	ConstructLibrary struct {
		Functions func(http.ResponseWriter, *http.Request)
		Triggers  func(http.ResponseWriter, *http.Request)
	}
)

func NewConstructLibrary(h ConstructLibraryAPI) *ConstructLibrary {
	return &ConstructLibrary{
		Functions: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewConstructLibraryFunctions()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Functions(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Triggers: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewConstructLibraryTriggers()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Triggers(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h ConstructLibrary) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/construct-library/functions", h.Functions)
		r.Get("/construct-library/triggers", h.Triggers)
	})
}

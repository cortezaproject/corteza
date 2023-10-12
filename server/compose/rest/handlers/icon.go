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
	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	IconAPI interface {
		List(context.Context, *request.IconList) (interface{}, error)
		Upload(context.Context, *request.IconUpload) (interface{}, error)
		Delete(context.Context, *request.IconDelete) (interface{}, error)
	}

	// HTTP API interface
	Icon struct {
		List   func(http.ResponseWriter, *http.Request)
		Upload func(http.ResponseWriter, *http.Request)
		Delete func(http.ResponseWriter, *http.Request)
	}
)

func NewIcon(h IconAPI) *Icon {
	return &Icon{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewIconList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Upload: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewIconUpload()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Upload(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewIconDelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Icon) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/icon/", h.List)
		r.Post("/icon/", h.Upload)
		r.Delete("/icon/{iconID}", h.Delete)
	})
}

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
	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	SearchAPI interface {
		Messages(context.Context, *request.SearchMessages) (interface{}, error)
		Threads(context.Context, *request.SearchThreads) (interface{}, error)
	}

	// HTTP API interface
	Search struct {
		Messages func(http.ResponseWriter, *http.Request)
		Threads  func(http.ResponseWriter, *http.Request)
	}
)

func NewSearch(h SearchAPI) *Search {
	return &Search{
		Messages: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSearchMessages()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Messages(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Threads: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSearchThreads()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Threads(r.Context(), params)
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
		r.Get("/search/messages", h.Messages)
		r.Get("/search/threads", h.Threads)
	})
}

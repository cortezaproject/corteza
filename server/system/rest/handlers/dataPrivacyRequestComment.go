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
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	DataPrivacyRequestCommentAPI interface {
		List(context.Context, *request.DataPrivacyRequestCommentList) (interface{}, error)
		Create(context.Context, *request.DataPrivacyRequestCommentCreate) (interface{}, error)
	}

	// HTTP API interface
	DataPrivacyRequestComment struct {
		List   func(http.ResponseWriter, *http.Request)
		Create func(http.ResponseWriter, *http.Request)
	}
)

func NewDataPrivacyRequestComment(h DataPrivacyRequestCommentAPI) *DataPrivacyRequestComment {
	return &DataPrivacyRequestComment{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestCommentList()
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
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestCommentCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h DataPrivacyRequestComment) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/data-privacy/requests/{requestID}/comments/", h.List)
		r.Post("/data-privacy/requests/{requestID}/comments/", h.Create)
	})
}

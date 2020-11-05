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
	ActivityAPI interface {
		Send(context.Context, *request.ActivitySend) (interface{}, error)
	}

	// HTTP API interface
	Activity struct {
		Send func(http.ResponseWriter, *http.Request)
	}
)

func NewActivity(h ActivityAPI) *Activity {
	return &Activity{
		Send: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewActivitySend()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Send(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Activity) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/activity/", h.Send)
	})
}

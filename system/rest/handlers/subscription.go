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
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	SubscriptionAPI interface {
		Current(context.Context, *request.SubscriptionCurrent) (interface{}, error)
	}

	// HTTP API interface
	Subscription struct {
		Current func(http.ResponseWriter, *http.Request)
	}
)

func NewSubscription(h SubscriptionAPI) *Subscription {
	return &Subscription{
		Current: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSubscriptionCurrent()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Current(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Subscription) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/subscription/", h.Current)
	})
}

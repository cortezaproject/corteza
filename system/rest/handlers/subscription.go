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
	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/rest/request"
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
				logger.LogParamError("Subscription.Current", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Current(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Subscription.Current", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Subscription.Current", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Subscription) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/subscription/", h.Current)
	})
}

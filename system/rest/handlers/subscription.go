package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `subscription.go`, `subscription.util.go` or `subscription_test.go` to
	implement your API calls, helper functions and tests. The file `subscription.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/rest/request"
)

// Internal API interface
type SubscriptionAPI interface {
	Current(context.Context, *request.SubscriptionCurrent) (interface{}, error)
}

// HTTP API interface
type Subscription struct {
	Current func(http.ResponseWriter, *http.Request)
}

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

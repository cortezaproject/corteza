package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `webhooks.go`, `webhooks.util.go` or `webhooks_test.go` to
	implement your API calls, helper functions and tests. The file `webhooks.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/messaging/rest/request"
)

// Internal API interface
type WebhooksAPI interface {
	List(context.Context, *request.WebhooksList) (interface{}, error)
	Create(context.Context, *request.WebhooksCreate) (interface{}, error)
	Update(context.Context, *request.WebhooksUpdate) (interface{}, error)
	Get(context.Context, *request.WebhooksGet) (interface{}, error)
	Delete(context.Context, *request.WebhooksDelete) (interface{}, error)
}

// HTTP API interface
type Webhooks struct {
	List   func(http.ResponseWriter, *http.Request)
	Create func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
	Get    func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
}

func NewWebhooks(wh WebhooksAPI) *Webhooks {
	return &Webhooks{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.Create(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.Update(r.Context(), params)
			})
		},
		Get: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksGet()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.Get(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.Delete(r.Context(), params)
			})
		},
	}
}

func (wh *Webhooks) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/webhooks/", wh.List)
		r.Post("/webhooks/", wh.Create)
		r.Post("/webhooks/{webhookID}", wh.Update)
		r.Get("/webhooks/{webhookID}", wh.Get)
		r.Delete("/webhooks/{webhookID}", wh.Delete)
	})
}

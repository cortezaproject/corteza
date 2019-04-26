package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `webhooks_public.go`, `webhooks_public.util.go` or `webhooks_public_test.go` to
	implement your API calls, helper functions and tests. The file `webhooks_public.go`
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
type WebhooksPublicAPI interface {
	Delete(context.Context, *request.WebhooksPublicDelete) (interface{}, error)
	Create(context.Context, *request.WebhooksPublicCreate) (interface{}, error)
}

// HTTP API interface
type WebhooksPublic struct {
	Delete func(http.ResponseWriter, *http.Request)
	Create func(http.ResponseWriter, *http.Request)
}

func NewWebhooksPublic(wh WebhooksPublicAPI) *WebhooksPublic {
	return &WebhooksPublic{
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksPublicDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.Delete(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksPublicCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.Create(r.Context(), params)
			})
		},
	}
}

func (wh *WebhooksPublic) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Delete("/webhooks/{webhookID}/{webhookToken}", wh.Delete)
		r.Post("/webhooks/{webhookID}/{webhookToken}", wh.Create)
	})
}

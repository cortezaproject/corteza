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

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
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

func NewWebhooks(h WebhooksAPI) *Webhooks {
	return &Webhooks{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Webhooks.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Webhooks.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Webhooks.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Webhooks.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Webhooks.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Webhooks.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Webhooks.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Webhooks.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Webhooks.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Get: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksGet()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Webhooks.Get", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Get(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Webhooks.Get", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Webhooks.Get", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWebhooksDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Webhooks.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Webhooks.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Webhooks.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Webhooks) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/webhooks/", h.List)
		r.Post("/webhooks/", h.Create)
		r.Post("/webhooks/{webhookID}", h.Update)
		r.Get("/webhooks/{webhookID}", h.Get)
		r.Delete("/webhooks/{webhookID}", h.Delete)
	})
}

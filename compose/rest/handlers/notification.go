package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `notification.go`, `notification.util.go` or `notification_test.go` to
	implement your API calls, helper functions and tests. The file `notification.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

// Internal API interface
type NotificationAPI interface {
	EmailSend(context.Context, *request.NotificationEmailSend) (interface{}, error)
}

// HTTP API interface
type Notification struct {
	EmailSend func(http.ResponseWriter, *http.Request)
}

func NewNotification(h NotificationAPI) *Notification {
	return &Notification{
		EmailSend: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNotificationEmailSend()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Notification.EmailSend", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.EmailSend(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Notification.EmailSend", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Notification.EmailSend", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Notification) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/notification/email", h.EmailSend)
	})
}

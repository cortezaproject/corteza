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

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

type (
	// Internal API interface
	NotificationAPI interface {
		EmailSend(context.Context, *request.NotificationEmailSend) (interface{}, error)
	}

	// HTTP API interface
	Notification struct {
		EmailSend func(http.ResponseWriter, *http.Request)
	}
)

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

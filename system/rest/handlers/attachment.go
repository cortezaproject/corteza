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
	AttachmentAPI interface {
		Read(context.Context, *request.AttachmentRead) (interface{}, error)
		Delete(context.Context, *request.AttachmentDelete) (interface{}, error)
		Original(context.Context, *request.AttachmentOriginal) (interface{}, error)
		Preview(context.Context, *request.AttachmentPreview) (interface{}, error)
	}

	// HTTP API interface
	Attachment struct {
		Read     func(http.ResponseWriter, *http.Request)
		Delete   func(http.ResponseWriter, *http.Request)
		Original func(http.ResponseWriter, *http.Request)
		Preview  func(http.ResponseWriter, *http.Request)
	}
)

func NewAttachment(h AttachmentAPI) *Attachment {
	return &Attachment{
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Attachment.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Attachment.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Attachment.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Attachment.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Attachment.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Attachment.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Original: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentOriginal()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Attachment.Original", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Original(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Attachment.Original", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Attachment.Original", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Preview: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentPreview()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Attachment.Preview", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Preview(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Attachment.Preview", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Attachment.Preview", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Attachment) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/attachment/{kind}/{attachmentID}", h.Read)
		r.Delete("/attachment/{kind}/{attachmentID}", h.Delete)
		r.Get("/attachment/{kind}/{attachmentID}/original/{name}", h.Original)
		r.Get("/attachment/{kind}/{attachmentID}/preview.{ext}", h.Preview)
	})
}

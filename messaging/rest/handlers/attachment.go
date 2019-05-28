package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `attachment.go`, `attachment.util.go` or `attachment_test.go` to
	implement your API calls, helper functions and tests. The file `attachment.go`
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
type AttachmentAPI interface {
	Original(context.Context, *request.AttachmentOriginal) (interface{}, error)
	Preview(context.Context, *request.AttachmentPreview) (interface{}, error)
}

// HTTP API interface
type Attachment struct {
	Original func(http.ResponseWriter, *http.Request)
	Preview  func(http.ResponseWriter, *http.Request)
}

func NewAttachment(h AttachmentAPI) *Attachment {
	return &Attachment{
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
		r.Get("/attachment/{attachmentID}/original/{name}", h.Original)
		r.Get("/attachment/{attachmentID}/preview.{ext}", h.Preview)
	})
}

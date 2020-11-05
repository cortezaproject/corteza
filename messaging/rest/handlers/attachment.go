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
	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	AttachmentAPI interface {
		Original(context.Context, *request.AttachmentOriginal) (interface{}, error)
		Preview(context.Context, *request.AttachmentPreview) (interface{}, error)
	}

	// HTTP API interface
	Attachment struct {
		Original func(http.ResponseWriter, *http.Request)
		Preview  func(http.ResponseWriter, *http.Request)
	}
)

func NewAttachment(h AttachmentAPI) *Attachment {
	return &Attachment{
		Original: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentOriginal()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Original(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Preview: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentPreview()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Preview(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
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

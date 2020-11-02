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
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/go-chi/chi"
	"net/http"
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
				api.Send(w, r, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentDelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
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
		r.Get("/attachment/{kind}/{attachmentID}", h.Read)
		r.Delete("/attachment/{kind}/{attachmentID}", h.Delete)
		r.Get("/attachment/{kind}/{attachmentID}/original/{name}", h.Original)
		r.Get("/attachment/{kind}/{attachmentID}/preview.{ext}", h.Preview)
	})
}

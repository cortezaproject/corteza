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

	"github.com/crusttech/crust/crm/rest/request"
)

// Internal API interface
type AttachmentAPI interface {
	List(context.Context, *request.AttachmentList) (interface{}, error)
	Details(context.Context, *request.AttachmentDetails) (interface{}, error)
	Original(context.Context, *request.AttachmentOriginal) (interface{}, error)
	Preview(context.Context, *request.AttachmentPreview) (interface{}, error)
}

// HTTP API interface
type Attachment struct {
	List     func(http.ResponseWriter, *http.Request)
	Details  func(http.ResponseWriter, *http.Request)
	Original func(http.ResponseWriter, *http.Request)
	Preview  func(http.ResponseWriter, *http.Request)
}

func NewAttachment(ah AttachmentAPI) *Attachment {
	return &Attachment{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.List(r.Context(), params)
			})
		},
		Details: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentDetails()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Details(r.Context(), params)
			})
		},
		Original: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentOriginal()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Original(r.Context(), params)
			})
		},
		Preview: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAttachmentPreview()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Preview(r.Context(), params)
			})
		},
	}
}

func (ah *Attachment) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/attachment/{kind}/", ah.List)
		r.Get("/attachment/{kind}/{attachmentID}", ah.Details)
		r.Get("/attachment/{kind}/{attachmentID}/original/{name}", ah.Original)
		r.Get("/attachment/{kind}/{attachmentID}/preview.{ext}", ah.Preview)
	})
}

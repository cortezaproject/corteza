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
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	DataPrivacyAPI interface {
		ConnectionList(context.Context, *request.DataPrivacyConnectionList) (interface{}, error)
		RequestList(context.Context, *request.DataPrivacyRequestList) (interface{}, error)
		RequestCreate(context.Context, *request.DataPrivacyRequestCreate) (interface{}, error)
		RequestRead(context.Context, *request.DataPrivacyRequestRead) (interface{}, error)
		RequestUpdateStatus(context.Context, *request.DataPrivacyRequestUpdateStatus) (interface{}, error)
		RequestCommentList(context.Context, *request.DataPrivacyRequestCommentList) (interface{}, error)
		RequestCommentCreate(context.Context, *request.DataPrivacyRequestCommentCreate) (interface{}, error)
	}

	// HTTP API interface
	DataPrivacy struct {
		ConnectionList       func(http.ResponseWriter, *http.Request)
		RequestList          func(http.ResponseWriter, *http.Request)
		RequestCreate        func(http.ResponseWriter, *http.Request)
		RequestRead          func(http.ResponseWriter, *http.Request)
		RequestUpdateStatus  func(http.ResponseWriter, *http.Request)
		RequestCommentList   func(http.ResponseWriter, *http.Request)
		RequestCommentCreate func(http.ResponseWriter, *http.Request)
	}
)

func NewDataPrivacy(h DataPrivacyAPI) *DataPrivacy {
	return &DataPrivacy{
		ConnectionList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyConnectionList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ConnectionList(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RequestList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RequestList(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RequestCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RequestCreate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RequestRead: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestRead()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RequestRead(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RequestUpdateStatus: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestUpdateStatus()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RequestUpdateStatus(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RequestCommentList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestCommentList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RequestCommentList(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RequestCommentCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestCommentCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RequestCommentCreate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h DataPrivacy) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/data-privacy/connection/", h.ConnectionList)
		r.Get("/data-privacy/requests/", h.RequestList)
		r.Post("/data-privacy/requests/", h.RequestCreate)
		r.Get("/data-privacy/requests/{requestID}", h.RequestRead)
		r.Patch("/data-privacy/requests/{requestID}/status/{status}", h.RequestUpdateStatus)
		r.Get("/data-privacy/requests/{requestID}/comments/", h.RequestCommentList)
		r.Post("/data-privacy/requests/{requestID}/comments/", h.RequestCommentCreate)
	})
}

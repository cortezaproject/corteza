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
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	DalConnectionAPI interface {
		List(context.Context, *request.DalConnectionList) (interface{}, error)
		Create(context.Context, *request.DalConnectionCreate) (interface{}, error)
		Update(context.Context, *request.DalConnectionUpdate) (interface{}, error)
		ReadPrimary(context.Context, *request.DalConnectionReadPrimary) (interface{}, error)
		Read(context.Context, *request.DalConnectionRead) (interface{}, error)
		Delete(context.Context, *request.DalConnectionDelete) (interface{}, error)
		Undelete(context.Context, *request.DalConnectionUndelete) (interface{}, error)
	}

	// HTTP API interface
	DalConnection struct {
		List        func(http.ResponseWriter, *http.Request)
		Create      func(http.ResponseWriter, *http.Request)
		Update      func(http.ResponseWriter, *http.Request)
		ReadPrimary func(http.ResponseWriter, *http.Request)
		Read        func(http.ResponseWriter, *http.Request)
		Delete      func(http.ResponseWriter, *http.Request)
		Undelete    func(http.ResponseWriter, *http.Request)
	}
)

func NewDalConnection(h DalConnectionAPI) *DalConnection {
	return &DalConnection{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDalConnectionList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDalConnectionCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDalConnectionUpdate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ReadPrimary: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDalConnectionReadPrimary()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReadPrimary(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDalConnectionRead()
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
			params := request.NewDalConnectionDelete()
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
		Undelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDalConnectionUndelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Undelete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h DalConnection) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/dal/connections/", h.List)
		r.Post("/dal/connections/", h.Create)
		r.Put("/dal/connections/{connectionID}", h.Update)
		r.Get("/dal/connections/primary", h.ReadPrimary)
		r.Get("/dal/connections/{connectionID}", h.Read)
		r.Delete("/dal/connections/{connectionID}", h.Delete)
		r.Post("/dal/connections/{connectionID}/undelete", h.Undelete)
	})
}

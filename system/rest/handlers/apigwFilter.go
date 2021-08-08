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
	ApigwFilterAPI interface {
		List(context.Context, *request.ApigwFilterList) (interface{}, error)
		Create(context.Context, *request.ApigwFilterCreate) (interface{}, error)
		Update(context.Context, *request.ApigwFilterUpdate) (interface{}, error)
		Read(context.Context, *request.ApigwFilterRead) (interface{}, error)
		Delete(context.Context, *request.ApigwFilterDelete) (interface{}, error)
		Undelete(context.Context, *request.ApigwFilterUndelete) (interface{}, error)
		DefFilter(context.Context, *request.ApigwFilterDefFilter) (interface{}, error)
		DefProxyAuth(context.Context, *request.ApigwFilterDefProxyAuth) (interface{}, error)
	}

	// HTTP API interface
	ApigwFilter struct {
		List         func(http.ResponseWriter, *http.Request)
		Create       func(http.ResponseWriter, *http.Request)
		Update       func(http.ResponseWriter, *http.Request)
		Read         func(http.ResponseWriter, *http.Request)
		Delete       func(http.ResponseWriter, *http.Request)
		Undelete     func(http.ResponseWriter, *http.Request)
		DefFilter    func(http.ResponseWriter, *http.Request)
		DefProxyAuth func(http.ResponseWriter, *http.Request)
	}
)

func NewApigwFilter(h ApigwFilterAPI) *ApigwFilter {
	return &ApigwFilter{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApigwFilterList()
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
			params := request.NewApigwFilterCreate()
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
			params := request.NewApigwFilterUpdate()
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
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApigwFilterRead()
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
			params := request.NewApigwFilterDelete()
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
			params := request.NewApigwFilterUndelete()
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
		DefFilter: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApigwFilterDefFilter()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.DefFilter(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		DefProxyAuth: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApigwFilterDefProxyAuth()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.DefProxyAuth(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h ApigwFilter) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/apigw/filter/", h.List)
		r.Put("/apigw/filter", h.Create)
		r.Post("/apigw/filter/{filterID}", h.Update)
		r.Get("/apigw/filter/{filterID}", h.Read)
		r.Delete("/apigw/filter/{filterID}", h.Delete)
		r.Post("/apigw/filter/{filterID}/undelete", h.Undelete)
		r.Get("/apigw/filter/def", h.DefFilter)
		r.Get("/apigw/filter/proxy_auth/def", h.DefProxyAuth)
	})
}

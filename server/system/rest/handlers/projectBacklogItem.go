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
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	ProjectBacklogItemAPI interface {
		List(context.Context, *request.ProjectBacklogItemList) (interface{}, error)
		Create(context.Context, *request.ProjectBacklogItemCreate) (interface{}, error)
		Read(context.Context, *request.ProjectBacklogItemRead) (interface{}, error)
		Update(context.Context, *request.ProjectBacklogItemUpdate) (interface{}, error)
		Delete(context.Context, *request.ProjectBacklogItemDelete) (interface{}, error)
	}

	// HTTP API interface
	ProjectBacklogItem struct {
		List   func(http.ResponseWriter, *http.Request)
		Create func(http.ResponseWriter, *http.Request)
		Read   func(http.ResponseWriter, *http.Request)
		Update func(http.ResponseWriter, *http.Request)
		Delete func(http.ResponseWriter, *http.Request)
	}
)

func NewProjectBacklogItem(h ProjectBacklogItemAPI) *ProjectBacklogItem {
	return &ProjectBacklogItem{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectBacklogItemList()
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
			params := request.NewProjectBacklogItemCreate()
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
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectBacklogItemRead()
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
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectBacklogItemUpdate()
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
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectBacklogItemDelete()
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
	}
}

func (h ProjectBacklogItem) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/project-backlog-items/", h.List)
		r.Post("/project-backlog-items", h.Create)
		r.Get("/project-backlog-items/{backlogItemID}", h.Read)
		r.Put("/project-backlog-items/{backlogItemID}", h.Update)
		r.Delete("/project-backlog-items/{backlogItemID}", h.Delete)
	})
}

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
	ProjectGroupAPI interface {
		List(context.Context, *request.ProjectGroupList) (interface{}, error)
		Create(context.Context, *request.ProjectGroupCreate) (interface{}, error)
		Read(context.Context, *request.ProjectGroupRead) (interface{}, error)
		Update(context.Context, *request.ProjectGroupUpdate) (interface{}, error)
		Delete(context.Context, *request.ProjectGroupDelete) (interface{}, error)
		EntryAdd(context.Context, *request.ProjectGroupEntryAdd) (interface{}, error)
		EntryRemove(context.Context, *request.ProjectGroupEntryRemove) (interface{}, error)
	}

	// HTTP API interface
	ProjectGroup struct {
		List        func(http.ResponseWriter, *http.Request)
		Create      func(http.ResponseWriter, *http.Request)
		Read        func(http.ResponseWriter, *http.Request)
		Update      func(http.ResponseWriter, *http.Request)
		Delete      func(http.ResponseWriter, *http.Request)
		EntryAdd    func(http.ResponseWriter, *http.Request)
		EntryRemove func(http.ResponseWriter, *http.Request)
	}
)

func NewProjectGroup(h ProjectGroupAPI) *ProjectGroup {
	return &ProjectGroup{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectGroupList()
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
			params := request.NewProjectGroupCreate()
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
			params := request.NewProjectGroupRead()
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
			params := request.NewProjectGroupUpdate()
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
			params := request.NewProjectGroupDelete()
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
		EntryAdd: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectGroupEntryAdd()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.EntryAdd(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		EntryRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectGroupEntryRemove()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.EntryRemove(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h ProjectGroup) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/projects/{projectID}/groups/", h.List)
		r.Post("/projects/{projectID}/groups/", h.Create)
		r.Get("/projects/{projectID}/groups/{projectGroupID}", h.Read)
		r.Put("/projects/{projectID}/groups/{projectGroupID}", h.Update)
		r.Delete("/projects/{projectID}/groups/{projectGroupID}", h.Delete)
		r.Post("/projects/{projectID}/groups/{projectGroupID}/entries/{resourceRef}", h.EntryAdd)
		r.Delete("/projects/{projectID}/groups/{projectGroupID}/entries/{resourceRef}", h.EntryRemove)
	})
}

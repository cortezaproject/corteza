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
	ProjectAiSystemAPI interface {
		List(context.Context, *request.ProjectAiSystemList) (interface{}, error)
		Create(context.Context, *request.ProjectAiSystemCreate) (interface{}, error)
		Read(context.Context, *request.ProjectAiSystemRead) (interface{}, error)
		Update(context.Context, *request.ProjectAiSystemUpdate) (interface{}, error)
		Delete(context.Context, *request.ProjectAiSystemDelete) (interface{}, error)
		EntryAdd(context.Context, *request.ProjectAiSystemEntryAdd) (interface{}, error)
		EntryRemove(context.Context, *request.ProjectAiSystemEntryRemove) (interface{}, error)
	}

	// HTTP API interface
	ProjectAiSystem struct {
		List        func(http.ResponseWriter, *http.Request)
		Create      func(http.ResponseWriter, *http.Request)
		Read        func(http.ResponseWriter, *http.Request)
		Update      func(http.ResponseWriter, *http.Request)
		Delete      func(http.ResponseWriter, *http.Request)
		EntryAdd    func(http.ResponseWriter, *http.Request)
		EntryRemove func(http.ResponseWriter, *http.Request)
	}
)

func NewProjectAiSystem(h ProjectAiSystemAPI) *ProjectAiSystem {
	return &ProjectAiSystem{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectAiSystemList()
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
			params := request.NewProjectAiSystemCreate()
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
			params := request.NewProjectAiSystemRead()
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
			params := request.NewProjectAiSystemUpdate()
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
			params := request.NewProjectAiSystemDelete()
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
			params := request.NewProjectAiSystemEntryAdd()
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
			params := request.NewProjectAiSystemEntryRemove()
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

func (h ProjectAiSystem) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/projects/{projectID}/ai-systems/", h.List)
		r.Post("/projects/{projectID}/ai-systems/", h.Create)
		r.Get("/projects/{projectID}/ai-systems/{projectAiSystemID}", h.Read)
		r.Put("/projects/{projectID}/ai-systems/{projectAiSystemID}", h.Update)
		r.Delete("/projects/{projectID}/ai-systems/{projectAiSystemID}", h.Delete)
		r.Post("/projects/{projectID}/ai-systems/{projectAiSystemID}/entries", h.EntryAdd)
		r.Delete("/projects/{projectID}/ai-systems/{projectAiSystemID}/entries", h.EntryRemove)
	})
}

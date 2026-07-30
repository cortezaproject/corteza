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
	ProjectFriaScenarioAPI interface {
		List(context.Context, *request.ProjectFriaScenarioList) (interface{}, error)
		Create(context.Context, *request.ProjectFriaScenarioCreate) (interface{}, error)
		Read(context.Context, *request.ProjectFriaScenarioRead) (interface{}, error)
		Update(context.Context, *request.ProjectFriaScenarioUpdate) (interface{}, error)
		Delete(context.Context, *request.ProjectFriaScenarioDelete) (interface{}, error)
	}

	// HTTP API interface
	ProjectFriaScenario struct {
		List   func(http.ResponseWriter, *http.Request)
		Create func(http.ResponseWriter, *http.Request)
		Read   func(http.ResponseWriter, *http.Request)
		Update func(http.ResponseWriter, *http.Request)
		Delete func(http.ResponseWriter, *http.Request)
	}
)

func NewProjectFriaScenario(h ProjectFriaScenarioAPI) *ProjectFriaScenario {
	return &ProjectFriaScenario{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectFriaScenarioList()
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
			params := request.NewProjectFriaScenarioCreate()
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
			params := request.NewProjectFriaScenarioRead()
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
			params := request.NewProjectFriaScenarioUpdate()
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
			params := request.NewProjectFriaScenarioDelete()
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

func (h ProjectFriaScenario) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/projects/{projectID}/fria-scenarios/", h.List)
		r.Post("/projects/{projectID}/fria-scenarios/", h.Create)
		r.Get("/projects/{projectID}/fria-scenarios/{projectFriaScenarioID}", h.Read)
		r.Put("/projects/{projectID}/fria-scenarios/{projectFriaScenarioID}", h.Update)
		r.Delete("/projects/{projectID}/fria-scenarios/{projectFriaScenarioID}", h.Delete)
	})
}

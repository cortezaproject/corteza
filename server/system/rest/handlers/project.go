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
	ProjectAPI interface {
		List(context.Context, *request.ProjectList) (interface{}, error)
		Create(context.Context, *request.ProjectCreate) (interface{}, error)
		Read(context.Context, *request.ProjectRead) (interface{}, error)
		Update(context.Context, *request.ProjectUpdate) (interface{}, error)
		Delete(context.Context, *request.ProjectDelete) (interface{}, error)
		Undelete(context.Context, *request.ProjectUndelete) (interface{}, error)
		ListMembers(context.Context, *request.ProjectListMembers) (interface{}, error)
		AddMember(context.Context, *request.ProjectAddMember) (interface{}, error)
		UpdateMember(context.Context, *request.ProjectUpdateMember) (interface{}, error)
		RemoveMember(context.Context, *request.ProjectRemoveMember) (interface{}, error)
	}

	// HTTP API interface
	Project struct {
		List         func(http.ResponseWriter, *http.Request)
		Create       func(http.ResponseWriter, *http.Request)
		Read         func(http.ResponseWriter, *http.Request)
		Update       func(http.ResponseWriter, *http.Request)
		Delete       func(http.ResponseWriter, *http.Request)
		Undelete     func(http.ResponseWriter, *http.Request)
		ListMembers  func(http.ResponseWriter, *http.Request)
		AddMember    func(http.ResponseWriter, *http.Request)
		UpdateMember func(http.ResponseWriter, *http.Request)
		RemoveMember func(http.ResponseWriter, *http.Request)
	}
)

func NewProject(h ProjectAPI) *Project {
	return &Project{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectList()
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
			params := request.NewProjectCreate()
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
			params := request.NewProjectRead()
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
			params := request.NewProjectUpdate()
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
			params := request.NewProjectDelete()
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
			params := request.NewProjectUndelete()
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
		ListMembers: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectListMembers()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListMembers(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		AddMember: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectAddMember()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.AddMember(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		UpdateMember: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectUpdateMember()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.UpdateMember(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RemoveMember: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectRemoveMember()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RemoveMember(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Project) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/projects/", h.List)
		r.Post("/projects", h.Create)
		r.Get("/projects/{projectID}", h.Read)
		r.Put("/projects/{projectID}", h.Update)
		r.Delete("/projects/{projectID}", h.Delete)
		r.Post("/projects/{projectID}/undelete", h.Undelete)
		r.Get("/projects/{projectID}/members", h.ListMembers)
		r.Post("/projects/{projectID}/members", h.AddMember)
		r.Put("/projects/{projectID}/members/{userID}", h.UpdateMember)
		r.Delete("/projects/{projectID}/members/{userID}", h.RemoveMember)
	})
}

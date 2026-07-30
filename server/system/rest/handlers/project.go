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
		Archive(context.Context, *request.ProjectArchive) (interface{}, error)
		Unarchive(context.Context, *request.ProjectUnarchive) (interface{}, error)
		ListMembers(context.Context, *request.ProjectListMembers) (interface{}, error)
		AddMember(context.Context, *request.ProjectAddMember) (interface{}, error)
		UpdateMember(context.Context, *request.ProjectUpdateMember) (interface{}, error)
		RemoveMember(context.Context, *request.ProjectRemoveMember) (interface{}, error)
		Graph(context.Context, *request.ProjectGraph) (interface{}, error)
		CreateRevision(context.Context, *request.ProjectCreateRevision) (interface{}, error)
		ListRevisions(context.Context, *request.ProjectListRevisions) (interface{}, error)
		RequestApproval(context.Context, *request.ProjectRequestApproval) (interface{}, error)
		GrantApproval(context.Context, *request.ProjectGrantApproval) (interface{}, error)
		RejectApproval(context.Context, *request.ProjectRejectApproval) (interface{}, error)
		GetDeploymentPlan(context.Context, *request.ProjectGetDeploymentPlan) (interface{}, error)
		Publish(context.Context, *request.ProjectPublish) (interface{}, error)
	}

	// HTTP API interface
	Project struct {
		List              func(http.ResponseWriter, *http.Request)
		Create            func(http.ResponseWriter, *http.Request)
		Read              func(http.ResponseWriter, *http.Request)
		Update            func(http.ResponseWriter, *http.Request)
		Delete            func(http.ResponseWriter, *http.Request)
		Undelete          func(http.ResponseWriter, *http.Request)
		Archive           func(http.ResponseWriter, *http.Request)
		Unarchive         func(http.ResponseWriter, *http.Request)
		ListMembers       func(http.ResponseWriter, *http.Request)
		AddMember         func(http.ResponseWriter, *http.Request)
		UpdateMember      func(http.ResponseWriter, *http.Request)
		RemoveMember      func(http.ResponseWriter, *http.Request)
		Graph             func(http.ResponseWriter, *http.Request)
		CreateRevision    func(http.ResponseWriter, *http.Request)
		ListRevisions     func(http.ResponseWriter, *http.Request)
		RequestApproval   func(http.ResponseWriter, *http.Request)
		GrantApproval     func(http.ResponseWriter, *http.Request)
		RejectApproval    func(http.ResponseWriter, *http.Request)
		GetDeploymentPlan func(http.ResponseWriter, *http.Request)
		Publish           func(http.ResponseWriter, *http.Request)
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
		Archive: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectArchive()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Archive(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Unarchive: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectUnarchive()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Unarchive(r.Context(), params)
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
		Graph: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectGraph()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Graph(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		CreateRevision: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectCreateRevision()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.CreateRevision(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ListRevisions: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectListRevisions()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListRevisions(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RequestApproval: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectRequestApproval()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RequestApproval(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		GrantApproval: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectGrantApproval()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.GrantApproval(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RejectApproval: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectRejectApproval()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RejectApproval(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		GetDeploymentPlan: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectGetDeploymentPlan()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.GetDeploymentPlan(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Publish: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectPublish()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Publish(r.Context(), params)
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
		r.Post("/projects/{projectID}/archive", h.Archive)
		r.Post("/projects/{projectID}/unarchive", h.Unarchive)
		r.Get("/projects/{projectID}/members", h.ListMembers)
		r.Post("/projects/{projectID}/members", h.AddMember)
		r.Put("/projects/{projectID}/members/{userID}", h.UpdateMember)
		r.Delete("/projects/{projectID}/members/{userID}", h.RemoveMember)
		r.Get("/projects/{projectID}/graph", h.Graph)
		r.Post("/projects/{projectID}/revision", h.CreateRevision)
		r.Get("/projects/{projectID}/revisions", h.ListRevisions)
		r.Post("/projects/{projectID}/approval/request", h.RequestApproval)
		r.Post("/projects/{projectID}/approval/grant", h.GrantApproval)
		r.Post("/projects/{projectID}/approval/reject", h.RejectApproval)
		r.Get("/projects/{projectID}/publish", h.GetDeploymentPlan)
		r.Post("/projects/{projectID}/publish", h.Publish)
	})
}

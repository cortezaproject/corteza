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
	TenantAPI interface {
		List(context.Context, *request.TenantList) (interface{}, error)
		Create(context.Context, *request.TenantCreate) (interface{}, error)
		Read(context.Context, *request.TenantRead) (interface{}, error)
		Update(context.Context, *request.TenantUpdate) (interface{}, error)
		Delete(context.Context, *request.TenantDelete) (interface{}, error)
		Undelete(context.Context, *request.TenantUndelete) (interface{}, error)
		Suspend(context.Context, *request.TenantSuspend) (interface{}, error)
		Activate(context.Context, *request.TenantActivate) (interface{}, error)
		Archive(context.Context, *request.TenantArchive) (interface{}, error)
		ListMembers(context.Context, *request.TenantListMembers) (interface{}, error)
		AddMember(context.Context, *request.TenantAddMember) (interface{}, error)
		UpdateMember(context.Context, *request.TenantUpdateMember) (interface{}, error)
		RemoveMember(context.Context, *request.TenantRemoveMember) (interface{}, error)
		SuspendMember(context.Context, *request.TenantSuspendMember) (interface{}, error)
		ActivateMember(context.Context, *request.TenantActivateMember) (interface{}, error)
	}

	// HTTP API interface
	Tenant struct {
		List           func(http.ResponseWriter, *http.Request)
		Create         func(http.ResponseWriter, *http.Request)
		Read           func(http.ResponseWriter, *http.Request)
		Update         func(http.ResponseWriter, *http.Request)
		Delete         func(http.ResponseWriter, *http.Request)
		Undelete       func(http.ResponseWriter, *http.Request)
		Suspend        func(http.ResponseWriter, *http.Request)
		Activate       func(http.ResponseWriter, *http.Request)
		Archive        func(http.ResponseWriter, *http.Request)
		ListMembers    func(http.ResponseWriter, *http.Request)
		AddMember      func(http.ResponseWriter, *http.Request)
		UpdateMember   func(http.ResponseWriter, *http.Request)
		RemoveMember   func(http.ResponseWriter, *http.Request)
		SuspendMember  func(http.ResponseWriter, *http.Request)
		ActivateMember func(http.ResponseWriter, *http.Request)
	}
)

func NewTenant(h TenantAPI) *Tenant {
	return &Tenant{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTenantList()
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
			params := request.NewTenantCreate()
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
			params := request.NewTenantRead()
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
			params := request.NewTenantUpdate()
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
			params := request.NewTenantDelete()
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
			params := request.NewTenantUndelete()
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
		Suspend: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTenantSuspend()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Suspend(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Activate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTenantActivate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Activate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Archive: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTenantArchive()
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
		ListMembers: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTenantListMembers()
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
			params := request.NewTenantAddMember()
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
			params := request.NewTenantUpdateMember()
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
			params := request.NewTenantRemoveMember()
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
		SuspendMember: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTenantSuspendMember()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.SuspendMember(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ActivateMember: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTenantActivateMember()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ActivateMember(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Tenant) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/tenants/", h.List)
		r.Post("/tenants", h.Create)
		r.Get("/tenants/{tenantID}", h.Read)
		r.Put("/tenants/{tenantID}", h.Update)
		r.Delete("/tenants/{tenantID}", h.Delete)
		r.Post("/tenants/{tenantID}/undelete", h.Undelete)
		r.Post("/tenants/{tenantID}/suspend", h.Suspend)
		r.Post("/tenants/{tenantID}/activate", h.Activate)
		r.Post("/tenants/{tenantID}/archive", h.Archive)
		r.Get("/tenants/{tenantID}/members", h.ListMembers)
		r.Post("/tenants/{tenantID}/members", h.AddMember)
		r.Put("/tenants/{tenantID}/members/{userID}", h.UpdateMember)
		r.Delete("/tenants/{tenantID}/members/{userID}", h.RemoveMember)
		r.Post("/tenants/{tenantID}/members/{userID}/suspend", h.SuspendMember)
		r.Post("/tenants/{tenantID}/members/{userID}/activate", h.ActivateMember)
	})
}

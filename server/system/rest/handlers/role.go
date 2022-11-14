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
	RoleAPI interface {
		List(context.Context, *request.RoleList) (interface{}, error)
		Create(context.Context, *request.RoleCreate) (interface{}, error)
		Update(context.Context, *request.RoleUpdate) (interface{}, error)
		Read(context.Context, *request.RoleRead) (interface{}, error)
		Delete(context.Context, *request.RoleDelete) (interface{}, error)
		Archive(context.Context, *request.RoleArchive) (interface{}, error)
		Unarchive(context.Context, *request.RoleUnarchive) (interface{}, error)
		Undelete(context.Context, *request.RoleUndelete) (interface{}, error)
		Move(context.Context, *request.RoleMove) (interface{}, error)
		Merge(context.Context, *request.RoleMerge) (interface{}, error)
		MemberList(context.Context, *request.RoleMemberList) (interface{}, error)
		MemberAdd(context.Context, *request.RoleMemberAdd) (interface{}, error)
		MemberRemove(context.Context, *request.RoleMemberRemove) (interface{}, error)
		TriggerScript(context.Context, *request.RoleTriggerScript) (interface{}, error)
	}

	// HTTP API interface
	Role struct {
		List          func(http.ResponseWriter, *http.Request)
		Create        func(http.ResponseWriter, *http.Request)
		Update        func(http.ResponseWriter, *http.Request)
		Read          func(http.ResponseWriter, *http.Request)
		Delete        func(http.ResponseWriter, *http.Request)
		Archive       func(http.ResponseWriter, *http.Request)
		Unarchive     func(http.ResponseWriter, *http.Request)
		Undelete      func(http.ResponseWriter, *http.Request)
		Move          func(http.ResponseWriter, *http.Request)
		Merge         func(http.ResponseWriter, *http.Request)
		MemberList    func(http.ResponseWriter, *http.Request)
		MemberAdd     func(http.ResponseWriter, *http.Request)
		MemberRemove  func(http.ResponseWriter, *http.Request)
		TriggerScript func(http.ResponseWriter, *http.Request)
	}
)

func NewRole(h RoleAPI) *Role {
	return &Role{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleList()
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
			params := request.NewRoleCreate()
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
			params := request.NewRoleUpdate()
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
			params := request.NewRoleRead()
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
			params := request.NewRoleDelete()
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
		Archive: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleArchive()
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
			params := request.NewRoleUnarchive()
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
		Undelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleUndelete()
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
		Move: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMove()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Move(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Merge: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMerge()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Merge(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MemberList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MemberList(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MemberAdd: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberAdd()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MemberAdd(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MemberRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberRemove()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MemberRemove(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		TriggerScript: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleTriggerScript()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.TriggerScript(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Role) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/roles/", h.List)
		r.Post("/roles/", h.Create)
		r.Put("/roles/{roleID}", h.Update)
		r.Get("/roles/{roleID}", h.Read)
		r.Delete("/roles/{roleID}", h.Delete)
		r.Post("/roles/{roleID}/archive", h.Archive)
		r.Post("/roles/{roleID}/unarchive", h.Unarchive)
		r.Post("/roles/{roleID}/undelete", h.Undelete)
		r.Post("/roles/{roleID}/move", h.Move)
		r.Post("/roles/{roleID}/merge", h.Merge)
		r.Get("/roles/{roleID}/members", h.MemberList)
		r.Post("/roles/{roleID}/member/{userID}", h.MemberAdd)
		r.Delete("/roles/{roleID}/member/{userID}", h.MemberRemove)
		r.Post("/roles/{roleID}/trigger", h.TriggerScript)
	})
}

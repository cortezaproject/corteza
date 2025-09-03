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
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	UserGroupAPI interface {
		List(context.Context, *request.UserGroupList) (interface{}, error)
		Create(context.Context, *request.UserGroupCreate) (interface{}, error)
		Update(context.Context, *request.UserGroupUpdate) (interface{}, error)
		Read(context.Context, *request.UserGroupRead) (interface{}, error)
		Delete(context.Context, *request.UserGroupDelete) (interface{}, error)
		Undelete(context.Context, *request.UserGroupUndelete) (interface{}, error)
		MemberList(context.Context, *request.UserGroupMemberList) (interface{}, error)
		MemberAdd(context.Context, *request.UserGroupMemberAdd) (interface{}, error)
	}

	// HTTP API interface
	UserGroup struct {
		List       func(http.ResponseWriter, *http.Request)
		Create     func(http.ResponseWriter, *http.Request)
		Update     func(http.ResponseWriter, *http.Request)
		Read       func(http.ResponseWriter, *http.Request)
		Delete     func(http.ResponseWriter, *http.Request)
		Undelete   func(http.ResponseWriter, *http.Request)
		MemberList func(http.ResponseWriter, *http.Request)
		MemberAdd  func(http.ResponseWriter, *http.Request)
	}
)

func NewUserGroup(h UserGroupAPI) *UserGroup {
	return &UserGroup{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserGroupList()
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
			params := request.NewUserGroupCreate()
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
			params := request.NewUserGroupUpdate()
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
			params := request.NewUserGroupRead()
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
			params := request.NewUserGroupDelete()
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
			params := request.NewUserGroupUndelete()
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
		MemberList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserGroupMemberList()
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
			params := request.NewUserGroupMemberAdd()
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
	}
}

func (h UserGroup) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/user-groups/", h.List)
		r.Post("/user-groups/", h.Create)
		r.Put("/user-groups/{userGroupID}", h.Update)
		r.Get("/user-groups/{userGroupID}", h.Read)
		r.Delete("/user-groups/{userGroupID}", h.Delete)
		r.Post("/user-groups/{userGroupID}/undelete", h.Undelete)
		r.Get("/user-groups/{userGroupID}/members", h.MemberList)
		r.Post("/user-groups/{userGroupID}/member/{userID}", h.MemberAdd)
	})
}

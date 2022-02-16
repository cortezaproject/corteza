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
	UserAPI interface {
		List(context.Context, *request.UserList) (interface{}, error)
		Create(context.Context, *request.UserCreate) (interface{}, error)
		Update(context.Context, *request.UserUpdate) (interface{}, error)
		PartialUpdate(context.Context, *request.UserPartialUpdate) (interface{}, error)
		Read(context.Context, *request.UserRead) (interface{}, error)
		Delete(context.Context, *request.UserDelete) (interface{}, error)
		Suspend(context.Context, *request.UserSuspend) (interface{}, error)
		Unsuspend(context.Context, *request.UserUnsuspend) (interface{}, error)
		Undelete(context.Context, *request.UserUndelete) (interface{}, error)
		SetPassword(context.Context, *request.UserSetPassword) (interface{}, error)
		MembershipList(context.Context, *request.UserMembershipList) (interface{}, error)
		MembershipAdd(context.Context, *request.UserMembershipAdd) (interface{}, error)
		MembershipRemove(context.Context, *request.UserMembershipRemove) (interface{}, error)
		TriggerScript(context.Context, *request.UserTriggerScript) (interface{}, error)
		SessionsRemove(context.Context, *request.UserSessionsRemove) (interface{}, error)
		Export(context.Context, *request.UserExport) (interface{}, error)
		Import(context.Context, *request.UserImport) (interface{}, error)
	}

	// HTTP API interface
	User struct {
		List             func(http.ResponseWriter, *http.Request)
		Create           func(http.ResponseWriter, *http.Request)
		Update           func(http.ResponseWriter, *http.Request)
		PartialUpdate    func(http.ResponseWriter, *http.Request)
		Read             func(http.ResponseWriter, *http.Request)
		Delete           func(http.ResponseWriter, *http.Request)
		Suspend          func(http.ResponseWriter, *http.Request)
		Unsuspend        func(http.ResponseWriter, *http.Request)
		Undelete         func(http.ResponseWriter, *http.Request)
		SetPassword      func(http.ResponseWriter, *http.Request)
		MembershipList   func(http.ResponseWriter, *http.Request)
		MembershipAdd    func(http.ResponseWriter, *http.Request)
		MembershipRemove func(http.ResponseWriter, *http.Request)
		TriggerScript    func(http.ResponseWriter, *http.Request)
		SessionsRemove   func(http.ResponseWriter, *http.Request)
		Export           func(http.ResponseWriter, *http.Request)
		Import           func(http.ResponseWriter, *http.Request)
	}
)

func NewUser(h UserAPI) *User {
	return &User{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserList()
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
			params := request.NewUserCreate()
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
			params := request.NewUserUpdate()
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
		PartialUpdate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserPartialUpdate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.PartialUpdate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserRead()
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
			params := request.NewUserDelete()
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
		Suspend: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserSuspend()
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
		Unsuspend: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserUnsuspend()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Unsuspend(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Undelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserUndelete()
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
		SetPassword: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserSetPassword()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.SetPassword(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MembershipList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserMembershipList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MembershipList(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MembershipAdd: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserMembershipAdd()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MembershipAdd(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MembershipRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserMembershipRemove()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MembershipRemove(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		TriggerScript: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserTriggerScript()
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
		SessionsRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserSessionsRemove()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.SessionsRemove(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Export: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserExport()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Export(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Import: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserImport()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Import(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h User) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/users/", h.List)
		r.Post("/users/", h.Create)
		r.Put("/users/{userID}", h.Update)
		r.Patch("/users/{userID}", h.PartialUpdate)
		r.Get("/users/{userID}", h.Read)
		r.Delete("/users/{userID}", h.Delete)
		r.Post("/users/{userID}/suspend", h.Suspend)
		r.Post("/users/{userID}/unsuspend", h.Unsuspend)
		r.Post("/users/{userID}/undelete", h.Undelete)
		r.Post("/users/{userID}/password", h.SetPassword)
		r.Get("/users/{userID}/membership", h.MembershipList)
		r.Post("/users/{userID}/membership/{roleID}", h.MembershipAdd)
		r.Delete("/users/{userID}/membership/{roleID}", h.MembershipRemove)
		r.Post("/users/{userID}/trigger", h.TriggerScript)
		r.Delete("/users/{userID}/sessions", h.SessionsRemove)
		r.Get("/users/export/{filename}.zip", h.Export)
		r.Post("/users/import", h.Import)
	})
}

package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `user.go`, `user.util.go` or `user_test.go` to
	implement your API calls, helper functions and tests. The file `user.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/rest/request"
)

// Internal API interface
type UserAPI interface {
	List(context.Context, *request.UserList) (interface{}, error)
	Create(context.Context, *request.UserCreate) (interface{}, error)
	Update(context.Context, *request.UserUpdate) (interface{}, error)
	Read(context.Context, *request.UserRead) (interface{}, error)
	Delete(context.Context, *request.UserDelete) (interface{}, error)
	Suspend(context.Context, *request.UserSuspend) (interface{}, error)
	Unsuspend(context.Context, *request.UserUnsuspend) (interface{}, error)
	SetPassword(context.Context, *request.UserSetPassword) (interface{}, error)
	MembershipList(context.Context, *request.UserMembershipList) (interface{}, error)
	MembershipAdd(context.Context, *request.UserMembershipAdd) (interface{}, error)
	MembershipRemove(context.Context, *request.UserMembershipRemove) (interface{}, error)
}

// HTTP API interface
type User struct {
	List             func(http.ResponseWriter, *http.Request)
	Create           func(http.ResponseWriter, *http.Request)
	Update           func(http.ResponseWriter, *http.Request)
	Read             func(http.ResponseWriter, *http.Request)
	Delete           func(http.ResponseWriter, *http.Request)
	Suspend          func(http.ResponseWriter, *http.Request)
	Unsuspend        func(http.ResponseWriter, *http.Request)
	SetPassword      func(http.ResponseWriter, *http.Request)
	MembershipList   func(http.ResponseWriter, *http.Request)
	MembershipAdd    func(http.ResponseWriter, *http.Request)
	MembershipRemove func(http.ResponseWriter, *http.Request)
}

func NewUser(h UserAPI) *User {
	return &User{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Suspend: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserSuspend()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.Suspend", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Suspend(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.Suspend", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.Suspend", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Unsuspend: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserUnsuspend()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.Unsuspend", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Unsuspend(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.Unsuspend", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.Unsuspend", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		SetPassword: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserSetPassword()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.SetPassword", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.SetPassword(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.SetPassword", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.SetPassword", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		MembershipList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserMembershipList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.MembershipList", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.MembershipList(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.MembershipList", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.MembershipList", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		MembershipAdd: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserMembershipAdd()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.MembershipAdd", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.MembershipAdd(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.MembershipAdd", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.MembershipAdd", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		MembershipRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserMembershipRemove()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("User.MembershipRemove", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.MembershipRemove(r.Context(), params)
			if err != nil {
				logger.LogControllerError("User.MembershipRemove", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("User.MembershipRemove", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h User) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/users/", h.List)
		r.Post("/users/", h.Create)
		r.Put("/users/{userID}", h.Update)
		r.Get("/users/{userID}", h.Read)
		r.Delete("/users/{userID}", h.Delete)
		r.Post("/users/{userID}/suspend", h.Suspend)
		r.Post("/users/{userID}/unsuspend", h.Unsuspend)
		r.Post("/users/{userID}/password", h.SetPassword)
		r.Get("/users/{userID}/membership", h.MembershipList)
		r.Post("/users/{userID}/membership/{roleID}", h.MembershipAdd)
		r.Delete("/users/{userID}/membership/{roleID}", h.MembershipRemove)
	})
}

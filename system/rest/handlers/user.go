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

	"github.com/crusttech/crust/system/rest/request"
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
}

// HTTP API interface
type User struct {
	List      func(http.ResponseWriter, *http.Request)
	Create    func(http.ResponseWriter, *http.Request)
	Update    func(http.ResponseWriter, *http.Request)
	Read      func(http.ResponseWriter, *http.Request)
	Delete    func(http.ResponseWriter, *http.Request)
	Suspend   func(http.ResponseWriter, *http.Request)
	Unsuspend func(http.ResponseWriter, *http.Request)
}

func NewUser(uh UserAPI) *User {
	return &User{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Create(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Update(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Read(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Delete(r.Context(), params)
			})
		},
		Suspend: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserSuspend()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Suspend(r.Context(), params)
			})
		},
		Unsuspend: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserUnsuspend()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Unsuspend(r.Context(), params)
			})
		},
	}
}

func (uh *User) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/users/", uh.List)
		r.Post("/users/", uh.Create)
		r.Put("/users/{userID}", uh.Update)
		r.Get("/users/{userID}", uh.Read)
		r.Delete("/users/{userID}", uh.Delete)
		r.Post("/users/{userID}/suspend", uh.Suspend)
		r.Post("/users/{userID}/unsuspend", uh.Unsuspend)
	})
}

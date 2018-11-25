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
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/system/rest/request"
)

// Internal API interface
type UserAPI interface {
	Login(context.Context, *request.UserLogin) (interface{}, error)
	Logout(context.Context, *request.UserLogout) (interface{}, error)
	List(context.Context, *request.UserList) (interface{}, error)
	Create(context.Context, *request.UserCreate) (interface{}, error)
	Edit(context.Context, *request.UserEdit) (interface{}, error)
	Read(context.Context, *request.UserRead) (interface{}, error)
	Remove(context.Context, *request.UserRemove) (interface{}, error)
	Suspend(context.Context, *request.UserSuspend) (interface{}, error)
	Unsuspend(context.Context, *request.UserUnsuspend) (interface{}, error)
}

// HTTP API interface
type User struct {
	Login     func(http.ResponseWriter, *http.Request)
	Logout    func(http.ResponseWriter, *http.Request)
	List      func(http.ResponseWriter, *http.Request)
	Create    func(http.ResponseWriter, *http.Request)
	Edit      func(http.ResponseWriter, *http.Request)
	Read      func(http.ResponseWriter, *http.Request)
	Remove    func(http.ResponseWriter, *http.Request)
	Suspend   func(http.ResponseWriter, *http.Request)
	Unsuspend func(http.ResponseWriter, *http.Request)
}

func NewUser(uh UserAPI) *User {
	return &User{
		Login: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserLogin()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Login(r.Context(), params)
			})
		},
		Logout: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserLogout()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Logout(r.Context(), params)
			})
		},
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
		Edit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserEdit()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Edit(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Read(r.Context(), params)
			})
		},
		Remove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserRemove()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Remove(r.Context(), params)
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
		r.Route("/users", func(r chi.Router) {
			r.Post("/login", uh.Login)
			r.Get("/logout", uh.Logout)
			r.Get("/", uh.List)
			r.Post("/", uh.Create)
			r.Put("/{userID}", uh.Edit)
			r.Get("/{userID}", uh.Read)
			r.Delete("/{userID}", uh.Remove)
			r.Post("/{userID}/suspend", uh.Suspend)
			r.Post("/{userID}/unsuspend", uh.Unsuspend)
		})
	})
}

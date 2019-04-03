package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `auth.go`, `auth.util.go` or `auth_test.go` to
	implement your API calls, helper functions and tests. The file `auth.go`
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
type AuthAPI interface {
	Check(context.Context, *request.AuthCheck) (interface{}, error)
	Login(context.Context, *request.AuthLogin) (interface{}, error)
	Logout(context.Context, *request.AuthLogout) (interface{}, error)
}

// HTTP API interface
type Auth struct {
	Check  func(http.ResponseWriter, *http.Request)
	Login  func(http.ResponseWriter, *http.Request)
	Logout func(http.ResponseWriter, *http.Request)
}

func NewAuth(ah AuthAPI) *Auth {
	return &Auth{
		Check: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthCheck()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Check(r.Context(), params)
			})
		},
		Login: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthLogin()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Login(r.Context(), params)
			})
		},
		Logout: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthLogout()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Logout(r.Context(), params)
			})
		},
	}
}

func (ah *Auth) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/auth/check", ah.Check)
		r.Post("/auth/login", ah.Login)
		r.Get("/auth/logout", ah.Logout)
	})
}

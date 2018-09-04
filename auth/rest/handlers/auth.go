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
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/auth/rest/request"
)

// Internal API interface
type AuthAPI interface {
	Login(context.Context, *request.AuthLogin) (interface{}, error)
	Create(context.Context, *request.AuthCreate) (interface{}, error)
}

// HTTP API interface
type Auth struct {
	Login  func(http.ResponseWriter, *http.Request)
	Create func(http.ResponseWriter, *http.Request)
}

func NewAuth(ah AuthAPI) *Auth {
	return &Auth{
		Login: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthLogin()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Login(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Create(r.Context(), params)
			})
		},
	}
}

func (ah *Auth) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", ah.Login)
			r.Post("/create", ah.Create)
		})
	})
}

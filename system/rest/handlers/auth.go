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
	Settings(context.Context, *request.AuthSettings) (interface{}, error)
	Check(context.Context, *request.AuthCheck) (interface{}, error)
	ExchangeAuthToken(context.Context, *request.AuthExchangeAuthToken) (interface{}, error)
	Logout(context.Context, *request.AuthLogout) (interface{}, error)
}

// HTTP API interface
type Auth struct {
	Settings          func(http.ResponseWriter, *http.Request)
	Check             func(http.ResponseWriter, *http.Request)
	ExchangeAuthToken func(http.ResponseWriter, *http.Request)
	Logout            func(http.ResponseWriter, *http.Request)
}

func NewAuth(ah AuthAPI) *Auth {
	return &Auth{
		Settings: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthSettings()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Settings(r.Context(), params)
			})
		},
		Check: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthCheck()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Check(r.Context(), params)
			})
		},
		ExchangeAuthToken: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthExchangeAuthToken()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.ExchangeAuthToken(r.Context(), params)
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
		r.Get("/auth/", ah.Settings)
		r.Get("/auth/check", ah.Check)
		r.Post("/auth/exchange", ah.ExchangeAuthToken)
		r.Get("/auth/logout", ah.Logout)
	})
}

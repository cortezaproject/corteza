package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `permissions.go`, `permissions.util.go` or `permissions_test.go` to
	implement your API calls, helper functions and tests. The file `permissions.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/messaging/rest/request"
)

// Internal API interface
type PermissionsAPI interface {
	List(context.Context, *request.PermissionsList) (interface{}, error)
	Get(context.Context, *request.PermissionsGet) (interface{}, error)
	Set(context.Context, *request.PermissionsSet) (interface{}, error)
	Scopes(context.Context, *request.PermissionsScopes) (interface{}, error)
}

// HTTP API interface
type Permissions struct {
	List   func(http.ResponseWriter, *http.Request)
	Get    func(http.ResponseWriter, *http.Request)
	Set    func(http.ResponseWriter, *http.Request)
	Scopes func(http.ResponseWriter, *http.Request)
}

func NewPermissions(ph PermissionsAPI) *Permissions {
	return &Permissions{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.List(r.Context(), params)
			})
		},
		Get: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsGet()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Get(r.Context(), params)
			})
		},
		Set: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsSet()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Set(r.Context(), params)
			})
		},
		Scopes: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsScopes()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Scopes(r.Context(), params)
			})
		},
	}
}

func (ph *Permissions) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/permissions", func(r chi.Router) {
			r.Get("/", ph.List)
			r.Get("/{teamID}", ph.Get)
			r.Post("/{teamID}", ph.Set)
			r.Get("/scopes/{scope}", ph.Scopes)
		})
	})
}

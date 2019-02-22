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

	"github.com/crusttech/crust/system/rest/request"
)

// Internal API interface
type PermissionsAPI interface {
	List(context.Context, *request.PermissionsList) (interface{}, error)
	Read(context.Context, *request.PermissionsRead) (interface{}, error)
	Delete(context.Context, *request.PermissionsDelete) (interface{}, error)
	Update(context.Context, *request.PermissionsUpdate) (interface{}, error)
}

// HTTP API interface
type Permissions struct {
	List   func(http.ResponseWriter, *http.Request)
	Read   func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
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
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Read(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Delete(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Update(r.Context(), params)
			})
		},
	}
}

func (ph *Permissions) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/permissions", func(r chi.Router) {
			r.Get("/", ph.List)
			r.Get("/{roleID}/rules", ph.Read)
			r.Delete("/{roleID}/rules", ph.Delete)
			r.Patch("/{roleID}/rules", ph.Update)
		})
	})
}

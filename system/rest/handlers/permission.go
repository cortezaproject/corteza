package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `permission.go`, `permission.util.go` or `permission_test.go` to
	implement your API calls, helper functions and tests. The file `permission.go`
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
type PermissionAPI interface {
	Read(context.Context, *request.PermissionRead) (interface{}, error)
	Delete(context.Context, *request.PermissionDelete) (interface{}, error)
	Update(context.Context, *request.PermissionUpdate) (interface{}, error)
}

// HTTP API interface
type Permission struct {
	Read   func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
}

func NewPermission(ph PermissionAPI) *Permission {
	return &Permission{
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Read(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Delete(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Update(r.Context(), params)
			})
		},
	}
}

func (ph *Permission) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/permissions", func(r chi.Router) {
			r.Get("/{roleID}/rules", ph.Read)
			r.Delete("/{roleID}/rules", ph.Delete)
			r.Patch("/{roleID}/rules", ph.Update)
		})
	})
}

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

	"github.com/crusttech/crust/sam/rest/request"
)

// Internal API interface
type PermissionsAPI interface {
	List(context.Context, *request.PermissionsList) (interface{}, error)
	GetTeam(context.Context, *request.PermissionsGetTeam) (interface{}, error)
	SetTeam(context.Context, *request.PermissionsSetTeam) (interface{}, error)
}

// HTTP API interface
type Permissions struct {
	List    func(http.ResponseWriter, *http.Request)
	GetTeam func(http.ResponseWriter, *http.Request)
	SetTeam func(http.ResponseWriter, *http.Request)
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
		GetTeam: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsGetTeam()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.GetTeam(r.Context(), params)
			})
		},
		SetTeam: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsSetTeam()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.SetTeam(r.Context(), params)
			})
		},
	}
}

func (ph *Permissions) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/permissions", func(r chi.Router) {
			r.Get("/permissions", ph.List)
			r.Get("/permissions/{team}", ph.GetTeam)
			r.Post("/permissions/{team}", ph.SetTeam)
		})
	})
}

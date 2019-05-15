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

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/compose/rest/request"
	"github.com/crusttech/crust/internal/logger"
)

// Internal API interface
type PermissionsAPI interface {
	List(context.Context, *request.PermissionsList) (interface{}, error)
	Effective(context.Context, *request.PermissionsEffective) (interface{}, error)
	Read(context.Context, *request.PermissionsRead) (interface{}, error)
	Delete(context.Context, *request.PermissionsDelete) (interface{}, error)
	Update(context.Context, *request.PermissionsUpdate) (interface{}, error)
}

// HTTP API interface
type Permissions struct {
	List      func(http.ResponseWriter, *http.Request)
	Effective func(http.ResponseWriter, *http.Request)
	Read      func(http.ResponseWriter, *http.Request)
	Delete    func(http.ResponseWriter, *http.Request)
	Update    func(http.ResponseWriter, *http.Request)
}

func NewPermissions(h PermissionsAPI) *Permissions {
	return &Permissions{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Permissions.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := h.List(r.Context(), params); err != nil {
				logger.LogControllerError("Permissions.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("Permissions.List", r, params.Auditable())
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Effective: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsEffective()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Permissions.Effective", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := h.Effective(r.Context(), params); err != nil {
				logger.LogControllerError("Permissions.Effective", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("Permissions.Effective", r, params.Auditable())
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Permissions.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := h.Read(r.Context(), params); err != nil {
				logger.LogControllerError("Permissions.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("Permissions.Read", r, params.Auditable())
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Permissions.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := h.Delete(r.Context(), params); err != nil {
				logger.LogControllerError("Permissions.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("Permissions.Delete", r, params.Auditable())
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Permissions.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := h.Update(r.Context(), params); err != nil {
				logger.LogControllerError("Permissions.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("Permissions.Update", r, params.Auditable())
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
	}
}

func (h Permissions) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/permissions/", h.List)
		r.Get("/permissions/effective", h.Effective)
		r.Get("/permissions/{roleID}/rules", h.Read)
		r.Delete("/permissions/{roleID}/rules", h.Delete)
		r.Patch("/permissions/{roleID}/rules", h.Update)
	})
}

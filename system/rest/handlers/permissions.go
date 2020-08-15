package handlers

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/rest/request"
)

type (
	// Internal API interface
	PermissionsAPI interface {
		List(context.Context, *request.PermissionsList) (interface{}, error)
		Effective(context.Context, *request.PermissionsEffective) (interface{}, error)
		Read(context.Context, *request.PermissionsRead) (interface{}, error)
		Delete(context.Context, *request.PermissionsDelete) (interface{}, error)
		Update(context.Context, *request.PermissionsUpdate) (interface{}, error)
	}

	// HTTP API interface
	Permissions struct {
		List      func(http.ResponseWriter, *http.Request)
		Effective func(http.ResponseWriter, *http.Request)
		Read      func(http.ResponseWriter, *http.Request)
		Delete    func(http.ResponseWriter, *http.Request)
		Update    func(http.ResponseWriter, *http.Request)
	}
)

func NewPermissions(h PermissionsAPI) *Permissions {
	return &Permissions{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Permissions.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Permissions.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Permissions.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Effective: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsEffective()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Permissions.Effective", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Effective(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Permissions.Effective", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Permissions.Effective", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Permissions.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Permissions.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Permissions.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Permissions.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Permissions.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Permissions.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPermissionsUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Permissions.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Permissions.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Permissions.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
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

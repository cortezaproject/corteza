package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `namespace.go`, `namespace.util.go` or `namespace_test.go` to
	implement your API calls, helper functions and tests. The file `namespace.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

// Internal API interface
type NamespaceAPI interface {
	List(context.Context, *request.NamespaceList) (interface{}, error)
	Create(context.Context, *request.NamespaceCreate) (interface{}, error)
	Read(context.Context, *request.NamespaceRead) (interface{}, error)
	Update(context.Context, *request.NamespaceUpdate) (interface{}, error)
	Delete(context.Context, *request.NamespaceDelete) (interface{}, error)
}

// HTTP API interface
type Namespace struct {
	List   func(http.ResponseWriter, *http.Request)
	Create func(http.ResponseWriter, *http.Request)
	Read   func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
}

func NewNamespace(h NamespaceAPI) *Namespace {
	return &Namespace{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Namespace.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Namespace.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Namespace.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Namespace.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Namespace.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Namespace.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Namespace.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Namespace.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Namespace.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Namespace.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Namespace.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Namespace.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Namespace.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Namespace.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Namespace.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Namespace) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/", h.List)
		r.Post("/namespace/", h.Create)
		r.Get("/namespace/{namespaceID}", h.Read)
		r.Post("/namespace/{namespaceID}", h.Update)
		r.Delete("/namespace/{namespaceID}", h.Delete)
	})
}

package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `page.go`, `page.util.go` or `page_test.go` to
	implement your API calls, helper functions and tests. The file `page.go`
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
type PageAPI interface {
	List(context.Context, *request.PageList) (interface{}, error)
	Create(context.Context, *request.PageCreate) (interface{}, error)
	Read(context.Context, *request.PageRead) (interface{}, error)
	Tree(context.Context, *request.PageTree) (interface{}, error)
	Update(context.Context, *request.PageUpdate) (interface{}, error)
	Reorder(context.Context, *request.PageReorder) (interface{}, error)
	Delete(context.Context, *request.PageDelete) (interface{}, error)
	Upload(context.Context, *request.PageUpload) (interface{}, error)
}

// HTTP API interface
type Page struct {
	List    func(http.ResponseWriter, *http.Request)
	Create  func(http.ResponseWriter, *http.Request)
	Read    func(http.ResponseWriter, *http.Request)
	Tree    func(http.ResponseWriter, *http.Request)
	Update  func(http.ResponseWriter, *http.Request)
	Reorder func(http.ResponseWriter, *http.Request)
	Delete  func(http.ResponseWriter, *http.Request)
	Upload  func(http.ResponseWriter, *http.Request)
}

func NewPage(h PageAPI) *Page {
	return &Page{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Page.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Page.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Page.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Page.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Page.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Page.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Page.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Page.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Page.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Tree: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageTree()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Page.Tree", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Tree(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Page.Tree", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Page.Tree", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Page.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Page.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Page.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Reorder: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageReorder()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Page.Reorder", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Reorder(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Page.Reorder", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Page.Reorder", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Page.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Page.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Page.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Upload: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageUpload()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Page.Upload", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Upload(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Page.Upload", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Page.Upload", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Page) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/{namespaceID}/page/", h.List)
		r.Post("/namespace/{namespaceID}/page/", h.Create)
		r.Get("/namespace/{namespaceID}/page/{pageID}", h.Read)
		r.Get("/namespace/{namespaceID}/page/tree", h.Tree)
		r.Post("/namespace/{namespaceID}/page/{pageID}", h.Update)
		r.Post("/namespace/{namespaceID}/page/{selfID}/reorder", h.Reorder)
		r.Delete("/namespace/{namespaceID}/page/{pageID}", h.Delete)
		r.Post("/namespace/{namespaceID}/page/{pageID}/attachment", h.Upload)
	})
}

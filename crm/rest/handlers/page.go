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
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/crm/rest/request"
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
}

func NewPage(ph PageAPI) *Page {
	return &Page{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Create(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Read(r.Context(), params)
			})
		},
		Tree: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageTree()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Tree(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Update(r.Context(), params)
			})
		},
		Reorder: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageReorder()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Reorder(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ph.Delete(r.Context(), params)
			})
		},
	}
}

func (ph *Page) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/page", func(r chi.Router) {
			r.Get("/", ph.List)
			r.Post("/", ph.Create)
			r.Get("/{pageID}", ph.Read)
			r.Get("/tree", ph.Tree)
			r.Post("/{pageID}", ph.Update)
			r.Post("/{selfID}/reorder", ph.Reorder)
			r.Delete("/{pageID}", ph.Delete)
		})
	})
}

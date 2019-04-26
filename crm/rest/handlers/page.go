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

func NewPage(ph PageAPI) *Page {
	return &Page{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageList()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ph.List(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageCreate()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ph.Create(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
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
			params := request.NewPageRead()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ph.Read(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Tree: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageTree()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ph.Tree(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
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
			params := request.NewPageUpdate()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ph.Update(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Reorder: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageReorder()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ph.Reorder(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
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
			params := request.NewPageDelete()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ph.Delete(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Upload: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageUpload()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ph.Upload(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
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

func (ph *Page) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/page/", ph.List)
		r.Post("/page/", ph.Create)
		r.Get("/page/{pageID}", ph.Read)
		r.Get("/page/tree", ph.Tree)
		r.Post("/page/{pageID}", ph.Update)
		r.Post("/page/{selfID}/reorder", ph.Reorder)
		r.Delete("/page/{pageID}", ph.Delete)
		r.Post("/page/{pageID}/attachment", ph.Upload)
	})
}

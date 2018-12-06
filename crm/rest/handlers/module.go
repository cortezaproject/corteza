package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `module.go`, `module.util.go` or `module_test.go` to
	implement your API calls, helper functions and tests. The file `module.go`
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
type ModuleAPI interface {
	List(context.Context, *request.ModuleList) (interface{}, error)
	Create(context.Context, *request.ModuleCreate) (interface{}, error)
	Read(context.Context, *request.ModuleRead) (interface{}, error)
	Chart(context.Context, *request.ModuleChart) (interface{}, error)
	Edit(context.Context, *request.ModuleEdit) (interface{}, error)
	Delete(context.Context, *request.ModuleDelete) (interface{}, error)
	ContentReport(context.Context, *request.ModuleContentReport) (interface{}, error)
	ContentList(context.Context, *request.ModuleContentList) (interface{}, error)
	ContentCreate(context.Context, *request.ModuleContentCreate) (interface{}, error)
	ContentRead(context.Context, *request.ModuleContentRead) (interface{}, error)
	ContentEdit(context.Context, *request.ModuleContentEdit) (interface{}, error)
	ContentDelete(context.Context, *request.ModuleContentDelete) (interface{}, error)
}

// HTTP API interface
type Module struct {
	List          func(http.ResponseWriter, *http.Request)
	Create        func(http.ResponseWriter, *http.Request)
	Read          func(http.ResponseWriter, *http.Request)
	Chart         func(http.ResponseWriter, *http.Request)
	Edit          func(http.ResponseWriter, *http.Request)
	Delete        func(http.ResponseWriter, *http.Request)
	ContentReport func(http.ResponseWriter, *http.Request)
	ContentList   func(http.ResponseWriter, *http.Request)
	ContentCreate func(http.ResponseWriter, *http.Request)
	ContentRead   func(http.ResponseWriter, *http.Request)
	ContentEdit   func(http.ResponseWriter, *http.Request)
	ContentDelete func(http.ResponseWriter, *http.Request)
}

func NewModule(mh ModuleAPI) *Module {
	return &Module{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Create(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Read(r.Context(), params)
			})
		},
		Chart: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleChart()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Chart(r.Context(), params)
			})
		},
		Edit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleEdit()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Edit(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Delete(r.Context(), params)
			})
		},
		ContentReport: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleContentReport()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.ContentReport(r.Context(), params)
			})
		},
		ContentList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleContentList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.ContentList(r.Context(), params)
			})
		},
		ContentCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleContentCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.ContentCreate(r.Context(), params)
			})
		},
		ContentRead: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleContentRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.ContentRead(r.Context(), params)
			})
		},
		ContentEdit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleContentEdit()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.ContentEdit(r.Context(), params)
			})
		},
		ContentDelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleContentDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.ContentDelete(r.Context(), params)
			})
		},
	}
}

func (mh *Module) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/module", func(r chi.Router) {
			r.Get("/", mh.List)
			r.Post("/", mh.Create)
			r.Get("/{moduleID}", mh.Read)
			r.Get("/{moduleID}/chart", mh.Chart)
			r.Post("/{moduleID}", mh.Edit)
			r.Delete("/{moduleID}", mh.Delete)
			r.Get("/{moduleID}/report", mh.ContentReport)
			r.Get("/{moduleID}/content", mh.ContentList)
			r.Post("/{moduleID}/content", mh.ContentCreate)
			r.Get("/{moduleID}/content/{contentID}", mh.ContentRead)
			r.Post("/{moduleID}/content/{contentID}", mh.ContentEdit)
			r.Delete("/{moduleID}/content/{contentID}", mh.ContentDelete)
		})
	})
}

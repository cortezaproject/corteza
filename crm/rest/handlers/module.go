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
	Update(context.Context, *request.ModuleUpdate) (interface{}, error)
	Delete(context.Context, *request.ModuleDelete) (interface{}, error)
	RecordReport(context.Context, *request.ModuleRecordReport) (interface{}, error)
	RecordList(context.Context, *request.ModuleRecordList) (interface{}, error)
	RecordCreate(context.Context, *request.ModuleRecordCreate) (interface{}, error)
	RecordRead(context.Context, *request.ModuleRecordRead) (interface{}, error)
	RecordUpdate(context.Context, *request.ModuleRecordUpdate) (interface{}, error)
	RecordDelete(context.Context, *request.ModuleRecordDelete) (interface{}, error)
}

// HTTP API interface
type Module struct {
	List         func(http.ResponseWriter, *http.Request)
	Create       func(http.ResponseWriter, *http.Request)
	Read         func(http.ResponseWriter, *http.Request)
	Update       func(http.ResponseWriter, *http.Request)
	Delete       func(http.ResponseWriter, *http.Request)
	RecordReport func(http.ResponseWriter, *http.Request)
	RecordList   func(http.ResponseWriter, *http.Request)
	RecordCreate func(http.ResponseWriter, *http.Request)
	RecordRead   func(http.ResponseWriter, *http.Request)
	RecordUpdate func(http.ResponseWriter, *http.Request)
	RecordDelete func(http.ResponseWriter, *http.Request)
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
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Update(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Delete(r.Context(), params)
			})
		},
		RecordReport: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleRecordReport()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.RecordReport(r.Context(), params)
			})
		},
		RecordList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleRecordList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.RecordList(r.Context(), params)
			})
		},
		RecordCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleRecordCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.RecordCreate(r.Context(), params)
			})
		},
		RecordRead: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleRecordRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.RecordRead(r.Context(), params)
			})
		},
		RecordUpdate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleRecordUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.RecordUpdate(r.Context(), params)
			})
		},
		RecordDelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleRecordDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.RecordDelete(r.Context(), params)
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
			r.Post("/{moduleID}", mh.Update)
			r.Delete("/{moduleID}", mh.Delete)
			r.Get("/{moduleID}/report", mh.RecordReport)
			r.Get("/{moduleID}/record", mh.RecordList)
			r.Post("/{moduleID}/record", mh.RecordCreate)
			r.Get("/{moduleID}/record/{recordID}", mh.RecordRead)
			r.Post("/{moduleID}/record/{recordID}", mh.RecordUpdate)
			r.Delete("/{moduleID}/record/{recordID}", mh.RecordDelete)
		})
	})
}

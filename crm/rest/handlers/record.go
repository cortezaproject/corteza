package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `record.go`, `record.util.go` or `record_test.go` to
	implement your API calls, helper functions and tests. The file `record.go`
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
type RecordAPI interface {
	Report(context.Context, *request.RecordReport) (interface{}, error)
	List(context.Context, *request.RecordList) (interface{}, error)
	Create(context.Context, *request.RecordCreate) (interface{}, error)
	Read(context.Context, *request.RecordRead) (interface{}, error)
	Update(context.Context, *request.RecordUpdate) (interface{}, error)
	Delete(context.Context, *request.RecordDelete) (interface{}, error)
	Upload(context.Context, *request.RecordUpload) (interface{}, error)
}

// HTTP API interface
type Record struct {
	Report func(http.ResponseWriter, *http.Request)
	List   func(http.ResponseWriter, *http.Request)
	Create func(http.ResponseWriter, *http.Request)
	Read   func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
	Upload func(http.ResponseWriter, *http.Request)
}

func NewRecord(rh RecordAPI) *Record {
	return &Record{
		Report: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordReport()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Report(r.Context(), params)
			})
		},
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Create(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Read(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Update(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Delete(r.Context(), params)
			})
		},
		Upload: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordUpload()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Upload(r.Context(), params)
			})
		},
	}
}

func (rh *Record) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/module/{moduleID}/record/report", rh.Report)
		r.Get("/module/{moduleID}/record/", rh.List)
		r.Post("/module/{moduleID}/record/", rh.Create)
		r.Get("/module/{moduleID}/record/{recordID}", rh.Read)
		r.Post("/module/{moduleID}/record/{recordID}", rh.Update)
		r.Delete("/module/{moduleID}/record/{recordID}", rh.Delete)
		r.Post("/module/{moduleID}/record/{recordID}/{fieldName}/attachment", rh.Upload)
	})
}

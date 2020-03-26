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

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

// Internal API interface
type RecordAPI interface {
	Report(context.Context, *request.RecordReport) (interface{}, error)
	List(context.Context, *request.RecordList) (interface{}, error)
	ImportInit(context.Context, *request.RecordImportInit) (interface{}, error)
	ImportRun(context.Context, *request.RecordImportRun) (interface{}, error)
	ImportProgress(context.Context, *request.RecordImportProgress) (interface{}, error)
	Export(context.Context, *request.RecordExport) (interface{}, error)
	Exec(context.Context, *request.RecordExec) (interface{}, error)
	Create(context.Context, *request.RecordCreate) (interface{}, error)
	Read(context.Context, *request.RecordRead) (interface{}, error)
	Update(context.Context, *request.RecordUpdate) (interface{}, error)
	BulkDelete(context.Context, *request.RecordBulkDelete) (interface{}, error)
	Delete(context.Context, *request.RecordDelete) (interface{}, error)
	Upload(context.Context, *request.RecordUpload) (interface{}, error)
	TriggerScript(context.Context, *request.RecordTriggerScript) (interface{}, error)
	TriggerScriptOnList(context.Context, *request.RecordTriggerScriptOnList) (interface{}, error)
}

// HTTP API interface
type Record struct {
	Report              func(http.ResponseWriter, *http.Request)
	List                func(http.ResponseWriter, *http.Request)
	ImportInit          func(http.ResponseWriter, *http.Request)
	ImportRun           func(http.ResponseWriter, *http.Request)
	ImportProgress      func(http.ResponseWriter, *http.Request)
	Export              func(http.ResponseWriter, *http.Request)
	Exec                func(http.ResponseWriter, *http.Request)
	Create              func(http.ResponseWriter, *http.Request)
	Read                func(http.ResponseWriter, *http.Request)
	Update              func(http.ResponseWriter, *http.Request)
	BulkDelete          func(http.ResponseWriter, *http.Request)
	Delete              func(http.ResponseWriter, *http.Request)
	Upload              func(http.ResponseWriter, *http.Request)
	TriggerScript       func(http.ResponseWriter, *http.Request)
	TriggerScriptOnList func(http.ResponseWriter, *http.Request)
}

func NewRecord(h RecordAPI) *Record {
	return &Record{
		Report: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordReport()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.Report", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Report(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.Report", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.Report", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ImportInit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordImportInit()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.ImportInit", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ImportInit(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.ImportInit", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.ImportInit", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ImportRun: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordImportRun()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.ImportRun", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ImportRun(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.ImportRun", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.ImportRun", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ImportProgress: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordImportProgress()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.ImportProgress", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ImportProgress(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.ImportProgress", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.ImportProgress", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Export: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordExport()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.Export", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Export(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.Export", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.Export", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Exec: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordExec()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.Exec", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Exec(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.Exec", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.Exec", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		BulkDelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordBulkDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.BulkDelete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.BulkDelete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.BulkDelete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.BulkDelete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Upload: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordUpload()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.Upload", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Upload(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.Upload", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.Upload", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		TriggerScript: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordTriggerScript()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.TriggerScript", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.TriggerScript(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.TriggerScript", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.TriggerScript", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		TriggerScriptOnList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordTriggerScriptOnList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Record.TriggerScriptOnList", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.TriggerScriptOnList(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Record.TriggerScriptOnList", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Record.TriggerScriptOnList", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Record) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/{namespaceID}/module/{moduleID}/record/report", h.Report)
		r.Get("/namespace/{namespaceID}/module/{moduleID}/record/", h.List)
		r.Post("/namespace/{namespaceID}/module/{moduleID}/record/import", h.ImportInit)
		r.Patch("/namespace/{namespaceID}/module/{moduleID}/record/import/{sessionID}", h.ImportRun)
		r.Get("/namespace/{namespaceID}/module/{moduleID}/record/import/{sessionID}", h.ImportProgress)
		r.Get("/namespace/{namespaceID}/module/{moduleID}/record/export{filename}.{ext}", h.Export)
		r.Post("/namespace/{namespaceID}/module/{moduleID}/record/exec/{procedure}", h.Exec)
		r.Post("/namespace/{namespaceID}/module/{moduleID}/record/", h.Create)
		r.Get("/namespace/{namespaceID}/module/{moduleID}/record/{recordID}", h.Read)
		r.Post("/namespace/{namespaceID}/module/{moduleID}/record/{recordID}", h.Update)
		r.Delete("/namespace/{namespaceID}/module/{moduleID}/record/", h.BulkDelete)
		r.Delete("/namespace/{namespaceID}/module/{moduleID}/record/{recordID}", h.Delete)
		r.Post("/namespace/{namespaceID}/module/{moduleID}/record/attachment", h.Upload)
		r.Post("/namespace/{namespaceID}/module/{moduleID}/record/{recordID}/trigger", h.TriggerScript)
		r.Post("/namespace/{namespaceID}/module/{moduleID}/record/trigger", h.TriggerScriptOnList)
	})
}

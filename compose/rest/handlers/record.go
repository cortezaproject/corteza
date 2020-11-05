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
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	RecordAPI interface {
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
	Record struct {
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
)

func NewRecord(h RecordAPI) *Record {
	return &Record{
		Report: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordReport()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Report(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ImportInit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordImportInit()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ImportInit(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ImportRun: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordImportRun()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ImportRun(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ImportProgress: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordImportProgress()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ImportProgress(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Export: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordExport()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Export(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Exec: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordExec()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Exec(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordRead()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordUpdate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		BulkDelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordBulkDelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.BulkDelete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordDelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Upload: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordUpload()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Upload(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		TriggerScript: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordTriggerScript()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.TriggerScript(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		TriggerScriptOnList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRecordTriggerScriptOnList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.TriggerScriptOnList(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
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

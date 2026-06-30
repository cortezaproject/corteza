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
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	DmlAPI interface {
		ConnectionList(context.Context, *request.DmlConnectionList) (interface{}, error)
		ConnectionCreate(context.Context, *request.DmlConnectionCreate) (interface{}, error)
		ConnectionRead(context.Context, *request.DmlConnectionRead) (interface{}, error)
		ConnectionModels(context.Context, *request.DmlConnectionModels) (interface{}, error)
		MappingCreate(context.Context, *request.DmlMappingCreate) (interface{}, error)
		MappingRead(context.Context, *request.DmlMappingRead) (interface{}, error)
		MappingUpdate(context.Context, *request.DmlMappingUpdate) (interface{}, error)
		MappingDelete(context.Context, *request.DmlMappingDelete) (interface{}, error)
		ImportRun(context.Context, *request.DmlImportRun) (interface{}, error)
		ImportRunRead(context.Context, *request.DmlImportRunRead) (interface{}, error)
	}

	// HTTP API interface
	Dml struct {
		ConnectionList   func(http.ResponseWriter, *http.Request)
		ConnectionCreate func(http.ResponseWriter, *http.Request)
		ConnectionRead   func(http.ResponseWriter, *http.Request)
		ConnectionModels func(http.ResponseWriter, *http.Request)
		MappingCreate    func(http.ResponseWriter, *http.Request)
		MappingRead      func(http.ResponseWriter, *http.Request)
		MappingUpdate    func(http.ResponseWriter, *http.Request)
		MappingDelete    func(http.ResponseWriter, *http.Request)
		ImportRun        func(http.ResponseWriter, *http.Request)
		ImportRunRead    func(http.ResponseWriter, *http.Request)
	}
)

func NewDml(h DmlAPI) *Dml {
	return &Dml{
		ConnectionList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDmlConnectionList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ConnectionList(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ConnectionCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDmlConnectionCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ConnectionCreate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ConnectionRead: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDmlConnectionRead()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ConnectionRead(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ConnectionModels: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDmlConnectionModels()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ConnectionModels(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MappingCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDmlMappingCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MappingCreate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MappingRead: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDmlMappingRead()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MappingRead(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MappingUpdate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDmlMappingUpdate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MappingUpdate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MappingDelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDmlMappingDelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MappingDelete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ImportRun: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDmlImportRun()
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
		ImportRunRead: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDmlImportRunRead()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ImportRunRead(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Dml) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/dml/connections", h.ConnectionList)
		r.Post("/dml/connections", h.ConnectionCreate)
		r.Get("/dml/connections/{connectionID}", h.ConnectionRead)
		r.Get("/dml/connections/{connectionID}/models", h.ConnectionModels)
		r.Post("/dml/connections/{connectionID}/mapping", h.MappingCreate)
		r.Get("/dml/connections/{connectionID}/mapping/{mappingID}", h.MappingRead)
		r.Put("/dml/connections/{connectionID}/mapping/{mappingID}", h.MappingUpdate)
		r.Delete("/dml/connections/{connectionID}/mapping/{mappingID}", h.MappingDelete)
		r.Post("/dml/connections/{connectionID}/mapping/{mappingID}/import", h.ImportRun)
		r.Get("/dml/connections/{connectionID}/mapping/{mappingID}/import/{runID}", h.ImportRunRead)
	})
}

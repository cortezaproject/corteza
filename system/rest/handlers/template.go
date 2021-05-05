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
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	TemplateAPI interface {
		List(context.Context, *request.TemplateList) (interface{}, error)
		Create(context.Context, *request.TemplateCreate) (interface{}, error)
		Read(context.Context, *request.TemplateRead) (interface{}, error)
		Update(context.Context, *request.TemplateUpdate) (interface{}, error)
		Delete(context.Context, *request.TemplateDelete) (interface{}, error)
		Undelete(context.Context, *request.TemplateUndelete) (interface{}, error)
		RenderDrivers(context.Context, *request.TemplateRenderDrivers) (interface{}, error)
		Render(context.Context, *request.TemplateRender) (interface{}, error)
	}

	// HTTP API interface
	Template struct {
		List          func(http.ResponseWriter, *http.Request)
		Create        func(http.ResponseWriter, *http.Request)
		Read          func(http.ResponseWriter, *http.Request)
		Update        func(http.ResponseWriter, *http.Request)
		Delete        func(http.ResponseWriter, *http.Request)
		Undelete      func(http.ResponseWriter, *http.Request)
		RenderDrivers func(http.ResponseWriter, *http.Request)
		Render        func(http.ResponseWriter, *http.Request)
	}
)

func NewTemplate(h TemplateAPI) *Template {
	return &Template{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTemplateList()
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
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTemplateCreate()
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
			params := request.NewTemplateRead()
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
			params := request.NewTemplateUpdate()
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
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTemplateDelete()
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
		Undelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTemplateUndelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Undelete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RenderDrivers: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTemplateRenderDrivers()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RenderDrivers(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Render: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTemplateRender()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Render(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Template) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/template/", h.List)
		r.Post("/template/", h.Create)
		r.Get("/template/{templateID}", h.Read)
		r.Put("/template/{templateID}", h.Update)
		r.Delete("/template/{templateID}", h.Delete)
		r.Post("/template/{templateID}/undelete", h.Undelete)
		r.Get("/template/render/drivers", h.RenderDrivers)
		r.Post("/template/{templateID}/render/{filename}.{ext}", h.Render)
	})
}

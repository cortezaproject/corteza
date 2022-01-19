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
	GigAPI interface {
		Create(context.Context, *request.GigCreate) (interface{}, error)
		Go(context.Context, *request.GigGo) (interface{}, error)
		Read(context.Context, *request.GigRead) (interface{}, error)
		Update(context.Context, *request.GigUpdate) (interface{}, error)
		Delete(context.Context, *request.GigDelete) (interface{}, error)
		Undelete(context.Context, *request.GigUndelete) (interface{}, error)
		AddSource(context.Context, *request.GigAddSource) (interface{}, error)
		RemoveSource(context.Context, *request.GigRemoveSource) (interface{}, error)
		Prepare(context.Context, *request.GigPrepare) (interface{}, error)
		Exec(context.Context, *request.GigExec) (interface{}, error)
		Output(context.Context, *request.GigOutput) (interface{}, error)
		OutputAll(context.Context, *request.GigOutputAll) (interface{}, error)
		OutputSpecific(context.Context, *request.GigOutputSpecific) (interface{}, error)
		State(context.Context, *request.GigState) (interface{}, error)
		Status(context.Context, *request.GigStatus) (interface{}, error)
		Complete(context.Context, *request.GigComplete) (interface{}, error)
		Workers(context.Context, *request.GigWorkers) (interface{}, error)
		Tasks(context.Context, *request.GigTasks) (interface{}, error)
	}

	// HTTP API interface
	Gig struct {
		Create         func(http.ResponseWriter, *http.Request)
		Go             func(http.ResponseWriter, *http.Request)
		Read           func(http.ResponseWriter, *http.Request)
		Update         func(http.ResponseWriter, *http.Request)
		Delete         func(http.ResponseWriter, *http.Request)
		Undelete       func(http.ResponseWriter, *http.Request)
		AddSource      func(http.ResponseWriter, *http.Request)
		RemoveSource   func(http.ResponseWriter, *http.Request)
		Prepare        func(http.ResponseWriter, *http.Request)
		Exec           func(http.ResponseWriter, *http.Request)
		Output         func(http.ResponseWriter, *http.Request)
		OutputAll      func(http.ResponseWriter, *http.Request)
		OutputSpecific func(http.ResponseWriter, *http.Request)
		State          func(http.ResponseWriter, *http.Request)
		Status         func(http.ResponseWriter, *http.Request)
		Complete       func(http.ResponseWriter, *http.Request)
		Workers        func(http.ResponseWriter, *http.Request)
		Tasks          func(http.ResponseWriter, *http.Request)
	}
)

func NewGig(h GigAPI) *Gig {
	return &Gig{
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigCreate()
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
		Go: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigGo()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Go(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigRead()
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
			params := request.NewGigUpdate()
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
			params := request.NewGigDelete()
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
			params := request.NewGigUndelete()
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
		AddSource: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigAddSource()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.AddSource(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RemoveSource: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigRemoveSource()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RemoveSource(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Prepare: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigPrepare()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Prepare(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Exec: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigExec()
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
		Output: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigOutput()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Output(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		OutputAll: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigOutputAll()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.OutputAll(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		OutputSpecific: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigOutputSpecific()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.OutputSpecific(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		State: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigState()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.State(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Status: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigStatus()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Status(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Complete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigComplete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Complete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Workers: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigWorkers()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Workers(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Tasks: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGigTasks()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Tasks(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Gig) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/gig/", h.Create)
		r.Post("/gig/go", h.Go)
		r.Get("/gig/{gigID}", h.Read)
		r.Put("/gig/{gigID}", h.Update)
		r.Delete("/gig/{gigID}", h.Delete)
		r.Post("/gig/{gigID}/undelete", h.Undelete)
		r.Patch("/gig/{gigID}/sources", h.AddSource)
		r.Delete("/gig/{gigID}/sources/{sourceID}", h.RemoveSource)
		r.Put("/gig/{gigID}/prepare", h.Prepare)
		r.Put("/gig/{gigID}/exec", h.Exec)
		r.Get("/gig/{gigID}/output", h.Output)
		r.Get("/gig/{gigID}/output/all", h.OutputAll)
		r.Get("/gig/{gigID}/output/{sourceID}", h.OutputSpecific)
		r.Get("/gig/{gigID}/state", h.State)
		r.Get("/gig/{gigID}/status", h.Status)
		r.Patch("/gig/{gigID}/complete", h.Complete)
		r.Get("/gig/workers", h.Workers)
		r.Get("/gig/tasks", h.Tasks)
	})
}

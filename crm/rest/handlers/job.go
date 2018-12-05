package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `job.go`, `job.util.go` or `job_test.go` to
	implement your API calls, helper functions and tests. The file `job.go`
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
type JobAPI interface {
	List(context.Context, *request.JobList) (interface{}, error)
	Run(context.Context, *request.JobRun) (interface{}, error)
	Get(context.Context, *request.JobGet) (interface{}, error)
	Logs(context.Context, *request.JobLogs) (interface{}, error)
	Update(context.Context, *request.JobUpdate) (interface{}, error)
	Delete(context.Context, *request.JobDelete) (interface{}, error)
}

// HTTP API interface
type Job struct {
	List   func(http.ResponseWriter, *http.Request)
	Run    func(http.ResponseWriter, *http.Request)
	Get    func(http.ResponseWriter, *http.Request)
	Logs   func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
}

func NewJob(jh JobAPI) *Job {
	return &Job{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewJobList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return jh.List(r.Context(), params)
			})
		},
		Run: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewJobRun()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return jh.Run(r.Context(), params)
			})
		},
		Get: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewJobGet()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return jh.Get(r.Context(), params)
			})
		},
		Logs: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewJobLogs()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return jh.Logs(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewJobUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return jh.Update(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewJobDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return jh.Delete(r.Context(), params)
			})
		},
	}
}

func (jh *Job) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/job", func(r chi.Router) {
			r.Get("/", jh.List)
			r.Post("/", jh.Run)
			r.Get("/{jobID}", jh.Get)
			r.Get("/{jobID}/logs", jh.Logs)
			r.Post("/{jobID}", jh.Update)
			r.Delete("/{jobID}", jh.Delete)
		})
	})
}

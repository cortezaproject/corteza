package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `organisation.go`, `organisation.util.go` or `organisation_test.go` to
	implement your API calls, helper functions and tests. The file `organisation.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/system/rest/request"
)

// Internal API interface
type OrganisationAPI interface {
	List(context.Context, *request.OrganisationList) (interface{}, error)
	Create(context.Context, *request.OrganisationCreate) (interface{}, error)
	Update(context.Context, *request.OrganisationUpdate) (interface{}, error)
	Delete(context.Context, *request.OrganisationDelete) (interface{}, error)
	Read(context.Context, *request.OrganisationRead) (interface{}, error)
	Archive(context.Context, *request.OrganisationArchive) (interface{}, error)
}

// HTTP API interface
type Organisation struct {
	List    func(http.ResponseWriter, *http.Request)
	Create  func(http.ResponseWriter, *http.Request)
	Update  func(http.ResponseWriter, *http.Request)
	Delete  func(http.ResponseWriter, *http.Request)
	Read    func(http.ResponseWriter, *http.Request)
	Archive func(http.ResponseWriter, *http.Request)
}

func NewOrganisation(oh OrganisationAPI) *Organisation {
	return &Organisation{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return oh.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return oh.Create(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return oh.Update(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return oh.Delete(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return oh.Read(r.Context(), params)
			})
		},
		Archive: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationArchive()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return oh.Archive(r.Context(), params)
			})
		},
	}
}

func (oh *Organisation) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/organisations", func(r chi.Router) {
			r.Get("/", oh.List)
			r.Post("/", oh.Create)
			r.Put("/{id}", oh.Update)
			r.Delete("/{id}", oh.Delete)
			r.Get("/{id}", oh.Read)
			r.Post("/{id}/archive", oh.Archive)
		})
	})
}

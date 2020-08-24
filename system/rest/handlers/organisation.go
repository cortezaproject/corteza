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
	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/rest/request"
)

type (
	// Internal API interface
	OrganisationAPI interface {
		List(context.Context, *request.OrganisationList) (interface{}, error)
		Create(context.Context, *request.OrganisationCreate) (interface{}, error)
		Update(context.Context, *request.OrganisationUpdate) (interface{}, error)
		Delete(context.Context, *request.OrganisationDelete) (interface{}, error)
		Read(context.Context, *request.OrganisationRead) (interface{}, error)
		Archive(context.Context, *request.OrganisationArchive) (interface{}, error)
	}

	// HTTP API interface
	Organisation struct {
		List    func(http.ResponseWriter, *http.Request)
		Create  func(http.ResponseWriter, *http.Request)
		Update  func(http.ResponseWriter, *http.Request)
		Delete  func(http.ResponseWriter, *http.Request)
		Read    func(http.ResponseWriter, *http.Request)
		Archive func(http.ResponseWriter, *http.Request)
	}
)

func NewOrganisation(h OrganisationAPI) *Organisation {
	return &Organisation{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Organisation.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Organisation.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Organisation.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Organisation.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Organisation.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Organisation.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Organisation.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Organisation.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Organisation.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Organisation.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Organisation.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Organisation.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Organisation.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Organisation.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Organisation.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Archive: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewOrganisationArchive()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Organisation.Archive", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Archive(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Organisation.Archive", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Organisation.Archive", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Organisation) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/organisations/", h.List)
		r.Post("/organisations/", h.Create)
		r.Put("/organisations/{id}", h.Update)
		r.Delete("/organisations/{id}", h.Delete)
		r.Get("/organisations/{id}", h.Read)
		r.Post("/organisations/{id}/archive", h.Archive)
	})
}

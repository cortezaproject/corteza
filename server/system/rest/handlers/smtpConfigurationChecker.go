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
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	SmtpConfigurationCheckerAPI interface {
		Check(context.Context, *request.SmtpConfigurationCheckerCheck) (interface{}, error)
	}

	// HTTP API interface
	SmtpConfigurationChecker struct {
		Check func(http.ResponseWriter, *http.Request)
	}
)

func NewSmtpConfigurationChecker(h SmtpConfigurationCheckerAPI) *SmtpConfigurationChecker {
	return &SmtpConfigurationChecker{
		Check: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSmtpConfigurationCheckerCheck()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Check(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h SmtpConfigurationChecker) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/smtp/configuration-checker/", h.Check)
	})
}

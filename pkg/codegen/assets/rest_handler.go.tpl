package handlers

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

import (
	"context"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/{{ .App }}/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

type (
    // Internal API interface
    {{ pubIdent $.Endpoint.Entrypoint }}API interface {
    {{- range $a := $.Endpoint.Apis }}
        {{ pubIdent $a.Name }}(context.Context, *request.{{ pubIdent $.Endpoint.Entrypoint $a.Name }}) (interface{}, error)
    {{- end }}
    }

    // HTTP API interface
    {{ pubIdent .Endpoint.Entrypoint }} struct {
    {{- range $a := .Endpoint.Apis }}
        {{ pubIdent $a.Name }} func(http.ResponseWriter, *http.Request)
    {{- end }}
    }
)


func {{ pubIdent "New" $.Endpoint.Entrypoint }}(h {{ pubIdent $.Endpoint.Entrypoint }}API) *{{ pubIdent $.Endpoint.Entrypoint }} {
	return &{{ pubIdent $.Endpoint.Entrypoint }}{
    {{- range $a := .Endpoint.Apis }}
		{{ pubIdent $a.Name }}: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.New{{ pubIdent $.Endpoint.Entrypoint $a.Name }}()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("{{ pubIdent $.Endpoint.Entrypoint }}.{{ pubIdent $a.Name }}", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.{{ pubIdent $a.Name }}(r.Context(), params)
			if err != nil {
				logger.LogControllerError("{{ pubIdent $.Endpoint.Entrypoint }}.{{ pubIdent $a.Name }}", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("{{ pubIdent $.Endpoint.Entrypoint }}.{{ pubIdent $a.Name }}", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
    {{- end }}
	}
}

func (h {{ pubIdent $.Endpoint.Entrypoint }}) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)

		{{- range $a := .Endpoint.Apis }}
		r.{{ pubIdent ( toLower $a.Method ) }}("{{ $.Endpoint.Path }}{{ $a.Path }}", h.{{ pubIdent $a.Name }})
		{{- end }}
	})
}

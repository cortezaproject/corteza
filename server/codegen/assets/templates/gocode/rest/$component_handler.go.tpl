package {{ .package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"context"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/crusttech/human/server/{{ .app }}/rest/request"
	"github.com/crusttech/human/server/pkg/api"
)

type (
    // Internal API interface
    {{ .entrypointExp }}API interface {
    {{- range $a := .apis }}
        {{- if $a.raw }}
        {{ $a.nameExp }}(http.ResponseWriter, *http.Request)
        {{- else }}
        {{ $a.nameExp }}(context.Context, *request.{{ $a.reqIdent }}) (interface{}, error)
        {{- end }}
    {{- end }}
    }

    // HTTP API interface
    {{ .entrypointExp }} struct {
    {{- range $a := .apis }}
        {{ $a.nameExp }} func(http.ResponseWriter, *http.Request)
    {{- end }}
    }
)

func New{{ .entrypointExp }}(h {{ .entrypointExp }}API) *{{ .entrypointExp }} {
	return &{{ .entrypointExp }}{
    {{- range $a := .apis }}
        {{- if $a.raw }}
		{{ $a.nameExp }}: h.{{ $a.nameExp }},
        {{- else }}
		{{ $a.nameExp }}: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.New{{ $a.reqIdent }}()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.{{ $a.nameExp }}(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
        {{- end }}
    {{- end }}
	}
}

func (h {{ .entrypointExp }}) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)

		{{- range $a := .apis }}
		r.{{ $a.methodRoute }}("{{ $a.routePath }}", h.{{ $a.nameExp }})
		{{- end }}
	})
}

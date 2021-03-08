package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

import (
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	"io"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
{{- range .Imports }}
    {{ normalizeImport . }}
{{- end }}
)

// dummy vars to prevent
// unused imports complain
var (
    _ = chi.URLParam
    _ = multipart.ErrMessageTooLarge
    _ = payload.ParseUint64s
    _ = strings.ToLower
    _ = io.EOF
    _ = fmt.Errorf
    _ = json.NewEncoder
)

type (
    // Internal API interface
    {{- range $a := $.Endpoint.Apis }}
    {{ export $.Endpoint.Entrypoint $a.Name }} struct {
    {{- range $p := $a.Params.All }}
        // {{ export $p.Name }} {{ $p.Origin }} parameter
        //
        // {{ $p.Title }}
        {{ export $p.Name }} {{ $p.Type }} {{ $p.FieldTag }}
    {{ end }}
    }
    {{ end }}
)

{{- range $a := $.Endpoint.Apis }}
// {{ export "New" $.Endpoint.Entrypoint $a.Name }} request
func {{ export "New" $.Endpoint.Entrypoint $a.Name }}() *{{ export $.Endpoint.Entrypoint $a.Name }} {
	return &{{ export $.Endpoint.Entrypoint $a.Name }}{}
}

// Auditable returns all auditable/loggable parameters
func (r {{ export $.Endpoint.Entrypoint $a.Name }}) Auditable() map[string]interface{} {
	return map[string]interface{}{
    {{- range $p := $a.Params.All }}
    {{- if not $p.Sensitive }}
    	"{{ $p.Name }}": r.{{ export $p.Name }},
    {{- end }}
    {{- end }}
	}
}

{{- range $p := $a.Params.All }}
// Auditable returns all auditable/loggable parameters
func (r {{ export $.Endpoint.Entrypoint $a.Name }}) Get{{ export $p.Name }}() {{ $p.Type }} {
	return r.{{ export $p.Name }}
}
{{- end }}



// Fill processes request and fills internal variables
func (r *{{ export $.Endpoint.Entrypoint $a.Name }}) Fill(req *http.Request) (err error) {
    {{ if $a.Params.Post }}
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}
    {{- end }}


    {{ if $a.Params.Get }}
    {
        // GET params
	    tmp := req.URL.Query()
	{{ range $p := $a.Params.Get }}
        {{- if or $p.IsSlice $p.HasExplicitParser }}
        if val, ok := tmp["{{ $p.Name }}[]"]; ok   {
            r.{{ export $p.Name }}, err = {{ $p.Parser "val" }}
            if err != nil {
                return err
            }
        } else if val, ok := tmp["{{ $p.Name }}"]; ok   {
            r.{{ export $p.Name }}, err = {{ $p.Parser "val" }}
            if err != nil {
                return err
            }
        }
        {{- else }}
        if val, ok := tmp["{{ $p.Name }}"]; ok && len(val) > 0  {
            r.{{ export $p.Name }}, err = {{ $p.Parser "val[0]" }}
            if err != nil {
                return err
            }
        }
        {{- end }}
    {{- end }}
	}
	{{- end }}

    {{ if $a.Params.Post }}
    {
	if err = req.ParseForm(); err != nil {
		return err
	}

        // POST params
        {{ range $p := $a.Params.Post }}
            {{ if $p.IsUpload }}
            if _, r.{{ export $p.Name }}, err = req.FormFile("{{ $p.Name }}"); err != nil {
                return fmt.Errorf("error processing uploaded file: %w", err)
            }
            {{ else }}
                {{- if or $p.HasExplicitParser }}
                if val, ok := req.Form["{{ $p.Name }}[]"]; ok {
                    r.{{ export $p.Name }}, err = {{ $p.Parser "val" }}
                    if err != nil {
                        return err
                    }
                } else if val, ok := req.Form["{{ $p.Name }}"]; ok   {
					r.{{ export $p.Name }}, err = {{ $p.Parser "val" }}
					if err != nil {
						return err
					}
				}
                {{- else if or $p.IsSlice }}
                //if val, ok := req.Form["{{ $p.Name }}[]"]; ok && len(val) > 0  {
                //    r.{{ export $p.Name }}, err = {{ $p.Parser "val" }}
                //    if err != nil {
                //        return err
                //    }
                //}
                {{- else }}
                if val, ok := req.Form["{{ $p.Name }}"]; ok && len(val) > 0  {
                    r.{{ export $p.Name }}, err = {{ $p.Parser "val[0]" }}
                    if err != nil {
                        return err
                    }
                }
                {{- end }}
            {{- end }}

        {{- end }}
	}
	{{ end }}

	{{ if $a.Params.Path }}
    {
        var val string
        // path params
	{{ range $p := $a.Params.Path }}
        val = chi.URLParam(req, "{{ $p.Name }}")
        r.{{ export $p.Name }}, err = {{ $p.Parser "val" }}
        if err != nil {
            return err
        }
	{{ end }}

	}
	{{ end }}

	return err
}

{{- end }}

package {{ .package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"github.com/crusttech/human/server/pkg/payload"
	"github.com/go-chi/chi/v5"
	"io"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
{{- range .importsNorm }}
    {{ . }}
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
    {{- range $a := .apis }}
    {{ $a.reqIdent }} struct {
    {{- range $p := $a.params.all }}
        // {{ $p.exportedName }} {{ $p.origin }} parameter
        //
        // {{ $p.title }}
        {{ $p.exportedName }} {{ $p.type }} {{ $p.fieldTag }}
    {{ end }}
    }
    {{ end }}
)

{{- range $a := .apis }}
// New{{ $a.reqIdent }} request
func New{{ $a.reqIdent }}() *{{ $a.reqIdent }} {
	return &{{ $a.reqIdent }}{}
}

// Auditable returns all auditable/loggable parameters
func (r {{ $a.reqIdent }}) Auditable() map[string]interface{} {
	return map[string]interface{}{
    {{- range $p := $a.params.all }}
    {{- if not $p.sensitive }}
    	"{{ $p.name }}": r.{{ $p.exportedName }},
    {{- end }}
    {{- end }}
	}
}

{{- range $p := $a.params.all }}
// Auditable returns all auditable/loggable parameters
func (r {{ $a.reqIdent }}) Get{{ $p.exportedName }}() {{ $p.type }} {
	return r.{{ $p.exportedName }}
}
{{- end }}



// Fill processes request and fills internal variables
func (r *{{ $a.reqIdent }}) Fill(req *http.Request) (err error) {
    {{ if $a.params.post }}
	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}
    {{- end }}


    {{ if $a.params.get }}
    {
        // GET params
	    tmp := req.URL.Query()
	{{ range $p := $a.params.get }}
        {{- if or $p.isSlice $p.hasExplicitParser }}
        if val, ok := tmp["{{ $p.name }}[]"]; ok   {
            r.{{ $p.exportedName }}, err = {{ $p.parserVal }}
            if err != nil {
                return err
            }
        } else if val, ok := tmp["{{ $p.name }}"]; ok   {
            r.{{ $p.exportedName }}, err = {{ $p.parserVal }}
            if err != nil {
                return err
            }
        }
        {{- else }}
        if val, ok := tmp["{{ $p.name }}"]; ok && len(val) > 0  {
            r.{{ $p.exportedName }}, err = {{ $p.parserVal0 }}
            if err != nil {
                return err
            }
        }
        {{- end }}
    {{- end }}
	}
	{{- end }}

    {{ if $a.params.post }}
    {
        // Caching 32MB to memory, the rest to disk
        if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
            return err
        } else if err == nil {
            // Multipart params
            {{ range $p := $a.params.post }}
                {{ if $p.isUpload }}
                // Ignoring {{ $p.name }} as its handled in the POST params section
                {{- else }}
                    {{- if or $p.hasExplicitParser }}
                    if val, ok := req.MultipartForm.Value["{{ $p.name }}[]"]; ok {
                        r.{{ $p.exportedName }}, err = {{ $p.parserVal }}
                        if err != nil {
                            return err
                        }
                    } else if val, ok := req.MultipartForm.Value["{{ $p.name }}"]; ok   {
                        r.{{ $p.exportedName }}, err = {{ $p.parserVal }}
                        if err != nil {
                            return err
                        }
                    }
                    {{- else if not $p.isSlice }}
                    if val, ok := req.MultipartForm.Value["{{ $p.name }}"]; ok && len(val) > 0  {
                        r.{{ $p.exportedName }}, err = {{ $p.parserVal0 }}
                        if err != nil {
                            return err
                        }
                    }
                    {{- end }}
                {{- end }}

            {{- end }}
        }
	}

    {
	if err = req.ParseForm(); err != nil {
		return err
	}

        // POST params
        {{ range $p := $a.params.post }}
            {{ if $p.isUpload }}
            if _, r.{{ $p.exportedName }}, err = req.FormFile("{{ $p.name }}"); err != nil {
                return fmt.Errorf("error processing uploaded file: %w", err)
            }
            {{ else }}
                {{- if or $p.hasExplicitParser }}
                if val, ok := req.Form["{{ $p.name }}[]"]; ok {
                    r.{{ $p.exportedName }}, err = {{ $p.parserVal }}
                    if err != nil {
                        return err
                    }
                } else if val, ok := req.Form["{{ $p.name }}"]; ok   {
					r.{{ $p.exportedName }}, err = {{ $p.parserVal }}
					if err != nil {
						return err
					}
				}
                {{- else if or $p.isSlice }}
                //if val, ok := req.Form["{{ $p.name }}[]"]; ok && len(val) > 0  {
                //    r.{{ $p.exportedName }}, err = {{ $p.parserVal }}
                //    if err != nil {
                //        return err
                //    }
                //}
                {{- else }}
                if val, ok := req.Form["{{ $p.name }}"]; ok && len(val) > 0  {
                    r.{{ $p.exportedName }}, err = {{ $p.parserVal0 }}
                    if err != nil {
                        return err
                    }
                }
                {{- end }}
            {{- end }}

        {{- end }}
	}
	{{ end }}

	{{ if $a.params.path }}
    {
        var val string
        // path params
	{{ range $p := $a.params.path }}
        val = chi.URLParam(req, "{{ $p.name }}")
        r.{{ $p.exportedName }}, err = {{ $p.parserVal }}
        if err != nil {
            return err
        }
	{{ end }}

	}
	{{ end }}

	return err
}

{{- end }}

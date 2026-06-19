package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"github.com/crusttech/human/server/pkg/resourceref"
)

// ResourceRefs returns configuration-level references to other resources
func (r {{ .expIdent }}) ResourceRefs() (out []resourceref.Ref) {
{{- range .refs }}
{{- if eq .iter "direct" }}
	out = resourceref.Append(out, resourceref.Make(
		resourceref.{{ .kind }}, r.{{ .path }}, resourceref.{{ .reason }},
	))
{{- else if eq .iter "sliceID" }}
	for _, id := range r.{{ .path }} {
		out = resourceref.Append(out, resourceref.Make(
			resourceref.{{ .kind }}, id, resourceref.{{ .reason }},
		))
	}
{{- else if eq .iter "sliceField" }}
	for _, item := range r.{{ .path }} {
		out = resourceref.Append(out, resourceref.Make(
			resourceref.{{ .kind }}, item.{{ .field }}, resourceref.{{ .reason }},
		))
	}
{{- end }}
{{- end }}
{{- if .extended }}
	return r.resourceRefsExt(out)
{{- else }}
	return
{{- end }}
}

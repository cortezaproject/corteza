package {{ .Package }}

{{ template "header-gentext.tpl" }}
{{ template "header-definitions.tpl" . }}

import (
{{- range .Imports }}
    {{ . }}
{{- end }}
)

{{- range .Def }}
{{- if gt (len .RBAC.Resource.References) 0 }}
// {{ export .Component .Resource }}RbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func {{ export .Component .Resource }}RbacReferences({{- range .RBAC.Resource.References }}{{ unexport .Resource }} string, {{- end }}) (res *Ref, pp []*Ref, err error) {
	{{- range .RBAC.Resource.References }}
	{{- if eq .Field "ID" }}
	if {{ unexport .Resource }} != "*" {
		res = &Ref{ResourceType: types.{{ export .Resource }}ResourceType, Identifiers: MakeIdentifiers({{ unexport .Resource }})}
	}
	{{- else }}
	if {{ unexport .Resource }} != "*" {
		pp = append(pp, &Ref{ResourceType: types.{{ export .Resource }}ResourceType, Identifiers: MakeIdentifiers({{ unexport .Resource }})})
	}
	{{- end }}
	{{- end }}

	return
}
{{- end }}
{{- end }}

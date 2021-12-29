package {{ .Package }}

{{ template "header-gentext.tpl" }}
{{ template "header-definitions.tpl" . }}

import (
{{- range .Imports }}
    {{ . }}
{{- end }}
)

{{- range .Def }}
{{- if .Locale }}
{{- if gt (len .Locale.Resource.References) 0 }}
// {{ export .Component .Resource }}ResourceTranslationReferences generates Locale references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func {{ export .Component .Resource }}ResourceTranslationReferences({{- range .Locale.Resource.References }}{{ unexport .Resource }} string, {{- end }}) (res *Ref, pp []*Ref, err error) {
	{{- range .Locale.Resource.References }}
	{{- if eq .Field "ID" }}
	res = &Ref{ResourceType: types.{{ export .Resource }}ResourceType, Identifiers: MakeIdentifiers({{ unexport .Resource }})}
	{{- else }}
	pp = append(pp, &Ref{ResourceType: types.{{ export .Resource }}ResourceType, Identifiers: MakeIdentifiers({{ unexport .Resource }})})
	{{- end }}
	{{- end }}

	return
}
{{- end }}
{{- end }}
{{- end }}

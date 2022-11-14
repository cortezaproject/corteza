package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
{{- range .imports }}
    {{ . }}
{{- end }}
)

{{- range .resources }}
// {{ .rbacRefFunc }} generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func {{ .rbacRefFunc }}({{- range .references }}{{ .param }} string, {{- end }}) (res *Ref, pp []*Ref, err error) {
	{{- range .references }}
		{{- if eq .refField "ID" }}
		if {{ .param }} != "*" {
			res = &Ref{ResourceType: types.{{ .expIdent }}ResourceType, Identifiers: MakeIdentifiers({{ .param }})}
		}
		{{- else }}
		if {{ .param }} != "*" {
			pp = append(pp, &Ref{ResourceType: types.{{ .expIdent }}ResourceType, Identifiers: MakeIdentifiers({{ .param }})})
		}
		{{- end }}
	{{- end }}

	return
}
{{- end }}

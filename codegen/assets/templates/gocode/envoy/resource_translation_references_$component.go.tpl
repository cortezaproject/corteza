package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
{{- range .imports }}
    {{ . }}
{{- end }}
)

{{- range .resources }}
// {{ .resTrRefFunc }} generates Locale references
//
// This function is auto-generated
func {{ .resTrRefFunc }}({{- range .references }}{{ .param }} string, {{- end }} self string) (res *Ref, pp []*Ref, err error) {
	res = &Ref{ResourceType: types.{{ .expIdent }}ResourceType, Identifiers: MakeIdentifiers(self)}

	{{- range .references }}
		pp = append(pp, &Ref{ResourceType: types.{{ .expIdent }}ResourceType, Identifiers: MakeIdentifiers({{ .param }})})
	{{- end }}

	return
}
{{- end }}

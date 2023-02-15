package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"github.com/cortezaproject/corteza/server/pkg/cast2"
)

{{ range .resources }}
	{{ if .model.omitGetterSetter }}
		{{continue}}
	{{ end }}

func (r {{.expIdent}}) GetID() (uint64) {
	{{- $ok := false -}}
	{{- range .model.attributes -}}
		{{- if eq .expIdent "ID" -}}
		return r.ID
			{{- $ok = true -}}
			{{break}}
		{{- end -}}
	{{- end -}}
	{{- if not $ok }}
	// The resource does not define an ID field
	return 0
	{{- end -}}
}

func (r *{{.expIdent}}) GetValue(name string, pos uint) (any, error) {
	switch name {
	{{ range .model.attributes -}}
		{{- if .omitGetter -}}
			{{continue}}
		{{- end -}}
		{{- $identAlias := .identAlias -}}
		case {{ range $i, $l := $identAlias -}}
			"{{ $l }}"{{if not (eq $i (sub (len $identAlias) 1))}},{{end}}
		{{- end}}:
			return r.{{.expIdent}}, nil
	{{ end }}

	{{ if .model.defaultGetter }}
	default:
		return r.getValue(name, pos)
	{{end}}
	}
	return nil, nil
}

func (r *{{.expIdent}}) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	{{ range .model.attributes -}}
		{{- if .omitSetter -}}
			{{continue}}
		{{- end -}}
		{{ $identAlias := .identAlias -}}
		case {{ range $i, $l := $identAlias -}}
			"{{ $l }}"{{if not (eq $i (sub (len $identAlias) 1))}},{{end}}
		{{- end}}:
			return cast2.{{ .goCastFnc }}(value, &r.{{ .expIdent }})
	{{ end }}

	{{ if .model.defaultSetter }}
	default:
		return r.setValue(name, pos, value)
	{{end}}
	}
	return nil
}

{{ end }}
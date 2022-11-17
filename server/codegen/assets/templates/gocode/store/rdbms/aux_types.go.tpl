package rdbms

{{ template "gocode/header-gentext.tpl" }}

import (
	"time"
	"github.com/cortezaproject/corteza/server/pkg/expr"
{{- range $path, $alias :=  .imports }}
    {{ $alias }} {{ printf "%q" $path }}
{{- end }}
)

type (
{{ range .types }}
	// {{ .auxIdent }} is an auxiliary structure used for transporting to/from RDBMS store
	{{ .auxIdent }} struct {
	{{ range .auxStruct }}
		{{ .expIdent }} {{ .goType }} {{ printf "`db:%q`" .name }}
	{{- end }}
	}
{{ end }}
)

{{- range .types }}
// encodes {{ .expIdent }} to {{ .auxIdent }}
//
// This function is auto-generated
func (aux *{{ .auxIdent }}) encode(res *{{ .goType }}) (_ error) {
	{{- range .auxStruct }}
	aux.{{ .expIdent }} = res.{{ .expIdent }}
	{{- end }}
	return
}

// decodes {{ .expIdent }} from {{ .auxIdent }}
//
// This function is auto-generated
func (aux {{ .auxIdent }}) decode() (res *{{ .goType }}, _ error) {
	res = new({{ .goType }})
	{{- range .auxStruct }}
		res.{{ .expIdent }} = aux.{{ .expIdent }}
	{{- end }}
	return
}

// scans row and fills {{ .auxIdent }} fields
//
// This function is auto-generated
func (aux *{{ .auxIdent }})scan(row scanner) (error) {
	return row.Scan(
	{{- range .auxStruct }}
		&aux.{{ .expIdent }},
	{{- end }}
	)
}

{{ end }}

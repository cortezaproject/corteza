package rdbms

{{ template "gocode/header-gentext.tpl" }}

import (
	"strings"
	{{- range $path, $alias :=  .imports }}
    {{ $alias }} {{ printf "%q" $path }}
	{{- end }}
	"github.com/doug-martin/goqu/v9"
)

type (
	// extendedFilters allows special per-resource
	// filters to be attached to store
	//
	// when optional filter is set, generated filter function is NOT called automatically
	// (but can be called from the optional filter)
	extendedFilters struct {
		// Filter extensions for search/query functions
	{{ range .types }}

		// optional {{ .ident }} filter function called after the generated function
		{{ .expIdent }} func({{ .goFilterType }}) ([]goqu.Expression, {{ .goFilterType }}, error)
	{{ end }}
	}
)

{{- range .types }}
// {{ .expIdent }}Filter returns logical expressions
//
// This function is called from Store.Query{{ .expIdentPlural }}() and can be extended
// by setting Store.Filters.{{ .expIdent }}. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func {{ .expIdent }}Filter(f {{ .goFilterType }})(ee []goqu.Expression, _ {{ .goFilterType }}, err error) {
	{{ range .filter.byNilState }}
		if expr := stateNilComparison({{ printf "%q" .storeIdent }}, f.{{ .expIdent }}); expr != nil {
			ee = append(ee, expr)
		}
	{{ end }}

	{{ range .filter.byFalseState }}
		if expr := stateFalseComparison({{ printf "%q" .storeIdent }}, f.{{ .expIdent }}); expr != nil {
			ee = append(ee, expr)
		}
	{{ end }}

	{{ range .filter.byValue }}
		{{ if eq .goType "string" }}
		if val := strings.TrimSpace(f.{{ .expIdent }}); len(val) > 0 {
			ee = append(ee, goqu.C({{ printf "%q" .storeIdent }}).Eq(f.{{ .expIdent }}))
		}
		{{ else if eq .goType "bool" }}
		if f.{{ .expIdent }} {
			ee = append(ee, goqu.C({{ printf "%q" .storeIdent }}).IsTrue())
		}
		{{ else if eq .goType "uint64" }}
		if f.{{ .expIdent }} > 0 {
			ee = append(ee, goqu.C({{ printf "%q" .storeIdent }}).Eq(f.{{ .expIdent }}))
		}
		{{ else if eq .goType "[]uint64" }}
		if len(f.{{ .expIdent }}) > 0 {
			ee = append(ee, goqu.C({{ printf "%q" .storeIdent }}).In(f.{{ .expIdent }}))
		}
		{{ else }}
		// @todo codegen warning: filtering by {{ .expIdent }} ({{ .goType }}) not supported,
		//       see rdbms.go.tpl and add an exception
		{{ end }}
	{{ end }}

	{{ if .filter.byLabel }}
	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}
	{{ end }}

	{{ if .filter.query }}
	if f.Query != "" {
		ee = append(ee, goqu.Or(
		{{- range .filter.query }}
			goqu.C({{ printf "%q" .storeIdent }}).ILike("%" + f.Query + "%"),
		{{- end }}
		))
	}
	{{ end }}

	return ee, f, err
}
{{ end }}

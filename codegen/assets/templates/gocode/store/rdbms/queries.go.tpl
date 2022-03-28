package rdbms

{{ template "gocode/header-gentext.tpl" }}

import (
	"github.com/doug-martin/goqu/v9"
{{- range $path, $alias :=  .imports }}
    {{ $alias }} {{ printf "%q" $path }}
{{- end }}
)

var (
{{- range .types }}
// {{ .ident }}Table represents {{ .identPlural }} store table
//
// This value is auto-generated
{{ .ident }}Table = goqu.T({{ printf "%q" .settings.rdbms.table }})

// {{ .ident }}SelectQuery assembles select query for fetching {{ .identPlural }}
//
// This function is auto-generated
{{ .ident }}SelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
	return d.Select(
		{{- range .struct }}
		{{ printf "%q" .storeIdent }},
		{{- end }}
	).From({{ .ident }}Table)
}

// {{ .ident }}InsertQuery assembles query inserting {{ .identPlural }}
//
// This function is auto-generated
{{ .ident }}InsertQuery = func(d goqu.DialectWrapper, res *{{ .goType }}) *goqu.InsertDataset {
	return d.Insert({{ .ident }}Table).
		Rows(goqu.Record{
		{{- range .struct }}
			{{ printf "%q" .storeIdent }}: res.{{ .expIdent }},
		{{- end }}
		})
}

// {{ .ident }}UpsertQuery assembles (insert+on-conflict) query for replacing {{ .identPlural }}
//
// This function is auto-generated
{{ .ident }}UpsertQuery = func(d goqu.DialectWrapper, res *{{ .goType }}) *goqu.InsertDataset {
	var target = `
	{{- range .struct -}}
		{{- if .primaryKey -}}
			,
			{{- if .ignoreCase -}}
			LOWER({{- .storeIdent -}})
			{{- else -}}
			{{- .storeIdent -}}
			{{- end -}}
		{{- end -}}
	{{- end -}}`

	return {{ .ident }}InsertQuery(d, res).
		OnConflict(
			goqu.DoUpdate(target[1:],
			goqu.Record{
			{{- range .struct }}
				{{- if not .primaryKey }}
					{{ printf "%q" .storeIdent }}: res.{{ .expIdent }},
				{{- end }}
			{{- end }}
			},
		),
	)
}

// {{ .ident }}UpdateQuery assembles query for updating {{ .identPlural }}
//
// This function is auto-generated
{{ .ident }}UpdateQuery = func(d goqu.DialectWrapper, res *{{ .goType }}) *goqu.UpdateDataset {
	return d.Update({{ .ident }}Table).
		Set(goqu.Record{
		{{- range .struct }}
			{{- if not .primaryKey }}
				{{ printf "%q" .storeIdent }}: res.{{ .expIdent }},
			{{- end }}
		{{- end }}
		}).
		Where({{ .ident }}PrimaryKeys(res))
}

// {{ .ident }}DeleteQuery assembles delete query for removing {{ .identPlural }}
//
// This function is auto-generated
{{ .ident }}DeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
	return d.Delete({{ .ident }}Table).Where(ee...)
}

// {{ .ident }}DeleteQuery assembles delete query for removing {{ .identPlural }}
//
// This function is auto-generated
{{ .ident }}TruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
	return d.Truncate({{ .ident }}Table)
}

// {{ .ident }}PrimaryKeys assembles set of conditions for all primary keys
//
// This function is auto-generated
{{ .ident }}PrimaryKeys = func(res *{{ .goType }}) goqu.Ex {
	return goqu.Ex{
	{{- range .struct }}
		{{- if .primaryKey }}
			{{ printf "%q" .storeIdent }}: res.{{ .expIdent }},
		{{- end }}
	{{- end }}
	}
}
{{ end }}
)

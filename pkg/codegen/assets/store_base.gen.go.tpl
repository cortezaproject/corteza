package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: {{ .Source }}
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
{{- range .Import }}
	{{ normalizeImport . }}
{{- end }}
)

type (
	{{- $Types := .Types }}
	{{- $Fields := .Fields }}

	{{ export .Types.Plural }} interface {

{{ if .Publish }}
	{{- if .Search.Enable }}
		Search{{ export $Types.Plural }}(ctx context.Context{{ template "extraArgsDef" . }}, f {{ $Types.GoFilterType }}) ({{ $Types.GoSetType }}, {{ $Types.GoFilterType }}, error)
	{{- end }}
{{- range .Lookups }}
		Lookup{{ export $Types.Singular }}By{{ export .Suffix }}(ctx context.Context{{ template "extraArgsDef" $ }}{{- range .Fields }}, {{ cc2underscore .Field }} {{ .Type  }}{{- end }}) (*{{ $Types.GoType }}, error)
{{- end }}
	{{ if .Create.Enable }}
		Create{{ export $Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $Types.GoType }}) error
	{{- end }}
	{{ if .Update.Enable }}
		Update{{ export $Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $Types.GoType }}) error
	{{- end }}
	{{ if .Upsert.Enable }}
		Upsert{{ export $Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $Types.GoType }}) error
	{{- end }}
	{{ if .Delete.Enable }}
		Delete{{ export $Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $Types.GoType }}) error
		Delete{{ export $Types.Singular }}By{{ template "primaryKeySuffix" $Fields }}(ctx context.Context{{ template "extraArgsDef" . }} {{ template "primaryKeyArgsDef" $Fields }}) error
	{{- end }}

		Truncate{{ export $Types.Plural }}(ctx context.Context{{ template "extraArgsDef" . }}) error
{{ end }}

{{- if .Functions}}
		// Additional custom functions
	{{- range .Functions }}

		// {{ .Name }} (custom function)
		{{ .Name }}(ctx context.Context{{ template "extraArgsDef" . }}) ({{ join ", " .Return }})
	{{- end }}
{{- end -}}
	}
)

{{/* convering scenario with non-exported main functions and no additional functions defined
     where we get "imported and not used error" */}}
var _ *{{ $Types.GoType }}
var _ context.Context

{{ if .Publish }}
{{- if .Search.Enable }}
// Search{{ export $.Types.Plural }} returns all matching {{ $.Types.Plural }} from store
func Search{{ export $Types.Plural }}(ctx context.Context, s {{ export $Types.Plural }}{{ template "extraArgsDef" . }}, f {{ $Types.GoFilterType }}) ({{ $Types.GoSetType }}, {{ $Types.GoFilterType }}, error) {
	return s.Search{{ export $Types.Plural }}(ctx{{ template "extraArgsCall" . }}, f)
}
{{- end -}}

{{ range .Lookups }}

// Lookup{{ export $.Types.Singular }}By{{ export .Suffix }} {{ comment .Description true -}}
func Lookup{{ export $Types.Singular }}By{{ export .Suffix }}(ctx context.Context, s {{ export $Types.Plural }}{{ template "extraArgsDef" $ }}{{- range .Fields }}, {{ cc2underscore .Field }} {{ .Type }}{{- end }}) (*{{ $Types.GoType }}, error) {
    return s.Lookup{{ export $Types.Singular }}By{{ export .Suffix }}(ctx{{ template "extraArgsCall" $ }}{{- range .Fields }}, {{ cc2underscore .Field }}{{- end }})
}
{{- end }}

{{ if .Create.Enable }}
// Create{{ export $.Types.Singular }} creates one or more {{ $.Types.Plural }} in store
func Create{{ export $Types.Singular }}(ctx context.Context, s {{ export $Types.Plural }}{{ template "extraArgsDef" . }}, rr ... *{{ $Types.GoType }}) error {
	return s.Create{{ export $Types.Singular }}(ctx{{ template "extraArgsCall" . }}, rr... )
}
{{- end }}

{{ if .Update.Enable }}
// Update{{ export $.Types.Singular }} updates one or more (existing) {{ $.Types.Plural }} in store
func Update{{ export $Types.Singular }}(ctx context.Context, s {{ export $Types.Plural }}{{ template "extraArgsDef" . }}, rr ... *{{ $Types.GoType }}) error {
	return s.Update{{ export $Types.Singular }}(ctx{{ template "extraArgsCall" . }}, rr... )
}
{{ end }}

{{ if .Upsert.Enable }}
// Upsert{{ export $.Types.Singular }} creates new or updates existing one or more {{ $.Types.Plural }} in store
func Upsert{{ export $Types.Singular }}(ctx context.Context, s {{ export $Types.Plural }}{{ template "extraArgsDef" . }}, rr ... *{{ $Types.GoType }}) error {
	return s.Upsert{{ export $Types.Singular }}(ctx{{ template "extraArgsCall" . }}, rr... )
}
{{ end }}

{{ if .Delete.Enable }}
// Delete{{ export $.Types.Singular }} Deletes one or more {{ $.Types.Plural }} from store
func Delete{{ export $Types.Singular }}(ctx context.Context, s {{ export $Types.Plural }}{{ template "extraArgsDef" . }}, rr ... *{{ $Types.GoType }}) error {
	return s.Delete{{ export $Types.Singular }}(ctx{{ template "extraArgsCall" . }}, rr...)
}

// Delete{{ export $.Types.Singular }}By{{ template "primaryKeySuffix" $.Fields }} Deletes {{ $.Types.Singular }} from store
func Delete{{ export $Types.Singular }}By{{ template "primaryKeySuffix" $Fields }}(ctx context.Context, s {{ export $Types.Plural }}{{ template "extraArgsDef" . }} {{ template "primaryKeyArgsDef" $Fields }}) error {
	return s.Delete{{ export $Types.Singular }}By{{ template "primaryKeySuffix" $Fields }}(ctx{{ template "extraArgsCall" . }}{{ template "primaryKeyArgsCall" $Fields }})
}
{{ end }}

// Truncate{{ export $.Types.Plural }} Deletes all {{ $.Types.Plural }} from store
func Truncate{{ export $Types.Plural }}(ctx context.Context, s {{ export $Types.Plural }}{{ template "extraArgsDef" . }}) error {
	return s.Truncate{{ export $Types.Plural }}(ctx{{ template "extraArgsCall" . }})
}
{{ end }}

{{ range .Functions }}
func {{ .Name }}(ctx context.Context, s {{ export $Types.Plural }}{{ template "extraArgsDef" . }}) ({{ join ", " .Return }}) {
	return s.{{ .Name }}(ctx{{ template "extraArgsCall" . }})
}
{{ end }}

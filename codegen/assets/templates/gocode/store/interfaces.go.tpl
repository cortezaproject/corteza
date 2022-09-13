package store

{{ template "gocode/header-gentext.tpl" }}

import (
	"context"
	"go.uber.org/zap"
{{- range $path, $alias :=  .imports }}
    {{ $alias }} {{ printf "%q" $path }}
{{- end }}
  "github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"golang.org/x/text/language"
)

{{ define "extraArgs" -}}
{{/*This is temporary solution until we properly implement Compose Record Store*/}}
{{- if eq . "composeRecord" }}mod *composeType.Module, {{ end -}}
{{- end }}
{{ define "extraParams" -}}
{{/*This is temporary solution until we properly implement Compose Record Store*/}}
{{- if eq . "composeRecord" }}mod, {{ end -}}
{{- end }}

type (
	// Storer interface combines interfaces of all supported store interfaces
	Storer interface {
		// SetLogger sets new logging facility
		//
		// Store facility should fallback to logger.Default when no logging facility is set
		//
		// Intentionally closely coupled with Zap logger since this is not some public lib
		// and it's highly unlikely we'll support different/multiple logging "backend"
		SetLogger(*zap.Logger)

		// Returns underlying store as DAL connection
		ToDalConn() dal.Connection

		// Tx is a transaction handler
		Tx(context.Context, func(context.Context, Storer) error) error

		// Upgrade store's schema to the latest version
		Upgrade(context.Context) error

		{{- range .types }}
			{{ .expIdentPlural }}
		{{- end }}
	}

{{ range .types }}
	{{ .expIdentPlural }} interface {
		Search{{ .expIdentPlural }}(ctx context.Context, {{ template "extraArgs" .ident }} f {{ .goFilterType }}) ({{ .goSetType }}, {{ .goFilterType }}, error)
		Create{{ .expIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) error
		Update{{ .expIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) error
		Upsert{{ .expIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) error
		Delete{{ .expIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) error
		{{ .api.deleteByPK.expFnIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} {{ range .api.deleteByPK.attributes }}{{ .ident }} {{ .goType }},{{ end }}) error
		Truncate{{ .expIdentPlural }}(ctx context.Context, {{ template "extraArgs" .ident }}) error

		{{- range .api.lookups }}
		{{ .expFnIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} {{ range .args }}{{ .ident }} {{ .goType }}, {{ end }}) (*{{ .returnType }}, error)
		{{- end }}
		{{- range .api.functions }}
		{{ .expFnIdent }}(ctx context.Context, {{ range .args }}{{ .ident }} {{ if .spread}}...{{ end }}{{ .goType }}, {{ end }}) ({{ range .return }}{{ . }}, {{ end }} error)
		{{- end }}
	}
{{ end }}
)

{{- range .types }}
	// Search{{ .expIdentPlural }} returns all matching {{ .expIdentPlural }} from store
	//
	// This function is auto-generated
	func Search{{ .expIdentPlural }}(ctx context.Context, s {{ .expIdentPlural }}, {{ template "extraArgs" .ident }} f {{ .goFilterType }}) ({{ .goSetType }}, {{ .goFilterType }}, error) {
		return s.Search{{ .expIdentPlural }}(ctx, {{ template "extraParams" .ident }} f)
	}

	// Create{{ .expIdent }} creates one or more {{ .expIdentPlural }} in store
	//
	// This function is auto-generated
	func Create{{ .expIdent }}(ctx context.Context, s {{ .expIdentPlural }}, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) error {
		return s.Create{{ .expIdent }}(ctx, {{ template "extraParams" .ident }}rr...)
	}

	// Update{{ .expIdent }} updates one or more (existing) {{ .expIdentPlural }} in store
	//
	// This function is auto-generated
	func Update{{ .expIdent }}(ctx context.Context, s {{ .expIdentPlural }}, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) error {
		return s.Update{{ .expIdent }}(ctx, {{ template "extraParams" .ident }}rr...)
	}

	// Upsert{{ .expIdent }} creates new or updates existing one or more {{ .expIdentPlural }} in store
	//
	// This function is auto-generated
	func Upsert{{ .expIdent }}(ctx context.Context, s {{ .expIdentPlural }}, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) error {
		return s.Upsert{{ .expIdent }}(ctx, {{ template "extraParams" .ident }}rr...)
	}

	// Delete{{ .expIdent }} deletes one or more {{ .expIdentPlural }} from store
	//
	// This function is auto-generated
	func Delete{{ .expIdent }}(ctx context.Context, s {{ .expIdentPlural }}, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) error {
		return s.Delete{{ .expIdent }}(ctx, {{ template "extraParams" .ident }}rr...)
	}

	// Delete{{ .expIdent }}ByID deletes one or more {{ .expIdentPlural }} from store
	//
	// This function is auto-generated
	func {{ .api.deleteByPK.expFnIdent }}(ctx context.Context, s {{ .expIdentPlural }}, {{ template "extraArgs" .ident }} {{ range .api.deleteByPK.attributes }}{{ .ident }} {{ .goType }},{{ end }}) error {
		return s.{{ .api.deleteByPK.expFnIdent }}(ctx, {{ template "extraParams" .ident }}{{ range .api.deleteByPK.attributes }}{{ .ident }},{{ end }})
	}

	// Truncate{{ .expIdentPlural }} Deletes all {{ .expIdentPlural }} from store
	//
	// This function is auto-generated
	func Truncate{{ .expIdentPlural }}(ctx context.Context, s {{ .expIdentPlural }}, {{ template "extraArgs" .ident }}) error {
		return s.Truncate{{ .expIdentPlural }}(ctx, {{ template "extraParams" .ident }})
	}

	{{- range .api.lookups }}
		{{ if .description	}}{{ .description }}{{ end }}
		//
		// This function is auto-generated
		func {{ .expFnIdent }}(ctx context.Context, s {{ .expStoreIdent }}, {{ template "extraArgs" .ident }} {{ range .args }}{{ .ident }} {{ .goType }}, {{ end }}) (*{{ .returnType }}, error) {
			return s.{{ .expFnIdent }}(ctx, {{ template "extraParams" .ident }}{{ range .args }}{{ .ident }}, {{ end }})
		}
	{{- end }}
	{{- range .api.functions }}
		{{ if .description	}}{{ .description }}{{ end }}
		//
		// This function is auto-generated
		func {{ .expFnIdent }}(ctx context.Context, s {{ .expStoreIdent }}, {{ range .args }}{{ .ident }} {{ if .spread}}...{{ end }}{{ .goType }}, {{ end }}) ({{ range .return }}{{ . }}, {{ end }} error) {
			return s.{{ .expFnIdent }}(ctx, {{ range .args }}{{ .ident }}{{ if .spread}}...{{ end }}, {{ end }})
		}
	{{- end }}
{{- end }}

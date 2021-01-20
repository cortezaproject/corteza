package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

import (
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
{{- range .Imports }}
  {{ normalizeImport . }}
{{- end }}
)

var _ wfexec.ExecResponse

type (
	{{ $.Name }}HandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h {{ $.Name }}Handler) register() {
	h.reg.AddFunctions(
	{{- range .Functions }}
		h.{{ export .Name }}(),
	{{- end }}
	)
}

{{ range .Functions }}
	{{ $REF     := unexport $.Prefix $.Name .Name }}
	{{ $ARGS    := unexport $.Name .Name "Args" }}
	{{ $RESULTS := unexport $.Name .Name "Results" }}


type (
	{{ $ARGS }} struct {
		{{ range .Params }}
			{{ $name    := .Name }}
			{{ $isArray := .IsArray }}
			has{{ export .Name }} bool
			{{- if gt (len .Types) 1 }}
				{{ export .Name }} interface{}
				{{- range .Types }}
					{{ $name }}{{ export .Suffix }} {{ if $isArray }}[]{{ end }}{{ .GoType }}
				{{- end }}
			{{- else -}}
				{{ range .Types }}
					{{ export $name }}{{ export .Suffix }} {{ if $isArray }}[]{{ end }}{{ .GoType }}
				{{- end }}
			{{- end -}}
		{{- end }}
	}

	{{ if .Results }}
	{{ unexport $.Name .Name }}Results struct {
		{{ range .Results }}
			{{ export .Name }} {{ if .IsArray }}[]{{ end }}{{ .GoType }}
		{{- end }}
	}
	{{- end }}
)




{{ range .Params }}
{{- if gt (len .Types) 1 }}
{{ $name := .Name }}
func (a {{ $ARGS }}) {{ export "get" $name }}() (bool, {{ range .Types }}{{ .GoType }},{{ end }}) {
	return a.has{{ export $name }}{{ range .Types }}, a.{{ $name }}{{ export .Suffix }}{{ end }}
}
{{- end }}
{{- end }}


// {{ export .Name }} function {{ .Meta.Short }}
//
// expects implementation of {{ .Name }} function:
// func (h {{ $.Name }}Handler) {{ .Name }}(ctx context.Context, args *{{ $ARGS }}) ({{ if .Results }}results *{{ $RESULTS }}, {{ end }}err error) {
//    return
// }
func (h {{ $.Name }}Handler) {{ export .Name }}() *atypes.Function {
	return &atypes.Function{
		Ref: {{ printf "%q" ( $REF ) }},
		Kind: {{ printf "%q" .Kind }},
		{{- if .Meta }}
		Meta: &atypes.FunctionMeta{
			{{- if .Meta.Short }}
			Short: {{ printf "%q" .Meta.Short }},
			{{- end }}
			{{- if .Meta.Description }}
			Description: {{ printf "%q" .Meta.Description }},
			{{- end }}
			{{- if .Meta.Visual }}
			Visual: {{ printf "%#v" .Meta.Visual }},
			{{- end }}
		},
		{{- end }}

		Parameters: []*atypes.Param{
		{{- range .Params }}
			{
				Name: {{ printf "%q" .Name }},
				Types: []string{ {{ range .Types }}{{ printf "%q" .WorkflowType }},{{ end }} },
				{{- if .Required }}Required: true,{{ end }}
				{{- if .IsArray }}IsArray: true,{{ end }}
				{{- if .Meta }}
				Meta: &atypes.ParamMeta{
					{{- if .Meta.Label }}
					Label: {{ printf "%#v" .Meta.Label }},
					{{- end }}
					{{- if .Meta.Description }}
					Description: {{ printf "%#v" .Meta.Description }},
					{{- end }}
					{{- if .Meta.Visual }}
					Visual: {{ printf "%#v" .Meta.Visual }},
					{{- end }}
				},
				{{ end }}
			},
		{{- end }}
		},


		{{ if .Results }}
		Results: []*atypes.Param{
		{{ range .Results }}
			{
				Name: {{ printf "%q" .Name }},
				Types: []string{ {{ printf "%q" .WorkflowType }} },
				{{ if .IsArray }}IsArray: true,{{ end -}}
				{{ if .Meta -}}
				Meta: &atypes.ParamMeta{
					{{- if .Meta.Label }}
					Label: {{ printf "%#v" .Meta.Label }},
					{{- end }}
					{{- if .Meta.Description }}
					Description: {{ printf "%#v" .Meta.Description }},
					{{- end }}
					{{- if .Meta.Visual }}
					Visual: {{ printf "%#v" .Meta.Visual }},
					{{- end }}
				},
				{{ end }}
			},
		{{ end }}
		},
		{{ end }}

		{{ if eq .Kind "iterator" }}
		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &{{ $ARGS }}{
				{{- range .Params }}
					has{{ export .Name }}: in.Has({{ printf "%q" .Name }}),
				{{- end }}
				}
			)

			{{ template "params" .Params }}

			return h.{{ .Name }}(ctx, args)
		},
		{{ else }}
		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &{{ $ARGS }}{
				{{- range .Params }}
					has{{ export .Name }}: in.Has({{ printf "%q" .Name }}),
				{{- end }}
				}
			)

			{{ template "params" .Params }}

			{{ if .Results }}
			var results *{{ $RESULTS }}
			if results, err = h.{{ .Name }}(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{{- range .Results }}
				{{ if .IsArray }}

				{
					// converting results.{{ export .Name }} ({{ .GoType }}) to Array (of {{ .WorkflowType }})
					var (
						tval expr.TypedValue
						tarr = make([]expr.TypedValue, len(results.{{ export .Name }}))
					)

					for i := range results.{{ export .Name }} {
						if tarr[i], err = h.reg.Type({{ printf "%q" .WorkflowType }}).Cast(results.{{ export .Name }}[i]); err != nil {
							return
						}
					}

					if tval, err = expr.NewArray(tarr); err != nil {
						return
					} else if err = expr.Assign(out, {{ printf "%q" .Name }}, tval); err != nil {
						return
					}
				}


				{{ else }}

				{
					// converting results.{{ export .Name }} ({{ .GoType }}) to {{ .WorkflowType }}
					var (
						tval expr.TypedValue
					)

					if tval, err = h.reg.Type({{ printf "%q" .WorkflowType }}).Cast(results.{{ export .Name }}); err != nil {
						return
					} else if err = expr.Assign(out, {{ printf "%q" .Name }}, tval); err != nil {
						return
					}
				}

				{{ end }}
			{{ end }}

			return
			{{- else }}
			return out, h.{{ .Name }}(ctx, args)
			{{- end }}
		},
		{{ end }}
	}
}
{{ end }}

{{ define "params" }}
	if err = in.Decode(args); err != nil {
		return
	}

	{{ range . }}
		{{ $name := .Name }}
		{{ if gt (len .Types) 1 }}
		// Converting {{ export .Name }} argument
		if args.has{{ export .Name }} {
			aux := expr.Must(expr.Select(in, {{ printf "%q" .Name }}))
			switch aux.Type() {
		{{- range .Types }}
			case h.reg.Type({{ printf "%q" .WorkflowType }}).Type():
				args.{{ $name }}{{ export .Suffix }} = aux.Get().({{ .GoType }})
		{{- end -}}
			}
		}

		{{- end }}
	{{ end }}
{{ end }}

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
{{- range .Imports }}
  {{ normalizeImport . }}
{{- end }}
)

var (
	{{ $.Name }} = &{{ $.Name }}Handler{}
)

func (h {{ $.Name }}Handler) register(reg func(*atypes.Function)) {
{{- range .Functions }}
	reg(h.{{ export .Name }}())
{{- end }}
}

{{ range .Functions }}
	{{ $REF     := unexport $.Prefix $.Name .Name }}
	{{ $ARGS    := unexport $.Name .Name "Args" }}
	{{ $RESULTS := unexport $.Name .Name "Results" }}


type (
	{{ $ARGS }} struct {
		{{ range .Params }}
			{{ $NAME := .Name }}
			has{{ export .Name }} bool
			{{ if gt (len .Types) 1 }}
				{{ export .Name }} interface{}
				{{- range .Types }}
					{{ $NAME }}{{ export .Suffix }} {{ .GoType }}
				{{- end }}
			{{- else -}}
				{{ range .Types }}
					{{ export $NAME }}{{ export .Suffix }} {{ .GoType }}
				{{- end }}
			{{- end -}}
		{{- end }}
	}

	{{ if .Results }}
	{{ unexport $.Name .Name }}Results struct {
		{{ range .Results }}
			{{ export .Name }} {{ .GoType }}
		{{- end }}
	}
	{{- end }}
)

//
//
// expects implementation of {{ .Name }} function:
// func (h {{ $.Name }}) {{ .Name }}(ctx context.Context, args *{{ $ARGS }}) (results *{{ $RESULTS }}, err error) {
//    return
// }
func (h {{ $.Name }}Handler) {{ export .Name }}() *atypes.Function {
	return &atypes.Function{
		Ref: {{ printf "%q" ( $REF ) }},
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
				Types: []string{ {{ range .Types }}({{ .WorkflowType }}{}).Type(),{{ end }} },
				{{- if .Required }}Required: true,{{ end }}
				{{- if .SetOf }}SetOf: true,{{ end }}
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
			atypes.NewParam({{ printf "%q" .Name }},
				atypes.Types(&{{ .WorkflowType }}{}),
			),
		{{ end }}
		},
		{{ end }}

		Handler: func(ctx context.Context, in expr.Vars) (out expr.Vars, err error) {
			var (
				args = &{{ $ARGS }}{
				{{- range .Params }}
					has{{ export .Name }}: in.Has({{ printf "%q" .Name }}),
				{{- end }}
				}

				{{ if .Results }}
				results *{{ $RESULTS }}
				{{ end }}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			{{ range .Params }}
				{{ $NAME := .Name }}
				{{ if gt (len .Types) 1 }}
				// Converting {{ export .Name }} to go type
				switch casted := args.{{ export .Name }}.(type) {
				{{- range .Types }}
					case {{ .GoType }}:
						args.{{ $NAME }}{{ export .Suffix }} = casted
				{{- end -}}
				}
				{{- end }}
			{{ end }}


			{{ if .Results }}
			if results, err = h.{{ .Name }}(ctx, args); err != nil {
				return
			}

			out = expr.Vars{}

			{{- range .Results }}
			if out[{{ printf "%q" .Name }}], err = ({{ .WorkflowType }}{}).Cast(results.{{ export .Name }}); err != nil {
				return nil, err
			}
			{{- end }}

			return
			{{- else }}
			return out, h.{{ .Name }}(ctx, args)
			{{- end }}
		},
	}
}
{{ end }}

package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

import (
	"encoding/json"
{{- range .Imports }}
  {{ normalizeImport . }}
{{- end }}
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

// dummy placing to simplify import generation logic
var _ = json.NewEncoder

type (
{{ range $r := $.Resources }}
// {{ camelCase $r.ResourceIdent "base" }}
//
// This type is auto-generated.
{{ camelCase $r.ResourceIdent "base" }} struct {
	immutable bool
{{- range $r.Properties }}
	{{ .Name }} {{ .Type }}
{{- end }}
}

{{ range $event := $r.Events }}
// {{ camelCase $r.ResourceIdent $event }}
//
// This type is auto-generated.
{{ camelCase $r.ResourceIdent $event }} struct {
	*{{ camelCase $r.ResourceIdent "base" }}
}
{{ end }}
{{ end }}
)


{{ range $r := $.Resources }}

// ResourceType returns "{{ $r.ResourceString }}"
//
// This function is auto-generated.
func ({{ camelCase .ResourceIdent "base" }}) ResourceType() string {
	return "{{ .ResourceString }}"
}

{{ range $event := $r.Events }}
// EventType on {{ camelCase $r.ResourceIdent $event }} returns "{{ $event }}"
//
// This function is auto-generated.
func ({{ camelCase $r.ResourceIdent $event }}) EventType() string {
	return "{{ $event }}"
}
{{ end }}

{{ range $event := $r.Events }}
// {{ camelCase "" $r.ResourceIdent $event }} creates {{ $event }} for {{ $r.ResourceString }} resource
//
// This function is auto-generated.
func {{ camelCase "" $r.ResourceIdent $event }}(
{{- range $r.Properties }}
	{{- if not .Internal }}
		{{ camelCase "arg" .Name }} {{ .Type }},
	{{- end -}}
{{- end}}
) *{{ camelCase $r.ResourceIdent $event }} {
	return &{{ camelCase $r.ResourceIdent $event }}{
		{{ camelCase $r.ResourceIdent "base" }}: &{{ camelCase $r.ResourceIdent "base" }}{
			immutable: false,
		{{- range $r.Properties }}
			{{- if not .Internal }}
				{{ .Name }}: {{ camelCase "arg" .Name }},
			{{- end -}}
		{{- end}}
		},
	}
}

// {{ camelCase "" $r.ResourceIdent $event "Immutable" }} creates {{ $event }} for {{ $r.ResourceString }} resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func {{ camelCase "" $r.ResourceIdent $event "Immutable" }}(
{{- range $r.Properties }}
	{{- if not .Internal }}
		{{ camelCase "arg" .Name }} {{ .Type }},
	{{- end -}}
{{- end}}
) *{{ camelCase $r.ResourceIdent $event }} {
	return &{{ camelCase $r.ResourceIdent $event }}{
		{{ camelCase $r.ResourceIdent "base" }}: &{{ camelCase $r.ResourceIdent "base" }}{
			immutable: true,
		{{- range $r.Properties }}
			{{- if not .Internal }}
				{{ .Name }}: {{ camelCase "arg" .Name }},
			{{- end -}}
		{{- end}}
		},
	}
}
{{ end }}



{{ range $r.Properties }}
{{ if not .Immutable }}
// {{ camelCase "Set" .Name }} sets new {{ .Name }} value
//
// This function is auto-generated.
func (res *{{ camelCase $r.ResourceIdent "base" }}) {{ camelCase "Set" .Name }}({{ camelCase "arg" .Name }} {{ .Type }}) {
	res.{{ .Name }} = {{ camelCase "arg" .Name }}
}
{{ end }}

// {{ camelCase "" .Name }} returns {{ .Name }}
//
// This function is auto-generated.
func (res {{ camelCase $r.ResourceIdent "base" }}) {{ camelCase "" .Name }}() {{ .Type }} {
	return res.{{ .Name }}
}
{{ end }}


// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res {{ camelCase .ResourceIdent "base" }}) Encode() (args map[string][]byte, err error) {
	{{- if $r.Properties }}
	args = make(map[string][]byte)

	{{ range $prop := $r.Properties }}
	if args["{{ $prop.Name }}"], err = json.Marshal(res.{{ $prop.Name }}); err != nil {
		return nil, err
	}
	{{ end }}
	{{ else }}
	// Handle argument encoding
	{{ end -}}
	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res {{ camelCase .ResourceIdent "base" }}) EncodeVars() (vars *expr.Vars, err error) {
	{{- if $r.Properties }}
	var (
		rvars = expr.RVars{}
	)

	{{ range $r.Properties }}
	{{- if .ExprType }}
	if rvars[{{ printf "%q" .Name }}], err = automation.{{ export "new" .ExprType }}(res.{{ .Name }}); err != nil {
		return nil, err
	}
	{{- else }}
	// Could not found expression-type counterpart for {{ .Type }}
	{{- end }}
	{{ end }}
	return rvars.Vars(), err
	{{ else }}
	return
	{{ end -}}
}

// Decode return values from Corredor script into struct props
func (res *{{ camelCase .ResourceIdent "base" }}) Decode(results map[string][]byte)( err error) {
	if res.immutable {
		// Respect immutability
		return
	}

	{{- if $r.Result }}
	if res.{{ $r.Result }} != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.{{ $r.Result }}); err != nil {
				return
			}
		}
	}
	{{ end -}}

	{{- range $prop := $r.Properties }}
	{{- if not $prop.Immutable }}
	if res.{{ $prop.Name }} != nil {
		if r, ok := results["{{ $prop.Name }}"]; ok {
			if err = json.Unmarshal(r, res.{{ $prop.Name }}); err != nil {
				return
			}
		}
	}
	{{ else }}
	// Do not decode {{ $prop.Name }}; marked as immutable
	{{ end -}}
	{{ end -}}

	return
}

func (res *{{ camelCase .ResourceIdent "base" }}) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}

	{{- range $r.Properties }}
	{{- if .Immutable }}
	// {{ .Name }} marked as immutable
	{{- else }}
	{{- if .ExprType }}
	if res.{{ .Name }} != nil && vars.Has({{ printf "%q" .Name }}) {
		var aux *automation.{{ export .ExprType }}
		aux, err = automation.{{ export "new" .ExprType }}(expr.Must(vars.Select({{ printf "%q" .Name }})))
		if err != nil {
			return
		}

		res.{{ .Name }} = aux.GetValue()
	}
	{{- else }}
	// Could not find expression-type counterpart for {{ .Type }}
	{{- end }}
	{{- end }}
	{{- end }}

	return
}

{{ end }}





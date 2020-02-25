package {{ .Package }}

// This file is auto-generated.
//
// YAML event definitions:
//   {{ .YAML }}
//
// Regenerate with:
//   {{ .Command }}
//

{{ if or .Imports $.Events.Properties }}
import (
{{ if $.Events.Properties }}
	"encoding/json"
{{ end }}
{{ range $i, $import := .Imports }}
  "{{ $import }}"
{{ end }}
)
{{ end }}


type (
// {{ camelCase $.ResourceIdent "base" }}
//
// This type is auto-generated.
{{ camelCase $.ResourceIdent "base" }} struct {
	immutable bool
{{- range $p := $.Events.Properties }}
	{{ $p.Name }} {{ $p.Type }}
{{- end }}
}


{{ range $event := makeEvents $.Events }}
// {{ camelCase $.ResourceIdent $event }}
//
// This type is auto-generated.
{{ camelCase $.ResourceIdent $event }} struct {
	*{{ camelCase $.ResourceIdent "base" }}
}
{{ end }}
)

// ResourceType returns "{{ .ResourceString }}"
//
// This function is auto-generated.
func ({{ camelCase .ResourceIdent "base" }}) ResourceType() string {
	return "{{ .ResourceString }}"
}

{{ range $event := makeEvents $.Events }}
// EventType on {{ camelCase $.ResourceIdent $event }} returns "{{ $event }}"
//
// This function is auto-generated.
func ({{ camelCase $.ResourceIdent $event }}) EventType() string {
	return "{{ $event }}"
}
{{ end }}


{{ range $event := makeEvents $.Events }}
// {{ camelCase "" $.ResourceIdent $event }} creates {{ $event }} for {{ $.ResourceString }} resource
//
// This function is auto-generated.
func {{ camelCase "" $.ResourceIdent $event }}(
{{- range $p := $.Events.Properties }}
	{{- if not $p.Internal }}
		{{ camelCase "arg" $p.Name }} {{ $p.Type }},
	{{- end -}}
{{- end}}
) *{{ camelCase $.ResourceIdent $event }} {
	return &{{ camelCase $.ResourceIdent $event }}{
		{{ camelCase $.ResourceIdent "base" }}: &{{ camelCase $.ResourceIdent "base" }}{
			immutable: false,
		{{- range $p := $.Events.Properties }}
			{{- if not $p.Internal }}
				{{ $p.Name }}: {{ camelCase "arg" $p.Name }},
			{{- end -}}
		{{- end}}
		},
	}
}

// {{ camelCase "" $.ResourceIdent $event "Immutable" }} creates {{ $event }} for {{ $.ResourceString }} resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func {{ camelCase "" $.ResourceIdent $event "Immutable" }}(
{{- range $p := $.Events.Properties }}
	{{- if not $p.Internal }}
		{{ camelCase "arg" $p.Name }} {{ $p.Type }},
	{{- end -}}
{{- end}}
) *{{ camelCase $.ResourceIdent $event }} {
	return &{{ camelCase $.ResourceIdent $event }}{
		{{ camelCase $.ResourceIdent "base" }}: &{{ camelCase $.ResourceIdent "base" }}{
			immutable: true,
		{{- range $p := $.Events.Properties }}
			{{- if not $p.Internal }}
				{{ $p.Name }}: {{ camelCase "arg" $p.Name }},
			{{- end -}}
		{{- end}}
		},
	}
}
{{ end }}

{{ range $p := $.Events.Properties }}
{{ if not $p.Immutable }}
// {{ camelCase "Set" $p.Name }} sets new {{ $p.Name }} value
//
// This function is auto-generated.
func (res *{{ camelCase $.ResourceIdent "base" }}) {{ camelCase "Set" $p.Name }}({{ camelCase "arg" $p.Name }} {{ $p.Type }}) {
	res.{{ $p.Name }} = {{ camelCase "arg" $p.Name }}
}
{{ end }}

// {{ camelCase "" $p.Name }} returns {{ $p.Name }}
//
// This function is auto-generated.
func (res {{ camelCase $.ResourceIdent "base" }}) {{ camelCase "" $p.Name }}() {{ $p.Type }} {
	return res.{{ $p.Name }}
}
{{ end }}


// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res {{ camelCase .ResourceIdent "base" }}) Encode() (args map[string][]byte, err error) {
	{{- if $.Events.Properties }}
	args = make(map[string][]byte)

	{{ range $prop := $.Events.Properties }}
	if args["{{ $prop.Name }}"], err = json.Marshal(res.{{ $prop.Name }}); err != nil {
		return nil, err
	}
	{{ end }}
	{{ else }}
	// Handle argument encoding
	{{ end -}}
	return
}

// Decode return values from Corredor script into struct props
func (res *{{ camelCase .ResourceIdent "base" }}) Decode(results map[string][]byte)( err error) {
	if res.immutable {
		// Respect immutability
		return
	}

	{{- if $.Events.Result }}
	if res.{{ $.Events.Result }} != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.{{ $.Events.Result }}); err != nil {
				return
			}
		}
	}
	{{ end -}}

	{{- range $prop := $.Events.Properties }}
	{{- if not $prop.Immutable }}
	if res.{{ $prop.Name }} != nil {
		if r, ok := results["{{ $prop.Name }}"]; ok {
			if err = json.Unmarshal(r, res.{{ $prop.Name }}); err != nil {
				return
			}
		}
	}
	{{ end -}}
	{{ end -}}

	return
}

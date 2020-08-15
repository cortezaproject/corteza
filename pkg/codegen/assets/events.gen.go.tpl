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
{{- range $i, $import := .Imports }}
  {{ normalizeImport $import }}
{{- end }}
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
{{- range $p := $r.Properties }}
	{{ $p.Name }} {{ $p.Type }}
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
{{- range $p := $r.Properties }}
	{{- if not $p.Internal }}
		{{ camelCase "arg" $p.Name }} {{ $p.Type }},
	{{- end -}}
{{- end}}
) *{{ camelCase $r.ResourceIdent $event }} {
	return &{{ camelCase $r.ResourceIdent $event }}{
		{{ camelCase $r.ResourceIdent "base" }}: &{{ camelCase $r.ResourceIdent "base" }}{
			immutable: false,
		{{- range $p := $r.Properties }}
			{{- if not $p.Internal }}
				{{ $p.Name }}: {{ camelCase "arg" $p.Name }},
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
{{- range $p := $r.Properties }}
	{{- if not $p.Internal }}
		{{ camelCase "arg" $p.Name }} {{ $p.Type }},
	{{- end -}}
{{- end}}
) *{{ camelCase $r.ResourceIdent $event }} {
	return &{{ camelCase $r.ResourceIdent $event }}{
		{{ camelCase $r.ResourceIdent "base" }}: &{{ camelCase $r.ResourceIdent "base" }}{
			immutable: true,
		{{- range $p := $r.Properties }}
			{{- if not $p.Internal }}
				{{ $p.Name }}: {{ camelCase "arg" $p.Name }},
			{{- end -}}
		{{- end}}
		},
	}
}
{{ end }}



{{ range $p := $r.Properties }}
{{ if not $p.Immutable }}
// {{ camelCase "Set" $p.Name }} sets new {{ $p.Name }} value
//
// This function is auto-generated.
func (res *{{ camelCase $r.ResourceIdent "base" }}) {{ camelCase "Set" $p.Name }}({{ camelCase "arg" $p.Name }} {{ $p.Type }}) {
	res.{{ $p.Name }} = {{ camelCase "arg" $p.Name }}
}
{{ end }}

// {{ camelCase "" $p.Name }} returns {{ $p.Name }}
//
// This function is auto-generated.
func (res {{ camelCase $r.ResourceIdent "base" }}) {{ camelCase "" $p.Name }}() {{ $p.Type }} {
	return res.{{ $p.Name }}
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


{{ end }}

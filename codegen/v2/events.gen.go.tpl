package {{ .Package }}

// This file is auto-generated.
//
// YAML event definitions:
//   {{ .YAML }}
//
// Regenerate with:
//   {{ .Command }}
//

{{ if .Imports }}
import (
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

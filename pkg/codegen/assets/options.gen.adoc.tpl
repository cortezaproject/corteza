// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}


{{ range .Definitions }}

= {{ .Docs.Title }}

{{ .Docs.Intro }}

{{ range .Properties }}

== *{{ toUpper .Env }}* `{{ .Type }}`

{{ if .Default }}
Default::
	`{{ .Default }}`
{{ end -}}
{{ if .Description }}
Description::
	{{ .Description }}
{{ end -}}

{{ end }}
{{ end }}


// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}

= ENV options
{{ range .Definitions }}
== {{ .Docs.Title }}
{{ if .Docs.Intro }}
{{ .Docs.Intro }}
{{ end }}
{{- range .Properties }}
=== *{{ toUpper .Env }}* `{{ .Type }}`

{{ if or .Default .Description -}}
{{ if .Default -}}
Default::
    `{{ .Default }}`
{{ end -}}
{{ if .Description -}}
Description::
    {{ .Description }}
{{ end -}}{{ end -}}
{{ end }}{{ end }}
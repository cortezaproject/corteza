// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}
include::ROOT:partial$variables.adoc[]

{{ range .Definitions }}
= {{ .Docs.Title }}
{{ if .Docs.Intro }}
{{ .Docs.Intro }}
{{ end }}
{{- range .Properties }}
== *{{ toUpper .Env }}*

=== Type

`{{ .Type }}`

{{ if or .Default .Description -}}
{{ if .Default -}}
=== Default

[source]
----
{{ .Default }}
----

{{ end -}}
{{ if .Description -}}
=== Description

{{ .Description }}
{{ end -}}{{ end -}}
{{ end }}{{ end }}
//
// This file is auto-generated.
//

include::ROOT:partial$variables.adoc[]

{{- range .groups }}
= {{ .title }}

{{- if .intro }}
{{ .intro }}
#
{{ end -}}

{{ range .options }}

== *{{ .env }}*

=== Type

`{{ .type }}`

{{ if .defaultValue }}
=== Default

[source]
----
{{ .defaultValue }}
----
{{- end }}

{{- if .description }}
=== Description

{{ .description }}
{{- end }}

{{- end }}
{{- end }}

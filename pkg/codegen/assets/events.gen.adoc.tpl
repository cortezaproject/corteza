// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}

{{- range .Definitions }}
{{- range .Resources }}

= {{ .ResourceString }}

== Events

{{- if .BeforeAfter }}
.Before/after events:
{{- range $ba := .BeforeAfter }}
* `before('{{ $ba }}')`
{{- end }}
{{- range $ba := .BeforeAfter }}
* `after('{{ $ba }}')`
{{- end }}
{{- end }}

{{ if .On -}}
.On events:
{{- range $on := .On }}
* `on('{{ $on }}')`
{{- end }}
{{- end }}

== Exec arguments

.Argument properties:
[%header, cols=3*]
|===
|Name|Type|Mutable
{{- range $p := .Properties }}
| `{{ $p.Name }}`
| `{{ $p.Type }}`
{{- if $p.Immutable }}
| no
{{ else }}
| yes
{{ end -}}

{{ end -}}
|===

{{- end }}
{{- end }}

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

.Events:
{{- range $ba := .BeforeAfter }}
* `before('{{ $ba }}')`
{{- end }}
{{- range $ba := .BeforeAfter }}
* `after('{{ $ba }}')`
{{- end }}
{{- range $on := .On }}
* `on('{{ $on }}')`
{{- end }}

== Argument properties

.Argument properties:
[%header, cols=3*]
|===
|Name|Type|Immutable
{{- range $p := .Properties }}
| `{{ $p.Name }}`
| `{{ $p.Type }}`
{{- if $p.Immutable }}
| yes
{{ else }}
| no
{{ end -}}

{{ end -}}
|===

{{- end }}
{{- end }}

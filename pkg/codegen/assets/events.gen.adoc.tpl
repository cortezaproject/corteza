// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

= {{ $.ResourceString }}


.List of events on `{{ $.ResourceString }}`
{{- range $event := $.Events.MakeEvents }}
- `{{ $event }}`
{{- end }}

.Event arguments for `{{ $.ResourceString }}`
[%header,cols=3*]
|===
|Name
|Type
|Immutable

{{- range $p := $.Events.Properties }}
|`{{ camelCase $p.Name }}`
|`{{ $p.Type }}`
|{{ $p.Immutable }}
{{- end }}
|===

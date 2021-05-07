// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}

= Types

[%header, cols="1m,2a,2m"]
|===
| Type | DataType | Function
{{- range .Definitions }}
{{- range $tName, $tDef := .Types }}
{{- if gt (len $tDef.Struct) 0}}
| *{{ $tName }}*
| `{{ $tDef.As }}`

{{ if $tDef.Struct }}
[source, go]
----
{{ $tDef.As }}{
{{- range $s := .Struct }}
   {{ $s.Name }} {{ $s.ExprType }}
{{- end }}
}
----
{{- end }}
| {{ $tDef.AssignerFn }}

{{- end }}

{{- end }}
{{- end }}
|===

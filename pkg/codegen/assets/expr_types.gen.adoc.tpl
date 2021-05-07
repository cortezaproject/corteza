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

[%header, cols="1,2a"]
|===
| Type | DataType
{{- range .Definitions }}
{{- range $tName, $tDef := .Types }}
| *{{ $tName }}*
| `{{ $tDef.As }}`

{{ if $tDef.Struct }}
[source, go]
----
{{ $tDef.As }}{
{{- range $s := .Struct }}
   {{ $s.Name }} {{ $s.GoType }}
{{- end }}
}
----
{{- end }}

{{- end }}
{{- end }}
|===

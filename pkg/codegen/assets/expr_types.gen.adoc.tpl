// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}

[cols="2m,3a"]
|===
| Type | Structure
{{- range .Definitions }}
{{- range $tName, $tDef := .Types }}
   {{- if gt (len $tDef.Struct) 0}}
| [#objref-{{ toLower $tName }}]#<<objref-{{ toLower $tName }},{{ $tName }}>>#
|
      {{- if $tDef.Struct }}
[source]
----
{
         {{- range $s := .Struct }}
   {{ $s.Name }} ({{ $s.ExprType }})
         {{- end }}
}
----
{{ end }}
   {{- end }}
{{- end }}
{{ end }}
|===

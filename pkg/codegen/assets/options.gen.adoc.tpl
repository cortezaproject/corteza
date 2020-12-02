// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

= {{ export $.Name }}

[cols="2,3,5a"]
|===
|Type|Default value|Description

{{- range $prop := $.Properties }}
3+| *{{ toUpper $prop.Env }}*
|`{{ $prop.Type }}`
|{{- if $prop.Default }}
   {{- $prop.Default -}}
 {{- end -}}
|{{ $prop.Description }}
{{- end }}
|===

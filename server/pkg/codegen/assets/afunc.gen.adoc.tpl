// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}

{{ range $d := .Definitions }}
= `{{ $d.Name }}`

[cols="2m,4a,3a"]
|===
| Name | Description | I/O

{{- range $f := .Functions }}

{{- if eq $f.Kind "function" }}

| [#fnc-{{ toLower $d.Name }}-{{ toLower $f.Name }}]#<<fnc-{{ toLower $d.Name }}-{{ toLower $f.Name }},{{ if $f.Meta.Short }}{{ $f.Meta.Short }}{{ else }}{{ $f.Name }}{{ end }}>>#
| {{ if $f.Meta.Description }}{{ $f.Meta.Description }}{{ end }}
| 
{{- if gt (len $f.Params) 0}}
.Parameters:
{{- range $p := $f.Params}}
* {{ if $p.Required }}#*# {{ end }}`{{ if $p.Meta }}
   {{- if $p.Meta.Label }}
      {{- $p.Meta.Label }}
   {{- else }}
      {{- $p.Name }}
   {{- end }}
{{- else }}
   {{- $p.Name }}
{{- end }}`
({{- range $pti, $pt := $p.Types }}
   `{{- $pt.WorkflowType }}`,
{{- end }})
{{- end }}
{{- end }}

{{- if gt (len $f.Results) 0}}

.Results:
{{- range $r := $f.Results}}
* {{ if $r.Meta }}
   {{- if $r.Meta.Label }}
      {{- $r.Meta.Label }}
   {{- else }}
      {{- $r.Name }}
   {{- end }}
{{- else }}
   {{- $r.Name }}
{{- end }} (`{{ $r.WorkflowType }}`)
{{- end }}
{{- end }}

{{- end }}
{{- end }}

|===
{{- end }}

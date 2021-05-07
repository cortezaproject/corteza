// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}

= Automatic functions

{{- range .Definitions }}
{{- range $f := .Functions }}

{{- if eq $f.Kind "function" }}

.*{{ $f.Name }}* {{ if $f.Meta.Description }} {{ $f.Meta.Description }} {{ else }} {{ $f.Meta.Short }} {{ end }}
{{ if gt (len $f.Params) 0}}
* It takes following parameters:
+
[cols="1m,3a,4s"]
|===
{{- range $p := $f.Params}}
|
{{ $p.Name }}
{{ if $p.Required }} `Required` {{ end }}
|
{{ if eq (len $p.Types) 1}}
   {{- (index $p.Types 0).WorkflowType }}
{{else if gt (len $p.Types) 1}}
[source, go]
----
struct {
{{- range $pti, $pt := $p.Types }}
   field{{ add $pti 1 }} {{ $pt.WorkflowType }}
{{- end }}
}
----
{{ end }}
|
{{ $pMeta := $p.Meta }}
{{- if $pMeta }}
{{ $mDescription := printf "%s" $pMeta.Description }}
{{ $mLabel := printf "%s" $pMeta.Label }}
{{- if $mDescription}} {{- $mDescription }} {{- else if $mLabel }} {{- $mLabel }} {{- end }}
{{- end }}

{{ end }}
|===
{{ end }}

{{ if gt (len $f.Results) 0}}
* It returns as follow:
+
[cols="1m,2m,4s"]
|===
{{- range $r := $f.Results}}
| {{ $r.Name }}
| {{ $r.WorkflowType }}
|
{{ $rMeta := $r.Meta }}
{{- if $rMeta }}
{{ $rmDescription := printf "%s" $rMeta.Description }}
{{ $rmLabel := printf "%s" $rMeta.Label }}
{{- if $rmDescription}} {{- $rmDescription }} {{- else if $rmLabel }} {{- $rmLabel }} {{- end }}
{{- end }}
{{ end }}
|===
{{ end }}

{{- end }}
{{- end }}
{{- end }}

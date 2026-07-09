{{- define "fieldVal" -}}
{{- $base := .field.goType | trimPrefix "*" | trimPrefix "[]" -}}
{{- if hasKey .defs $base -}}{{- template "obj" (dict "entries" (index .defs $base).fields "defs" .defs) -}}{{- else -}}"{{$base}}"{{- end -}}
{{- end -}}

{{- define "obj" -}}
{
{{- $first := true -}}
{{- $defs := .defs -}}
{{- range .entries -}}
{{- if not $first }},{{ end -}}
{{- $first = false }}
  "{{.jsonTag | trimPrefix "json:\"" | trimSuffix "\"" | splitList "," | first}}": {{template "fieldVal" (dict "field" . "defs" $defs)}}
{{- end }}
}{{- end -}}

[
{{- $res := list -}}
{{- range .components -}}{{- range .resources -}}{{- $res = append $res . -}}{{- end -}}{{- end -}}
{{- $first := true -}}
{{- range $res }}
{{- if not $first }},{{ end -}}
{{- $first = false -}}
{{- $defs := dict -}}
{{- range .structDefs -}}{{- $defs = set $defs .name . -}}{{- end }}
  { "{{.handle}}": {{template "obj" (dict "entries" .fields "defs" $defs)}} }
{{- end }}
]

{{- define "primaryKeyArgs" -}}
    {{- range $field := . -}}
        {{- if $field.IsPrimaryKey -}}
           , {{ $field.Arg }} {{ camelCase $field.Type }}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{- define "primaryKeySuffix" -}}
    {{- range $field := . }}{{ if $field.IsPrimaryKey }}{{ $field.Field }}{{ end }}{{ end -}}
{{- end -}}

{{- define "partialUpdateArgs" -}}
    {{- range .Args -}}
	   , {{ .Arg }} {{ .Type }}
    {{- end -}}
{{- end -}}

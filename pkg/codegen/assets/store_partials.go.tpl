{{- define "primaryKeyArgsDef" -}}
    {{- range $field := . -}}
        {{- if $field.IsPrimaryKey -}}
           , {{ $field.Arg }} {{ camelCase $field.Type }}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{- define "primaryKeyArgsCall" -}}
    {{- range $field := . -}}
        {{- if $field.IsPrimaryKey -}}
           , {{ $field.Arg }}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{- define "primaryKeySuffix" -}}
    {{- range $field := . }}{{ if $field.IsPrimaryKey }}{{ $field.Field }}{{ end }}{{ end -}}
{{- end -}}

{{ define "extraArgsDefFirst" }}
    {{- range .Arguments -}}
	   _{{ .Name }} {{ .Type }},
    {{- end -}}
{{ end }}

{{ define "extraArgsDef" }}
    {{- range .Arguments -}}
	   , _{{ .Name }}  {{ .Type }}
    {{- end -}}
{{ end }}

{{ define "extraArgsDefTypesOnly" }}
    {{- range .Arguments -}}
	   , {{ .Type }}
    {{- end -}}
{{ end }}

{{ define "extraArgsCallFirst" }}
    {{- range .Arguments -}}
	   _{{ .Name }},
    {{- end -}}
{{ end }}

{{ define "extraArgsCall" }}
    {{- range .Arguments -}}
	   , _{{ .Name }}{{ if (hasPrefix "..." .Type) }}...{{ end }}
    {{- end -}}
{{ end }}

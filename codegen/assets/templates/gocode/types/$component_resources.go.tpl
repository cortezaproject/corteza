package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

const (
{{- range .types }}
	{{ .const }} = {{ printf "%q" .type }}
{{- end }}
)

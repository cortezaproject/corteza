package {{ .Package }}

{{ template "header-gentext.tpl" }}
{{ template "header-definitions.tpl" . }}

import (
	"fmt"
	"strconv"
	"github.com/cortezaproject/corteza-server/pkg/locale"
)

type (
	LocaleKey struct {
		Name          string
		Resource      string
		Path          string
		CustomHandler string
	}
)

// Types and stuff
const (
{{- range .Def }}
	{{ .Resource }}ResourceTranslationType = "{{ .Locale.ResourceType }}"
{{- end }}
)

var (
{{- range .Def }}
{{- $Resource := .Resource }}
	{{- range .Locale.Keys}}
	LocaleKey{{ $Resource }}{{coalesce (export .Name) (export .Path) }} = LocaleKey{
		Name: "{{.Name}}",
		Resource: {{ $Resource }}ResourceTranslationType,
		Path: "{{.Path }}",{{ if .CustomHandler }}
		CustomHandler: "{{ .CustomHandler }}",
		{{- end }}
	}
	{{- end}}
{{- end }}
)

{{- range .Def }}
{{ $Resource := .Resource }}
{{ $GoType   := printf "types.%s" .Resource }}


// ResourceTranslation returns string representation of Locale resource for {{ .Resource }} by calling {{ .Resource }}ResourceTranslation fn
//
// Locale resource is in "{{ .Locale.ResourceType }}/..." format
//
// This function is auto-generated
func (r {{ .Resource }}) ResourceTranslation() string {
	return {{ .Resource }}ResourceTranslation({{ if .Locale.Resource }}{{ range .Locale.Resource.References }}r.{{ export .Field }},{{ end }}{{ end }})
}

// {{ .Resource }}ResourceTranslation returns string representation of Locale resource for {{ .Resource }}
//
// Locale resource is in the {{ .Locale.ResourceType }}/{{- if .Locale.Resource.References }}...{{ end }} format
//
// This function is auto-generated
func {{ .Resource }}ResourceTranslation({{ if .Locale.Resource }}{{ range .Locale.Resource.References }}{{ unexport .Field }} uint64,{{ end }}{{ end }}) string {
	{{- if .Locale.Resource.References }}
	cpts := []interface{{"{}"}}{{"{"}}{{ .Resource }}ResourceTranslationType{{"}"}}
	cpts = append(cpts, {{range .Locale.Resource.References -}}
		strconv.FormatUint({{ unexport .Field }}, 10),
		{{- end }})

	return fmt.Sprintf({{ .Resource }}ResourceTranslationTpl(), cpts...)
	{{- end }}
}

// @todo template
func {{ .Resource }}ResourceTranslationTpl() string {
	{{- if .Locale.Resource.References }}
	return "%s
	{{- range .Locale.Resource.References }}/%s{{- end }}"

	{{- else }}
	return "%s"
	{{- end }}
}

func (r *{{ .Resource }}) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

{{- range .Locale.Keys}}
{{- if not .Custom }}
	if aux = tt.FindByKey(LocaleKey{{ $Resource }}{{coalesce (export .Name) (export .Path) }}.Path); aux != nil {
		r.{{ .Field }} = aux.Msg
	}
{{- end}}
{{- end}}

{{- range .Locale.Keys}}
{{- if and .Custom .CustomHandler }}
	r.decodeTranslations{{export .CustomHandler }}(tt)
{{- end}}
{{- end}}

{{- if .Locale.Extended }}

	r.decodeTranslations(tt)
{{- end }}
}

func (r *{{ .Resource }}) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}
{{- range .Locale.Keys}}
{{- if not .Custom }}
	if r.{{ .Field }} != "" {
		out = append(out, &locale.ResourceTranslation{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKey{{ $Resource }}{{coalesce (export .Name) (export .Path) }}.Path,
			Msg:      locale.SanitizeMessage(r.{{ .Field }}),
		})
	}

{{- end}}
{{- end}}

{{range .Locale.Keys}}
{{- if and .Custom .CustomHandler }}
	out = append(out, r.encodeTranslations{{export .CustomHandler}}()...)
{{- end}}
{{- end}}

{{- if .Locale.Extended }}
	out = append(out, r.encodeTranslations()...)
{{- end }}

	return out
}

{{- end }}

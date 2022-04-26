package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

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
{{- range .resources }}
	{{ .const }} = "{{ .type }}"
{{- end }}
)

var (
 // @todo can we remove LocaleKey struct for string constant?
{{- range .resources }}
{{- range .keys }}
	{{ .struct }} = LocaleKey{ Path: {{ printf "%q" .path }} }
{{- end }}
{{- end }}
)

{{- range .resources }}

// ResourceTranslation returns string representation of Locale resource for {{ .expIdent }} by calling {{ .expIdent }}ResourceTranslation fn
//
// Locale resource is in "{{ .type }}/..." format
//
// This function is auto-generated
func (r {{ .expIdent }}) ResourceTranslation() string {
	return {{ .expIdent }}ResourceTranslation({{ range .references }}r.{{ .refField }},{{ end }})
}

// {{ .expIdent }}ResourceTranslation returns string representation of Locale resource for {{ .expIdent }}
//
// Locale resource is in the {{ .type }}/{{- if .references }}...{{ end }} format
//
// This function is auto-generated
func {{ .expIdent }}ResourceTranslation({{ range .references }}{{ .refField }} uint64,{{ end }}) string {
	{{- if .references }}
	cpts := []interface{{"{}"}}{
		{{ .expIdent }}ResourceTranslationType,
	  {{- range .references }}
		strconv.FormatUint({{ .refField }}, 10),
		{{- end }}
	}

	return fmt.Sprintf({{ .expIdent }}ResourceTranslationTpl(), cpts...)
	{{- end }}
}

func {{ .expIdent }}ResourceTranslationTpl() string {
	{{- if .references }}
	return "%s
	{{- range .references }}/%s{{- end }}"

	{{- else }}
	return "%s"
	{{- end }}
}

func (r *{{ .expIdent }}) DecodeTranslations(tt locale.ResourceTranslationIndex) {
    {{- $decodeFuncExist := false }}
    {{- range .keys }}
        {{- if not .decodeFunc }}
            {{- $decodeFuncExist = true }}
        {{- end }}
    {{- end}}

    {{- if $decodeFuncExist }}
        var aux *locale.ResourceTranslation
    {{- end }}

	{{- range .keys }}
		{{ if .decodeFunc }}
			{{ if not .extended }}
			r.{{ .decodeFunc }}(tt)
			{{- end}}
		{{ else }}
		if aux = tt.FindByKey({{ .struct }}.Path); aux != nil {
			r.{{ .fieldPath }} = aux.Msg
		}
		{{- end}}
	{{- end}}

	{{- if .extended }}
		r.decodeTranslations(tt)
	{{- end }}
}

func (r *{{ .expIdent }}) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}
	{{- range .keys }}
		{{ if .encodeFunc }}
			{{ if not .extended }}
			out = append(out, r.{{ .encodeFunc }}()...)
			{{- end}}
		{{ else }}
		out = append(out, &locale.ResourceTranslation{
			Resource: r.ResourceTranslation(),
			Key:      {{ .struct }}.Path,
			Msg:      locale.SanitizeMessage(r.{{ .fieldPath }}),
		})
		{{- end}}
	{{- end}}

	{{- if .extended }}
		out = append(out, r.encodeTranslations()...)
	{{- end }}

	return out
}

{{- end }}

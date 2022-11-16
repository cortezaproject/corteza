package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"fmt"
	"strconv"
	"github.com/cortezaproject/corteza/server/pkg/locale"
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
	return {{ .expIdent }}ResourceTranslation({{ if .references }}{{ range .references }}r.{{ . }},{{ end }}{{ end }})
}

// {{ .expIdent }}ResourceTranslation returns string representation of Locale resource for {{ .expIdent }}
//
// Locale resource is in the {{ .type }}/{{- if .references }}...{{ end }} format
//
// This function is auto-generated
func {{ .expIdent }}ResourceTranslation({{ if .references }}{{ range .references }}{{ . }} uint64,{{ end }}{{ end }}) string {
	{{- if .references }}
	cpts := []interface{{"{}"}}{
		{{ .expIdent }}ResourceTranslationType,
	  {{- range .references }}
		strconv.FormatUint({{ . }}, 10),
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
	var aux *locale.ResourceTranslation

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
		if r.{{ .fieldPath }} != "" {
			out = append(out, &locale.ResourceTranslation{
				Resource: r.ResourceTranslation(),
				Key:      {{ .struct }}.Path,
				Msg:      locale.SanitizeMessage(r.{{ .fieldPath }}),
			})
		}
		{{- end}}
	{{- end}}

	{{- if .extended }}
		out = append(out, r.encodeTranslations()...)
	{{- end }}

	return out
}

{{- end }}

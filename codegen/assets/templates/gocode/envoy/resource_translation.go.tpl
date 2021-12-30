package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

{{- range .resources }}
func (r *{{ .expIdent }}) EncodeTranslations() ([]*ResourceTranslation, error) {
	out := make([]*ResourceTranslation, 0, 10)

	rr := r.Res.EncodeTranslations()
	rr.SetLanguage(defaultLanguage)
	res, ref, pp := r.ResourceTranslationParts()
	out = append(out, NewResourceTranslation(systemTypes.FromLocale(rr), res, ref, pp...))
{{ if .extended }}
	tmp, err := r.encodeTranslations()
	return append(out, tmp...), err
{{ else }}
	return out, nil
{{- end }}
}
{{- end }}


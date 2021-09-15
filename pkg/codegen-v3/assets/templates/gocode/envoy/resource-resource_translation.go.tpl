package {{ .Package }}

{{ template "header-gentext.tpl" }}
{{ template "header-definitions.tpl" . }}

import (
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
{{- range .Imports }}
    {{ . }}
{{- end }}
)

{{- range .Def }}
{{ $Component := .Component }}
{{ $Resource := .Resource }}
{{ $GoType   := printf "types.%s" .Resource }}
func (r *{{export $Component}}{{$Resource}}) EncodeTranslations() ([]*ResourceTranslation, error) {
	out := make([]*ResourceTranslation, 0, 10)

	rr := r.Res.EncodeTranslations()
	rr.SetLanguage(defaultLanguage)
	res, ref, pp := r.ResourceTranslationParts()
	out = append(out, NewResourceTranslation(systemTypes.FromLocale(rr), res, ref, pp...))
{{ if .Locale.Extended }}
	tmp, err := r.encodeTranslations()
	return append(out, tmp...), err
{{ else }}
	return out, nil
{{- end }}
}
{{- end }}

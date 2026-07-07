package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}
{{- if .types }}
import (
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)
{{- range .types }}

func (m *{{ .expIdent }}) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m {{ .expIdent }}) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func ({{ .expIdent }}) LabelResourceKind() string {
	return {{ printf "%q" .labelResourceType }}
}

func (m {{ .expIdent }}) LabelResourceID() uint64 {
	return m.ID
}
{{- end }}
{{- end }}

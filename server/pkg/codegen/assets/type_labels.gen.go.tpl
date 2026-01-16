package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}
 import (
  	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
  )


{{ range $name, $set := .Types }}
{{ if $set.LabelResourceType }}
// SetLabel adds new label to label map
func (m *{{ $name }}) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m {{ $name }}) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

// GetLabels adds new label to label map
func ({{ $name }}) LabelResourceKind() string {
	return {{ printf "%q" $set.LabelResourceType }}
}

// GetLabels adds new label to label map
func (m {{ $name }}) LabelResourceID() uint64 {
	return m.ID
}
{{ end }}
{{ end }}

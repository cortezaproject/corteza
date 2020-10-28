package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}


{{ range $name, $set := .Types }}
{{ if $set.LabelResourceType }}
// SetLabel adds new label to label map
func (m *{{ $name }}) SetLabel(key string, value string) {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m {{ $name }}) GetLabels() map[string]string {
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

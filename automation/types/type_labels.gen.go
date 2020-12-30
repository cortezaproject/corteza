package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/types/types.yaml

// SetLabel adds new label to label map
func (m *Trigger) SetLabel(key string, value string) {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Trigger) GetLabels() map[string]string {
	return m.Labels
}

// GetLabels adds new label to label map
func (Trigger) LabelResourceKind() string {
	return "trigger"
}

// GetLabels adds new label to label map
func (m Trigger) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *Workflow) SetLabel(key string, value string) {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Workflow) GetLabels() map[string]string {
	return m.Labels
}

// GetLabels adds new label to label map
func (Workflow) LabelResourceKind() string {
	return "workflow"
}

// GetLabels adds new label to label map
func (m Workflow) LabelResourceID() uint64 {
	return m.ID
}

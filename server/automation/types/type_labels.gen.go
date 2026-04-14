package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/types/types.yaml
import (
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
)

// SetLabel adds new label to label map
func (m *NgAutomation) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m NgAutomation) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

// GetLabels adds new label to label map
func (NgAutomation) LabelResourceKind() string {
	return "executable"
}

// GetLabels adds new label to label map
func (m NgAutomation) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *Trigger) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Trigger) GetLabels() map[string]labelTypes.LabelValue {
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
func (m *TriggerDefinition) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m TriggerDefinition) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

// GetLabels adds new label to label map
func (TriggerDefinition) LabelResourceKind() string {
	return "trigger-definition"
}

// GetLabels adds new label to label map
func (m TriggerDefinition) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *Workflow) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Workflow) GetLabels() map[string]labelTypes.LabelValue {
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

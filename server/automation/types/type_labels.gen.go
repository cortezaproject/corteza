package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

func (m *Workflow) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Workflow) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Workflow) LabelResourceKind() string {
	return "workflow"
}

func (m Workflow) LabelResourceID() uint64 {
	return m.ID
}

func (m *Trigger) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Trigger) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Trigger) LabelResourceKind() string {
	return "trigger"
}

func (m Trigger) LabelResourceID() uint64 {
	return m.ID
}

func (m *NgAutomation) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m NgAutomation) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (NgAutomation) LabelResourceKind() string {
	return "executable"
}

func (m NgAutomation) LabelResourceID() uint64 {
	return m.ID
}

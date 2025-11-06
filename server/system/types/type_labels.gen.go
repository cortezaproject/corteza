package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/types/types.yaml
import (
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
)

// SetLabel adds new label to label map
func (m *Application) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Application) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

// GetLabels adds new label to label map
func (Application) LabelResourceKind() string {
	return "application"
}

// GetLabels adds new label to label map
func (m Application) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *AuthClient) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m AuthClient) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

// GetLabels adds new label to label map
func (AuthClient) LabelResourceKind() string {
	return "authClient"
}

// GetLabels adds new label to label map
func (m AuthClient) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *Report) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Report) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

// GetLabels adds new label to label map
func (Report) LabelResourceKind() string {
	return "report"
}

// GetLabels adds new label to label map
func (m Report) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *Role) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Role) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

// GetLabels adds new label to label map
func (Role) LabelResourceKind() string {
	return "role"
}

// GetLabels adds new label to label map
func (m Role) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *Template) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Template) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

// GetLabels adds new label to label map
func (Template) LabelResourceKind() string {
	return "template"
}

// GetLabels adds new label to label map
func (m Template) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *User) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m User) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

// GetLabels adds new label to label map
func (User) LabelResourceKind() string {
	return "user"
}

// GetLabels adds new label to label map
func (m User) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *UserGroup) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m UserGroup) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

// GetLabels adds new label to label map
func (UserGroup) LabelResourceKind() string {
	return "userGroup"
}

// GetLabels adds new label to label map
func (m UserGroup) LabelResourceID() uint64 {
	return m.ID
}

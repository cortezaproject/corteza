package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

func (m *Application) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Application) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Application) LabelResourceKind() string {
	return "application"
}

func (m Application) LabelResourceID() uint64 {
	return m.ID
}

func (m *AuthClient) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m AuthClient) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (AuthClient) LabelResourceKind() string {
	return "authClient"
}

func (m AuthClient) LabelResourceID() uint64 {
	return m.ID
}

func (m *Report) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Report) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Report) LabelResourceKind() string {
	return "report"
}

func (m Report) LabelResourceID() uint64 {
	return m.ID
}

func (m *Role) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Role) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Role) LabelResourceKind() string {
	return "role"
}

func (m Role) LabelResourceID() uint64 {
	return m.ID
}

func (m *UserGroup) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m UserGroup) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (UserGroup) LabelResourceKind() string {
	return "userGroup"
}

func (m UserGroup) LabelResourceID() uint64 {
	return m.ID
}

func (m *Template) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Template) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Template) LabelResourceKind() string {
	return "template"
}

func (m Template) LabelResourceID() uint64 {
	return m.ID
}

func (m *User) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m User) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (User) LabelResourceKind() string {
	return "user"
}

func (m User) LabelResourceID() uint64 {
	return m.ID
}

func (m *Connection) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Connection) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Connection) LabelResourceKind() string {
	return "connection"
}

func (m Connection) LabelResourceID() uint64 {
	return m.ID
}

func (m *ConfiguredConnection) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m ConfiguredConnection) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (ConfiguredConnection) LabelResourceKind() string {
	return "configuredConnection"
}

func (m ConfiguredConnection) LabelResourceID() uint64 {
	return m.ID
}

func (m *Agent) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Agent) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Agent) LabelResourceKind() string {
	return "agent"
}

func (m Agent) LabelResourceID() uint64 {
	return m.ID
}

func (m *Chatbot) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Chatbot) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Chatbot) LabelResourceKind() string {
	return "chatbot"
}

func (m Chatbot) LabelResourceID() uint64 {
	return m.ID
}

func (m *Tenant) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Tenant) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Tenant) LabelResourceKind() string {
	return "tenant"
}

func (m Tenant) LabelResourceID() uint64 {
	return m.ID
}

func (m *Project) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Project) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Project) LabelResourceKind() string {
	return "project"
}

func (m Project) LabelResourceID() uint64 {
	return m.ID
}

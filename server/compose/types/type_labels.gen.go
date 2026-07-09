package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

func (m *Chart) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Chart) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Chart) LabelResourceKind() string {
	return "compose:chart"
}

func (m Chart) LabelResourceID() uint64 {
	return m.ID
}

func (m *Module) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Module) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Module) LabelResourceKind() string {
	return "compose:module"
}

func (m Module) LabelResourceID() uint64 {
	return m.ID
}

func (m *ModuleField) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m ModuleField) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (ModuleField) LabelResourceKind() string {
	return "compose:module:field"
}

func (m ModuleField) LabelResourceID() uint64 {
	return m.ID
}

func (m *Namespace) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Namespace) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Namespace) LabelResourceKind() string {
	return "compose:namespace"
}

func (m Namespace) LabelResourceID() uint64 {
	return m.ID
}

func (m *Page) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m Page) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (Page) LabelResourceKind() string {
	return "compose:page"
}

func (m Page) LabelResourceID() uint64 {
	return m.ID
}

func (m *PageLayout) SetLabel(key string, value labelTypes.LabelValue) {
	if m.Labels == nil {
		m.Labels = make(map[string]labelTypes.LabelValue)
	}

	m.Labels[key] = value
}

func (m PageLayout) GetLabels() map[string]labelTypes.LabelValue {
	return m.Labels
}

func (PageLayout) LabelResourceKind() string {
	return "compose:page-layout"
}

func (m PageLayout) LabelResourceID() uint64 {
	return m.ID
}

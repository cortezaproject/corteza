package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/types/types.yaml

// SetLabel adds new label to label map
func (m *Chart) SetLabel(key string, value string) {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Chart) GetLabels() map[string]string {
	return m.Labels
}

// GetLabels adds new label to label map
func (Chart) LabelResourceKind() string {
	return "compose:chart"
}

// GetLabels adds new label to label map
func (m Chart) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *Module) SetLabel(key string, value string) {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Module) GetLabels() map[string]string {
	return m.Labels
}

// GetLabels adds new label to label map
func (Module) LabelResourceKind() string {
	return "compose:module"
}

// GetLabels adds new label to label map
func (m Module) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (f *ModuleField) SetLabel(key string, value string) {
	if f.Labels == nil {
		f.Labels = make(map[string]string)
	}

	f.Labels[key] = value
}

// GetLabels adds new label to label map
func (f ModuleField) GetLabels() map[string]string {
	return f.Labels
}

// GetLabels adds new label to label map
func (ModuleField) LabelResourceKind() string {
	return "compose:module:field"
}

// GetLabels adds new label to label map
func (f ModuleField) LabelResourceID() uint64 {
	return f.ID
}

// SetLabel adds new label to label map
func (m *Namespace) SetLabel(key string, value string) {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Namespace) GetLabels() map[string]string {
	return m.Labels
}

// GetLabels adds new label to label map
func (Namespace) LabelResourceKind() string {
	return "compose:namespace"
}

// GetLabels adds new label to label map
func (m Namespace) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *Page) SetLabel(key string, value string) {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Page) GetLabels() map[string]string {
	return m.Labels
}

// GetLabels adds new label to label map
func (Page) LabelResourceKind() string {
	return "compose:page"
}

// GetLabels adds new label to label map
func (m Page) LabelResourceID() uint64 {
	return m.ID
}

// SetLabel adds new label to label map
func (m *Record) SetLabel(key string, value string) {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}

	m.Labels[key] = value
}

// GetLabels adds new label to label map
func (m Record) GetLabels() map[string]string {
	return m.Labels
}

// GetLabels adds new label to label map
func (Record) LabelResourceKind() string {
	return "compose:record"
}

// GetLabels adds new label to label map
func (m Record) LabelResourceID() uint64 {
	return m.ID
}

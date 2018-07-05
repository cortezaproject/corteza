package crm

import (
	"github.com/jmoiron/sqlx/types"
)

type (
	// Modules
	Module struct {
		ID   uint64
		Name string

		changed []string
	}

	// Modules
	ModuleContentRow struct {
		ID       uint64         `db:"id"`
		ModuleID uint64         `db:"module_id"`
		Fields   types.JSONText `db:"address"`

		changed []string
	}
)

/* Constructors */
func (Module) new() *Module {
	return &Module{}
}
func (ModuleContentRow) new() *ModuleContentRow {
	return &ModuleContentRow{}
}

/* Getters/setters */
func (m *Module) GetID() uint64 {
	return m.ID
}

func (m *Module) SetID(value uint64) *Module {
	if m.ID != value {
		m.changed = append(m.changed, "ID")
		m.ID = value
	}
	return m
}
func (m *Module) GetName() string {
	return m.Name
}

func (m *Module) SetName(value string) *Module {
	if m.Name != value {
		m.changed = append(m.changed, "Name")
		m.Name = value
	}
	return m
}
func (m *ModuleContentRow) GetID() uint64 {
	return m.ID
}

func (m *ModuleContentRow) SetID(value uint64) *ModuleContentRow {
	if m.ID != value {
		m.changed = append(m.changed, "ID")
		m.ID = value
	}
	return m
}
func (m *ModuleContentRow) GetModuleID() uint64 {
	return m.ModuleID
}

func (m *ModuleContentRow) SetModuleID(value uint64) *ModuleContentRow {
	if m.ModuleID != value {
		m.changed = append(m.changed, "ModuleID")
		m.ModuleID = value
	}
	return m
}
func (m *ModuleContentRow) GetFields() types.JSONText {
	return m.Fields
}

func (m *ModuleContentRow) SetFields(value types.JSONText) *ModuleContentRow {
	m.Fields = value
	return m
}

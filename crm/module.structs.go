package crm

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `module.go`, `module.util.go` or `module_test.go` to
	implement your API calls, helper functions and tests. The file `module.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/jmoiron/sqlx/types"
)

type (
	// Modules
	Module struct {
		ID   uint64 `db:"id"`
		Name string `db:"name"`

		changed []string
	}

	// Modules
	ModuleContentRow struct {
		ID       uint64         `db:"id" db:"id"`
		ModuleID uint64         `db:"module_id" db:"module_id"`
		Fields   types.JSONText `db:"address" db:"fields"`

		changed []string
	}
)

/* Constructors */
func (Module) New() *Module {
	return &Module{}
}
func (ModuleContentRow) New() *ModuleContentRow {
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

package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"

	"github.com/crusttech/crust/internal/rules"
)

type (
	// Record is a stored row in the `record` table
	Record struct {
		ID       uint64 `json:"recordID,string" db:"id"`
		ModuleID uint64 `json:"moduleID,string" db:"module_id"`

		Values RecordValueSet `json:"values,omitempty" db:"-"`

		OwnedBy   uint64     `db:"owned_by"   json:"ownedBy,string"`
		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		CreatedBy uint64     `db:"created_by" json:"createdBy,string" `
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty,omitempty"`
		UpdatedBy uint64     `db:"updated_by" json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy uint64     `db:"deleted_by" json:"deletedBy,string,omitempty" `
	}

	// RecordValue is a stored row in the `record_value` table
	RecordValue struct {
		RecordID  uint64     `db:"record_id"  json:"-"`
		Name      string     `db:"name"       json:"name"`
		Value     string     `db:"value"      json:"value,omitempty"`
		Ref       uint64     `db:"ref"        json:"-"`
		Place     uint       `db:"place"      json:"-"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	// Modules - CRM module definitions
	Module struct {
		ID     uint64         `json:"moduleID,string" db:"id"`
		Name   string         `json:"name" db:"name"`
		Meta   types.JSONText `json:"meta" db:"json"`
		Fields ModuleFieldSet `json:"fields" db:"-"`
		Page   *Page          `json:"page,omitempty"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	// Modules - CRM module definitions
	ModuleField struct {
		ModuleID uint64 `json:"moduleID,string" db:"module_id"`
		Place    int    `json:"-" db:"place"`

		Kind  string `json:"kind" db:"kind"`
		Name  string `json:"name" db:"name"`
		Label string `json:"label" db:"label"`

		Options types.JSONText `json:"options" db:"json"`

		Private  bool `json:"isPrivate" db:"is_private"`
		Required bool `json:"isRequired" db:"is_required"`
		Visible  bool `json:"isVisible" db:"is_visible"`
		Multi    bool `json:"isMulti" db:"is_multi"`
	}

	// Page - page structure
	Page struct {
		ID     uint64 `json:"pageID,string" db:"id"`
		SelfID uint64 `json:"selfID,string" db:"self_id"`

		ModuleID uint64  `json:"moduleID,string" db:"module_id"`
		Module   *Module `json:"module,omitempty" db:"-"`

		Title       string `json:"title" db:"title"`
		Description string `json:"description" db:"description"`

		Blocks types.JSONText `json:"blocks" db:"blocks"`

		Children PageSet `json:"children,omitempty" db:"-"`

		Visible bool `json:"visible" db:"visible"`
		Weight  int  `json:"-" db:"weight"`
	}

	// Block - value of Page.Blocks ([]Block)
	Block struct {
		Title       string         `json:"title"`
		Description string         `json:"description"`
		Options     types.JSONText `json:"options"`
		Kind        string         `json:"kind"`
		X           int            `json:"x"`
		Y           int            `json:"y"`
		Width       int            `json:"width"`
		Height      int            `json:"height"`
	}
)

func (set *ModuleFieldSet) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, set)
	}
	return nil
}

func (set ModuleFieldSet) Value() (driver.Value, error) {
	return json.Marshal(set)
}

func (set ModuleFieldSet) Names() (names []string) {
	names = make([]string, len(set))

	for i := range set {
		names[i] = set[i].Name
	}

	return
}

func (set ModuleFieldSet) HasName(name string) bool {
	for i := range set {
		if name == set[i].Name {
			return true
		}
	}

	return false
}

func (set ModuleFieldSet) FindByName(name string) *ModuleField {
	for i := range set {
		if name == set[i].Name {
			return set[i]
		}
	}

	return nil
}

func (set ModuleFieldSet) FilterByModule(moduleID uint64) (ff ModuleFieldSet) {
	for i := range set {
		if set[i].ModuleID == moduleID {
			ff = append(ff, set[i])
		}
	}

	return
}

// IsRef tells us if value of this field be a reference to something (another record, user)?
func (f ModuleField) IsRef() bool {
	return f.Kind == "Record" || f.Kind == "Owner" || f.Kind == "File"
}

// UserIDs returns a slice of user IDs from all items in the set
//
// This function is auto-generated.
func (set RecordSet) UserIDs() (IDs []uint64) {
	IDs = make([]uint64, 0)

loop:
	for i := range set {
		for _, id := range IDs {
			if id == set[i].OwnedBy {
				continue loop
			}
		}

		IDs = append(IDs, set[i].OwnedBy)
	}

	return
}

func (set RecordValueSet) FilterByName(name string) (vv RecordValueSet) {
	for i := range set {
		if set[i].Name == name {
			vv = append(vv, set[i])
		}
	}

	return
}

func (set RecordValueSet) FilterByRecordID(recordID uint64) (vv RecordValueSet) {
	for i := range set {
		if set[i].RecordID == recordID {
			vv = append(vv, set[i])
		}
	}

	return
}

// Resource returns a system resource ID for this type
func (r *Module) Resource() rules.Resource {
	resource := rules.Resource{
		Service: "compose",
		Scope:   "module",
		ID:      r.ID,
	}

	return resource
}

// Resource returns a system resource ID for this type
func (r *Record) Resource() rules.Resource {
	resource := rules.Resource{
		Service: "compose",
		Scope:   "module", // intentionally using module here so we can use Record's resource
		ID:      r.ModuleID,
	}

	return resource
}

// Resource returns a system resource ID for this type
func (r *Page) Resource() rules.Resource {
	resource := rules.Resource{
		Service: "compose",
		Scope:   "page",
		ID:      r.ID,
	}

	return resource
}

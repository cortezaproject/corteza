package types

import (
	"database/sql/driver"
	"encoding/json"
	"sort"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	// Modules - CRM module definitions
	ModuleField struct {
		ID       uint64 `json:"fieldID,string" db:"id"`
		ModuleID uint64 `json:"moduleID,string" db:"rel_module"`
		Place    int    `json:"-" db:"place"`

		Kind  string `json:"kind"  db:"kind"`
		Name  string `json:"name"  db:"name"`
		Label string `json:"label" db:"label"`

		Options ModuleFieldOptions `json:"options" db:"options"`

		Private      bool           `json:"isPrivate" db:"is_private"`
		Required     bool           `json:"isRequired" db:"is_required"`
		Visible      bool           `json:"isVisible" db:"is_visible"`
		Multi        bool           `json:"isMulti" db:"is_multi"`
		DefaultValue RecordValueSet `json:"defaultValue" db:"default_value"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	ModuleFieldOptions map[string]interface{}
)

var (
	_ sort.Interface = &ModuleFieldSet{}
)

func (mfo *ModuleFieldOptions) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*mfo = ModuleFieldOptions{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, mfo); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into ModuleFieldOptions", string(b))
		}
	}

	return nil
}

func (mfo ModuleFieldOptions) Value() (driver.Value, error) {
	return json.Marshal(mfo)
}

// Resource returns a system resource ID for this type
func (m ModuleField) PermissionResource() permissions.Resource {
	return ModuleFieldPermissionResource.AppendID(m.ID)
}

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

func (set ModuleFieldSet) Len() int {
	return len(set)
}

func (set ModuleFieldSet) Less(i, j int) bool {
	return set[i].Place < set[j].Place
}

func (set ModuleFieldSet) Swap(i, j int) {
	set[i], set[j] = set[j], set[i]
}

// IsRef tells us if value of this field be a reference to something (another record, user)?
func (f ModuleField) IsRef() bool {
	return f.Kind == "Record" || f.Kind == "Owner" || f.Kind == "File"
}

func (f ModuleField) IsNumeric() bool {
	return f.Kind == "Number"
}

func (f ModuleField) IsDateTime() bool {
	return f.Kind == "DateTime"
}

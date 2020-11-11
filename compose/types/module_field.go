package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"sort"
	"time"
)

type (
	// Modules - CRM module definitions
	ModuleField struct {
		ID       uint64 `json:"fieldID,string"`
		ModuleID uint64 `json:"moduleID,string"`
		Place    int    `json:"-"`

		Kind  string `json:"kind"`
		Name  string `json:"name"`
		Label string `json:"label"`

		Options ModuleFieldOptions `json:"options"`

		Private      bool           `json:"isPrivate"`
		Required     bool           `json:"isRequired"`
		Visible      bool           `json:"isVisible"`
		Multi        bool           `json:"isMulti"`
		DefaultValue RecordValueSet `json:"defaultValue"`

		Expressions ModuleFieldExpr `json:"expressions"`

		Labels map[string]string `json:"labels,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	ModuleFieldFilter struct {
		ModuleID []uint64
		Deleted  filter.State
	}

	ModuleFieldExpr struct {
		Value      string   `json:"value,omitempty"`
		Sanitizers []string `json:"sanitizers,omitempty"`

		Validators               []ModuleFieldValidator `json:"validators,omitempty"`
		DisableDefaultValidators bool                   `json:"disableDefaultValidators,omitempty"`

		Formatters               []string `json:"formatters,omitempty"`
		DisableDefaultFormatters bool     `json:"disableDefaultFormatters,omitempty"`
	}

	ModuleFieldValidator struct {
		Test  string `json:"test,omitempty"`
		Error string `json:"error,omitempty"`
	}
)

var (
	_ sort.Interface = &ModuleFieldSet{}
)

// Resource returns a system resource ID for this type
func (m ModuleField) RBACResource() rbac.Resource {
	return ModuleFieldRBACResource.AppendID(m.ID)
}

func (m ModuleField) DynamicRoles(userID uint64) []uint64 {
	return nil
}

func (m ModuleField) Clone() *ModuleField {
	return &m
}

func (set ModuleFieldSet) Clone() (out ModuleFieldSet) {
	out = make([]*ModuleField, len(set))
	for i := range set {
		out[i] = set[i].Clone()
	}

	return out
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

func (f ModuleField) IsBoolean() bool {
	return f.Kind == "Bool"
}

func (f ModuleField) IsNumeric() bool {
	return f.Kind == "Number"
}

func (f ModuleField) IsDateTime() bool {
	return f.Kind == "DateTime"
}

// IsRef tells us if value of this field be a reference to something
// (another record, file , user)?
func (f ModuleField) IsRef() bool {
	return f.Kind == "Record" || f.Kind == "User" || f.Kind == "File"
}

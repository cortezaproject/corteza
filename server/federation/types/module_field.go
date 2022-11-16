package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza/server/pkg/sql"
	"strings"
)

type (
	ModuleFieldSet []*ModuleField

	ModuleField struct {
		Kind    string `json:"kind"`
		Name    string `json:"name"`
		Label   string `json:"label"`
		IsMulti bool   `json:"isMulti"`
	}
)

// HasField checks if the fieldset has a value by name
// keeping error field to better match the existing Filter()
// methods in generated sets
func (list ModuleFieldSet) HasField(name string) (bool, error) {
	for _, f := range list {
		if strings.ToLower(f.Name) == strings.ToLower(name) {
			return true, nil
		}
	}

	return false, nil
}

func (list *ModuleFieldSet) Scan(src any) error          { return sql.ParseJSON(src, list) }
func (list ModuleFieldSet) Value() (driver.Value, error) { return json.Marshal(list) }

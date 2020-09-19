package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type (
	ModuleFieldList []*ModuleField

	ModuleField struct {
		Kind    string `json:"kind"`
		Name    string `json:"name"`
		Label   string `json:"label"`
		IsMulti bool   `json:"isMulti"`
	}
)

func (list ModuleFieldList) Value() (driver.Value, error) {
	return json.Marshal(list)
}

func (list *ModuleFieldList) Scan(value interface{}) error {
	switch value.(type) {
	case nil:
		*list = ModuleFieldList{}
	case []uint8:
		if err := json.Unmarshal(value.([]byte), list); err != nil {
			return errors.New(fmt.Sprintf("Can not scan '%v' into ModuleFieldList", value))
		}
	}

	return nil
}

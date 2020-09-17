package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type (
	ModuleFieldMappingList []*ModuleFieldMapping

	ModuleFieldMapping struct {
		Kind    string `json:"kind"`
		Name    string `json:"name"`
		Label   string `json:"label"`
		IsMulti bool   `json:"isMulti"`
	}
)

func (list ModuleFieldMappingList) Value() (driver.Value, error) {
	return json.Marshal(list)
}

func (list *ModuleFieldMappingList) Scan(value interface{}) error {
	switch value.(type) {
	case nil:
		*list = ModuleFieldMappingList{}
	case []uint8:
		if err := json.Unmarshal(value.([]byte), list); err != nil {
			return errors.New(fmt.Sprintf("Can not scan '%v' into RecordValueSet", value))
		}
	}

	return nil
}

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
		Origin      ModuleField `json:"origin"`
		Destination ModuleField `json:"destination"`
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
			return errors.New(fmt.Sprintf("Can not scan '%v' into ModuleFieldMappingList", value))
		}
	}

	return nil
}

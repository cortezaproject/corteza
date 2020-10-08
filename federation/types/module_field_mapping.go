package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type (
	ModuleFieldMappingSet []*ModuleFieldMapping

	ModuleFieldMapping struct {
		Origin      ModuleField `json:"origin"`
		Destination ModuleField `json:"destination"`
	}
)

func (list ModuleFieldMappingSet) Value() (driver.Value, error) {
	return json.Marshal(list)
}

func (list *ModuleFieldMappingSet) Scan(value interface{}) error {
	switch value.(type) {
	case nil:
		*list = ModuleFieldMappingSet{}
	case []uint8:
		if err := json.Unmarshal(value.([]byte), list); err != nil {
			return errors.New(fmt.Sprintf("Can not scan '%v' into ModuleFieldMappingSet", value))
		}
	}

	return nil
}

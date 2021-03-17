package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	ModuleFieldMappingSetFindTypeOrigin ModuleFieldMappingSetFindType = iota
	ModuleFieldMappingSetFindTypeDestination
)

type (
	ModuleFieldMappingSet []*ModuleFieldMapping

	ModuleFieldMapping struct {
		Origin      ModuleField `json:"origin"`
		Destination ModuleField `json:"destination"`
	}

	ModuleFieldMappingSetFindType int
)

// Find looks up a origin or destination mapping
func (list *ModuleFieldMappingSet) FindByName(name string, findType ModuleFieldMappingSetFindType) (*ModuleFieldMapping, error) {
	for _, mfm := range *list {
		switch findType {

		case ModuleFieldMappingSetFindTypeOrigin:
			if mfm.Origin.Name == name {
				return mfm, nil
			}

		case ModuleFieldMappingSetFindTypeDestination:
			if mfm.Destination.Name == name {
				return mfm, nil
			}
		}
	}

	return nil, nil
}

func (list ModuleFieldMappingSet) Value() (driver.Value, error) {
	return json.Marshal(list)
}

func (list *ModuleFieldMappingSet) Scan(value interface{}) error {
	switch value.(type) {
	case nil:
		*list = ModuleFieldMappingSet{}
	case []uint8:
		if err := json.Unmarshal(value.([]byte), list); err != nil {
			return errors.New(fmt.Sprintf("cannot scan '%v' into ModuleFieldMappingSet", value))
		}
	}

	return nil
}

package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)

type (
	ModuleFieldOptions map[string]interface{}
)

const (
	moduleFieldOptionIsUnique           = "isUnique"
	moduleFieldOptionIsUniqueMultiValue = "isUniqueMultiValue"
)

func (opt *ModuleFieldOptions) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*opt = ModuleFieldOptions{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, opt); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into ModuleFieldOptions", string(b))
		}
	}

	return nil
}

func (opt ModuleFieldOptions) Value() (driver.Value, error) {
	return json.Marshal(opt)
}

// Bool returns option value for key as boolean true or false
//
// Invalid, non-existing are returned as false
func (opt ModuleFieldOptions) Bool(key string) bool {
	if _, has := opt[key]; has {
		if v, ok := opt[key].(bool); ok {
			return v
		}
	}

	return false
}

// Strings returns option value for key as slice of strings
//
// Invalid, non-existing are returned as nil
func (opt ModuleFieldOptions) Strings(key string) []string {
	if _, has := opt[key]; has {
		if v, ok := opt[key].([]string); ok {
			return v
		}
	}

	return nil
}

// IsUnique - should value in this field be unique across records?
func (opt ModuleFieldOptions) IsUnique() bool {
	return opt.Bool(moduleFieldOptionIsUnique)
}

// SetIsUnique - should value in this field be unique across records?
func (opt ModuleFieldOptions) SetIsUnique(value bool) {
	opt[moduleFieldOptionIsUnique] = value
}

// IsUniqueMultiValue - should value in this field be unique in the multi-value set?
func (opt ModuleFieldOptions) IsUniqueMultiValue() bool {
	return opt.Bool(moduleFieldOptionIsUniqueMultiValue)
}

func (opt ModuleFieldOptions) SetIsUniqueMultiValue(value bool) {
	// SetIsUniqueMultiValue - should value in this field be unique in the multi-value set?
	opt[moduleFieldOptionIsUniqueMultiValue] = value
}

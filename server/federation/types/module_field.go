package types

import (
	"strings"
)

type (
	ModuleFieldSet []*ModuleField
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

package types

import (
	"encoding/json"

	"github.com/crusttech/human/server/pkg/id"
)

type (
	settingValueAlias SettingValue

	// settingValueJSON is a SettingValue as JSON carries it: the user ID
	// written as a string, read from a string or a number.
	settingValueJSON struct {
		*settingValueAlias
		UpdatedBy id.Uint64 `json:"updatedBy"`
	}
)

func (v SettingValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(settingValueJSON{
		settingValueAlias: (*settingValueAlias)(&v),
		UpdatedBy:         id.Uint64(v.UpdatedBy),
	})
}

func (v *SettingValue) UnmarshalJSON(data []byte) error {
	aux := settingValueJSON{settingValueAlias: (*settingValueAlias)(v)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	v.UpdatedBy = uint64(aux.UpdatedBy)
	return nil
}

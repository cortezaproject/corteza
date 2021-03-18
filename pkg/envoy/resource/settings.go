package resource

import (
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Settings struct {
		*base
		Res types.SettingValueSet
	}

	Setting struct {
		*base
		Res types.SettingValue
	}
)

func NewSettings(vv map[string]interface{}) *Settings {
	r := &Settings{base: &base{}}
	r.SetResourceType(SETTINGS_RESOURCE_TYPE)

	r.Res = make(types.SettingValueSet, 0, len(vv))
	for k, v := range vv {
		sv := &types.SettingValue{
			Name: k,
		}
		sv.SetValue(v)
		r.Res = append(r.Res, sv)
	}

	return r
}

func NewSetting(s *types.SettingValue) *Setting {
	r := &Setting{base: &base{}}
	r.SetResourceType(SETTINGS_RESOURCE_TYPE)
	r.Res = *s

	r.AddIdentifier(identifiers(s.Name)...)

	// Set initial stamps
	r.SetTimestamps(MakeTimestampsCUDA(nil, &s.UpdatedAt, nil, nil))
	r.SetUserstamps(&Userstamps{
		UpdatedBy: &Userstamp{UserID: s.UpdatedBy},
	})

	return r
}

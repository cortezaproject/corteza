package tests

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/settings.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	settingsStore interface {
		SearchSettings(ctx context.Context, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error)
		LookupSettingByNameOwnedBy(ctx context.Context, name string, owned_by uint64) (*types.SettingValue, error)
		CreateSetting(ctx context.Context, rr ...*types.SettingValue) error
		UpdateSetting(ctx context.Context, rr ...*types.SettingValue) error
		PartialUpdateSetting(ctx context.Context, onlyColumns []string, rr ...*types.SettingValue) error
		RemoveSetting(ctx context.Context, rr ...*types.SettingValue) error
		RemoveSettingByNameOwnedBy(ctx context.Context, name string, ownedBy uint64) error

		TruncateSettings(ctx context.Context) error
	}
)

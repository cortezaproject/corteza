package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/settings.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Settings interface {
		SearchSettings(ctx context.Context, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error)
		LookupSettingByNameOwnedBy(ctx context.Context, name string, owned_by uint64) (*types.SettingValue, error)

		CreateSetting(ctx context.Context, rr ...*types.SettingValue) error

		UpdateSetting(ctx context.Context, rr ...*types.SettingValue) error

		UpsertSetting(ctx context.Context, rr ...*types.SettingValue) error

		DeleteSetting(ctx context.Context, rr ...*types.SettingValue) error
		DeleteSettingByNameOwnedBy(ctx context.Context, name string, ownedBy uint64) error

		TruncateSettings(ctx context.Context) error
	}
)

var _ *types.SettingValue
var _ context.Context

// SearchSettings returns all matching Settings from store
func SearchSettings(ctx context.Context, s Settings, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error) {
	return s.SearchSettings(ctx, f)
}

// LookupSettingByNameOwnedBy searches for settings by name and owner
func LookupSettingByNameOwnedBy(ctx context.Context, s Settings, name string, owned_by uint64) (*types.SettingValue, error) {
	return s.LookupSettingByNameOwnedBy(ctx, name, owned_by)
}

// CreateSetting creates one or more Settings in store
func CreateSetting(ctx context.Context, s Settings, rr ...*types.SettingValue) error {
	return s.CreateSetting(ctx, rr...)
}

// UpdateSetting updates one or more (existing) Settings in store
func UpdateSetting(ctx context.Context, s Settings, rr ...*types.SettingValue) error {
	return s.UpdateSetting(ctx, rr...)
}

// UpsertSetting creates new or updates existing one or more Settings in store
func UpsertSetting(ctx context.Context, s Settings, rr ...*types.SettingValue) error {
	return s.UpsertSetting(ctx, rr...)
}

// DeleteSetting Deletes one or more Settings from store
func DeleteSetting(ctx context.Context, s Settings, rr ...*types.SettingValue) error {
	return s.DeleteSetting(ctx, rr...)
}

// DeleteSettingByNameOwnedBy Deletes Setting from store
func DeleteSettingByNameOwnedBy(ctx context.Context, s Settings, name string, ownedBy uint64) error {
	return s.DeleteSettingByNameOwnedBy(ctx, name, ownedBy)
}

// TruncateSettings Deletes all Settings from store
func TruncateSettings(ctx context.Context, s Settings) error {
	return s.TruncateSettings(ctx)
}

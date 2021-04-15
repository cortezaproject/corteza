package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messagebus_queue_settings.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
)

type (
	MessagebusQueueSettings interface {
		SearchMessagebusQueueSettings(ctx context.Context, f messagebus.QueueSettingsFilter) (messagebus.QueueSettingsSet, messagebus.QueueSettingsFilter, error)
		LookupMessagebusQueueSettingByID(ctx context.Context, id uint64) (*messagebus.QueueSettings, error)
		LookupMessagebusQueueSettingByQueue(ctx context.Context, queue string) (*messagebus.QueueSettings, error)

		CreateMessagebusQueueSetting(ctx context.Context, rr ...*messagebus.QueueSettings) error

		UpdateMessagebusQueueSetting(ctx context.Context, rr ...*messagebus.QueueSettings) error

		UpsertMessagebusQueueSetting(ctx context.Context, rr ...*messagebus.QueueSettings) error

		DeleteMessagebusQueueSetting(ctx context.Context, rr ...*messagebus.QueueSettings) error
		DeleteMessagebusQueueSettingByID(ctx context.Context, ID uint64) error

		TruncateMessagebusQueueSettings(ctx context.Context) error
	}
)

var _ *messagebus.QueueSettings
var _ context.Context

// SearchMessagebusQueueSettings returns all matching MessagebusQueueSettings from store
func SearchMessagebusQueueSettings(ctx context.Context, s MessagebusQueueSettings, f messagebus.QueueSettingsFilter) (messagebus.QueueSettingsSet, messagebus.QueueSettingsFilter, error) {
	return s.SearchMessagebusQueueSettings(ctx, f)
}

// LookupMessagebusQueueSettingByID searches for queue by ID
func LookupMessagebusQueueSettingByID(ctx context.Context, s MessagebusQueueSettings, id uint64) (*messagebus.QueueSettings, error) {
	return s.LookupMessagebusQueueSettingByID(ctx, id)
}

// LookupMessagebusQueueSettingByQueue searches for queue by queue name
func LookupMessagebusQueueSettingByQueue(ctx context.Context, s MessagebusQueueSettings, queue string) (*messagebus.QueueSettings, error) {
	return s.LookupMessagebusQueueSettingByQueue(ctx, queue)
}

// CreateMessagebusQueueSetting creates one or more MessagebusQueueSettings in store
func CreateMessagebusQueueSetting(ctx context.Context, s MessagebusQueueSettings, rr ...*messagebus.QueueSettings) error {
	return s.CreateMessagebusQueueSetting(ctx, rr...)
}

// UpdateMessagebusQueueSetting updates one or more (existing) MessagebusQueueSettings in store
func UpdateMessagebusQueueSetting(ctx context.Context, s MessagebusQueueSettings, rr ...*messagebus.QueueSettings) error {
	return s.UpdateMessagebusQueueSetting(ctx, rr...)
}

// UpsertMessagebusQueueSetting creates new or updates existing one or more MessagebusQueueSettings in store
func UpsertMessagebusQueueSetting(ctx context.Context, s MessagebusQueueSettings, rr ...*messagebus.QueueSettings) error {
	return s.UpsertMessagebusQueueSetting(ctx, rr...)
}

// DeleteMessagebusQueueSetting Deletes one or more MessagebusQueueSettings from store
func DeleteMessagebusQueueSetting(ctx context.Context, s MessagebusQueueSettings, rr ...*messagebus.QueueSettings) error {
	return s.DeleteMessagebusQueueSetting(ctx, rr...)
}

// DeleteMessagebusQueueSettingByID Deletes MessagebusQueueSetting from store
func DeleteMessagebusQueueSettingByID(ctx context.Context, s MessagebusQueueSettings, ID uint64) error {
	return s.DeleteMessagebusQueueSettingByID(ctx, ID)
}

// TruncateMessagebusQueueSettings Deletes all MessagebusQueueSettings from store
func TruncateMessagebusQueueSettings(ctx context.Context, s MessagebusQueueSettings) error {
	return s.TruncateMessagebusQueueSettings(ctx)
}

package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messagebus_queuesettings.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
)

type (
	MessagebusQueuesettings interface {
		SearchMessagebusQueuesettings(ctx context.Context, f messagebus.QueueSettingsFilter) (messagebus.QueueSettingsSet, messagebus.QueueSettingsFilter, error)

		CreateMessagebusQueuesetting(ctx context.Context, rr ...*messagebus.QueueSettings) error

		UpdateMessagebusQueuesetting(ctx context.Context, rr ...*messagebus.QueueSettings) error

		UpsertMessagebusQueuesetting(ctx context.Context, rr ...*messagebus.QueueSettings) error

		DeleteMessagebusQueuesetting(ctx context.Context, rr ...*messagebus.QueueSettings) error
		DeleteMessagebusQueuesettingByID(ctx context.Context, ID uint64) error

		TruncateMessagebusQueuesettings(ctx context.Context) error
	}
)

var _ *messagebus.QueueSettings
var _ context.Context

// SearchMessagebusQueuesettings returns all matching MessagebusQueuesettings from store
func SearchMessagebusQueuesettings(ctx context.Context, s MessagebusQueuesettings, f messagebus.QueueSettingsFilter) (messagebus.QueueSettingsSet, messagebus.QueueSettingsFilter, error) {
	return s.SearchMessagebusQueuesettings(ctx, f)
}

// CreateMessagebusQueuesetting creates one or more MessagebusQueuesettings in store
func CreateMessagebusQueuesetting(ctx context.Context, s MessagebusQueuesettings, rr ...*messagebus.QueueSettings) error {
	return s.CreateMessagebusQueuesetting(ctx, rr...)
}

// UpdateMessagebusQueuesetting updates one or more (existing) MessagebusQueuesettings in store
func UpdateMessagebusQueuesetting(ctx context.Context, s MessagebusQueuesettings, rr ...*messagebus.QueueSettings) error {
	return s.UpdateMessagebusQueuesetting(ctx, rr...)
}

// UpsertMessagebusQueuesetting creates new or updates existing one or more MessagebusQueuesettings in store
func UpsertMessagebusQueuesetting(ctx context.Context, s MessagebusQueuesettings, rr ...*messagebus.QueueSettings) error {
	return s.UpsertMessagebusQueuesetting(ctx, rr...)
}

// DeleteMessagebusQueuesetting Deletes one or more MessagebusQueuesettings from store
func DeleteMessagebusQueuesetting(ctx context.Context, s MessagebusQueuesettings, rr ...*messagebus.QueueSettings) error {
	return s.DeleteMessagebusQueuesetting(ctx, rr...)
}

// DeleteMessagebusQueuesettingByID Deletes MessagebusQueuesetting from store
func DeleteMessagebusQueuesettingByID(ctx context.Context, s MessagebusQueuesettings, ID uint64) error {
	return s.DeleteMessagebusQueuesettingByID(ctx, ID)
}

// TruncateMessagebusQueuesettings Deletes all MessagebusQueuesettings from store
func TruncateMessagebusQueuesettings(ctx context.Context, s MessagebusQueuesettings) error {
	return s.TruncateMessagebusQueuesettings(ctx)
}

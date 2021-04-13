package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messagebus_queuemessage.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
)

type (
	MessagebusQueuemessages interface {
		SearchMessagebusQueuemessages(ctx context.Context, f messagebus.QueueMessageFilter) (messagebus.QueueMessageSet, messagebus.QueueMessageFilter, error)

		CreateMessagebusQueuemessage(ctx context.Context, rr ...*messagebus.QueueMessage) error

		UpdateMessagebusQueuemessage(ctx context.Context, rr ...*messagebus.QueueMessage) error

		DeleteMessagebusQueuemessage(ctx context.Context, rr ...*messagebus.QueueMessage) error
		DeleteMessagebusQueuemessageByID(ctx context.Context, ID uint64) error

		TruncateMessagebusQueuemessages(ctx context.Context) error
	}
)

var _ *messagebus.QueueMessage
var _ context.Context

// SearchMessagebusQueuemessages returns all matching MessagebusQueuemessages from store
func SearchMessagebusQueuemessages(ctx context.Context, s MessagebusQueuemessages, f messagebus.QueueMessageFilter) (messagebus.QueueMessageSet, messagebus.QueueMessageFilter, error) {
	return s.SearchMessagebusQueuemessages(ctx, f)
}

// CreateMessagebusQueuemessage creates one or more MessagebusQueuemessages in store
func CreateMessagebusQueuemessage(ctx context.Context, s MessagebusQueuemessages, rr ...*messagebus.QueueMessage) error {
	return s.CreateMessagebusQueuemessage(ctx, rr...)
}

// UpdateMessagebusQueuemessage updates one or more (existing) MessagebusQueuemessages in store
func UpdateMessagebusQueuemessage(ctx context.Context, s MessagebusQueuemessages, rr ...*messagebus.QueueMessage) error {
	return s.UpdateMessagebusQueuemessage(ctx, rr...)
}

// DeleteMessagebusQueuemessage Deletes one or more MessagebusQueuemessages from store
func DeleteMessagebusQueuemessage(ctx context.Context, s MessagebusQueuemessages, rr ...*messagebus.QueueMessage) error {
	return s.DeleteMessagebusQueuemessage(ctx, rr...)
}

// DeleteMessagebusQueuemessageByID Deletes MessagebusQueuemessage from store
func DeleteMessagebusQueuemessageByID(ctx context.Context, s MessagebusQueuemessages, ID uint64) error {
	return s.DeleteMessagebusQueuemessageByID(ctx, ID)
}

// TruncateMessagebusQueuemessages Deletes all MessagebusQueuemessages from store
func TruncateMessagebusQueuemessages(ctx context.Context, s MessagebusQueuemessages) error {
	return s.TruncateMessagebusQueuemessages(ctx)
}

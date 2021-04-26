package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messagebus_queue_message.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
)

type (
	MessagebusQueueMessages interface {
		SearchMessagebusQueueMessages(ctx context.Context, f messagebus.QueueMessageFilter) (messagebus.QueueMessageSet, messagebus.QueueMessageFilter, error)

		CreateMessagebusQueueMessage(ctx context.Context, rr ...*messagebus.QueueMessage) error

		UpdateMessagebusQueueMessage(ctx context.Context, rr ...*messagebus.QueueMessage) error

		DeleteMessagebusQueueMessage(ctx context.Context, rr ...*messagebus.QueueMessage) error
		DeleteMessagebusQueueMessageByID(ctx context.Context, ID uint64) error

		TruncateMessagebusQueueMessages(ctx context.Context) error
	}
)

var _ *messagebus.QueueMessage
var _ context.Context

// SearchMessagebusQueueMessages returns all matching MessagebusQueueMessages from store
func SearchMessagebusQueueMessages(ctx context.Context, s MessagebusQueueMessages, f messagebus.QueueMessageFilter) (messagebus.QueueMessageSet, messagebus.QueueMessageFilter, error) {
	return s.SearchMessagebusQueueMessages(ctx, f)
}

// CreateMessagebusQueueMessage creates one or more MessagebusQueueMessages in store
func CreateMessagebusQueueMessage(ctx context.Context, s MessagebusQueueMessages, rr ...*messagebus.QueueMessage) error {
	return s.CreateMessagebusQueueMessage(ctx, rr...)
}

// UpdateMessagebusQueueMessage updates one or more (existing) MessagebusQueueMessages in store
func UpdateMessagebusQueueMessage(ctx context.Context, s MessagebusQueueMessages, rr ...*messagebus.QueueMessage) error {
	return s.UpdateMessagebusQueueMessage(ctx, rr...)
}

// DeleteMessagebusQueueMessage Deletes one or more MessagebusQueueMessages from store
func DeleteMessagebusQueueMessage(ctx context.Context, s MessagebusQueueMessages, rr ...*messagebus.QueueMessage) error {
	return s.DeleteMessagebusQueueMessage(ctx, rr...)
}

// DeleteMessagebusQueueMessageByID Deletes MessagebusQueueMessage from store
func DeleteMessagebusQueueMessageByID(ctx context.Context, s MessagebusQueueMessages, ID uint64) error {
	return s.DeleteMessagebusQueueMessageByID(ctx, ID)
}

// TruncateMessagebusQueueMessages Deletes all MessagebusQueueMessages from store
func TruncateMessagebusQueueMessages(ctx context.Context, s MessagebusQueueMessages) error {
	return s.TruncateMessagebusQueueMessages(ctx)
}

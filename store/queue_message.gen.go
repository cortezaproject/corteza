package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/queue_message.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	QueueMessages interface {
		SearchQueueMessages(ctx context.Context, f types.QueueMessageFilter) (types.QueueMessageSet, types.QueueMessageFilter, error)

		CreateQueueMessage(ctx context.Context, rr ...*types.QueueMessage) error

		UpdateQueueMessage(ctx context.Context, rr ...*types.QueueMessage) error

		DeleteQueueMessage(ctx context.Context, rr ...*types.QueueMessage) error
		DeleteQueueMessageByID(ctx context.Context, ID uint64) error

		TruncateQueueMessages(ctx context.Context) error
	}
)

var _ *types.QueueMessage
var _ context.Context

// SearchQueueMessages returns all matching QueueMessages from store
func SearchQueueMessages(ctx context.Context, s QueueMessages, f types.QueueMessageFilter) (types.QueueMessageSet, types.QueueMessageFilter, error) {
	return s.SearchQueueMessages(ctx, f)
}

// CreateQueueMessage creates one or more QueueMessages in store
func CreateQueueMessage(ctx context.Context, s QueueMessages, rr ...*types.QueueMessage) error {
	return s.CreateQueueMessage(ctx, rr...)
}

// UpdateQueueMessage updates one or more (existing) QueueMessages in store
func UpdateQueueMessage(ctx context.Context, s QueueMessages, rr ...*types.QueueMessage) error {
	return s.UpdateQueueMessage(ctx, rr...)
}

// DeleteQueueMessage Deletes one or more QueueMessages from store
func DeleteQueueMessage(ctx context.Context, s QueueMessages, rr ...*types.QueueMessage) error {
	return s.DeleteQueueMessage(ctx, rr...)
}

// DeleteQueueMessageByID Deletes QueueMessage from store
func DeleteQueueMessageByID(ctx context.Context, s QueueMessages, ID uint64) error {
	return s.DeleteQueueMessageByID(ctx, ID)
}

// TruncateQueueMessages Deletes all QueueMessages from store
func TruncateQueueMessages(ctx context.Context, s QueueMessages) error {
	return s.TruncateQueueMessages(ctx)
}

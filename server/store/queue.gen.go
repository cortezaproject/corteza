package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/queue.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	Queues interface {
		SearchQueues(ctx context.Context, f types.QueueFilter) (types.QueueSet, types.QueueFilter, error)
		LookupQueueByID(ctx context.Context, id uint64) (*types.Queue, error)
		LookupQueueByQueue(ctx context.Context, queue string) (*types.Queue, error)

		CreateQueue(ctx context.Context, rr ...*types.Queue) error

		UpdateQueue(ctx context.Context, rr ...*types.Queue) error

		UpsertQueue(ctx context.Context, rr ...*types.Queue) error

		DeleteQueue(ctx context.Context, rr ...*types.Queue) error
		DeleteQueueByID(ctx context.Context, ID uint64) error

		TruncateQueues(ctx context.Context) error
	}
)

var _ *types.Queue
var _ context.Context

// SearchQueues returns all matching Queues from store
func SearchQueues(ctx context.Context, s Queues, f types.QueueFilter) (types.QueueSet, types.QueueFilter, error) {
	return s.SearchQueues(ctx, f)
}

// LookupQueueByID searches for queue by ID
func LookupQueueByID(ctx context.Context, s Queues, id uint64) (*types.Queue, error) {
	return s.LookupQueueByID(ctx, id)
}

// LookupQueueByQueue searches for queue by queue name
func LookupQueueByQueue(ctx context.Context, s Queues, queue string) (*types.Queue, error) {
	return s.LookupQueueByQueue(ctx, queue)
}

// CreateQueue creates one or more Queues in store
func CreateQueue(ctx context.Context, s Queues, rr ...*types.Queue) error {
	return s.CreateQueue(ctx, rr...)
}

// UpdateQueue updates one or more (existing) Queues in store
func UpdateQueue(ctx context.Context, s Queues, rr ...*types.Queue) error {
	return s.UpdateQueue(ctx, rr...)
}

// UpsertQueue creates new or updates existing one or more Queues in store
func UpsertQueue(ctx context.Context, s Queues, rr ...*types.Queue) error {
	return s.UpsertQueue(ctx, rr...)
}

// DeleteQueue Deletes one or more Queues from store
func DeleteQueue(ctx context.Context, s Queues, rr ...*types.Queue) error {
	return s.DeleteQueue(ctx, rr...)
}

// DeleteQueueByID Deletes Queue from store
func DeleteQueueByID(ctx context.Context, s Queues, ID uint64) error {
	return s.DeleteQueueByID(ctx, ID)
}

// TruncateQueues Deletes all Queues from store
func TruncateQueues(ctx context.Context, s Queues) error {
	return s.TruncateQueues(ctx)
}

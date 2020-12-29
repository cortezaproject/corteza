package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/triggers.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Triggers interface {
		SearchTriggers(ctx context.Context, f types.TriggerFilter) (types.TriggerSet, types.TriggerFilter, error)
		LookupTriggerByID(ctx context.Context, id uint64) (*types.Trigger, error)

		CreateTrigger(ctx context.Context, rr ...*types.Trigger) error

		UpdateTrigger(ctx context.Context, rr ...*types.Trigger) error

		UpsertTrigger(ctx context.Context, rr ...*types.Trigger) error

		DeleteTrigger(ctx context.Context, rr ...*types.Trigger) error
		DeleteTriggerByID(ctx context.Context, ID uint64) error

		TruncateTriggers(ctx context.Context) error
	}
)

var _ *types.Trigger
var _ context.Context

// SearchTriggers returns all matching Triggers from store
func SearchTriggers(ctx context.Context, s Triggers, f types.TriggerFilter) (types.TriggerSet, types.TriggerFilter, error) {
	return s.SearchTriggers(ctx, f)
}

// LookupTriggerByID searches for trigger by ID
//
// It returns trigger even if deleted or suspended
func LookupTriggerByID(ctx context.Context, s Triggers, id uint64) (*types.Trigger, error) {
	return s.LookupTriggerByID(ctx, id)
}

// CreateTrigger creates one or more Triggers in store
func CreateTrigger(ctx context.Context, s Triggers, rr ...*types.Trigger) error {
	return s.CreateTrigger(ctx, rr...)
}

// UpdateTrigger updates one or more (existing) Triggers in store
func UpdateTrigger(ctx context.Context, s Triggers, rr ...*types.Trigger) error {
	return s.UpdateTrigger(ctx, rr...)
}

// UpsertTrigger creates new or updates existing one or more Triggers in store
func UpsertTrigger(ctx context.Context, s Triggers, rr ...*types.Trigger) error {
	return s.UpsertTrigger(ctx, rr...)
}

// DeleteTrigger Deletes one or more Triggers from store
func DeleteTrigger(ctx context.Context, s Triggers, rr ...*types.Trigger) error {
	return s.DeleteTrigger(ctx, rr...)
}

// DeleteTriggerByID Deletes Trigger from store
func DeleteTriggerByID(ctx context.Context, s Triggers, ID uint64) error {
	return s.DeleteTriggerByID(ctx, ID)
}

// TruncateTriggers Deletes all Triggers from store
func TruncateTriggers(ctx context.Context, s Triggers) error {
	return s.TruncateTriggers(ctx)
}

package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/automation_triggers.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/types"
)

type (
	AutomationTriggers interface {
		SearchAutomationTriggers(ctx context.Context, f types.TriggerFilter) (types.TriggerSet, types.TriggerFilter, error)
		LookupAutomationTriggerByID(ctx context.Context, id uint64) (*types.Trigger, error)

		CreateAutomationTrigger(ctx context.Context, rr ...*types.Trigger) error

		UpdateAutomationTrigger(ctx context.Context, rr ...*types.Trigger) error

		UpsertAutomationTrigger(ctx context.Context, rr ...*types.Trigger) error

		DeleteAutomationTrigger(ctx context.Context, rr ...*types.Trigger) error
		DeleteAutomationTriggerByID(ctx context.Context, ID uint64) error

		TruncateAutomationTriggers(ctx context.Context) error
	}
)

var _ *types.Trigger
var _ context.Context

// SearchAutomationTriggers returns all matching AutomationTriggers from store
func SearchAutomationTriggers(ctx context.Context, s AutomationTriggers, f types.TriggerFilter) (types.TriggerSet, types.TriggerFilter, error) {
	return s.SearchAutomationTriggers(ctx, f)
}

// LookupAutomationTriggerByID searches for trigger by ID
//
// It returns trigger even if deleted
func LookupAutomationTriggerByID(ctx context.Context, s AutomationTriggers, id uint64) (*types.Trigger, error) {
	return s.LookupAutomationTriggerByID(ctx, id)
}

// CreateAutomationTrigger creates one or more AutomationTriggers in store
func CreateAutomationTrigger(ctx context.Context, s AutomationTriggers, rr ...*types.Trigger) error {
	return s.CreateAutomationTrigger(ctx, rr...)
}

// UpdateAutomationTrigger updates one or more (existing) AutomationTriggers in store
func UpdateAutomationTrigger(ctx context.Context, s AutomationTriggers, rr ...*types.Trigger) error {
	return s.UpdateAutomationTrigger(ctx, rr...)
}

// UpsertAutomationTrigger creates new or updates existing one or more AutomationTriggers in store
func UpsertAutomationTrigger(ctx context.Context, s AutomationTriggers, rr ...*types.Trigger) error {
	return s.UpsertAutomationTrigger(ctx, rr...)
}

// DeleteAutomationTrigger Deletes one or more AutomationTriggers from store
func DeleteAutomationTrigger(ctx context.Context, s AutomationTriggers, rr ...*types.Trigger) error {
	return s.DeleteAutomationTrigger(ctx, rr...)
}

// DeleteAutomationTriggerByID Deletes AutomationTrigger from store
func DeleteAutomationTriggerByID(ctx context.Context, s AutomationTriggers, ID uint64) error {
	return s.DeleteAutomationTriggerByID(ctx, ID)
}

// TruncateAutomationTriggers Deletes all AutomationTriggers from store
func TruncateAutomationTriggers(ctx context.Context, s AutomationTriggers) error {
	return s.TruncateAutomationTriggers(ctx)
}

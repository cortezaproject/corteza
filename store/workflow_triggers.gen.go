package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/workflow_triggers.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	WorkflowTriggers interface {
		SearchWorkflowTriggers(ctx context.Context, f types.WorkflowTriggerFilter) (types.WorkflowTriggerSet, types.WorkflowTriggerFilter, error)

		CreateWorkflowTrigger(ctx context.Context, rr ...*types.WorkflowTrigger) error

		UpdateWorkflowTrigger(ctx context.Context, rr ...*types.WorkflowTrigger) error

		UpsertWorkflowTrigger(ctx context.Context, rr ...*types.WorkflowTrigger) error

		DeleteWorkflowTrigger(ctx context.Context, rr ...*types.WorkflowTrigger) error
		DeleteWorkflowTriggerByID(ctx context.Context, ID uint64) error

		TruncateWorkflowTriggers(ctx context.Context) error
	}
)

var _ *types.WorkflowTrigger
var _ context.Context

// SearchWorkflowTriggers returns all matching WorkflowTriggers from store
func SearchWorkflowTriggers(ctx context.Context, s WorkflowTriggers, f types.WorkflowTriggerFilter) (types.WorkflowTriggerSet, types.WorkflowTriggerFilter, error) {
	return s.SearchWorkflowTriggers(ctx, f)
}

// CreateWorkflowTrigger creates one or more WorkflowTriggers in store
func CreateWorkflowTrigger(ctx context.Context, s WorkflowTriggers, rr ...*types.WorkflowTrigger) error {
	return s.CreateWorkflowTrigger(ctx, rr...)
}

// UpdateWorkflowTrigger updates one or more (existing) WorkflowTriggers in store
func UpdateWorkflowTrigger(ctx context.Context, s WorkflowTriggers, rr ...*types.WorkflowTrigger) error {
	return s.UpdateWorkflowTrigger(ctx, rr...)
}

// UpsertWorkflowTrigger creates new or updates existing one or more WorkflowTriggers in store
func UpsertWorkflowTrigger(ctx context.Context, s WorkflowTriggers, rr ...*types.WorkflowTrigger) error {
	return s.UpsertWorkflowTrigger(ctx, rr...)
}

// DeleteWorkflowTrigger Deletes one or more WorkflowTriggers from store
func DeleteWorkflowTrigger(ctx context.Context, s WorkflowTriggers, rr ...*types.WorkflowTrigger) error {
	return s.DeleteWorkflowTrigger(ctx, rr...)
}

// DeleteWorkflowTriggerByID Deletes WorkflowTrigger from store
func DeleteWorkflowTriggerByID(ctx context.Context, s WorkflowTriggers, ID uint64) error {
	return s.DeleteWorkflowTriggerByID(ctx, ID)
}

// TruncateWorkflowTriggers Deletes all WorkflowTriggers from store
func TruncateWorkflowTriggers(ctx context.Context, s WorkflowTriggers) error {
	return s.TruncateWorkflowTriggers(ctx)
}

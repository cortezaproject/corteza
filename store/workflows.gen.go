package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/workflows.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Workflows interface {
		SearchWorkflows(ctx context.Context, f types.WorkflowFilter) (types.WorkflowSet, types.WorkflowFilter, error)
		LookupWorkflowByID(ctx context.Context, id uint64) (*types.Workflow, error)
		LookupWorkflowByHandle(ctx context.Context, handle string) (*types.Workflow, error)

		CreateWorkflow(ctx context.Context, rr ...*types.Workflow) error

		UpdateWorkflow(ctx context.Context, rr ...*types.Workflow) error

		UpsertWorkflow(ctx context.Context, rr ...*types.Workflow) error

		DeleteWorkflow(ctx context.Context, rr ...*types.Workflow) error
		DeleteWorkflowByID(ctx context.Context, ID uint64) error

		TruncateWorkflows(ctx context.Context) error
	}
)

var _ *types.Workflow
var _ context.Context

// SearchWorkflows returns all matching Workflows from store
func SearchWorkflows(ctx context.Context, s Workflows, f types.WorkflowFilter) (types.WorkflowSet, types.WorkflowFilter, error) {
	return s.SearchWorkflows(ctx, f)
}

// LookupWorkflowByID searches for workflow by ID
//
// It returns workflow even if deleted or suspended
func LookupWorkflowByID(ctx context.Context, s Workflows, id uint64) (*types.Workflow, error) {
	return s.LookupWorkflowByID(ctx, id)
}

// LookupWorkflowByHandle searches for workflow by their handle
//
// It returns only valid workflows (not deleted, not suspended)
func LookupWorkflowByHandle(ctx context.Context, s Workflows, handle string) (*types.Workflow, error) {
	return s.LookupWorkflowByHandle(ctx, handle)
}

// CreateWorkflow creates one or more Workflows in store
func CreateWorkflow(ctx context.Context, s Workflows, rr ...*types.Workflow) error {
	return s.CreateWorkflow(ctx, rr...)
}

// UpdateWorkflow updates one or more (existing) Workflows in store
func UpdateWorkflow(ctx context.Context, s Workflows, rr ...*types.Workflow) error {
	return s.UpdateWorkflow(ctx, rr...)
}

// UpsertWorkflow creates new or updates existing one or more Workflows in store
func UpsertWorkflow(ctx context.Context, s Workflows, rr ...*types.Workflow) error {
	return s.UpsertWorkflow(ctx, rr...)
}

// DeleteWorkflow Deletes one or more Workflows from store
func DeleteWorkflow(ctx context.Context, s Workflows, rr ...*types.Workflow) error {
	return s.DeleteWorkflow(ctx, rr...)
}

// DeleteWorkflowByID Deletes Workflow from store
func DeleteWorkflowByID(ctx context.Context, s Workflows, ID uint64) error {
	return s.DeleteWorkflowByID(ctx, ID)
}

// TruncateWorkflows Deletes all Workflows from store
func TruncateWorkflows(ctx context.Context, s Workflows) error {
	return s.TruncateWorkflows(ctx)
}

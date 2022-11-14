package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/automation_workflows.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/types"
)

type (
	AutomationWorkflows interface {
		SearchAutomationWorkflows(ctx context.Context, f types.WorkflowFilter) (types.WorkflowSet, types.WorkflowFilter, error)
		LookupAutomationWorkflowByID(ctx context.Context, id uint64) (*types.Workflow, error)
		LookupAutomationWorkflowByHandle(ctx context.Context, handle string) (*types.Workflow, error)

		CreateAutomationWorkflow(ctx context.Context, rr ...*types.Workflow) error

		UpdateAutomationWorkflow(ctx context.Context, rr ...*types.Workflow) error

		UpsertAutomationWorkflow(ctx context.Context, rr ...*types.Workflow) error

		DeleteAutomationWorkflow(ctx context.Context, rr ...*types.Workflow) error
		DeleteAutomationWorkflowByID(ctx context.Context, ID uint64) error

		TruncateAutomationWorkflows(ctx context.Context) error
	}
)

var _ *types.Workflow
var _ context.Context

// SearchAutomationWorkflows returns all matching AutomationWorkflows from store
func SearchAutomationWorkflows(ctx context.Context, s AutomationWorkflows, f types.WorkflowFilter) (types.WorkflowSet, types.WorkflowFilter, error) {
	return s.SearchAutomationWorkflows(ctx, f)
}

// LookupAutomationWorkflowByID searches for workflow by ID
//
// It returns workflow even if deleted
func LookupAutomationWorkflowByID(ctx context.Context, s AutomationWorkflows, id uint64) (*types.Workflow, error) {
	return s.LookupAutomationWorkflowByID(ctx, id)
}

// LookupAutomationWorkflowByHandle searches for workflow by their handle
//
// It returns only valid workflows
func LookupAutomationWorkflowByHandle(ctx context.Context, s AutomationWorkflows, handle string) (*types.Workflow, error) {
	return s.LookupAutomationWorkflowByHandle(ctx, handle)
}

// CreateAutomationWorkflow creates one or more AutomationWorkflows in store
func CreateAutomationWorkflow(ctx context.Context, s AutomationWorkflows, rr ...*types.Workflow) error {
	return s.CreateAutomationWorkflow(ctx, rr...)
}

// UpdateAutomationWorkflow updates one or more (existing) AutomationWorkflows in store
func UpdateAutomationWorkflow(ctx context.Context, s AutomationWorkflows, rr ...*types.Workflow) error {
	return s.UpdateAutomationWorkflow(ctx, rr...)
}

// UpsertAutomationWorkflow creates new or updates existing one or more AutomationWorkflows in store
func UpsertAutomationWorkflow(ctx context.Context, s AutomationWorkflows, rr ...*types.Workflow) error {
	return s.UpsertAutomationWorkflow(ctx, rr...)
}

// DeleteAutomationWorkflow Deletes one or more AutomationWorkflows from store
func DeleteAutomationWorkflow(ctx context.Context, s AutomationWorkflows, rr ...*types.Workflow) error {
	return s.DeleteAutomationWorkflow(ctx, rr...)
}

// DeleteAutomationWorkflowByID Deletes AutomationWorkflow from store
func DeleteAutomationWorkflowByID(ctx context.Context, s AutomationWorkflows, ID uint64) error {
	return s.DeleteAutomationWorkflowByID(ctx, ID)
}

// TruncateAutomationWorkflows Deletes all AutomationWorkflows from store
func TruncateAutomationWorkflows(ctx context.Context, s AutomationWorkflows) error {
	return s.TruncateAutomationWorkflows(ctx)
}

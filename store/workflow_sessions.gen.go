package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/workflow_sessions.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	WorkflowSessions interface {
		SearchWorkflowSessions(ctx context.Context, f types.WorkflowSessionFilter) (types.WorkflowSessionSet, types.WorkflowSessionFilter, error)

		CreateWorkflowSession(ctx context.Context, rr ...*types.WorkflowSession) error

		UpdateWorkflowSession(ctx context.Context, rr ...*types.WorkflowSession) error

		UpsertWorkflowSession(ctx context.Context, rr ...*types.WorkflowSession) error

		DeleteWorkflowSession(ctx context.Context, rr ...*types.WorkflowSession) error
		DeleteWorkflowSessionByID(ctx context.Context, ID uint64) error

		TruncateWorkflowSessions(ctx context.Context) error
	}
)

var _ *types.WorkflowSession
var _ context.Context

// SearchWorkflowSessions returns all matching WorkflowSessions from store
func SearchWorkflowSessions(ctx context.Context, s WorkflowSessions, f types.WorkflowSessionFilter) (types.WorkflowSessionSet, types.WorkflowSessionFilter, error) {
	return s.SearchWorkflowSessions(ctx, f)
}

// CreateWorkflowSession creates one or more WorkflowSessions in store
func CreateWorkflowSession(ctx context.Context, s WorkflowSessions, rr ...*types.WorkflowSession) error {
	return s.CreateWorkflowSession(ctx, rr...)
}

// UpdateWorkflowSession updates one or more (existing) WorkflowSessions in store
func UpdateWorkflowSession(ctx context.Context, s WorkflowSessions, rr ...*types.WorkflowSession) error {
	return s.UpdateWorkflowSession(ctx, rr...)
}

// UpsertWorkflowSession creates new or updates existing one or more WorkflowSessions in store
func UpsertWorkflowSession(ctx context.Context, s WorkflowSessions, rr ...*types.WorkflowSession) error {
	return s.UpsertWorkflowSession(ctx, rr...)
}

// DeleteWorkflowSession Deletes one or more WorkflowSessions from store
func DeleteWorkflowSession(ctx context.Context, s WorkflowSessions, rr ...*types.WorkflowSession) error {
	return s.DeleteWorkflowSession(ctx, rr...)
}

// DeleteWorkflowSessionByID Deletes WorkflowSession from store
func DeleteWorkflowSessionByID(ctx context.Context, s WorkflowSessions, ID uint64) error {
	return s.DeleteWorkflowSessionByID(ctx, ID)
}

// TruncateWorkflowSessions Deletes all WorkflowSessions from store
func TruncateWorkflowSessions(ctx context.Context, s WorkflowSessions) error {
	return s.TruncateWorkflowSessions(ctx)
}

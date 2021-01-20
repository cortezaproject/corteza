package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/automation_sessions.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/types"
)

type (
	AutomationSessions interface {
		SearchAutomationSessions(ctx context.Context, f types.SessionFilter) (types.SessionSet, types.SessionFilter, error)
		LookupAutomationSessionByID(ctx context.Context, id uint64) (*types.Session, error)

		CreateAutomationSession(ctx context.Context, rr ...*types.Session) error

		UpdateAutomationSession(ctx context.Context, rr ...*types.Session) error

		UpsertAutomationSession(ctx context.Context, rr ...*types.Session) error

		DeleteAutomationSession(ctx context.Context, rr ...*types.Session) error
		DeleteAutomationSessionByID(ctx context.Context, ID uint64) error

		TruncateAutomationSessions(ctx context.Context) error
	}
)

var _ *types.Session
var _ context.Context

// SearchAutomationSessions returns all matching AutomationSessions from store
func SearchAutomationSessions(ctx context.Context, s AutomationSessions, f types.SessionFilter) (types.SessionSet, types.SessionFilter, error) {
	return s.SearchAutomationSessions(ctx, f)
}

// LookupAutomationSessionByID searches for session by ID
//
// It returns session even if deleted
func LookupAutomationSessionByID(ctx context.Context, s AutomationSessions, id uint64) (*types.Session, error) {
	return s.LookupAutomationSessionByID(ctx, id)
}

// CreateAutomationSession creates one or more AutomationSessions in store
func CreateAutomationSession(ctx context.Context, s AutomationSessions, rr ...*types.Session) error {
	return s.CreateAutomationSession(ctx, rr...)
}

// UpdateAutomationSession updates one or more (existing) AutomationSessions in store
func UpdateAutomationSession(ctx context.Context, s AutomationSessions, rr ...*types.Session) error {
	return s.UpdateAutomationSession(ctx, rr...)
}

// UpsertAutomationSession creates new or updates existing one or more AutomationSessions in store
func UpsertAutomationSession(ctx context.Context, s AutomationSessions, rr ...*types.Session) error {
	return s.UpsertAutomationSession(ctx, rr...)
}

// DeleteAutomationSession Deletes one or more AutomationSessions from store
func DeleteAutomationSession(ctx context.Context, s AutomationSessions, rr ...*types.Session) error {
	return s.DeleteAutomationSession(ctx, rr...)
}

// DeleteAutomationSessionByID Deletes AutomationSession from store
func DeleteAutomationSessionByID(ctx context.Context, s AutomationSessions, ID uint64) error {
	return s.DeleteAutomationSessionByID(ctx, ID)
}

// TruncateAutomationSessions Deletes all AutomationSessions from store
func TruncateAutomationSessions(ctx context.Context, s AutomationSessions) error {
	return s.TruncateAutomationSessions(ctx)
}

package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/auth_sessions.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	AuthSessions interface {
		SearchAuthSessions(ctx context.Context, f types.AuthSessionFilter) (types.AuthSessionSet, types.AuthSessionFilter, error)
		LookupAuthSessionByID(ctx context.Context, id string) (*types.AuthSession, error)

		CreateAuthSession(ctx context.Context, rr ...*types.AuthSession) error

		UpdateAuthSession(ctx context.Context, rr ...*types.AuthSession) error

		UpsertAuthSession(ctx context.Context, rr ...*types.AuthSession) error

		DeleteAuthSession(ctx context.Context, rr ...*types.AuthSession) error
		DeleteAuthSessionByID(ctx context.Context, ID string) error

		TruncateAuthSessions(ctx context.Context) error

		// Additional custom functions

		// DeleteExpiredAuthSessions (custom function)
		DeleteExpiredAuthSessions(ctx context.Context) error

		// DeleteAuthSessionsByUserID (custom function)
		DeleteAuthSessionsByUserID(ctx context.Context, _userID uint64) error
	}
)

var _ *types.AuthSession
var _ context.Context

// SearchAuthSessions returns all matching AuthSessions from store
func SearchAuthSessions(ctx context.Context, s AuthSessions, f types.AuthSessionFilter) (types.AuthSessionSet, types.AuthSessionFilter, error) {
	return s.SearchAuthSessions(ctx, f)
}

// LookupAuthSessionByID
func LookupAuthSessionByID(ctx context.Context, s AuthSessions, id string) (*types.AuthSession, error) {
	return s.LookupAuthSessionByID(ctx, id)
}

// CreateAuthSession creates one or more AuthSessions in store
func CreateAuthSession(ctx context.Context, s AuthSessions, rr ...*types.AuthSession) error {
	return s.CreateAuthSession(ctx, rr...)
}

// UpdateAuthSession updates one or more (existing) AuthSessions in store
func UpdateAuthSession(ctx context.Context, s AuthSessions, rr ...*types.AuthSession) error {
	return s.UpdateAuthSession(ctx, rr...)
}

// UpsertAuthSession creates new or updates existing one or more AuthSessions in store
func UpsertAuthSession(ctx context.Context, s AuthSessions, rr ...*types.AuthSession) error {
	return s.UpsertAuthSession(ctx, rr...)
}

// DeleteAuthSession Deletes one or more AuthSessions from store
func DeleteAuthSession(ctx context.Context, s AuthSessions, rr ...*types.AuthSession) error {
	return s.DeleteAuthSession(ctx, rr...)
}

// DeleteAuthSessionByID Deletes AuthSession from store
func DeleteAuthSessionByID(ctx context.Context, s AuthSessions, ID string) error {
	return s.DeleteAuthSessionByID(ctx, ID)
}

// TruncateAuthSessions Deletes all AuthSessions from store
func TruncateAuthSessions(ctx context.Context, s AuthSessions) error {
	return s.TruncateAuthSessions(ctx)
}

func DeleteExpiredAuthSessions(ctx context.Context, s AuthSessions) error {
	return s.DeleteExpiredAuthSessions(ctx)
}

func DeleteAuthSessionsByUserID(ctx context.Context, s AuthSessions, _userID uint64) error {
	return s.DeleteAuthSessionsByUserID(ctx, _userID)
}

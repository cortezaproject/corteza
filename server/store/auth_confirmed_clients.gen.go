package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/auth_confirmed_clients.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	AuthConfirmedClients interface {
		SearchAuthConfirmedClients(ctx context.Context, f types.AuthConfirmedClientFilter) (types.AuthConfirmedClientSet, types.AuthConfirmedClientFilter, error)
		LookupAuthConfirmedClientByUserIDClientID(ctx context.Context, user_id uint64, client_id uint64) (*types.AuthConfirmedClient, error)

		CreateAuthConfirmedClient(ctx context.Context, rr ...*types.AuthConfirmedClient) error

		UpdateAuthConfirmedClient(ctx context.Context, rr ...*types.AuthConfirmedClient) error

		UpsertAuthConfirmedClient(ctx context.Context, rr ...*types.AuthConfirmedClient) error

		DeleteAuthConfirmedClient(ctx context.Context, rr ...*types.AuthConfirmedClient) error
		DeleteAuthConfirmedClientByUserIDClientID(ctx context.Context, userID uint64, clientID uint64) error

		TruncateAuthConfirmedClients(ctx context.Context) error
	}
)

var _ *types.AuthConfirmedClient
var _ context.Context

// SearchAuthConfirmedClients returns all matching AuthConfirmedClients from store
func SearchAuthConfirmedClients(ctx context.Context, s AuthConfirmedClients, f types.AuthConfirmedClientFilter) (types.AuthConfirmedClientSet, types.AuthConfirmedClientFilter, error) {
	return s.SearchAuthConfirmedClients(ctx, f)
}

// LookupAuthConfirmedClientByUserIDClientID
func LookupAuthConfirmedClientByUserIDClientID(ctx context.Context, s AuthConfirmedClients, user_id uint64, client_id uint64) (*types.AuthConfirmedClient, error) {
	return s.LookupAuthConfirmedClientByUserIDClientID(ctx, user_id, client_id)
}

// CreateAuthConfirmedClient creates one or more AuthConfirmedClients in store
func CreateAuthConfirmedClient(ctx context.Context, s AuthConfirmedClients, rr ...*types.AuthConfirmedClient) error {
	return s.CreateAuthConfirmedClient(ctx, rr...)
}

// UpdateAuthConfirmedClient updates one or more (existing) AuthConfirmedClients in store
func UpdateAuthConfirmedClient(ctx context.Context, s AuthConfirmedClients, rr ...*types.AuthConfirmedClient) error {
	return s.UpdateAuthConfirmedClient(ctx, rr...)
}

// UpsertAuthConfirmedClient creates new or updates existing one or more AuthConfirmedClients in store
func UpsertAuthConfirmedClient(ctx context.Context, s AuthConfirmedClients, rr ...*types.AuthConfirmedClient) error {
	return s.UpsertAuthConfirmedClient(ctx, rr...)
}

// DeleteAuthConfirmedClient Deletes one or more AuthConfirmedClients from store
func DeleteAuthConfirmedClient(ctx context.Context, s AuthConfirmedClients, rr ...*types.AuthConfirmedClient) error {
	return s.DeleteAuthConfirmedClient(ctx, rr...)
}

// DeleteAuthConfirmedClientByUserIDClientID Deletes AuthConfirmedClient from store
func DeleteAuthConfirmedClientByUserIDClientID(ctx context.Context, s AuthConfirmedClients, userID uint64, clientID uint64) error {
	return s.DeleteAuthConfirmedClientByUserIDClientID(ctx, userID, clientID)
}

// TruncateAuthConfirmedClients Deletes all AuthConfirmedClients from store
func TruncateAuthConfirmedClients(ctx context.Context, s AuthConfirmedClients) error {
	return s.TruncateAuthConfirmedClients(ctx)
}

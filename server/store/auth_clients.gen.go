package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/auth_clients.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	AuthClients interface {
		SearchAuthClients(ctx context.Context, f types.AuthClientFilter) (types.AuthClientSet, types.AuthClientFilter, error)
		LookupAuthClientByID(ctx context.Context, id uint64) (*types.AuthClient, error)
		LookupAuthClientByHandle(ctx context.Context, handle string) (*types.AuthClient, error)

		CreateAuthClient(ctx context.Context, rr ...*types.AuthClient) error

		UpdateAuthClient(ctx context.Context, rr ...*types.AuthClient) error

		UpsertAuthClient(ctx context.Context, rr ...*types.AuthClient) error

		DeleteAuthClient(ctx context.Context, rr ...*types.AuthClient) error
		DeleteAuthClientByID(ctx context.Context, ID uint64) error

		TruncateAuthClients(ctx context.Context) error
	}
)

var _ *types.AuthClient
var _ context.Context

// SearchAuthClients returns all matching AuthClients from store
func SearchAuthClients(ctx context.Context, s AuthClients, f types.AuthClientFilter) (types.AuthClientSet, types.AuthClientFilter, error) {
	return s.SearchAuthClients(ctx, f)
}

// LookupAuthClientByID searches for auth client by ID
//
// It returns auth client even if deleted
func LookupAuthClientByID(ctx context.Context, s AuthClients, id uint64) (*types.AuthClient, error) {
	return s.LookupAuthClientByID(ctx, id)
}

// LookupAuthClientByHandle searches for auth client by ID
//
// It returns auth client even if deleted
func LookupAuthClientByHandle(ctx context.Context, s AuthClients, handle string) (*types.AuthClient, error) {
	return s.LookupAuthClientByHandle(ctx, handle)
}

// CreateAuthClient creates one or more AuthClients in store
func CreateAuthClient(ctx context.Context, s AuthClients, rr ...*types.AuthClient) error {
	return s.CreateAuthClient(ctx, rr...)
}

// UpdateAuthClient updates one or more (existing) AuthClients in store
func UpdateAuthClient(ctx context.Context, s AuthClients, rr ...*types.AuthClient) error {
	return s.UpdateAuthClient(ctx, rr...)
}

// UpsertAuthClient creates new or updates existing one or more AuthClients in store
func UpsertAuthClient(ctx context.Context, s AuthClients, rr ...*types.AuthClient) error {
	return s.UpsertAuthClient(ctx, rr...)
}

// DeleteAuthClient Deletes one or more AuthClients from store
func DeleteAuthClient(ctx context.Context, s AuthClients, rr ...*types.AuthClient) error {
	return s.DeleteAuthClient(ctx, rr...)
}

// DeleteAuthClientByID Deletes AuthClient from store
func DeleteAuthClientByID(ctx context.Context, s AuthClients, ID uint64) error {
	return s.DeleteAuthClientByID(ctx, ID)
}

// TruncateAuthClients Deletes all AuthClients from store
func TruncateAuthClients(ctx context.Context, s AuthClients) error {
	return s.TruncateAuthClients(ctx)
}

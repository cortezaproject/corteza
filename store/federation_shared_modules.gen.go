package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/federation_shared_modules.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/federation/types"
)

type (
	FederationSharedModules interface {
		SearchFederationSharedModules(ctx context.Context, f types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error)
		LookupFederationSharedModuleByID(ctx context.Context, id uint64) (*types.SharedModule, error)

		CreateFederationSharedModule(ctx context.Context, rr ...*types.SharedModule) error

		UpdateFederationSharedModule(ctx context.Context, rr ...*types.SharedModule) error

		UpsertFederationSharedModule(ctx context.Context, rr ...*types.SharedModule) error

		DeleteFederationSharedModule(ctx context.Context, rr ...*types.SharedModule) error
		DeleteFederationSharedModuleByID(ctx context.Context, ID uint64) error

		TruncateFederationSharedModules(ctx context.Context) error
	}
)

var _ *types.SharedModule
var _ context.Context

// SearchFederationSharedModules returns all matching FederationSharedModules from store
func SearchFederationSharedModules(ctx context.Context, s FederationSharedModules, f types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error) {
	return s.SearchFederationSharedModules(ctx, f)
}

// LookupFederationSharedModuleByID searches for shared federation module by ID
//
// It returns shared federation module
func LookupFederationSharedModuleByID(ctx context.Context, s FederationSharedModules, id uint64) (*types.SharedModule, error) {
	return s.LookupFederationSharedModuleByID(ctx, id)
}

// CreateFederationSharedModule creates one or more FederationSharedModules in store
func CreateFederationSharedModule(ctx context.Context, s FederationSharedModules, rr ...*types.SharedModule) error {
	return s.CreateFederationSharedModule(ctx, rr...)
}

// UpdateFederationSharedModule updates one or more (existing) FederationSharedModules in store
func UpdateFederationSharedModule(ctx context.Context, s FederationSharedModules, rr ...*types.SharedModule) error {
	return s.UpdateFederationSharedModule(ctx, rr...)
}

// UpsertFederationSharedModule creates new or updates existing one or more FederationSharedModules in store
func UpsertFederationSharedModule(ctx context.Context, s FederationSharedModules, rr ...*types.SharedModule) error {
	return s.UpsertFederationSharedModule(ctx, rr...)
}

// DeleteFederationSharedModule Deletes one or more FederationSharedModules from store
func DeleteFederationSharedModule(ctx context.Context, s FederationSharedModules, rr ...*types.SharedModule) error {
	return s.DeleteFederationSharedModule(ctx, rr...)
}

// DeleteFederationSharedModuleByID Deletes FederationSharedModule from store
func DeleteFederationSharedModuleByID(ctx context.Context, s FederationSharedModules, ID uint64) error {
	return s.DeleteFederationSharedModuleByID(ctx, ID)
}

// TruncateFederationSharedModules Deletes all FederationSharedModules from store
func TruncateFederationSharedModules(ctx context.Context, s FederationSharedModules) error {
	return s.TruncateFederationSharedModules(ctx)
}

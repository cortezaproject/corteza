package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/federation_exposed_modules.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/federation/types"
)

type (
	FederationExposedModules interface {
		SearchFederationExposedModules(ctx context.Context, f types.ExposedModuleFilter) (types.ExposedModuleSet, types.ExposedModuleFilter, error)
		LookupFederationExposedModuleByID(ctx context.Context, id uint64) (*types.ExposedModule, error)

		CreateFederationExposedModule(ctx context.Context, rr ...*types.ExposedModule) error

		UpdateFederationExposedModule(ctx context.Context, rr ...*types.ExposedModule) error

		UpsertFederationExposedModule(ctx context.Context, rr ...*types.ExposedModule) error

		DeleteFederationExposedModule(ctx context.Context, rr ...*types.ExposedModule) error
		DeleteFederationExposedModuleByID(ctx context.Context, ID uint64) error

		TruncateFederationExposedModules(ctx context.Context) error
	}
)

var _ *types.ExposedModule
var _ context.Context

// SearchFederationExposedModules returns all matching FederationExposedModules from store
func SearchFederationExposedModules(ctx context.Context, s FederationExposedModules, f types.ExposedModuleFilter) (types.ExposedModuleSet, types.ExposedModuleFilter, error) {
	return s.SearchFederationExposedModules(ctx, f)
}

// LookupFederationExposedModuleByID searches for federation module by ID
//
// It returns federation module
func LookupFederationExposedModuleByID(ctx context.Context, s FederationExposedModules, id uint64) (*types.ExposedModule, error) {
	return s.LookupFederationExposedModuleByID(ctx, id)
}

// CreateFederationExposedModule creates one or more FederationExposedModules in store
func CreateFederationExposedModule(ctx context.Context, s FederationExposedModules, rr ...*types.ExposedModule) error {
	return s.CreateFederationExposedModule(ctx, rr...)
}

// UpdateFederationExposedModule updates one or more (existing) FederationExposedModules in store
func UpdateFederationExposedModule(ctx context.Context, s FederationExposedModules, rr ...*types.ExposedModule) error {
	return s.UpdateFederationExposedModule(ctx, rr...)
}

// UpsertFederationExposedModule creates new or updates existing one or more FederationExposedModules in store
func UpsertFederationExposedModule(ctx context.Context, s FederationExposedModules, rr ...*types.ExposedModule) error {
	return s.UpsertFederationExposedModule(ctx, rr...)
}

// DeleteFederationExposedModule Deletes one or more FederationExposedModules from store
func DeleteFederationExposedModule(ctx context.Context, s FederationExposedModules, rr ...*types.ExposedModule) error {
	return s.DeleteFederationExposedModule(ctx, rr...)
}

// DeleteFederationExposedModuleByID Deletes FederationExposedModule from store
func DeleteFederationExposedModuleByID(ctx context.Context, s FederationExposedModules, ID uint64) error {
	return s.DeleteFederationExposedModuleByID(ctx, ID)
}

// TruncateFederationExposedModules Deletes all FederationExposedModules from store
func TruncateFederationExposedModules(ctx context.Context, s FederationExposedModules) error {
	return s.TruncateFederationExposedModules(ctx)
}

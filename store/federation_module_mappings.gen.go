package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/federation_module_mappings.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/federation/types"
)

type (
	FederationModuleMappings interface {
		SearchFederationModuleMappings(ctx context.Context, f types.ModuleMappingFilter) (types.ModuleMappingSet, types.ModuleMappingFilter, error)
		LookupFederationModuleMappingByFederationModuleIDComposeModuleID(ctx context.Context, federation_module_id uint64, compose_module_id uint64) (*types.ModuleMapping, error)
		LookupFederationModuleMappingByFederationModuleID(ctx context.Context, federation_module_id uint64) (*types.ModuleMapping, error)

		CreateFederationModuleMapping(ctx context.Context, rr ...*types.ModuleMapping) error

		UpdateFederationModuleMapping(ctx context.Context, rr ...*types.ModuleMapping) error

		UpsertFederationModuleMapping(ctx context.Context, rr ...*types.ModuleMapping) error

		DeleteFederationModuleMapping(ctx context.Context, rr ...*types.ModuleMapping) error
		DeleteFederationModuleMappingBy(ctx context.Context) error

		TruncateFederationModuleMappings(ctx context.Context) error
	}
)

var _ *types.ModuleMapping
var _ context.Context

// SearchFederationModuleMappings returns all matching FederationModuleMappings from store
func SearchFederationModuleMappings(ctx context.Context, s FederationModuleMappings, f types.ModuleMappingFilter) (types.ModuleMappingSet, types.ModuleMappingFilter, error) {
	return s.SearchFederationModuleMappings(ctx, f)
}

// LookupFederationModuleMappingByFederationModuleIDComposeModuleID searches for module mapping by federation module id and compose module id
//
// It returns module mapping
func LookupFederationModuleMappingByFederationModuleIDComposeModuleID(ctx context.Context, s FederationModuleMappings, federation_module_id uint64, compose_module_id uint64) (*types.ModuleMapping, error) {
	return s.LookupFederationModuleMappingByFederationModuleIDComposeModuleID(ctx, federation_module_id, compose_module_id)
}

// LookupFederationModuleMappingByFederationModuleID searches for module mapping by federation module id
//
// It returns module mapping
func LookupFederationModuleMappingByFederationModuleID(ctx context.Context, s FederationModuleMappings, federation_module_id uint64) (*types.ModuleMapping, error) {
	return s.LookupFederationModuleMappingByFederationModuleID(ctx, federation_module_id)
}

// CreateFederationModuleMapping creates one or more FederationModuleMappings in store
func CreateFederationModuleMapping(ctx context.Context, s FederationModuleMappings, rr ...*types.ModuleMapping) error {
	return s.CreateFederationModuleMapping(ctx, rr...)
}

// UpdateFederationModuleMapping updates one or more (existing) FederationModuleMappings in store
func UpdateFederationModuleMapping(ctx context.Context, s FederationModuleMappings, rr ...*types.ModuleMapping) error {
	return s.UpdateFederationModuleMapping(ctx, rr...)
}

// UpsertFederationModuleMapping creates new or updates existing one or more FederationModuleMappings in store
func UpsertFederationModuleMapping(ctx context.Context, s FederationModuleMappings, rr ...*types.ModuleMapping) error {
	return s.UpsertFederationModuleMapping(ctx, rr...)
}

// DeleteFederationModuleMapping Deletes one or more FederationModuleMappings from store
func DeleteFederationModuleMapping(ctx context.Context, s FederationModuleMappings, rr ...*types.ModuleMapping) error {
	return s.DeleteFederationModuleMapping(ctx, rr...)
}

// DeleteFederationModuleMappingBy Deletes FederationModuleMapping from store
func DeleteFederationModuleMappingBy(ctx context.Context, s FederationModuleMappings) error {
	return s.DeleteFederationModuleMappingBy(ctx)
}

// TruncateFederationModuleMappings Deletes all FederationModuleMappings from store
func TruncateFederationModuleMappings(ctx context.Context, s FederationModuleMappings) error {
	return s.TruncateFederationModuleMappings(ctx)
}

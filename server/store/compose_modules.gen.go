package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/compose_modules.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeModules interface {
		SearchComposeModules(ctx context.Context, f types.ModuleFilter) (types.ModuleSet, types.ModuleFilter, error)
		LookupComposeModuleByNamespaceIDHandle(ctx context.Context, namespace_id uint64, handle string) (*types.Module, error)
		LookupComposeModuleByNamespaceIDName(ctx context.Context, namespace_id uint64, name string) (*types.Module, error)
		LookupComposeModuleByID(ctx context.Context, id uint64) (*types.Module, error)

		CreateComposeModule(ctx context.Context, rr ...*types.Module) error

		UpdateComposeModule(ctx context.Context, rr ...*types.Module) error

		UpsertComposeModule(ctx context.Context, rr ...*types.Module) error

		DeleteComposeModule(ctx context.Context, rr ...*types.Module) error
		DeleteComposeModuleByID(ctx context.Context, ID uint64) error

		TruncateComposeModules(ctx context.Context) error
	}
)

var _ *types.Module
var _ context.Context

// SearchComposeModules returns all matching ComposeModules from store
func SearchComposeModules(ctx context.Context, s ComposeModules, f types.ModuleFilter) (types.ModuleSet, types.ModuleFilter, error) {
	return s.SearchComposeModules(ctx, f)
}

// LookupComposeModuleByNamespaceIDHandle searches for compose module by handle (case-insensitive)
func LookupComposeModuleByNamespaceIDHandle(ctx context.Context, s ComposeModules, namespace_id uint64, handle string) (*types.Module, error) {
	return s.LookupComposeModuleByNamespaceIDHandle(ctx, namespace_id, handle)
}

// LookupComposeModuleByNamespaceIDName searches for compose module by name (case-insensitive)
func LookupComposeModuleByNamespaceIDName(ctx context.Context, s ComposeModules, namespace_id uint64, name string) (*types.Module, error) {
	return s.LookupComposeModuleByNamespaceIDName(ctx, namespace_id, name)
}

// LookupComposeModuleByID searches for compose module by ID
//
// It returns compose module even if deleted
func LookupComposeModuleByID(ctx context.Context, s ComposeModules, id uint64) (*types.Module, error) {
	return s.LookupComposeModuleByID(ctx, id)
}

// CreateComposeModule creates one or more ComposeModules in store
func CreateComposeModule(ctx context.Context, s ComposeModules, rr ...*types.Module) error {
	return s.CreateComposeModule(ctx, rr...)
}

// UpdateComposeModule updates one or more (existing) ComposeModules in store
func UpdateComposeModule(ctx context.Context, s ComposeModules, rr ...*types.Module) error {
	return s.UpdateComposeModule(ctx, rr...)
}

// UpsertComposeModule creates new or updates existing one or more ComposeModules in store
func UpsertComposeModule(ctx context.Context, s ComposeModules, rr ...*types.Module) error {
	return s.UpsertComposeModule(ctx, rr...)
}

// DeleteComposeModule Deletes one or more ComposeModules from store
func DeleteComposeModule(ctx context.Context, s ComposeModules, rr ...*types.Module) error {
	return s.DeleteComposeModule(ctx, rr...)
}

// DeleteComposeModuleByID Deletes ComposeModule from store
func DeleteComposeModuleByID(ctx context.Context, s ComposeModules, ID uint64) error {
	return s.DeleteComposeModuleByID(ctx, ID)
}

// TruncateComposeModules Deletes all ComposeModules from store
func TruncateComposeModules(ctx context.Context, s ComposeModules) error {
	return s.TruncateComposeModules(ctx)
}

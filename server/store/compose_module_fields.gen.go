package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/compose_module_fields.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/compose/types"
)

type (
	ComposeModuleFields interface {
		SearchComposeModuleFields(ctx context.Context, f types.ModuleFieldFilter) (types.ModuleFieldSet, types.ModuleFieldFilter, error)
		LookupComposeModuleFieldByModuleIDName(ctx context.Context, module_id uint64, name string) (*types.ModuleField, error)

		CreateComposeModuleField(ctx context.Context, rr ...*types.ModuleField) error

		UpdateComposeModuleField(ctx context.Context, rr ...*types.ModuleField) error

		UpsertComposeModuleField(ctx context.Context, rr ...*types.ModuleField) error

		DeleteComposeModuleField(ctx context.Context, rr ...*types.ModuleField) error
		DeleteComposeModuleFieldByID(ctx context.Context, ID uint64) error

		TruncateComposeModuleFields(ctx context.Context) error
	}
)

var _ *types.ModuleField
var _ context.Context

// SearchComposeModuleFields returns all matching ComposeModuleFields from store
func SearchComposeModuleFields(ctx context.Context, s ComposeModuleFields, f types.ModuleFieldFilter) (types.ModuleFieldSet, types.ModuleFieldFilter, error) {
	return s.SearchComposeModuleFields(ctx, f)
}

// LookupComposeModuleFieldByModuleIDName searches for compose module field by name (case-insensitive)
func LookupComposeModuleFieldByModuleIDName(ctx context.Context, s ComposeModuleFields, module_id uint64, name string) (*types.ModuleField, error) {
	return s.LookupComposeModuleFieldByModuleIDName(ctx, module_id, name)
}

// CreateComposeModuleField creates one or more ComposeModuleFields in store
func CreateComposeModuleField(ctx context.Context, s ComposeModuleFields, rr ...*types.ModuleField) error {
	return s.CreateComposeModuleField(ctx, rr...)
}

// UpdateComposeModuleField updates one or more (existing) ComposeModuleFields in store
func UpdateComposeModuleField(ctx context.Context, s ComposeModuleFields, rr ...*types.ModuleField) error {
	return s.UpdateComposeModuleField(ctx, rr...)
}

// UpsertComposeModuleField creates new or updates existing one or more ComposeModuleFields in store
func UpsertComposeModuleField(ctx context.Context, s ComposeModuleFields, rr ...*types.ModuleField) error {
	return s.UpsertComposeModuleField(ctx, rr...)
}

// DeleteComposeModuleField Deletes one or more ComposeModuleFields from store
func DeleteComposeModuleField(ctx context.Context, s ComposeModuleFields, rr ...*types.ModuleField) error {
	return s.DeleteComposeModuleField(ctx, rr...)
}

// DeleteComposeModuleFieldByID Deletes ComposeModuleField from store
func DeleteComposeModuleFieldByID(ctx context.Context, s ComposeModuleFields, ID uint64) error {
	return s.DeleteComposeModuleFieldByID(ctx, ID)
}

// TruncateComposeModuleFields Deletes all ComposeModuleFields from store
func TruncateComposeModuleFields(ctx context.Context, s ComposeModuleFields) error {
	return s.TruncateComposeModuleFields(ctx)
}

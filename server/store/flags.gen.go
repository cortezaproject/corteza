package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/flags.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/flag/types"
)

type (
	Flags interface {
		SearchFlags(ctx context.Context, f types.FlagFilter) (types.FlagSet, types.FlagFilter, error)
		LookupFlagByKindResourceIDName(ctx context.Context, kind string, resource_id uint64, name string) (*types.Flag, error)
		LookupFlagByKindResourceID(ctx context.Context, kind string, resource_id uint64) (*types.Flag, error)
		LookupFlagByKindResourceIDOwnedBy(ctx context.Context, kind string, resource_id uint64, owned_by uint64) (*types.Flag, error)
		LookupFlagByKindResourceIDOwnedByName(ctx context.Context, kind string, resource_id uint64, owned_by uint64, name string) (*types.Flag, error)

		CreateFlag(ctx context.Context, rr ...*types.Flag) error

		UpdateFlag(ctx context.Context, rr ...*types.Flag) error

		UpsertFlag(ctx context.Context, rr ...*types.Flag) error

		DeleteFlag(ctx context.Context, rr ...*types.Flag) error
		DeleteFlagByKindResourceIDOwnedByName(ctx context.Context, kind string, resourceID uint64, ownedBy uint64, name string) error

		TruncateFlags(ctx context.Context) error
	}
)

var _ *types.Flag
var _ context.Context

// SearchFlags returns all matching Flags from store
func SearchFlags(ctx context.Context, s Flags, f types.FlagFilter) (types.FlagSet, types.FlagFilter, error) {
	return s.SearchFlags(ctx, f)
}

// LookupFlagByKindResourceIDName Flag lookup by kind, resource, name
func LookupFlagByKindResourceIDName(ctx context.Context, s Flags, kind string, resource_id uint64, name string) (*types.Flag, error) {
	return s.LookupFlagByKindResourceIDName(ctx, kind, resource_id, name)
}

// LookupFlagByKindResourceID Flag lookup by kind, resource
func LookupFlagByKindResourceID(ctx context.Context, s Flags, kind string, resource_id uint64) (*types.Flag, error) {
	return s.LookupFlagByKindResourceID(ctx, kind, resource_id)
}

// LookupFlagByKindResourceIDOwnedBy Flag lookup by kind, resource, owner
func LookupFlagByKindResourceIDOwnedBy(ctx context.Context, s Flags, kind string, resource_id uint64, owned_by uint64) (*types.Flag, error) {
	return s.LookupFlagByKindResourceIDOwnedBy(ctx, kind, resource_id, owned_by)
}

// LookupFlagByKindResourceIDOwnedByName Flag lookup by kind, resource, owner, name
func LookupFlagByKindResourceIDOwnedByName(ctx context.Context, s Flags, kind string, resource_id uint64, owned_by uint64, name string) (*types.Flag, error) {
	return s.LookupFlagByKindResourceIDOwnedByName(ctx, kind, resource_id, owned_by, name)
}

// CreateFlag creates one or more Flags in store
func CreateFlag(ctx context.Context, s Flags, rr ...*types.Flag) error {
	return s.CreateFlag(ctx, rr...)
}

// UpdateFlag updates one or more (existing) Flags in store
func UpdateFlag(ctx context.Context, s Flags, rr ...*types.Flag) error {
	return s.UpdateFlag(ctx, rr...)
}

// UpsertFlag creates new or updates existing one or more Flags in store
func UpsertFlag(ctx context.Context, s Flags, rr ...*types.Flag) error {
	return s.UpsertFlag(ctx, rr...)
}

// DeleteFlag Deletes one or more Flags from store
func DeleteFlag(ctx context.Context, s Flags, rr ...*types.Flag) error {
	return s.DeleteFlag(ctx, rr...)
}

// DeleteFlagByKindResourceIDOwnedByName Deletes Flag from store
func DeleteFlagByKindResourceIDOwnedByName(ctx context.Context, s Flags, kind string, resourceID uint64, ownedBy uint64, name string) error {
	return s.DeleteFlagByKindResourceIDOwnedByName(ctx, kind, resourceID, ownedBy, name)
}

// TruncateFlags Deletes all Flags from store
func TruncateFlags(ctx context.Context, s Flags) error {
	return s.TruncateFlags(ctx)
}

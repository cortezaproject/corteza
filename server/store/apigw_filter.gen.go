package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/apigw_filter.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	ApigwFilters interface {
		SearchApigwFilters(ctx context.Context, f types.ApigwFilterFilter) (types.ApigwFilterSet, types.ApigwFilterFilter, error)
		LookupApigwFilterByID(ctx context.Context, id uint64) (*types.ApigwFilter, error)
		LookupApigwFilterByRoute(ctx context.Context, route string) (*types.ApigwFilter, error)

		CreateApigwFilter(ctx context.Context, rr ...*types.ApigwFilter) error

		UpdateApigwFilter(ctx context.Context, rr ...*types.ApigwFilter) error

		DeleteApigwFilter(ctx context.Context, rr ...*types.ApigwFilter) error
		DeleteApigwFilterByID(ctx context.Context, ID uint64) error

		TruncateApigwFilters(ctx context.Context) error
	}
)

var _ *types.ApigwFilter
var _ context.Context

// SearchApigwFilters returns all matching ApigwFilters from store
func SearchApigwFilters(ctx context.Context, s ApigwFilters, f types.ApigwFilterFilter) (types.ApigwFilterSet, types.ApigwFilterFilter, error) {
	return s.SearchApigwFilters(ctx, f)
}

// LookupApigwFilterByID searches for filter by ID
func LookupApigwFilterByID(ctx context.Context, s ApigwFilters, id uint64) (*types.ApigwFilter, error) {
	return s.LookupApigwFilterByID(ctx, id)
}

// LookupApigwFilterByRoute searches for filter by route
func LookupApigwFilterByRoute(ctx context.Context, s ApigwFilters, route string) (*types.ApigwFilter, error) {
	return s.LookupApigwFilterByRoute(ctx, route)
}

// CreateApigwFilter creates one or more ApigwFilters in store
func CreateApigwFilter(ctx context.Context, s ApigwFilters, rr ...*types.ApigwFilter) error {
	return s.CreateApigwFilter(ctx, rr...)
}

// UpdateApigwFilter updates one or more (existing) ApigwFilters in store
func UpdateApigwFilter(ctx context.Context, s ApigwFilters, rr ...*types.ApigwFilter) error {
	return s.UpdateApigwFilter(ctx, rr...)
}

// DeleteApigwFilter Deletes one or more ApigwFilters from store
func DeleteApigwFilter(ctx context.Context, s ApigwFilters, rr ...*types.ApigwFilter) error {
	return s.DeleteApigwFilter(ctx, rr...)
}

// DeleteApigwFilterByID Deletes ApigwFilter from store
func DeleteApigwFilterByID(ctx context.Context, s ApigwFilters, ID uint64) error {
	return s.DeleteApigwFilterByID(ctx, ID)
}

// TruncateApigwFilters Deletes all ApigwFilters from store
func TruncateApigwFilters(ctx context.Context, s ApigwFilters) error {
	return s.TruncateApigwFilters(ctx)
}

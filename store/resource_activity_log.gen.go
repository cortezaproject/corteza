package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/resource_activity_log.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/discovery/types"
)

type (
	ResourceActivityLogs interface {
		SearchResourceActivityLogs(ctx context.Context, f types.ResourceActivityFilter) (types.ResourceActivitySet, types.ResourceActivityFilter, error)

		CreateResourceActivityLog(ctx context.Context, rr ...*types.ResourceActivity) error

		TruncateResourceActivityLogs(ctx context.Context) error
	}
)

var _ *types.ResourceActivity
var _ context.Context

// SearchResourceActivityLogs returns all matching ResourceActivityLogs from store
func SearchResourceActivityLogs(ctx context.Context, s ResourceActivityLogs, f types.ResourceActivityFilter) (types.ResourceActivitySet, types.ResourceActivityFilter, error) {
	return s.SearchResourceActivityLogs(ctx, f)
}

// CreateResourceActivityLog creates one or more ResourceActivityLogs in store
func CreateResourceActivityLog(ctx context.Context, s ResourceActivityLogs, rr ...*types.ResourceActivity) error {
	return s.CreateResourceActivityLog(ctx, rr...)
}

// TruncateResourceActivityLogs Deletes all ResourceActivityLogs from store
func TruncateResourceActivityLogs(ctx context.Context, s ResourceActivityLogs) error {
	return s.TruncateResourceActivityLogs(ctx)
}

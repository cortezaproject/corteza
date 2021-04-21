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
	"github.com/cortezaproject/corteza-server/pkg/discovery"
)

type (
	ResourceActivityLogs interface {
		SearchResourceActivityLogs(ctx context.Context, f discovery.Filter) (discovery.ResourceActivitySet, discovery.Filter, error)
		LookupResourceActivityLogByID(ctx context.Context, id uint64) (*discovery.ResourceActivity, error)

		CreateResourceActivityLog(ctx context.Context, rr ...*discovery.ResourceActivity) error

		UpdateResourceActivityLog(ctx context.Context, rr ...*discovery.ResourceActivity) error

		UpsertResourceActivityLog(ctx context.Context, rr ...*discovery.ResourceActivity) error

		DeleteResourceActivityLog(ctx context.Context, rr ...*discovery.ResourceActivity) error
		DeleteResourceActivityLogByID(ctx context.Context, ID uint64) error

		TruncateResourceActivityLogs(ctx context.Context) error
	}
)

var _ *discovery.ResourceActivity
var _ context.Context

// SearchResourceActivityLogs returns all matching ResourceActivityLogs from store
func SearchResourceActivityLogs(ctx context.Context, s ResourceActivityLogs, f discovery.Filter) (discovery.ResourceActivitySet, discovery.Filter, error) {
	return s.SearchResourceActivityLogs(ctx, f)
}

// LookupResourceActivityLogByID searches for corteza resource activity by ID
// It returns corteza resource activity even if deleted
func LookupResourceActivityLogByID(ctx context.Context, s ResourceActivityLogs, id uint64) (*discovery.ResourceActivity, error) {
	return s.LookupResourceActivityLogByID(ctx, id)
}

// CreateResourceActivityLog creates one or more ResourceActivityLogs in store
func CreateResourceActivityLog(ctx context.Context, s ResourceActivityLogs, rr ...*discovery.ResourceActivity) error {
	return s.CreateResourceActivityLog(ctx, rr...)
}

// UpdateResourceActivityLog updates one or more (existing) ResourceActivityLogs in store
func UpdateResourceActivityLog(ctx context.Context, s ResourceActivityLogs, rr ...*discovery.ResourceActivity) error {
	return s.UpdateResourceActivityLog(ctx, rr...)
}

// UpsertResourceActivityLog creates new or updates existing one or more ResourceActivityLogs in store
func UpsertResourceActivityLog(ctx context.Context, s ResourceActivityLogs, rr ...*discovery.ResourceActivity) error {
	return s.UpsertResourceActivityLog(ctx, rr...)
}

// DeleteResourceActivityLog Deletes one or more ResourceActivityLogs from store
func DeleteResourceActivityLog(ctx context.Context, s ResourceActivityLogs, rr ...*discovery.ResourceActivity) error {
	return s.DeleteResourceActivityLog(ctx, rr...)
}

// DeleteResourceActivityLogByID Deletes ResourceActivityLog from store
func DeleteResourceActivityLogByID(ctx context.Context, s ResourceActivityLogs, ID uint64) error {
	return s.DeleteResourceActivityLogByID(ctx, ID)
}

// TruncateResourceActivityLogs Deletes all ResourceActivityLogs from store
func TruncateResourceActivityLogs(ctx context.Context, s ResourceActivityLogs) error {
	return s.TruncateResourceActivityLogs(ctx)
}

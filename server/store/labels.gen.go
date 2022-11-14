package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/labels.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/label/types"
)

type (
	Labels interface {
		SearchLabels(ctx context.Context, f types.LabelFilter) (types.LabelSet, types.LabelFilter, error)
		LookupLabelByKindResourceIDName(ctx context.Context, kind string, resource_id uint64, name string) (*types.Label, error)

		CreateLabel(ctx context.Context, rr ...*types.Label) error

		UpdateLabel(ctx context.Context, rr ...*types.Label) error

		UpsertLabel(ctx context.Context, rr ...*types.Label) error

		DeleteLabel(ctx context.Context, rr ...*types.Label) error
		DeleteLabelByKindResourceIDName(ctx context.Context, kind string, resourceID uint64, name string) error

		TruncateLabels(ctx context.Context) error

		// Additional custom functions

		// DeleteExtraLabels (custom function)
		DeleteExtraLabels(ctx context.Context, _kind string, _resourceID uint64, _names ...string) error
	}
)

var _ *types.Label
var _ context.Context

// SearchLabels returns all matching Labels from store
func SearchLabels(ctx context.Context, s Labels, f types.LabelFilter) (types.LabelSet, types.LabelFilter, error) {
	return s.SearchLabels(ctx, f)
}

// LookupLabelByKindResourceIDName Label lookup by kind, resource, name
func LookupLabelByKindResourceIDName(ctx context.Context, s Labels, kind string, resource_id uint64, name string) (*types.Label, error) {
	return s.LookupLabelByKindResourceIDName(ctx, kind, resource_id, name)
}

// CreateLabel creates one or more Labels in store
func CreateLabel(ctx context.Context, s Labels, rr ...*types.Label) error {
	return s.CreateLabel(ctx, rr...)
}

// UpdateLabel updates one or more (existing) Labels in store
func UpdateLabel(ctx context.Context, s Labels, rr ...*types.Label) error {
	return s.UpdateLabel(ctx, rr...)
}

// UpsertLabel creates new or updates existing one or more Labels in store
func UpsertLabel(ctx context.Context, s Labels, rr ...*types.Label) error {
	return s.UpsertLabel(ctx, rr...)
}

// DeleteLabel Deletes one or more Labels from store
func DeleteLabel(ctx context.Context, s Labels, rr ...*types.Label) error {
	return s.DeleteLabel(ctx, rr...)
}

// DeleteLabelByKindResourceIDName Deletes Label from store
func DeleteLabelByKindResourceIDName(ctx context.Context, s Labels, kind string, resourceID uint64, name string) error {
	return s.DeleteLabelByKindResourceIDName(ctx, kind, resourceID, name)
}

// TruncateLabels Deletes all Labels from store
func TruncateLabels(ctx context.Context, s Labels) error {
	return s.TruncateLabels(ctx)
}

func DeleteExtraLabels(ctx context.Context, s Labels, _kind string, _resourceID uint64, _names ...string) error {
	return s.DeleteExtraLabels(ctx, _kind, _resourceID, _names...)
}

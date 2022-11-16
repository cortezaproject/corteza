package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/resource_translation.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/system/types"
	"golang.org/x/text/language"
)

type (
	ResourceTranslations interface {
		SearchResourceTranslations(ctx context.Context, f types.ResourceTranslationFilter) (types.ResourceTranslationSet, types.ResourceTranslationFilter, error)
		LookupResourceTranslationByID(ctx context.Context, id uint64) (*types.ResourceTranslation, error)

		CreateResourceTranslation(ctx context.Context, rr ...*types.ResourceTranslation) error

		UpdateResourceTranslation(ctx context.Context, rr ...*types.ResourceTranslation) error

		UpsertResourceTranslation(ctx context.Context, rr ...*types.ResourceTranslation) error

		DeleteResourceTranslation(ctx context.Context, rr ...*types.ResourceTranslation) error
		DeleteResourceTranslationByID(ctx context.Context, ID uint64) error

		TruncateResourceTranslations(ctx context.Context) error

		// Additional custom functions

		// TransformResource (custom function)
		TransformResource(ctx context.Context, _lang language.Tag) (map[string]map[string]*locale.ResourceTranslation, error)
	}
)

var _ *types.ResourceTranslation
var _ context.Context

// SearchResourceTranslations returns all matching ResourceTranslations from store
func SearchResourceTranslations(ctx context.Context, s ResourceTranslations, f types.ResourceTranslationFilter) (types.ResourceTranslationSet, types.ResourceTranslationFilter, error) {
	return s.SearchResourceTranslations(ctx, f)
}

// LookupResourceTranslationByID searches for resource translation by ID
// It also returns deleted resource translations.
func LookupResourceTranslationByID(ctx context.Context, s ResourceTranslations, id uint64) (*types.ResourceTranslation, error) {
	return s.LookupResourceTranslationByID(ctx, id)
}

// CreateResourceTranslation creates one or more ResourceTranslations in store
func CreateResourceTranslation(ctx context.Context, s ResourceTranslations, rr ...*types.ResourceTranslation) error {
	return s.CreateResourceTranslation(ctx, rr...)
}

// UpdateResourceTranslation updates one or more (existing) ResourceTranslations in store
func UpdateResourceTranslation(ctx context.Context, s ResourceTranslations, rr ...*types.ResourceTranslation) error {
	return s.UpdateResourceTranslation(ctx, rr...)
}

// UpsertResourceTranslation creates new or updates existing one or more ResourceTranslations in store
func UpsertResourceTranslation(ctx context.Context, s ResourceTranslations, rr ...*types.ResourceTranslation) error {
	return s.UpsertResourceTranslation(ctx, rr...)
}

// DeleteResourceTranslation Deletes one or more ResourceTranslations from store
func DeleteResourceTranslation(ctx context.Context, s ResourceTranslations, rr ...*types.ResourceTranslation) error {
	return s.DeleteResourceTranslation(ctx, rr...)
}

// DeleteResourceTranslationByID Deletes ResourceTranslation from store
func DeleteResourceTranslationByID(ctx context.Context, s ResourceTranslations, ID uint64) error {
	return s.DeleteResourceTranslationByID(ctx, ID)
}

// TruncateResourceTranslations Deletes all ResourceTranslations from store
func TruncateResourceTranslations(ctx context.Context, s ResourceTranslations) error {
	return s.TruncateResourceTranslations(ctx)
}

func TransformResource(ctx context.Context, s ResourceTranslations, _lang language.Tag) (map[string]map[string]*locale.ResourceTranslation, error) {
	return s.TransformResource(ctx, _lang)
}

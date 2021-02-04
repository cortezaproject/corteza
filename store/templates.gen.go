package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/templates.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Templates interface {
		SearchTemplates(ctx context.Context, f types.TemplateFilter) (types.TemplateSet, types.TemplateFilter, error)
		LookupTemplateByID(ctx context.Context, id uint64) (*types.Template, error)
		LookupTemplateByHandle(ctx context.Context, handle string) (*types.Template, error)

		CreateTemplate(ctx context.Context, rr ...*types.Template) error

		UpdateTemplate(ctx context.Context, rr ...*types.Template) error

		UpsertTemplate(ctx context.Context, rr ...*types.Template) error

		DeleteTemplate(ctx context.Context, rr ...*types.Template) error
		DeleteTemplateByID(ctx context.Context, ID uint64) error

		TruncateTemplates(ctx context.Context) error
	}
)

var _ *types.Template
var _ context.Context

// SearchTemplates returns all matching Templates from store
func SearchTemplates(ctx context.Context, s Templates, f types.TemplateFilter) (types.TemplateSet, types.TemplateFilter, error) {
	return s.SearchTemplates(ctx, f)
}

// LookupTemplateByID searches for template by ID
//
// It also returns deleted templates.
func LookupTemplateByID(ctx context.Context, s Templates, id uint64) (*types.Template, error) {
	return s.LookupTemplateByID(ctx, id)
}

// LookupTemplateByHandle searches for template by the handle
//
// It returns only valid templates (not deleted)
func LookupTemplateByHandle(ctx context.Context, s Templates, handle string) (*types.Template, error) {
	return s.LookupTemplateByHandle(ctx, handle)
}

// CreateTemplate creates one or more Templates in store
func CreateTemplate(ctx context.Context, s Templates, rr ...*types.Template) error {
	return s.CreateTemplate(ctx, rr...)
}

// UpdateTemplate updates one or more (existing) Templates in store
func UpdateTemplate(ctx context.Context, s Templates, rr ...*types.Template) error {
	return s.UpdateTemplate(ctx, rr...)
}

// UpsertTemplate creates new or updates existing one or more Templates in store
func UpsertTemplate(ctx context.Context, s Templates, rr ...*types.Template) error {
	return s.UpsertTemplate(ctx, rr...)
}

// DeleteTemplate Deletes one or more Templates from store
func DeleteTemplate(ctx context.Context, s Templates, rr ...*types.Template) error {
	return s.DeleteTemplate(ctx, rr...)
}

// DeleteTemplateByID Deletes Template from store
func DeleteTemplateByID(ctx context.Context, s Templates, ID uint64) error {
	return s.DeleteTemplateByID(ctx, ID)
}

// TruncateTemplates Deletes all Templates from store
func TruncateTemplates(ctx context.Context, s Templates) error {
	return s.TruncateTemplates(ctx)
}

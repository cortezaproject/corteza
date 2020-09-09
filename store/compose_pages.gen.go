package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/compose_pages.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposePages interface {
		SearchComposePages(ctx context.Context, f types.PageFilter) (types.PageSet, types.PageFilter, error)
		LookupComposePageByNamespaceIDHandle(ctx context.Context, namespace_id uint64, handle string) (*types.Page, error)
		LookupComposePageByNamespaceIDModuleID(ctx context.Context, namespace_id uint64, module_id uint64) (*types.Page, error)
		LookupComposePageByID(ctx context.Context, id uint64) (*types.Page, error)

		CreateComposePage(ctx context.Context, rr ...*types.Page) error

		UpdateComposePage(ctx context.Context, rr ...*types.Page) error

		UpsertComposePage(ctx context.Context, rr ...*types.Page) error

		DeleteComposePage(ctx context.Context, rr ...*types.Page) error
		DeleteComposePageByID(ctx context.Context, ID uint64) error

		TruncateComposePages(ctx context.Context) error

		// Additional custom functions

		// ReorderComposePages (custom function)
		ReorderComposePages(ctx context.Context, _namespaceID uint64, _parentID uint64, _pageIDs []uint64) error
	}
)

var _ *types.Page
var _ context.Context

// SearchComposePages returns all matching ComposePages from store
func SearchComposePages(ctx context.Context, s ComposePages, f types.PageFilter) (types.PageSet, types.PageFilter, error) {
	return s.SearchComposePages(ctx, f)
}

// LookupComposePageByNamespaceIDHandle searches for page by handle (case-insensitive)
func LookupComposePageByNamespaceIDHandle(ctx context.Context, s ComposePages, namespace_id uint64, handle string) (*types.Page, error) {
	return s.LookupComposePageByNamespaceIDHandle(ctx, namespace_id, handle)
}

// LookupComposePageByNamespaceIDModuleID searches for page by moduleID
func LookupComposePageByNamespaceIDModuleID(ctx context.Context, s ComposePages, namespace_id uint64, module_id uint64) (*types.Page, error) {
	return s.LookupComposePageByNamespaceIDModuleID(ctx, namespace_id, module_id)
}

// LookupComposePageByID searches for compose page by ID
//
// It returns compose page even if deleted
func LookupComposePageByID(ctx context.Context, s ComposePages, id uint64) (*types.Page, error) {
	return s.LookupComposePageByID(ctx, id)
}

// CreateComposePage creates one or more ComposePages in store
func CreateComposePage(ctx context.Context, s ComposePages, rr ...*types.Page) error {
	return s.CreateComposePage(ctx, rr...)
}

// UpdateComposePage updates one or more (existing) ComposePages in store
func UpdateComposePage(ctx context.Context, s ComposePages, rr ...*types.Page) error {
	return s.UpdateComposePage(ctx, rr...)
}

// UpsertComposePage creates new or updates existing one or more ComposePages in store
func UpsertComposePage(ctx context.Context, s ComposePages, rr ...*types.Page) error {
	return s.UpsertComposePage(ctx, rr...)
}

// DeleteComposePage Deletes one or more ComposePages from store
func DeleteComposePage(ctx context.Context, s ComposePages, rr ...*types.Page) error {
	return s.DeleteComposePage(ctx, rr...)
}

// DeleteComposePageByID Deletes ComposePage from store
func DeleteComposePageByID(ctx context.Context, s ComposePages, ID uint64) error {
	return s.DeleteComposePageByID(ctx, ID)
}

// TruncateComposePages Deletes all ComposePages from store
func TruncateComposePages(ctx context.Context, s ComposePages) error {
	return s.TruncateComposePages(ctx)
}

func ReorderComposePages(ctx context.Context, s ComposePages, _namespaceID uint64, _parentID uint64, _pageIDs []uint64) error {
	return s.ReorderComposePages(ctx, _namespaceID, _parentID, _pageIDs)
}

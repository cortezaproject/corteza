package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/compose_namespaces.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeNamespaces interface {
		SearchComposeNamespaces(ctx context.Context, f types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error)
		LookupComposeNamespaceBySlug(ctx context.Context, slug string) (*types.Namespace, error)
		LookupComposeNamespaceByID(ctx context.Context, id uint64) (*types.Namespace, error)

		CreateComposeNamespace(ctx context.Context, rr ...*types.Namespace) error

		UpdateComposeNamespace(ctx context.Context, rr ...*types.Namespace) error

		UpsertComposeNamespace(ctx context.Context, rr ...*types.Namespace) error

		DeleteComposeNamespace(ctx context.Context, rr ...*types.Namespace) error
		DeleteComposeNamespaceByID(ctx context.Context, ID uint64) error

		TruncateComposeNamespaces(ctx context.Context) error
	}
)

var _ *types.Namespace
var _ context.Context

// SearchComposeNamespaces returns all matching ComposeNamespaces from store
func SearchComposeNamespaces(ctx context.Context, s ComposeNamespaces, f types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error) {
	return s.SearchComposeNamespaces(ctx, f)
}

// LookupComposeNamespaceBySlug searches for namespace by slug (case-insensitive)
func LookupComposeNamespaceBySlug(ctx context.Context, s ComposeNamespaces, slug string) (*types.Namespace, error) {
	return s.LookupComposeNamespaceBySlug(ctx, slug)
}

// LookupComposeNamespaceByID searches for compose namespace by ID
//
// It returns compose namespace even if deleted
func LookupComposeNamespaceByID(ctx context.Context, s ComposeNamespaces, id uint64) (*types.Namespace, error) {
	return s.LookupComposeNamespaceByID(ctx, id)
}

// CreateComposeNamespace creates one or more ComposeNamespaces in store
func CreateComposeNamespace(ctx context.Context, s ComposeNamespaces, rr ...*types.Namespace) error {
	return s.CreateComposeNamespace(ctx, rr...)
}

// UpdateComposeNamespace updates one or more (existing) ComposeNamespaces in store
func UpdateComposeNamespace(ctx context.Context, s ComposeNamespaces, rr ...*types.Namespace) error {
	return s.UpdateComposeNamespace(ctx, rr...)
}

// UpsertComposeNamespace creates new or updates existing one or more ComposeNamespaces in store
func UpsertComposeNamespace(ctx context.Context, s ComposeNamespaces, rr ...*types.Namespace) error {
	return s.UpsertComposeNamespace(ctx, rr...)
}

// DeleteComposeNamespace Deletes one or more ComposeNamespaces from store
func DeleteComposeNamespace(ctx context.Context, s ComposeNamespaces, rr ...*types.Namespace) error {
	return s.DeleteComposeNamespace(ctx, rr...)
}

// DeleteComposeNamespaceByID Deletes ComposeNamespace from store
func DeleteComposeNamespaceByID(ctx context.Context, s ComposeNamespaces, ID uint64) error {
	return s.DeleteComposeNamespaceByID(ctx, ID)
}

// TruncateComposeNamespaces Deletes all ComposeNamespaces from store
func TruncateComposeNamespaces(ctx context.Context, s ComposeNamespaces) error {
	return s.TruncateComposeNamespaces(ctx)
}

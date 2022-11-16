package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/credentials.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	Credentials interface {
		SearchCredentials(ctx context.Context, f types.CredentialsFilter) (types.CredentialsSet, types.CredentialsFilter, error)
		LookupCredentialsByID(ctx context.Context, id uint64) (*types.Credentials, error)

		CreateCredentials(ctx context.Context, rr ...*types.Credentials) error

		UpdateCredentials(ctx context.Context, rr ...*types.Credentials) error

		UpsertCredentials(ctx context.Context, rr ...*types.Credentials) error

		DeleteCredentials(ctx context.Context, rr ...*types.Credentials) error
		DeleteCredentialsByID(ctx context.Context, ID uint64) error

		TruncateCredentials(ctx context.Context) error
	}
)

var _ *types.Credentials
var _ context.Context

// SearchCredentials returns all matching Credentials from store
func SearchCredentials(ctx context.Context, s Credentials, f types.CredentialsFilter) (types.CredentialsSet, types.CredentialsFilter, error) {
	return s.SearchCredentials(ctx, f)
}

// LookupCredentialsByID searches for credentials by ID
//
// It returns credentials even if deleted
func LookupCredentialsByID(ctx context.Context, s Credentials, id uint64) (*types.Credentials, error) {
	return s.LookupCredentialsByID(ctx, id)
}

// CreateCredentials creates one or more Credentials in store
func CreateCredentials(ctx context.Context, s Credentials, rr ...*types.Credentials) error {
	return s.CreateCredentials(ctx, rr...)
}

// UpdateCredentials updates one or more (existing) Credentials in store
func UpdateCredentials(ctx context.Context, s Credentials, rr ...*types.Credentials) error {
	return s.UpdateCredentials(ctx, rr...)
}

// UpsertCredentials creates new or updates existing one or more Credentials in store
func UpsertCredentials(ctx context.Context, s Credentials, rr ...*types.Credentials) error {
	return s.UpsertCredentials(ctx, rr...)
}

// DeleteCredentials Deletes one or more Credentials from store
func DeleteCredentials(ctx context.Context, s Credentials, rr ...*types.Credentials) error {
	return s.DeleteCredentials(ctx, rr...)
}

// DeleteCredentialsByID Deletes Credentials from store
func DeleteCredentialsByID(ctx context.Context, s Credentials, ID uint64) error {
	return s.DeleteCredentialsByID(ctx, ID)
}

// TruncateCredentials Deletes all Credentials from store
func TruncateCredentials(ctx context.Context, s Credentials) error {
	return s.TruncateCredentials(ctx)
}

package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/credentials.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	credentialsStore interface {
		SearchCredentials(ctx context.Context, f types.CredentialsFilter) (types.CredentialsSet, types.CredentialsFilter, error)
		LookupCredentialsByID(ctx context.Context, id uint64) (*types.Credentials, error)
		CreateCredentials(ctx context.Context, rr ...*types.Credentials) error
		UpdateCredentials(ctx context.Context, rr ...*types.Credentials) error
		PartialUpdateCredentials(ctx context.Context, onlyColumns []string, rr ...*types.Credentials) error
		RemoveCredentials(ctx context.Context, rr ...*types.Credentials) error
		RemoveCredentialsByID(ctx context.Context, ID uint64) error

		TruncateCredentials(ctx context.Context) error
	}
)

package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/applications.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	applicationsStore interface {
		SearchApplications(ctx context.Context, f types.ApplicationFilter) (types.ApplicationSet, types.ApplicationFilter, error)
		LookupApplicationByID(ctx context.Context, id uint64) (*types.Application, error)
		CreateApplication(ctx context.Context, rr ...*types.Application) error
		UpdateApplication(ctx context.Context, rr ...*types.Application) error
		PartialUpdateApplication(ctx context.Context, onlyColumns []string, rr ...*types.Application) error
		RemoveApplication(ctx context.Context, rr ...*types.Application) error
		RemoveApplicationByID(ctx context.Context, ID uint64) error

		TruncateApplications(ctx context.Context) error
	}
)

package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/compose_modules.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	composeModulesStore interface {
		SearchComposeModules(ctx context.Context, f types.ModuleFilter) (types.ModuleSet, types.ModuleFilter, error)
		LookupComposeModuleByHandle(ctx context.Context, handle string) (*types.Module, error)
		LookupComposeModuleByID(ctx context.Context, id uint64) (*types.Module, error)
		CreateComposeModule(ctx context.Context, rr ...*types.Module) error
		UpdateComposeModule(ctx context.Context, rr ...*types.Module) error
		PartialUpdateComposeModule(ctx context.Context, onlyColumns []string, rr ...*types.Module) error
		RemoveComposeModule(ctx context.Context, rr ...*types.Module) error
		RemoveComposeModuleByID(ctx context.Context, ID uint64) error

		TruncateComposeModules(ctx context.Context) error
	}
)

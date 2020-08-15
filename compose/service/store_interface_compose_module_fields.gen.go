package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/compose_module_fields.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	composeModuleFieldsStore interface {
		CreateComposeModuleField(ctx context.Context, rr ...*types.ModuleField) error
		UpdateComposeModuleField(ctx context.Context, rr ...*types.ModuleField) error
		PartialUpdateComposeModuleField(ctx context.Context, onlyColumns []string, rr ...*types.ModuleField) error
		RemoveComposeModuleField(ctx context.Context, rr ...*types.ModuleField) error
		RemoveComposeModuleFieldByID(ctx context.Context, ID uint64) error

		TruncateComposeModuleFields(ctx context.Context) error
	}
)

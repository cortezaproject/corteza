package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/compose_module_fields.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	composeModuleFieldCreate struct {
		Done chan struct{}
		res  *types.ModuleField
		err  error
	}

	composeModuleFieldUpdate struct {
		Done chan struct{}
		res  *types.ModuleField
		err  error
	}

	composeModuleFieldRemove struct {
		Done chan struct{}
		res  *types.ModuleField
		err  error
	}
)

// CreateComposeModuleField creates a new ComposeModuleField
// create job that can be pushed to store's transaction handler
func CreateComposeModuleField(res *types.ModuleField) *composeModuleFieldCreate {
	return &composeModuleFieldCreate{res: res}
}

// Do Executes composeModuleFieldCreate job
func (j *composeModuleFieldCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateComposeModuleField(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateComposeModuleField creates a new ComposeModuleField
// update job that can be pushed to store's transaction handler
func UpdateComposeModuleField(res *types.ModuleField) *composeModuleFieldUpdate {
	return &composeModuleFieldUpdate{res: res}
}

// Do Executes composeModuleFieldUpdate job
func (j *composeModuleFieldUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateComposeModuleField(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveComposeModuleField creates a new ComposeModuleField
// remove job that can be pushed to store's transaction handler
func RemoveComposeModuleField(res *types.ModuleField) *composeModuleFieldRemove {
	return &composeModuleFieldRemove{res: res}
}

// Do Executes composeModuleFieldRemove job
func (j *composeModuleFieldRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveComposeModuleField(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

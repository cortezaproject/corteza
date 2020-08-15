package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/compose_modules.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	composeModuleCreate struct {
		Done chan struct{}
		res  *types.Module
		err  error
	}

	composeModuleUpdate struct {
		Done chan struct{}
		res  *types.Module
		err  error
	}

	composeModuleRemove struct {
		Done chan struct{}
		res  *types.Module
		err  error
	}
)

// CreateComposeModule creates a new ComposeModule
// create job that can be pushed to store's transaction handler
func CreateComposeModule(res *types.Module) *composeModuleCreate {
	return &composeModuleCreate{res: res}
}

// Do Executes composeModuleCreate job
func (j *composeModuleCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateComposeModule(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateComposeModule creates a new ComposeModule
// update job that can be pushed to store's transaction handler
func UpdateComposeModule(res *types.Module) *composeModuleUpdate {
	return &composeModuleUpdate{res: res}
}

// Do Executes composeModuleUpdate job
func (j *composeModuleUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateComposeModule(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveComposeModule creates a new ComposeModule
// remove job that can be pushed to store's transaction handler
func RemoveComposeModule(res *types.Module) *composeModuleRemove {
	return &composeModuleRemove{res: res}
}

// Do Executes composeModuleRemove job
func (j *composeModuleRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveComposeModule(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

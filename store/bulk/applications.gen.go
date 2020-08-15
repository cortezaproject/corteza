package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/applications.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	applicationCreate struct {
		Done chan struct{}
		res  *types.Application
		err  error
	}

	applicationUpdate struct {
		Done chan struct{}
		res  *types.Application
		err  error
	}

	applicationRemove struct {
		Done chan struct{}
		res  *types.Application
		err  error
	}
)

// CreateApplication creates a new Application
// create job that can be pushed to store's transaction handler
func CreateApplication(res *types.Application) *applicationCreate {
	return &applicationCreate{res: res}
}

// Do Executes applicationCreate job
func (j *applicationCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateApplication(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateApplication creates a new Application
// update job that can be pushed to store's transaction handler
func UpdateApplication(res *types.Application) *applicationUpdate {
	return &applicationUpdate{res: res}
}

// Do Executes applicationUpdate job
func (j *applicationUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateApplication(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveApplication creates a new Application
// remove job that can be pushed to store's transaction handler
func RemoveApplication(res *types.Application) *applicationRemove {
	return &applicationRemove{res: res}
}

// Do Executes applicationRemove job
func (j *applicationRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveApplication(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

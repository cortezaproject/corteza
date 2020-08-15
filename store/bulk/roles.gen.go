package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/roles.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	roleCreate struct {
		Done chan struct{}
		res  *types.Role
		err  error
	}

	roleUpdate struct {
		Done chan struct{}
		res  *types.Role
		err  error
	}

	roleRemove struct {
		Done chan struct{}
		res  *types.Role
		err  error
	}
)

// CreateRole creates a new Role
// create job that can be pushed to store's transaction handler
func CreateRole(res *types.Role) *roleCreate {
	return &roleCreate{res: res}
}

// Do Executes roleCreate job
func (j *roleCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateRole(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateRole creates a new Role
// update job that can be pushed to store's transaction handler
func UpdateRole(res *types.Role) *roleUpdate {
	return &roleUpdate{res: res}
}

// Do Executes roleUpdate job
func (j *roleUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateRole(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveRole creates a new Role
// remove job that can be pushed to store's transaction handler
func RemoveRole(res *types.Role) *roleRemove {
	return &roleRemove{res: res}
}

// Do Executes roleRemove job
func (j *roleRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveRole(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

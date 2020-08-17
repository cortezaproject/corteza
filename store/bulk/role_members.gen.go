package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/role_members.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	roleMemberCreate struct {
		Done chan struct{}
		res  *types.RoleMember
		err  error
	}

	roleMemberUpdate struct {
		Done chan struct{}
		res  *types.RoleMember
		err  error
	}

	roleMemberRemove struct {
		Done chan struct{}
		res  *types.RoleMember
		err  error
	}
)

// CreateRoleMember creates a new RoleMember
// create job that can be pushed to store's transaction handler
func CreateRoleMember(res *types.RoleMember) *roleMemberCreate {
	return &roleMemberCreate{res: res}
}

// Do Executes roleMemberCreate job
func (j *roleMemberCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateRoleMember(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateRoleMember creates a new RoleMember
// update job that can be pushed to store's transaction handler
func UpdateRoleMember(res *types.RoleMember) *roleMemberUpdate {
	return &roleMemberUpdate{res: res}
}

// Do Executes roleMemberUpdate job
func (j *roleMemberUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateRoleMember(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveRoleMember creates a new RoleMember
// remove job that can be pushed to store's transaction handler
func RemoveRoleMember(res *types.RoleMember) *roleMemberRemove {
	return &roleMemberRemove{res: res}
}

// Do Executes roleMemberRemove job
func (j *roleMemberRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveRoleMember(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

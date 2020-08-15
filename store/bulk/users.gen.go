package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/users.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	userCreate struct {
		Done chan struct{}
		res  *types.User
		err  error
	}

	userUpdate struct {
		Done chan struct{}
		res  *types.User
		err  error
	}

	userRemove struct {
		Done chan struct{}
		res  *types.User
		err  error
	}
)

// CreateUser creates a new User
// create job that can be pushed to store's transaction handler
func CreateUser(res *types.User) *userCreate {
	return &userCreate{res: res}
}

// Do Executes userCreate job
func (j *userCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateUser(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateUser creates a new User
// update job that can be pushed to store's transaction handler
func UpdateUser(res *types.User) *userUpdate {
	return &userUpdate{res: res}
}

// Do Executes userUpdate job
func (j *userUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateUser(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveUser creates a new User
// remove job that can be pushed to store's transaction handler
func RemoveUser(res *types.User) *userRemove {
	return &userRemove{res: res}
}

// Do Executes userRemove job
func (j *userRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveUser(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

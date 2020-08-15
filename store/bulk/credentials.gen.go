package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/credentials.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	credentialsCreate struct {
		Done chan struct{}
		res  *types.Credentials
		err  error
	}

	credentialsUpdate struct {
		Done chan struct{}
		res  *types.Credentials
		err  error
	}

	credentialsRemove struct {
		Done chan struct{}
		res  *types.Credentials
		err  error
	}
)

// CreateCredentials creates a new Credentials
// create job that can be pushed to store's transaction handler
func CreateCredentials(res *types.Credentials) *credentialsCreate {
	return &credentialsCreate{res: res}
}

// Do Executes credentialsCreate job
func (j *credentialsCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateCredentials(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateCredentials creates a new Credentials
// update job that can be pushed to store's transaction handler
func UpdateCredentials(res *types.Credentials) *credentialsUpdate {
	return &credentialsUpdate{res: res}
}

// Do Executes credentialsUpdate job
func (j *credentialsUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateCredentials(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveCredentials creates a new Credentials
// remove job that can be pushed to store's transaction handler
func RemoveCredentials(res *types.Credentials) *credentialsRemove {
	return &credentialsRemove{res: res}
}

// Do Executes credentialsRemove job
func (j *credentialsRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveCredentials(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

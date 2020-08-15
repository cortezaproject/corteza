package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/actionlog.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
)

type (
	actionlogCreate struct {
		Done chan struct{}
		res  *actionlog.Action
		err  error
	}

	actionlogUpdate struct {
		Done chan struct{}
		res  *actionlog.Action
		err  error
	}

	actionlogRemove struct {
		Done chan struct{}
		res  *actionlog.Action
		err  error
	}
)

// CreateActionlog creates a new Actionlog
// create job that can be pushed to store's transaction handler
func CreateActionlog(res *actionlog.Action) *actionlogCreate {
	return &actionlogCreate{res: res}
}

// Do Executes actionlogCreate job
func (j *actionlogCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateActionlog(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateActionlog creates a new Actionlog
// update job that can be pushed to store's transaction handler
func UpdateActionlog(res *actionlog.Action) *actionlogUpdate {
	return &actionlogUpdate{res: res}
}

// Do Executes actionlogUpdate job
func (j *actionlogUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateActionlog(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveActionlog creates a new Actionlog
// remove job that can be pushed to store's transaction handler
func RemoveActionlog(res *actionlog.Action) *actionlogRemove {
	return &actionlogRemove{res: res}
}

// Do Executes actionlogRemove job
func (j *actionlogRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveActionlog(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

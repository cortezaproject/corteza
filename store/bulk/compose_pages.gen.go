package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/compose_pages.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	composePageCreate struct {
		Done chan struct{}
		res  *types.Page
		err  error
	}

	composePageUpdate struct {
		Done chan struct{}
		res  *types.Page
		err  error
	}

	composePageRemove struct {
		Done chan struct{}
		res  *types.Page
		err  error
	}
)

// CreateComposePage creates a new ComposePage
// create job that can be pushed to store's transaction handler
func CreateComposePage(res *types.Page) *composePageCreate {
	return &composePageCreate{res: res}
}

// Do Executes composePageCreate job
func (j *composePageCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateComposePage(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateComposePage creates a new ComposePage
// update job that can be pushed to store's transaction handler
func UpdateComposePage(res *types.Page) *composePageUpdate {
	return &composePageUpdate{res: res}
}

// Do Executes composePageUpdate job
func (j *composePageUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateComposePage(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveComposePage creates a new ComposePage
// remove job that can be pushed to store's transaction handler
func RemoveComposePage(res *types.Page) *composePageRemove {
	return &composePageRemove{res: res}
}

// Do Executes composePageRemove job
func (j *composePageRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveComposePage(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

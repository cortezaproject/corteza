package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/system_attachments.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	attachmentCreate struct {
		Done chan struct{}
		res  *types.Attachment
		err  error
	}

	attachmentUpdate struct {
		Done chan struct{}
		res  *types.Attachment
		err  error
	}

	attachmentRemove struct {
		Done chan struct{}
		res  *types.Attachment
		err  error
	}
)

// CreateAttachment creates a new Attachment
// create job that can be pushed to store's transaction handler
func CreateAttachment(res *types.Attachment) *attachmentCreate {
	return &attachmentCreate{res: res}
}

// Do Executes attachmentCreate job
func (j *attachmentCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateAttachment(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateAttachment creates a new Attachment
// update job that can be pushed to store's transaction handler
func UpdateAttachment(res *types.Attachment) *attachmentUpdate {
	return &attachmentUpdate{res: res}
}

// Do Executes attachmentUpdate job
func (j *attachmentUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateAttachment(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveAttachment creates a new Attachment
// remove job that can be pushed to store's transaction handler
func RemoveAttachment(res *types.Attachment) *attachmentRemove {
	return &attachmentRemove{res: res}
}

// Do Executes attachmentRemove job
func (j *attachmentRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveAttachment(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

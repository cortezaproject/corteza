package rest

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"context"
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/types"
)

func (ctrl *Queue) List(ctx context.Context, r *request.QueuesList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Queue) Create(ctx context.Context, r *request.QueuesCreate) (interface{}, error) {
	res := &types.Queue{
		Queue:    r.Queue,
		Consumer: r.Consumer,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Queue) Read(ctx context.Context, r *request.QueuesRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.QueueID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Queue) Update(ctx context.Context, r *request.QueuesUpdate) (interface{}, error) {
	res := &types.Queue{
		ID:        r.QueueID,
		Queue:     r.Queue,
		Consumer:  r.Consumer,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Queue) Delete(ctx context.Context, r *request.QueuesDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.QueueID)
}

func (ctrl *Queue) Undelete(ctx context.Context, r *request.QueuesUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.QueueID)
}

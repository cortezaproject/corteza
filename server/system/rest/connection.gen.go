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

func (ctrl *Connection) List(ctx context.Context, r *request.ConnectionList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Connection) Create(ctx context.Context, r *request.ConnectionCreate) (interface{}, error) {
	res := &types.Connection{
		Handle: r.Handle,
		Labels: r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Connection) Update(ctx context.Context, r *request.ConnectionUpdate) (interface{}, error) {
	res := &types.Connection{
		ID:        r.ConnectionID,
		Handle:    r.Handle,
		Labels:    r.Labels,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Connection) Read(ctx context.Context, r *request.ConnectionRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.ConnectionID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Connection) Delete(ctx context.Context, r *request.ConnectionDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ConnectionID)
}

func (ctrl *Connection) Undelete(ctx context.Context, r *request.ConnectionUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.ConnectionID)
}

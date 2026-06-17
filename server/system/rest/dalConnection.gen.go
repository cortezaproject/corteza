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

func (ctrl *DalConnection) List(ctx context.Context, r *request.DalConnectionList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *DalConnection) Create(ctx context.Context, r *request.DalConnectionCreate) (interface{}, error) {
	res := &types.DalConnection{
		Handle: r.Handle,
		Type:   r.Type,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *DalConnection) Update(ctx context.Context, r *request.DalConnectionUpdate) (interface{}, error) {
	res := &types.DalConnection{
		ID:        r.ConnectionID,
		Handle:    r.Handle,
		Type:      r.Type,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *DalConnection) Read(ctx context.Context, r *request.DalConnectionRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.ConnectionID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *DalConnection) Delete(ctx context.Context, r *request.DalConnectionDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ConnectionID)
}

func (ctrl *DalConnection) Undelete(ctx context.Context, r *request.DalConnectionUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.ConnectionID)
}

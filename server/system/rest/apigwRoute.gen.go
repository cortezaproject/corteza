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

func (ctrl *ApigwRoute) List(ctx context.Context, r *request.ApigwRouteList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ApigwRoute) Create(ctx context.Context, r *request.ApigwRouteCreate) (interface{}, error) {
	res := &types.ApigwRoute{
		Endpoint: r.Endpoint,
		Method:   r.Method,
		Enabled:  r.Enabled,
		Group:    r.Group,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ApigwRoute) Update(ctx context.Context, r *request.ApigwRouteUpdate) (interface{}, error) {
	res := &types.ApigwRoute{
		ID:        r.RouteID,
		Endpoint:  r.Endpoint,
		Method:    r.Method,
		Enabled:   r.Enabled,
		Group:     r.Group,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ApigwRoute) Read(ctx context.Context, r *request.ApigwRouteRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.RouteID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ApigwRoute) Delete(ctx context.Context, r *request.ApigwRouteDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.RouteID)
}

func (ctrl *ApigwRoute) Undelete(ctx context.Context, r *request.ApigwRouteUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.RouteID)
}

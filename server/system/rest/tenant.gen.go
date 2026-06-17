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

func (ctrl *Tenant) List(ctx context.Context, r *request.TenantList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Tenant) Create(ctx context.Context, r *request.TenantCreate) (interface{}, error) {
	res := &types.Tenant{
		Handle: r.Handle,
		Labels: r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Tenant) Read(ctx context.Context, r *request.TenantRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.TenantID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Tenant) Update(ctx context.Context, r *request.TenantUpdate) (interface{}, error) {
	res := &types.Tenant{
		ID:        r.TenantID,
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

func (ctrl *Tenant) Delete(ctx context.Context, r *request.TenantDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.TenantID)
}

func (ctrl *Tenant) Undelete(ctx context.Context, r *request.TenantUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.TenantID)
}

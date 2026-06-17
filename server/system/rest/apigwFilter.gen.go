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

func (ctrl *ApigwFilter) List(ctx context.Context, r *request.ApigwFilterList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ApigwFilter) Create(ctx context.Context, r *request.ApigwFilterCreate) (interface{}, error) {
	res := &types.ApigwFilter{
		Weight:  r.Weight,
		Kind:    r.Kind,
		Ref:     r.Ref,
		Enabled: r.Enabled,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ApigwFilter) Update(ctx context.Context, r *request.ApigwFilterUpdate) (interface{}, error) {
	res := &types.ApigwFilter{
		ID:        r.FilterID,
		Weight:    r.Weight,
		Kind:      r.Kind,
		Ref:       r.Ref,
		Enabled:   r.Enabled,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ApigwFilter) Read(ctx context.Context, r *request.ApigwFilterRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.FilterID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ApigwFilter) Delete(ctx context.Context, r *request.ApigwFilterDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.FilterID)
}

func (ctrl *ApigwFilter) Undelete(ctx context.Context, r *request.ApigwFilterUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.FilterID)
}

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

func (ctrl *Application) List(ctx context.Context, r *request.ApplicationList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.application.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Application) Create(ctx context.Context, r *request.ApplicationCreate) (interface{}, error) {
	res := &types.Application{
		Name:    r.Name,
		Enabled: r.Enabled,
		Weight:  r.Weight,
		Labels:  r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.application.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Application) Update(ctx context.Context, r *request.ApplicationUpdate) (interface{}, error) {
	res := &types.Application{
		ID:        r.ApplicationID,
		Name:      r.Name,
		Enabled:   r.Enabled,
		Weight:    r.Weight,
		Labels:    r.Labels,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.application.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Application) Delete(ctx context.Context, r *request.ApplicationDelete) (interface{}, error) {
	return api.OK(), ctrl.application.DeleteByID(ctx, r.ApplicationID)
}

func (ctrl *Application) Undelete(ctx context.Context, r *request.ApplicationUndelete) (interface{}, error) {
	return api.OK(), ctrl.application.UndeleteByID(ctx, r.ApplicationID)
}

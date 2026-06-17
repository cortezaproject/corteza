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

func (ctrl *Project) List(ctx context.Context, r *request.ProjectList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Project) Create(ctx context.Context, r *request.ProjectCreate) (interface{}, error) {
	res := &types.Project{
		Handle: r.Handle,
		Labels: r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Project) Read(ctx context.Context, r *request.ProjectRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.ProjectID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Project) Update(ctx context.Context, r *request.ProjectUpdate) (interface{}, error) {
	res := &types.Project{
		ID:        r.ProjectID,
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

func (ctrl *Project) Delete(ctx context.Context, r *request.ProjectDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ProjectID)
}

func (ctrl *Project) Undelete(ctx context.Context, r *request.ProjectUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.ProjectID)
}

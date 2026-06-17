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

func (ctrl *ProjectGroup) List(ctx context.Context, r *request.ProjectGroupList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.projectGroup.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ProjectGroup) Create(ctx context.Context, r *request.ProjectGroupCreate) (interface{}, error) {
	res := &types.ProjectGroup{
		Handle: r.Handle,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectGroup.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectGroup) Read(ctx context.Context, r *request.ProjectGroupRead) (interface{}, error) {
	res, err := ctrl.projectGroup.FindByID(ctx, r.ProjectGroupID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectGroup) Update(ctx context.Context, r *request.ProjectGroupUpdate) (interface{}, error) {
	res := &types.ProjectGroup{
		ID:        r.ProjectGroupID,
		Handle:    r.Handle,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectGroup.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectGroup) Delete(ctx context.Context, r *request.ProjectGroupDelete) (interface{}, error) {
	return api.OK(), ctrl.projectGroup.DeleteByID(ctx, r.ProjectGroupID)
}

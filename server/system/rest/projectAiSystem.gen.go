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

func (ctrl *ProjectAiSystem) List(ctx context.Context, r *request.ProjectAiSystemList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.projectAiSystem.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ProjectAiSystem) Create(ctx context.Context, r *request.ProjectAiSystemCreate) (interface{}, error) {
	res := &types.ProjectAiSystem{
		Handle:    r.Handle,
		RiskClass: r.RiskClass,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectAiSystem.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectAiSystem) Read(ctx context.Context, r *request.ProjectAiSystemRead) (interface{}, error) {
	res, err := ctrl.projectAiSystem.FindByID(ctx, r.ProjectAiSystemID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectAiSystem) Update(ctx context.Context, r *request.ProjectAiSystemUpdate) (interface{}, error) {
	res := &types.ProjectAiSystem{
		ID:        r.ProjectAiSystemID,
		Handle:    r.Handle,
		RiskClass: r.RiskClass,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectAiSystem.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectAiSystem) Delete(ctx context.Context, r *request.ProjectAiSystemDelete) (interface{}, error) {
	return api.OK(), ctrl.projectAiSystem.DeleteByID(ctx, r.ProjectAiSystemID)
}

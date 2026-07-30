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

func (ctrl *ProjectFriaScenario) List(ctx context.Context, r *request.ProjectFriaScenarioList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.projectFriaScenario.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ProjectFriaScenario) Create(ctx context.Context, r *request.ProjectFriaScenarioCreate) (interface{}, error) {
	res := &types.ProjectFriaScenario{
		AiSystemID: r.AiSystemID,
		Title:      r.Title,
		Severity:   r.Severity,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectFriaScenario.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectFriaScenario) Read(ctx context.Context, r *request.ProjectFriaScenarioRead) (interface{}, error) {
	res, err := ctrl.projectFriaScenario.FindByID(ctx, r.ProjectFriaScenarioID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectFriaScenario) Update(ctx context.Context, r *request.ProjectFriaScenarioUpdate) (interface{}, error) {
	res := &types.ProjectFriaScenario{
		ID:         r.ProjectFriaScenarioID,
		AiSystemID: r.AiSystemID,
		Title:      r.Title,
		Severity:   r.Severity,
		UpdatedAt:  r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectFriaScenario.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectFriaScenario) Delete(ctx context.Context, r *request.ProjectFriaScenarioDelete) (interface{}, error) {
	return api.OK(), ctrl.projectFriaScenario.DeleteByID(ctx, r.ProjectFriaScenarioID)
}

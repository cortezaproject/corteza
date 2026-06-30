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

func (ctrl *Agent) List(ctx context.Context, r *request.AgentList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.agent.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Agent) Create(ctx context.Context, r *request.AgentCreate) (interface{}, error) {
	res := &types.Agent{
		Handle:    r.Handle,
		ProjectID: r.ProjectID,
		Status:    r.Status,
		Labels:    r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.agent.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Agent) Read(ctx context.Context, r *request.AgentRead) (interface{}, error) {
	res, err := ctrl.agent.FindByID(ctx, r.AgentID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Agent) Update(ctx context.Context, r *request.AgentUpdate) (interface{}, error) {
	res := &types.Agent{
		ID:        r.AgentID,
		Handle:    r.Handle,
		Status:    r.Status,
		Labels:    r.Labels,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.agent.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Agent) Delete(ctx context.Context, r *request.AgentDelete) (interface{}, error) {
	return api.OK(), ctrl.agent.DeleteByID(ctx, r.AgentID)
}

func (ctrl *Agent) Undelete(ctx context.Context, r *request.AgentUndelete) (interface{}, error) {
	return api.OK(), ctrl.agent.UndeleteByID(ctx, r.AgentID)
}

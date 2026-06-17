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
	"github.com/crusttech/human/server/automation/rest/request"
	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/api"
)

func (ctrl *Workflow) List(ctx context.Context, r *request.WorkflowList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.workflow.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Workflow) Create(ctx context.Context, r *request.WorkflowCreate) (interface{}, error) {
	res := &types.Workflow{
		Handle:       r.Handle,
		Labels:       r.Labels,
		Enabled:      r.Enabled,
		Trace:        r.Trace,
		KeepSessions: r.KeepSessions,
		RunAs:        r.RunAs,
		OwnedBy:      r.OwnedBy,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.workflow.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Workflow) Update(ctx context.Context, r *request.WorkflowUpdate) (interface{}, error) {
	res := &types.Workflow{
		ID:           r.WorkflowID,
		Handle:       r.Handle,
		Labels:       r.Labels,
		Enabled:      r.Enabled,
		Trace:        r.Trace,
		KeepSessions: r.KeepSessions,
		RunAs:        r.RunAs,
		OwnedBy:      r.OwnedBy,
		UpdatedAt:    r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.workflow.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Workflow) Read(ctx context.Context, r *request.WorkflowRead) (interface{}, error) {
	res, err := ctrl.workflow.FindByID(ctx, r.WorkflowID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Workflow) Delete(ctx context.Context, r *request.WorkflowDelete) (interface{}, error) {
	return api.OK(), ctrl.workflow.DeleteByID(ctx, r.WorkflowID)
}

func (ctrl *Workflow) Undelete(ctx context.Context, r *request.WorkflowUndelete) (interface{}, error) {
	return api.OK(), ctrl.workflow.UndeleteByID(ctx, r.WorkflowID)
}

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

func (ctrl *Trigger) List(ctx context.Context, r *request.TriggerList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.trigger.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Trigger) Create(ctx context.Context, r *request.TriggerCreate) (interface{}, error) {
	res := &types.Trigger{
		EventType:    r.EventType,
		ResourceType: r.ResourceType,
		Enabled:      r.Enabled,
		WorkflowID:   r.WorkflowID,
		Labels:       r.Labels,
		OwnedBy:      r.OwnedBy,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.trigger.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Trigger) Update(ctx context.Context, r *request.TriggerUpdate) (interface{}, error) {
	res := &types.Trigger{
		ID:           r.TriggerID,
		EventType:    r.EventType,
		ResourceType: r.ResourceType,
		Enabled:      r.Enabled,
		WorkflowID:   r.WorkflowID,
		Labels:       r.Labels,
		OwnedBy:      r.OwnedBy,
		UpdatedAt:    r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.trigger.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Trigger) Read(ctx context.Context, r *request.TriggerRead) (interface{}, error) {
	res, err := ctrl.trigger.FindByID(ctx, r.TriggerID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Trigger) Delete(ctx context.Context, r *request.TriggerDelete) (interface{}, error) {
	return api.OK(), ctrl.trigger.DeleteByID(ctx, r.TriggerID)
}

func (ctrl *Trigger) Undelete(ctx context.Context, r *request.TriggerUndelete) (interface{}, error) {
	return api.OK(), ctrl.trigger.UndeleteByID(ctx, r.TriggerID)
}

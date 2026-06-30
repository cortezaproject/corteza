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

func (ctrl *NgAutomation) List(ctx context.Context, r *request.NgAutomationList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.ngAutomation.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *NgAutomation) Create(ctx context.Context, r *request.NgAutomationCreate) (interface{}, error) {
	res := &types.NgAutomation{
		Handle:    r.Handle,
		ProjectID: r.ProjectID,
		Labels:    r.Labels,
		Enabled:   r.Enabled,
		RunAs:     r.RunAs,
		OwnedBy:   r.OwnedBy,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.ngAutomation.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *NgAutomation) Update(ctx context.Context, r *request.NgAutomationUpdate) (interface{}, error) {
	res := &types.NgAutomation{
		ID:        r.AutomationID,
		Handle:    r.Handle,
		Labels:    r.Labels,
		Enabled:   r.Enabled,
		RunAs:     r.RunAs,
		OwnedBy:   r.OwnedBy,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.ngAutomation.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *NgAutomation) Read(ctx context.Context, r *request.NgAutomationRead) (interface{}, error) {
	res, err := ctrl.ngAutomation.FindByID(ctx, r.AutomationID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *NgAutomation) Delete(ctx context.Context, r *request.NgAutomationDelete) (interface{}, error) {
	return api.OK(), ctrl.ngAutomation.DeleteByID(ctx, r.AutomationID)
}

func (ctrl *NgAutomation) Undelete(ctx context.Context, r *request.NgAutomationUndelete) (interface{}, error) {
	return api.OK(), ctrl.ngAutomation.UndeleteByID(ctx, r.AutomationID)
}

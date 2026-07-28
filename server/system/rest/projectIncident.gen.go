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

func (ctrl *ProjectIncident) List(ctx context.Context, r *request.ProjectIncidentList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.projectIncident.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ProjectIncident) Create(ctx context.Context, r *request.ProjectIncidentCreate) (interface{}, error) {
	res := &types.ProjectIncident{
		ProjectID:        r.ProjectID,
		RevisionID:       r.RevisionID,
		Title:            r.Title,
		Description:      r.Description,
		IncidentType:     r.IncidentType,
		GroupSystem:      r.GroupSystem,
		Status:           r.Status,
		Severity:         r.Severity,
		Risk:             r.Risk,
		IssueOwner:       r.IssueOwner,
		ChangeOwner:      r.ChangeOwner,
		ChangeApprovedBy: r.ChangeApprovedBy,
		RiskIssue:        r.RiskIssue,
		ChangeRequired:   r.ChangeRequired,
		RiskChange:       r.RiskChange,
		DateDue:          r.DateDue,
		CompletedDate:    r.CompletedDate,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectIncident.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectIncident) Read(ctx context.Context, r *request.ProjectIncidentRead) (interface{}, error) {
	res, err := ctrl.projectIncident.FindByID(ctx, r.IncidentID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectIncident) Update(ctx context.Context, r *request.ProjectIncidentUpdate) (interface{}, error) {
	res := &types.ProjectIncident{
		ID:               r.IncidentID,
		RevisionID:       r.RevisionID,
		Title:            r.Title,
		Description:      r.Description,
		IncidentType:     r.IncidentType,
		GroupSystem:      r.GroupSystem,
		Status:           r.Status,
		Severity:         r.Severity,
		Risk:             r.Risk,
		IssueOwner:       r.IssueOwner,
		ChangeOwner:      r.ChangeOwner,
		ChangeApprovedBy: r.ChangeApprovedBy,
		RiskIssue:        r.RiskIssue,
		ChangeRequired:   r.ChangeRequired,
		RiskChange:       r.RiskChange,
		DateDue:          r.DateDue,
		CompletedDate:    r.CompletedDate,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectIncident.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectIncident) Delete(ctx context.Context, r *request.ProjectIncidentDelete) (interface{}, error) {
	return api.OK(), ctrl.projectIncident.DeleteByID(ctx, r.IncidentID)
}

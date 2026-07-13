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

func (ctrl *ProjectPrivacy) List(ctx context.Context, r *request.ProjectPrivacyList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.projectPrivacy.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ProjectPrivacy) Create(ctx context.Context, r *request.ProjectPrivacyCreate) (interface{}, error) {
	res := &types.ProjectPrivacy{
		ProjectID:        r.ProjectID,
		Title:            r.Title,
		Description:      r.Description,
		RequestType:      r.RequestType,
		Status:           r.Status,
		Severity:         r.Severity,
		Risk:             r.Risk,
		RequestOwner:     r.RequestOwner,
		ChangeOwner:      r.ChangeOwner,
		ChangeApprovedBy: r.ChangeApprovedBy,
		RiskAssessment:   r.RiskAssessment,
		ChangeRequired:   r.ChangeRequired,
		RiskChange:       r.RiskChange,
		Backlog:          r.Backlog,
		DateDue:          r.DateDue,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectPrivacy.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectPrivacy) Read(ctx context.Context, r *request.ProjectPrivacyRead) (interface{}, error) {
	res, err := ctrl.projectPrivacy.FindByID(ctx, r.PrivacyID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectPrivacy) Update(ctx context.Context, r *request.ProjectPrivacyUpdate) (interface{}, error) {
	res := &types.ProjectPrivacy{
		ID:               r.PrivacyID,
		Title:            r.Title,
		Description:      r.Description,
		RequestType:      r.RequestType,
		Status:           r.Status,
		Severity:         r.Severity,
		Risk:             r.Risk,
		RequestOwner:     r.RequestOwner,
		ChangeOwner:      r.ChangeOwner,
		ChangeApprovedBy: r.ChangeApprovedBy,
		RiskAssessment:   r.RiskAssessment,
		ChangeRequired:   r.ChangeRequired,
		RiskChange:       r.RiskChange,
		Backlog:          r.Backlog,
		DateDue:          r.DateDue,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectPrivacy.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectPrivacy) Delete(ctx context.Context, r *request.ProjectPrivacyDelete) (interface{}, error) {
	return api.OK(), ctrl.projectPrivacy.DeleteByID(ctx, r.PrivacyID)
}

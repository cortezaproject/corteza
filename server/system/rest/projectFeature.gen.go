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

func (ctrl *ProjectFeature) List(ctx context.Context, r *request.ProjectFeatureList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.projectFeature.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ProjectFeature) Create(ctx context.Context, r *request.ProjectFeatureCreate) (interface{}, error) {
	res := &types.ProjectFeature{
		ProjectID:        r.ProjectID,
		Title:            r.Title,
		Description:      r.Description,
		FeatureType:      r.FeatureType,
		Status:           r.Status,
		Severity:         r.Severity,
		Risk:             r.Risk,
		FeatureOwner:     r.FeatureOwner,
		ChangeOwner:      r.ChangeOwner,
		ChangeApprovedBy: r.ChangeApprovedBy,
		RiskFeature:      r.RiskFeature,
		ChangeRequired:   r.ChangeRequired,
		RiskChange:       r.RiskChange,
		DateDue:          r.DateDue,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectFeature.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectFeature) Read(ctx context.Context, r *request.ProjectFeatureRead) (interface{}, error) {
	res, err := ctrl.projectFeature.FindByID(ctx, r.FeatureID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectFeature) Update(ctx context.Context, r *request.ProjectFeatureUpdate) (interface{}, error) {
	res := &types.ProjectFeature{
		ID:               r.FeatureID,
		Title:            r.Title,
		Description:      r.Description,
		FeatureType:      r.FeatureType,
		Status:           r.Status,
		Severity:         r.Severity,
		Risk:             r.Risk,
		FeatureOwner:     r.FeatureOwner,
		ChangeOwner:      r.ChangeOwner,
		ChangeApprovedBy: r.ChangeApprovedBy,
		RiskFeature:      r.RiskFeature,
		ChangeRequired:   r.ChangeRequired,
		RiskChange:       r.RiskChange,
		DateDue:          r.DateDue,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectFeature.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectFeature) Delete(ctx context.Context, r *request.ProjectFeatureDelete) (interface{}, error) {
	return api.OK(), ctrl.projectFeature.DeleteByID(ctx, r.FeatureID)
}

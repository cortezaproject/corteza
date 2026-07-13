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

func (ctrl *ProjectReview) List(ctx context.Context, r *request.ProjectReviewList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.projectReview.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ProjectReview) Create(ctx context.Context, r *request.ProjectReviewCreate) (interface{}, error) {
	res := &types.ProjectReview{
		ProjectID:       r.ProjectID,
		Title:           r.Title,
		Description:     r.Description,
		ReviewType:      r.ReviewType,
		ReviewFrequency: r.ReviewFrequency,
		Scope:           r.Scope,
		Status:          r.Status,
		Reviewer:        r.Reviewer,
		ApprovedBy:      r.ApprovedBy,
		DateDue:         r.DateDue,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectReview.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectReview) Read(ctx context.Context, r *request.ProjectReviewRead) (interface{}, error) {
	res, err := ctrl.projectReview.FindByID(ctx, r.ReviewID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectReview) Update(ctx context.Context, r *request.ProjectReviewUpdate) (interface{}, error) {
	res := &types.ProjectReview{
		ID:              r.ReviewID,
		Title:           r.Title,
		Description:     r.Description,
		ReviewType:      r.ReviewType,
		ReviewFrequency: r.ReviewFrequency,
		Scope:           r.Scope,
		Status:          r.Status,
		Reviewer:        r.Reviewer,
		ApprovedBy:      r.ApprovedBy,
		DateDue:         r.DateDue,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectReview.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectReview) Delete(ctx context.Context, r *request.ProjectReviewDelete) (interface{}, error) {
	return api.OK(), ctrl.projectReview.DeleteByID(ctx, r.ReviewID)
}

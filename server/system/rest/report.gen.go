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

func (ctrl *Report) List(ctx context.Context, r *request.ReportList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.report.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Report) Create(ctx context.Context, r *request.ReportCreate) (interface{}, error) {
	res := &types.Report{
		Handle: r.Handle,
		Labels: r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.report.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Report) Update(ctx context.Context, r *request.ReportUpdate) (interface{}, error) {
	res := &types.Report{
		ID:        r.ReportID,
		Handle:    r.Handle,
		Labels:    r.Labels,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.report.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Report) Read(ctx context.Context, r *request.ReportRead) (interface{}, error) {
	res, err := ctrl.report.FindByID(ctx, r.ReportID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Report) Delete(ctx context.Context, r *request.ReportDelete) (interface{}, error) {
	return api.OK(), ctrl.report.DeleteByID(ctx, r.ReportID)
}

func (ctrl *Report) Undelete(ctx context.Context, r *request.ReportUndelete) (interface{}, error) {
	return api.OK(), ctrl.report.UndeleteByID(ctx, r.ReportID)
}

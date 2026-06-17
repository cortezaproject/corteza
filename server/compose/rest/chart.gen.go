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
	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/api"
)

func (ctrl *Chart) List(ctx context.Context, r *request.ChartList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.chart.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Chart) Create(ctx context.Context, r *request.ChartCreate) (interface{}, error) {
	res := &types.Chart{
		NamespaceID: r.NamespaceID,
		Name:        r.Name,
		Handle:      r.Handle,
		Labels:      r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.chart.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Chart) Read(ctx context.Context, r *request.ChartRead) (interface{}, error) {
	res, err := ctrl.chart.FindByID(ctx, r.NamespaceID, r.ChartID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Chart) Update(ctx context.Context, r *request.ChartUpdate) (interface{}, error) {
	res := &types.Chart{
		ID:          r.ChartID,
		NamespaceID: r.NamespaceID,
		Name:        r.Name,
		Handle:      r.Handle,
		Labels:      r.Labels,
		UpdatedAt:   r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.chart.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Chart) Delete(ctx context.Context, r *request.ChartDelete) (interface{}, error) {
	return api.OK(), ctrl.chart.DeleteByID(ctx, r.NamespaceID, r.ChartID)
}

package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/crm/internal/service"
	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/types"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Chart struct {
	module service.ModuleService
	chart  service.ChartService
}

func (Chart) New() *Chart {
	return &Chart{
		module: service.DefaultModule,
		chart:  service.DefaultChart,
	}
}

func (ctrl *Chart) List(ctx context.Context, r *request.ChartList) (interface{}, error) {
	return ctrl.chart.With(ctx).Find()
}

func (ctrl *Chart) Create(ctx context.Context, r *request.ChartCreate) (interface{}, error) {
	chart := &types.Chart{
		Name:   r.Name,
		Config: r.Config,
	}

	return ctrl.chart.With(ctx).Create(chart)
}

func (ctrl *Chart) Read(ctx context.Context, r *request.ChartRead) (interface{}, error) {
	return ctrl.chart.With(ctx).FindByID(r.ChartID)
}

func (ctrl *Chart) Update(ctx context.Context, r *request.ChartUpdate) (interface{}, error) {
	chart := &types.Chart{
		ID:     r.ChartID,
		Name:   r.Name,
		Config: r.Config,
	}

	return ctrl.chart.With(ctx).Update(chart)
}

func (ctrl *Chart) Delete(ctx context.Context, r *request.ChartDelete) (interface{}, error) {
	return resputil.OK(), ctrl.chart.With(ctx).DeleteByID(r.ChartID)
}

package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/compose/internal/service"
	"github.com/crusttech/crust/compose/rest/request"
	"github.com/crusttech/crust/compose/types"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	chartPayload struct {
		*types.Chart

		CanUpdateChart bool `json:"canUpdateChart"`
		CanDeleteChart bool `json:"canDeleteChart"`
	}

	chartSetPayload struct {
		Filter types.ChartFilter `json:"filter"`
		Set    []*chartPayload   `json:"set"`
	}

	Chart struct {
		chart       service.ChartService
		permissions service.PermissionsService
	}
)

func (Chart) New() *Chart {
	return &Chart{
		chart:       service.DefaultChart,
		permissions: service.DefaultPermissions,
	}
}

func (ctrl Chart) List(ctx context.Context, r *request.ChartList) (interface{}, error) {
	f := types.ChartFilter{
		Query:   r.Query,
		PerPage: r.PerPage,
		Page:    r.Page,
	}

	set, filter, err := ctrl.chart.With(ctx).Find(f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Chart) Create(ctx context.Context, r *request.ChartCreate) (interface{}, error) {
	var err error
	ns := &types.Chart{
		NamespaceID: r.NamespaceID,
		Name:        r.Name,
		Config:      r.Config,
	}

	ns, err = ctrl.chart.With(ctx).Create(ns)
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Chart) Read(ctx context.Context, r *request.ChartRead) (interface{}, error) {
	return ctrl.chart.With(ctx).FindByID(r.NamespaceID, r.ChartID)
}

func (ctrl Chart) Update(ctx context.Context, r *request.ChartUpdate) (interface{}, error) {
	var (
		ns  = &types.Chart{}
		err error
	)

	ns.ID = r.ChartID
	ns.Name = r.Name
	ns.Config = r.Config
	ns.NamespaceID = r.NamespaceID
	ns.UpdatedAt = r.UpdatedAt

	ns, err = ctrl.chart.With(ctx).Update(ns)
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Chart) Delete(ctx context.Context, r *request.ChartDelete) (interface{}, error) {
	_, err := ctrl.chart.With(ctx).FindByID(r.NamespaceID, r.ChartID)
	if err != nil {
		return nil, err
	}

	return resputil.OK(), ctrl.chart.With(ctx).DeleteByID(r.NamespaceID, r.ChartID)
}

func (ctrl Chart) makePayload(ctx context.Context, t *types.Chart, err error) (*chartPayload, error) {
	if err != nil || t == nil {
		return nil, err
	}

	perm := ctrl.permissions.With(ctx)

	return &chartPayload{
		Chart: t,

		CanUpdateChart: perm.CanUpdateChart(t),
		CanDeleteChart: perm.CanDeleteChart(t),
	}, nil
}

func (ctrl Chart) makeFilterPayload(ctx context.Context, nn types.ChartSet, f types.ChartFilter, err error) (*chartSetPayload, error) {
	if err != nil {
		return nil, err
	}

	nsp := &chartSetPayload{Filter: f, Set: make([]*chartPayload, len(nn))}

	for i := range nn {
		nsp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return nsp, nil
}

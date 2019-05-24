package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/types"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	chartPayload struct {
		*types.Chart

		CanGrant       bool `json:"canGrant"`
		CanUpdateChart bool `json:"canUpdateChart"`
		CanDeleteChart bool `json:"canDeleteChart"`
	}

	chartSetPayload struct {
		Filter types.ChartFilter `json:"filter"`
		Set    []*chartPayload   `json:"set"`
	}

	Chart struct {
		chart service.ChartService
		ac    chartAccessController
	}

	chartAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateChart(context.Context, *types.Chart) bool
		CanDeleteChart(context.Context, *types.Chart) bool
	}
)

func (Chart) New() *Chart {
	return &Chart{
		chart: service.DefaultChart,
		ac:    service.DefaultAccessControl,
	}
}

func (ctrl Chart) List(ctx context.Context, r *request.ChartList) (interface{}, error) {
	f := types.ChartFilter{
		NamespaceID: r.NamespaceID,

		Query:   r.Query,
		PerPage: r.PerPage,
		Page:    r.Page,
	}

	set, filter, err := ctrl.chart.With(ctx).Find(f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Chart) Create(ctx context.Context, r *request.ChartCreate) (interface{}, error) {
	var err error
	mod := &types.Chart{
		NamespaceID: r.NamespaceID,
		Name:        r.Name,
		Config:      r.Config,
	}

	mod, err = ctrl.chart.With(ctx).Create(mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl Chart) Read(ctx context.Context, r *request.ChartRead) (interface{}, error) {
	mod, err := ctrl.chart.With(ctx).FindByID(r.NamespaceID, r.ChartID)
	return ctrl.makePayload(ctx, mod, err)

}

func (ctrl Chart) Update(ctx context.Context, r *request.ChartUpdate) (interface{}, error) {
	var (
		mod = &types.Chart{}
		err error
	)

	mod.ID = r.ChartID
	mod.Name = r.Name
	mod.Config = r.Config
	mod.NamespaceID = r.NamespaceID
	mod.UpdatedAt = r.UpdatedAt

	mod, err = ctrl.chart.With(ctx).Update(mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl Chart) Delete(ctx context.Context, r *request.ChartDelete) (interface{}, error) {
	_, err := ctrl.chart.With(ctx).FindByID(r.NamespaceID, r.ChartID)
	if err != nil {
		return nil, err
	}

	return resputil.OK(), ctrl.chart.With(ctx).DeleteByID(r.NamespaceID, r.ChartID)
}

func (ctrl Chart) makePayload(ctx context.Context, c *types.Chart, err error) (*chartPayload, error) {
	if err != nil || c == nil {
		return nil, err
	}

	return &chartPayload{
		Chart: c,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateChart: ctrl.ac.CanUpdateChart(ctx, c),
		CanDeleteChart: ctrl.ac.CanDeleteChart(ctx, c),
	}, nil
}

func (ctrl Chart) makeFilterPayload(ctx context.Context, nn types.ChartSet, f types.ChartFilter, err error) (*chartSetPayload, error) {
	if err != nil {
		return nil, err
	}

	modp := &chartSetPayload{Filter: f, Set: make([]*chartPayload, len(nn))}

	for i := range nn {
		modp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return modp, nil
}

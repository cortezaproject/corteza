package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

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
		chart interface {
			FindByID(ctx context.Context, namespaceID, chartID uint64) (*types.Chart, error)
			FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*types.Chart, error)
			Find(ctx context.Context, filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error)

			Create(ctx context.Context, chart *types.Chart) (*types.Chart, error)
			Update(ctx context.Context, chart *types.Chart) (*types.Chart, error)
			DeleteByID(ctx context.Context, namespaceID, chartID uint64) error
		}
		locale service.ResourceTranslationsManagerService
		ac     chartAccessController
	}

	chartAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateChart(context.Context, *types.Chart) bool
		CanDeleteChart(context.Context, *types.Chart) bool
	}
)

func (Chart) New() *Chart {
	return &Chart{
		chart:  service.DefaultChart,
		locale: service.DefaultResourceTranslation,
		ac:     service.DefaultAccessControl,
	}
}

func (ctrl Chart) List(ctx context.Context, r *request.ChartList) (interface{}, error) {
	var (
		err error
		f   = types.ChartFilter{
			NamespaceID: r.NamespaceID,

			Handle: r.Handle,
			Query:  r.Query,
			Labels: r.Labels,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.chart.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Chart) Create(ctx context.Context, r *request.ChartCreate) (interface{}, error) {
	var err error
	mod := &types.Chart{
		NamespaceID: r.NamespaceID,
		Name:        r.Name,
		Handle:      r.Handle,
		Labels:      r.Labels,
	}

	if len(r.Config) > 2 {
		if err = r.Config.Unmarshal(&mod.Config); err != nil {
			return nil, err
		}
	}
	mod, err = ctrl.chart.Create(ctx, mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl Chart) Read(ctx context.Context, r *request.ChartRead) (interface{}, error) {
	mod, err := ctrl.chart.FindByID(ctx, r.NamespaceID, r.ChartID)
	return ctrl.makePayload(ctx, mod, err)

}

func (ctrl Chart) Update(ctx context.Context, r *request.ChartUpdate) (interface{}, error) {
	var (
		err error
		mod = &types.Chart{
			ID:          r.ChartID,
			Name:        r.Name,
			Handle:      r.Handle,
			NamespaceID: r.NamespaceID,
			UpdatedAt:   r.UpdatedAt,
			Labels:      r.Labels,
		}
	)

	if len(r.Config) > 2 {
		if err = r.Config.Unmarshal(&mod.Config); err != nil {
			return nil, err
		}
	}
	mod, err = ctrl.chart.Update(ctx, mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl Chart) Delete(ctx context.Context, r *request.ChartDelete) (interface{}, error) {
	_, err := ctrl.chart.FindByID(ctx, r.NamespaceID, r.ChartID)
	if err != nil {
		return nil, err
	}

	return api.OK(), ctrl.chart.DeleteByID(ctx, r.NamespaceID, r.ChartID)
}

func (ctrl Chart) ListTranslations(ctx context.Context, r *request.ChartListTranslations) (interface{}, error) {
	return ctrl.locale.Chart(ctx, r.NamespaceID, r.ChartID)
}

func (ctrl Chart) UpdateTranslations(ctx context.Context, r *request.ChartUpdateTranslations) (interface{}, error) {
	return api.OK(), ctrl.locale.Upsert(ctx, r.Translations)
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

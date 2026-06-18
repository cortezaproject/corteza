package rest

import (
	"context"

	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/compose/service"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/filter"
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
			Search(ctx context.Context, filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error)

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

// makeFilter builds the search filter for the generated List controller.
func (ctrl Chart) makeFilter(ctx context.Context, r *request.ChartList) (types.ChartFilter, error) {
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
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params and the compound (namespace) id.
func (ctrl Chart) beforeCreate(ctx context.Context, res *types.Chart, r *request.ChartCreate) error {
	if len(r.Config) > 2 {
		if err := r.Config.Unmarshal(&res.Config); err != nil {
			return err
		}
	}

	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params, the compound (namespace) id and
// the resource ID/UpdatedAt.
func (ctrl Chart) beforeUpdate(ctx context.Context, res *types.Chart, r *request.ChartUpdate) error {
	if len(r.Config) > 2 {
		if err := r.Config.Unmarshal(&res.Config); err != nil {
			return err
		}
	}

	return nil
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

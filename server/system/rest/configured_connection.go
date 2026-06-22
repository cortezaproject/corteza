package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ConfiguredConnection struct {
		configuredConnection configuredConnectionService
		ac                   configuredConnectionAccessController
	}

	configuredConnectionPayload struct {
		*types.ConfiguredConnection

		CanUpdateConfiguredConnection bool `json:"canUpdateConfiguredConnection"`
		CanDeleteConfiguredConnection bool `json:"canDeleteConfiguredConnection"`
	}

	configuredConnectionSetPayload struct {
		Filter types.ConfiguredConnectionFilter `json:"filter"`
		Set    []*configuredConnectionPayload   `json:"set"`
	}

	configuredConnectionAccessController interface {
		CanCreateConfiguredConnection(context.Context) bool
		CanUpdateConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
		CanDeleteConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
	}

	configuredConnectionService interface {
		FindByID(ctx context.Context, ID uint64) (*types.ConfiguredConnection, error)
		DeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.ConfiguredConnectionFilter) (types.ConfiguredConnectionSet, types.ConfiguredConnectionFilter, error)
		Enable(ctx context.Context, ID uint64) (*types.ConfiguredConnection, error)
		Check(ctx context.Context, ID uint64) (*types.ConfiguredConnectionCheckResult, error)
		RefreshDiscovery(ctx context.Context, ID uint64) (map[string]any, error)
	}
)

func (ConfiguredConnection) New() *ConfiguredConnection {
	return &ConfiguredConnection{
		configuredConnection: service.DefaultConfiguredConnection,
		ac:                   service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl ConfiguredConnection) makeFilter(ctx context.Context, r *request.ConfiguredConnectionList) (types.ConfiguredConnectionFilter, error) {
	var (
		err error
		f   = types.ConfiguredConnectionFilter{
			ConnectionID: r.ConnectionID,
			ProjectID:    r.ProjectID,
			Status:       r.Status,
			Query:        r.Query,
			Deleted:      filter.State(r.Deleted),
		}
	)

	if f.Deleted == 0 {
		f.Deleted = filter.StateExcluded
	}

	f.IncTotal = r.IncTotal

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

func (ctrl ConfiguredConnection) Enable(ctx context.Context, r *request.ConfiguredConnectionEnable) (interface{}, error) {
	res, err := ctrl.configuredConnection.Enable(ctx, r.ConnectionID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl ConfiguredConnection) Check(ctx context.Context, r *request.ConfiguredConnectionCheck) (interface{}, error) {
	return ctrl.configuredConnection.Check(ctx, r.ConnectionID)
}

func (ctrl ConfiguredConnection) RefreshDiscovery(ctx context.Context, r *request.ConfiguredConnectionCheck) (interface{}, error) {
	return ctrl.configuredConnection.RefreshDiscovery(ctx, r.ConnectionID)
}

func (ctrl ConfiguredConnection) makeFilterPayload(ctx context.Context, set types.ConfiguredConnectionSet, f types.ConfiguredConnectionFilter, err error) (*configuredConnectionSetPayload, error) {
	if err != nil {
		return nil, err
	}

	out := &configuredConnectionSetPayload{
		Filter: f,
		Set:    make([]*configuredConnectionPayload, 0, len(set)),
	}

	for _, c := range set {
		p, _ := ctrl.makePayload(ctx, c, nil)
		out.Set = append(out.Set, p)
	}

	return out, nil
}

func (ctrl ConfiguredConnection) makePayload(ctx context.Context, c *types.ConfiguredConnection, err error) (*configuredConnectionPayload, error) {
	if err != nil || c == nil {
		return nil, err
	}

	return &configuredConnectionPayload{
		ConfiguredConnection:          c,
		CanUpdateConfiguredConnection: ctrl.ac.CanUpdateConfiguredConnection(ctx, c),
		CanDeleteConfiguredConnection: ctrl.ac.CanDeleteConfiguredConnection(ctx, c),
	}, nil
}

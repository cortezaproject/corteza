package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ConfiguredConnection struct {
		svc configuredConnectionService
		ac  configuredConnectionAccessController
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
		svc: service.DefaultConfiguredConnection,
		ac:  service.DefaultAccessControl,
	}
}

func (ctrl ConfiguredConnection) List(ctx context.Context, r *request.ConfiguredConnectionList) (interface{}, error) {
	var (
		err error
		set types.ConfiguredConnectionSet

		f = types.ConfiguredConnectionFilter{
			ConnectionID: r.ConnectionID,
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
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err = ctrl.svc.Search(ctx, f)
	if err != nil {
		return nil, err
	}

	return ctrl.makeFilterPayload(ctx, set, f)
}

func (ctrl ConfiguredConnection) Read(ctx context.Context, r *request.ConfiguredConnectionRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.ConnectionID)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, res), nil
}

func (ctrl ConfiguredConnection) Delete(ctx context.Context, r *request.ConfiguredConnectionDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ConnectionID)
}

func (ctrl ConfiguredConnection) Enable(ctx context.Context, r *request.ConfiguredConnectionEnable) (interface{}, error) {
	res, err := ctrl.svc.Enable(ctx, r.ConnectionID)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, res), nil
}

func (ctrl ConfiguredConnection) Check(ctx context.Context, r *request.ConfiguredConnectionCheck) (interface{}, error) {
	return ctrl.svc.Check(ctx, r.ConnectionID)
}

func (ctrl ConfiguredConnection) RefreshDiscovery(ctx context.Context, r *request.ConfiguredConnectionCheck) (interface{}, error) {
	return ctrl.svc.RefreshDiscovery(ctx, r.ConnectionID)
}

func (ctrl ConfiguredConnection) makeFilterPayload(ctx context.Context, set types.ConfiguredConnectionSet, f types.ConfiguredConnectionFilter) (*configuredConnectionSetPayload, error) {
	out := &configuredConnectionSetPayload{
		Filter: f,
		Set:    make([]*configuredConnectionPayload, 0, len(set)),
	}

	for _, c := range set {
		out.Set = append(out.Set, ctrl.makePayload(ctx, c))
	}

	return out, nil
}

func (ctrl ConfiguredConnection) makePayload(ctx context.Context, c *types.ConfiguredConnection) *configuredConnectionPayload {
	return &configuredConnectionPayload{
		ConfiguredConnection:          c,
		CanUpdateConfiguredConnection: ctrl.ac.CanUpdateConfiguredConnection(ctx, c),
		CanDeleteConfiguredConnection: ctrl.ac.CanDeleteConfiguredConnection(ctx, c),
	}
}

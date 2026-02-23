package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	ConfiguredConnection struct {
		svc configuredConnectionService
		ac  configuredConnectionAccessController
	}

	configuredConnectionPayload struct {
		*types.ConfiguredConnection

		CanDeleteConfiguredConnection bool `json:"canDeleteConfiguredConnection"`
	}

	configuredConnectionSetPayload struct {
		Filter types.ConfiguredConnectionFilter `json:"filter"`
		Set    []*configuredConnectionPayload   `json:"set"`
	}

	configuredConnectionAccessController interface {
		CanCreateConfiguredConnection(context.Context) bool
		CanDeleteConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
	}

	configuredConnectionService interface {
		FindByID(ctx context.Context, ID uint64) (*types.ConfiguredConnection, error)
		Install(ctx context.Context, new *types.ConfiguredConnection) (*types.ConfiguredConnection, error)
		DeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.ConfiguredConnectionFilter) (types.ConfiguredConnectionSet, types.ConfiguredConnectionFilter, error)
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

	set, f, err = ctrl.svc.Search(ctx, f)
	if err != nil {
		return nil, err
	}

	return ctrl.makeFilterPayload(ctx, set, f)
}

func (ctrl ConfiguredConnection) Install(ctx context.Context, r *request.ConfiguredConnectionInstall) (interface{}, error) {
	conn := &types.ConfiguredConnection{
		ConnectionID: r.ConnectionID,
		Name:         r.Name,
		Config:       r.Config,
		Status:       "active",
	}

	res, err := ctrl.svc.Install(ctx, conn)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, res), nil
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
		CanDeleteConfiguredConnection: ctrl.ac.CanDeleteConfiguredConnection(ctx, c),
	}
}

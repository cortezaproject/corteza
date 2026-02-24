package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	configuredConnection struct {
		actionlog     actionlog.Recorder
		store         store.Storer
		ac            configuredConnectionAccessController
		connectionSvc *connection
	}

	configuredConnectionAccessController interface {
		CanSearchConfiguredConnections(ctx context.Context) bool

		CanCreateConfiguredConnection(context.Context) bool
		CanReadConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
		CanDeleteConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
	}
)

func ConfiguredConnectionSvc() *configuredConnection {
	return &configuredConnection{
		actionlog:     DefaultActionlog,
		store:         DefaultStore,
		ac:            DefaultAccessControl,
		connectionSvc: DefaultConnection,
	}
}

func (svc *configuredConnection) FindByID(ctx context.Context, ID uint64) (res *types.ConfiguredConnection, err error) {
	res, err = store.LookupConfiguredConnectionByID(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadConfiguredConnection(ctx, res) {
		return nil, errors.Unauthorized("connection read denied")
	}

	return
}

func (svc *configuredConnection) DeleteByID(ctx context.Context, ID uint64) (err error) {
	return ConfiguredConnectionErrDeletionNotSupported()
}

func (svc *configuredConnection) Search(ctx context.Context, filter types.ConfiguredConnectionFilter) (set types.ConfiguredConnectionSet, f types.ConfiguredConnectionFilter, err error) {
	if !svc.ac.CanSearchConfiguredConnections(ctx) {
		return nil, f, errors.Unauthorized("connection search denied")
	}

	set, f, err = store.SearchConfiguredConnections(ctx, svc.store, filter)
	return
}

func loadConfiguredConnection(ctx context.Context, s store.ConfiguredConnections, ID uint64) (res *types.ConfiguredConnection, err error) {
	if ID == 0 {
		return nil, errors.NotFound("connection not found")
	}

	if res, err = store.LookupConfiguredConnectionByID(ctx, s, ID); errors.IsNotFound(err) {
		err = errors.NotFound("connection not found")
	}

	return
}

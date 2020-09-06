package service

import (
	"context"

	composeService "github.com/cortezaproject/corteza-server/compose/service"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	module struct {
		ctx     context.Context
		fetcher composeService.ModuleService
		// actionlog actionlog.Recorder
		// ac        moduleAccessController
		store store.Storable
	}

	ModuleService interface {
		Find(ctx context.Context, filter types.ModuleFilter) (composeTypes.ModuleSet, error)
	}
)

func Module() ModuleService {
	return &module{
		ctx: context.Background(),
		// ac:       DefaultAccessControl,
		// eventbus: eventbus.Service(),
		store: composeService.DefaultNgStore,
	}
}

func (svc module) Find(ctx context.Context, filter types.ModuleFilter) (set composeTypes.ModuleSet, err error) {
	return composeTypes.ModuleSet{}, nil
}

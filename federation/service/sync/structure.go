package sync

import (
	"context"

	composeService "github.com/cortezaproject/corteza-server/compose/service"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	module struct {
		ctx        context.Context
		compose    composeService.ModuleService
		federation service.ModuleService
		store      store.Storable
	}

	ModuleService interface {
		FindForNode(ctx context.Context, filter types.ModuleFilter) (composeTypes.ModuleSet, error)
	}
)

func Module() ModuleService {
	return &module{
		ctx:        context.Background(),
		store:      composeService.DefaultNgStore,
		compose:    composeService.Module(),
		federation: service.Module(),
	}
}

func (svc module) FindForNode(ctx context.Context, filter types.ModuleFilter) (set composeTypes.ModuleSet, err error) {
	// get all modules per-node
	// feed the id's into the compose moduleservice
	// get the data
	// transform (but not here)
	return composeTypes.ModuleSet{}, nil
}

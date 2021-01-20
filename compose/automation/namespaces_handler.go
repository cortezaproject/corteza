package automation

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	moduleService interface {
		FindByID(ctx context.Context, namespaceID, moduleID uint64) (*types.Module, error)
		FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*types.Module, error)
		Find(ctx context.Context, filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error)

		Create(ctx context.Context, module *types.Module) (*types.Module, error)
		Update(ctx context.Context, module *types.Module) (*types.Module, error)

		DeleteByID(ctx context.Context, namespaceID uint64, moduleID ...uint64) error
	}

	moduleNamespaceService interface {
		FindByID(ctx context.Context, namespaceID uint64) (*types.Namespace, error)
		FindByHandle(ctx context.Context, handle string) (*types.Namespace, error)
	}

	modulesHandlers struct {
		reg modulesHandlerRegistry
		ns  moduleNamespaceService
		mod moduleService
	}
)

func ModulesHandlers(reg modulesHandlerRegistry, ns moduleNamespaceService, rec moduleService) *modulesHandlers {
	h := &modulesHandlers{
		reg: reg,
		ns:  ns,
		mod: rec,
	}

	h.register()
	return h
}

func (h modulesHandlers) resolveNamespace(ctx context.Context, id *uint64, handle string, res *types.Namespace) (err error) {
	if *id == 0 {
		if len(handle) > 0 {
			if res, err = h.ns.FindByHandle(ctx, handle); err != nil {
				return
			}
		}

		if res != nil {
			*id = res.ID
		}
	}

	return
}

func (h modulesHandlers) lookup(ctx context.Context, args *modulesLookupArgs) (results *modulesLookupResults, err error) {
	results = &modulesLookupResults{}

	if err = h.resolveNamespace(ctx, &args.namespaceID, args.namespaceHandle, args.namespaceRes); err != nil {
		return nil, err
	}

	if args.moduleID > 0 {
		results.Module, err = h.mod.FindByID(ctx, args.namespaceID, args.moduleID)
	} else {
		results.Module, err = h.mod.FindByHandle(ctx, args.namespaceID, args.moduleHandle)
	}

	return
}

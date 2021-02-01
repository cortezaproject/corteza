package automation

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	moduleService interface {
		FindByID(ctx context.Context, namespaceID, moduleID uint64) (*types.Module, error)
		FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*types.Module, error)
		Find(ctx context.Context, filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error)

		Create(ctx context.Context, module *types.Module) (*types.Module, error)
		Update(ctx context.Context, module *types.Module) (*types.Module, error)

		DeleteByID(ctx context.Context, namespaceID uint64, moduleID uint64) error
	}

	modulesHandler struct {
		reg modulesHandlerRegistry
		ns  namespaceService
		mod moduleService
	}

	moduleLookup interface {
		GetModule() (bool, uint64, string, *types.Module)
	}
)

func ModulesHandler(reg modulesHandlerRegistry, ns namespaceService, mod moduleService) *modulesHandler {
	h := &modulesHandler{
		reg: reg,
		ns:  ns,
		mod: mod,
	}

	h.register()
	return h
}

func (h modulesHandler) lookup(ctx context.Context, args *modulesLookupArgs) (results *modulesLookupResults, err error) {
	results = &modulesLookupResults{}
	results.Module, err = lookupModule(ctx, h.ns, h.mod, args)
	return
}

func getModuleID(ctx context.Context, nsSvc namespaceService, modSvc moduleService, args moduleLookup) (namespaceID uint64, moduleID uint64, err error) {
	namespaceID, err = getNamespaceID(ctx, nsSvc, args.(namespaceLookup))
	if err != nil {
		return
	}

	var mod *types.Module
	if _, moduleID, _, _ = args.GetModule(); moduleID > 0 {
		return
	} else if mod, err = lookupModule(ctx, nsSvc, modSvc, args); err != nil {
		return
	} else {
		return namespaceID, mod.ID, nil
	}
}

func lookupModule(ctx context.Context, nsSvc namespaceService, modSvc moduleService, args moduleLookup) (*types.Module, error) {
	namespaceID, err := getNamespaceID(ctx, nsSvc, args.(namespaceLookup))
	if err != nil {
		return nil, fmt.Errorf("could not load namespace: %w", err)
	}

	_, ID, handle, module := args.GetModule()

	switch {
	case module != nil:
		return module, nil
	case ID > 0:
		return modSvc.FindByID(ctx, namespaceID, ID)
	case len(handle) > 0:
		return modSvc.FindByHandle(ctx, namespaceID, handle)
	}

	return nil, fmt.Errorf("empty module lookup params")
}

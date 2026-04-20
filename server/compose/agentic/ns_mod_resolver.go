package agentic

import (
	"context"
	"fmt"

	cmpService "github.com/crusttech/human/server/compose/service"
	agenticRuntime "github.com/crusttech/human/server/system/agentic/runtime"
)

type nsModResolver struct{}

func NsModResolver() *nsModResolver {
	return &nsModResolver{}
}

func (r *nsModResolver) Resolve(ctx context.Context, namespace, module string) (nsID, modID uint64, err error) {
	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, namespace)
	if err != nil {
		return 0, 0, fmt.Errorf("namespace lookup failed: %w", err)
	}

	if module == "" {
		return ns.ID, 0, nil
	}

	mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, module)
	if err != nil {
		return 0, 0, fmt.Errorf("module lookup failed: %w", err)
	}

	return ns.ID, mod.ID, nil
}

func (r *nsModResolver) LookupNamespace(ctx context.Context, id uint64) (agenticRuntime.NsHandle, error) {
	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, id)
	if err != nil {
		return agenticRuntime.NsHandle{}, err
	}
	handle := ns.Slug
	if handle == "" {
		handle = ns.Name
	}
	return agenticRuntime.NsHandle{ID: ns.ID, Handle: handle}, nil
}

func (r *nsModResolver) LookupModule(ctx context.Context, nsID, modID uint64) (agenticRuntime.ModHandle, error) {
	mod, err := cmpService.DefaultModule.FindByAny(ctx, nsID, modID)
	if err != nil {
		return agenticRuntime.ModHandle{}, err
	}
	modHandle := mod.Handle
	if modHandle == "" {
		modHandle = mod.Name
	}
	return agenticRuntime.ModHandle{ID: mod.ID, NamespaceID: mod.NamespaceID, Handle: modHandle}, nil
}

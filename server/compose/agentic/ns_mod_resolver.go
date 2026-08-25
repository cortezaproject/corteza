package agentic

import (
	"context"
	"fmt"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
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
	return agenticRuntime.NsHandle{ID: ns.ID, Handle: handle, Name: ns.Name, CreatedByAgent: ns.CreatedByAgent}, nil
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
	return agenticRuntime.ModHandle{ID: mod.ID, NamespaceID: mod.NamespaceID, Handle: modHandle, Name: mod.Name, CreatedByAgent: mod.CreatedByAgent}, nil
}

// lookupNamespaceArg resolves the "namespace" argument, naming it when it is
// simply absent.
//
// FindByAny is given the raw argument, so a missing one arrives as nil and
// comes back as "invalid ID" — which reads as though a namespace was supplied
// and rejected. A model told that supplies an ID it invented, rather than the
// namespace it forgot.
func lookupNamespaceArg(ctx context.Context, args map[string]any) (*cmpTypes.Namespace, error) {
	if v, ok := args["namespace"]; !ok || v == nil || v == "" {
		return nil, fmt.Errorf("namespace is required — its handle, slug, name, or ID as a string")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, toolkit.Errf("namespace lookup", err)
	}

	return ns, nil
}

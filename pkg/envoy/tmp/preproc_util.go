package tmp

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

func filterComposeNamespaceResources(rr []resource.Interface) []*resource.ComposeNamespace {
	nn := make([]*resource.ComposeNamespace, 0, len(rr))
	for _, r := range rr {
		switch n := r.(type) {
		case *resource.ComposeNamespace:
			nn = append(nn, n)
		}
	}

	return nn
}

func filterComposeModuleResources(rr []resource.Interface) []*resource.ComposeModule {
	nn := make([]*resource.ComposeModule, 0, len(rr))
	for _, r := range rr {
		switch n := r.(type) {
		case *resource.ComposeModule:
			nn = append(nn, n)
		}
	}

	return nn
}

func findMissingComposeNamespace(ctx context.Context, s store.Storer, missing resource.RefSet) (*types.Namespace, resource.Identifiers, error) {
	// Check the store if we can find it
	nsDd := missing.FilterByResourceType("compose:namespace")
	if len(nsDd) > 0 {
		nsD := nsDd[0]
		ns, err := loadComposeNamespace(ctx, s, namespaceFilterFromGeneric(GetGenericFilter(nsD.Identifiers)))
		if err != nil {
			return nil, nil, err
		}
		if ns != nil {
			return ns, nsD.Identifiers, nil
		}
	}

	return nil, nil, nil
}

func findMissingComposeModule(ctx context.Context, s store.Storer, missing resource.RefSet) (*types.Module, resource.Identifiers, error) {
	// Check the store if we can find it
	modDd := missing.FilterByResourceType("compose:module")
	if len(modDd) > 0 {
		modD := modDd[0]
		mod, err := loadComposeModule(ctx, s, moduleFilterFromGeneric(GetGenericFilter(modD.Identifiers)))
		if err != nil {
			return nil, nil, err
		}
		if mod != nil {
			return mod, modD.Identifiers, nil
		}
	}

	return nil, nil, nil
}

package automation

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	namespaceService interface {
		FindByID(ctx context.Context, namespaceID uint64) (*types.Namespace, error)
		FindByHandle(ctx context.Context, handle string) (*types.Namespace, error)
		Find(ctx context.Context, filter types.NamespaceFilter) (set types.NamespaceSet, f types.NamespaceFilter, err error)

		Create(ctx context.Context, namespace *types.Namespace) (*types.Namespace, error)
		Update(ctx context.Context, namespace *types.Namespace) (*types.Namespace, error)

		DeleteByID(ctx context.Context, namespaceID uint64) error
	}

	namespacesHandler struct {
		reg namespacesHandlerRegistry
		ns  namespaceService
	}

	namespaceLookup interface {
		GetNamespace() (bool, uint64, string, *types.Namespace)
	}
)

func NamespacesHandler(reg namespacesHandlerRegistry, ns namespaceService) *namespacesHandler {
	h := &namespacesHandler{
		reg: reg,
		ns:  ns,
	}

	h.register()
	return h
}

func (h namespacesHandler) lookup(ctx context.Context, args *namespacesLookupArgs) (results *namespacesLookupResults, err error) {
	results = &namespacesLookupResults{}

	results.Namespace, err = lookupNamespace(ctx, h.ns, args)
	return
}

func getNamespaceID(ctx context.Context, svc namespaceService, args namespaceLookup) (uint64, error) {
	_, ID, _, _ := args.GetNamespace()
	if ID > 0 {
		return ID, nil
	} else if ns, err := lookupNamespace(ctx, svc, args); err != nil {
		return 0, err
	} else {
		return ns.ID, nil
	}
}

func lookupNamespace(ctx context.Context, svc namespaceService, args namespaceLookup) (*types.Namespace, error) {
	_, ID, handle, namespace := args.GetNamespace()

	switch {
	case namespace != nil:
		return namespace, nil
	case ID > 0:
		return svc.FindByID(ctx, ID)
	case len(handle) > 0:
		return svc.FindByHandle(ctx, handle)
	}

	return nil, fmt.Errorf("empty namespace lookup params")
}

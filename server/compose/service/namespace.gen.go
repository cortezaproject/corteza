package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/label"
)

func (svc *namespace) FindByID(ctx context.Context, ID uint64) (res *types.Namespace, err error) {
	var (
		aProps = &namespaceActionProps{namespace: &types.Namespace{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, NamespaceActionLookup, err)
}

func (svc *namespace) Search(ctx context.Context, filter types.NamespaceFilter) (set types.NamespaceSet, f types.NamespaceFilter, err error) {
	var (
		aProps = &namespaceActionProps{filter: &filter}
	)

	err = func() error {
		if !svc.ac.CanSearchNamespaces(ctx) {
			return NamespaceErrNotAllowedToSearch()
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, NamespaceActionSearch, err)
}

func (svc *namespace) Create(ctx context.Context, new *types.Namespace) (res *types.Namespace, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &namespaceActionProps{namespace: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, NamespaceActionCreate, err)
}

func (svc *namespace) Update(ctx context.Context, upd *types.Namespace) (res *types.Namespace, err error) {
	var (
		aProps = &namespaceActionProps{changed: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, NamespaceActionUpdate, err)
}

func (svc *namespace) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &namespaceActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, NamespaceActionDelete, err)
}

// toLabeledNamespaces converts to []label.LabeledResource
func toLabeledNamespaces(set []*types.Namespace) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

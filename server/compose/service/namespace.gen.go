package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	types "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
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
		old    *types.Namespace
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadNamespace(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setNamespace(res)
		aProps.setChanged(res)
		old = res.Clone()

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return NamespaceErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		return svc.onUpdate(ctx, s, upd, res, aProps, before, after)
	})

	return res, svc.recordAction(ctx, aProps, NamespaceActionUpdate, err, old, res)
}

func (svc *namespace) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &namespaceActionProps{}
		res    *types.Namespace
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadNamespace(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setNamespace(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, NamespaceActionDelete, err)
}

func loadNamespace(ctx context.Context, s store.ComposeNamespaces, ID uint64) (res *types.Namespace, err error) {
	if ID == 0 {
		return nil, NamespaceErrInvalidID()
	}

	if res, err = store.LookupComposeNamespaceByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, NamespaceErrNotFound()
	}

	return
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

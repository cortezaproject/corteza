package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	types "github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
)

func (svc *node) Search(ctx context.Context, filter types.NodeFilter) (set types.NodeSet, f types.NodeFilter, err error) {
	var (
		aProps = &nodeActionProps{filter: &filter}
	)

	err = func() error {
		if !svc.ac.CanSearchNodes(ctx) {
			return NodeErrNotAllowedToSearch()
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, NodeActionSearch, err)
}

func (svc *node) Create(ctx context.Context, new *types.Node) (res *types.Node, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &nodeActionProps{node: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateNode(ctx) {
			return NodeErrNotAllowedToCreate()
		}
		before := func() error { return nil }
		after := func() error { return nil }
		res = new
		return svc.onCreate(ctx, new, before, after)
	}()

	return res, svc.recordAction(ctx, aProps, NodeActionCreate, err)
}

func (svc *node) Update(ctx context.Context, upd *types.Node) (res *types.Node, err error) {
	var (
		aProps = &nodeActionProps{node: upd}
		old    *types.Node
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadNode(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setNode(res)
		aProps.setNode(res)
		old = res.Clone()

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return NodeErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Name = upd.Name
		res.BaseURL = upd.BaseURL
		res.Status = upd.Status
		res.Contact = upd.Contact
		res.PairToken = upd.PairToken
		res.AuthToken = upd.AuthToken
		res.CreatedBy = upd.CreatedBy
		res.UpdatedBy = upd.UpdatedBy
		res.DeletedBy = upd.DeletedBy
		res.UpdatedAt = now()

		if err = store.UpdateFederationNode(ctx, s, res); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, NodeActionUpdate, err, old, res)
}

func (svc *node) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &nodeActionProps{}
		res    *types.Node
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadNode(ctx, s, ID); err != nil {
			return
		}

		aProps.setNode(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, NodeActionDelete, err)
}

func (svc *node) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &nodeActionProps{}
		res    *types.Node
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadNode(ctx, s, ID); err != nil {
			return
		}

		aProps.setNode(res)

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, NodeActionUndelete, err)
}

func loadNode(ctx context.Context, s store.FederationNodes, ID uint64) (res *types.Node, err error) {
	if ID == 0 {
		return nil, NodeErrInvalidID()
	}

	if res, err = store.LookupFederationNodeByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, NodeErrNotFound()
	}

	return
}

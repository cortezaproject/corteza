package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/federation/types"
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
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, NodeActionCreate, err)
}

func (svc *node) Update(ctx context.Context, upd *types.Node) (res *types.Node, err error) {
	var (
		aProps = &nodeActionProps{node: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, NodeActionUpdate, err)
}

func (svc *node) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &nodeActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, NodeActionDelete, err)
}

func (svc *node) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &nodeActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, NodeActionUndelete, err)
}

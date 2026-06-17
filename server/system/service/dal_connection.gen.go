package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *dalConnection) FindByID(ctx context.Context, ID uint64) (res *types.DalConnection, err error) {
	var (
		aProps = &dalConnectionActionProps{connection: &types.DalConnection{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, DalConnectionActionLookup, err)
}

func (svc *dalConnection) Search(ctx context.Context, filter types.DalConnectionFilter) (set types.DalConnectionSet, f types.DalConnectionFilter, err error) {
	var (
		aProps = &dalConnectionActionProps{search: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, DalConnectionActionSearch, err)
}

func (svc *dalConnection) Create(ctx context.Context, new *types.DalConnection) (res *types.DalConnection, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &dalConnectionActionProps{connection: new, new: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, DalConnectionActionCreate, err)
}

func (svc *dalConnection) Update(ctx context.Context, upd *types.DalConnection) (res *types.DalConnection, err error) {
	var (
		aProps = &dalConnectionActionProps{update: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, DalConnectionActionUpdate, err)
}

func (svc *dalConnection) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &dalConnectionActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, DalConnectionActionDelete, err)
}

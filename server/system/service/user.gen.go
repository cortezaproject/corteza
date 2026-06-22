package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *user) FindByID(ctx context.Context, ID uint64) (res *types.User, err error) {
	var (
		aProps = &userActionProps{user: &types.User{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, UserActionLookup, err)
}

func (svc *user) Search(ctx context.Context, filter types.UserFilter) (set types.UserSet, f types.UserFilter, err error) {
	var (
		aProps = &userActionProps{filter: &filter}
	)

	err = func() error {
		if !svc.ac.CanSearchUsers(ctx) {
			return UserErrNotAllowedToSearch()
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, UserActionSearch, err)
}

func (svc *user) Create(ctx context.Context, new *types.User) (res *types.User, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &userActionProps{user: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateUser(ctx) {
			return UserErrNotAllowedToCreate()
		}
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, UserActionCreate, err)
}

func (svc *user) Update(ctx context.Context, upd *types.User) (res *types.User, err error) {
	var (
		aProps = &userActionProps{update: upd}
		old    *types.User
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadUser(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setUser(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return UserErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return UserErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		return svc.onUpdate(ctx, s, upd, res, aProps, before, after)
	})

	return res, svc.recordAction(ctx, aProps, UserActionUpdate, err, old, res)
}

func (svc *user) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &userActionProps{}
		res    *types.User
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadUser(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setUser(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, UserActionDelete, err)
}

func (svc *user) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &userActionProps{}
		res    *types.User
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadUser(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setUser(res)

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, UserActionUndelete, err)
}

func loadUser(ctx context.Context, s store.Users, ID uint64) (res *types.User, err error) {
	if ID == 0 {
		return nil, UserErrInvalidID()
	}

	if res, err = store.LookupUserByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, UserErrNotFound()
	}

	return
}

// toLabeledUsers converts to []label.LabeledResource
func toLabeledUsers(set []*types.User) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

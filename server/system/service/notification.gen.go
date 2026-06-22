package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *notification) FindByID(ctx context.Context, ID uint64) (res *types.Notification, err error) {
	var (
		aProps = &notificationActionProps{notification: &types.Notification{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, NotificationActionLookup, err)
}

func (svc *notification) Search(ctx context.Context, filter types.NotificationFilter) (set types.NotificationSet, f types.NotificationFilter, err error) {
	var (
		aProps = &notificationActionProps{filter: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, NotificationActionSearch, err)
}

func (svc *notification) Create(ctx context.Context, new *types.Notification) (res *types.Notification, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &notificationActionProps{notification: new, new: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, NotificationActionCreate, err)
}

func (svc *notification) Update(ctx context.Context, upd *types.Notification) (res *types.Notification, err error) {
	var (
		aProps = &notificationActionProps{updated: upd}
		old    *types.Notification
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadNotification(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setNotification(res)
		aProps.setUpdated(res)
		old = res.Clone()

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return NotificationErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		return svc.onUpdate(ctx, s, upd, res, aProps, before, after)
	})

	return res, svc.recordAction(ctx, aProps, NotificationActionUpdate, err, old, res)
}

func (svc *notification) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &notificationActionProps{}
		res    *types.Notification
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadNotification(ctx, s, ID); err != nil {
			return
		}

		aProps.setNotification(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, NotificationActionDelete, err)
}

func loadNotification(ctx context.Context, s store.Notifications, ID uint64) (res *types.Notification, err error) {
	if ID == 0 {
		return nil, NotificationErrInvalidID()
	}

	if res, err = store.LookupNotificationByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, NotificationErrNotFound()
	}

	return
}

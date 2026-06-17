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

func (svc *reminder) FindByID(ctx context.Context, ID uint64) (res *types.Reminder, err error) {
	var (
		aProps = &reminderActionProps{reminder: &types.Reminder{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ReminderActionLookup, err)
}

func (svc *reminder) Search(ctx context.Context, filter types.ReminderFilter) (set types.ReminderSet, f types.ReminderFilter, err error) {
	var (
		aProps = &reminderActionProps{filter: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ReminderActionSearch, err)
}

func (svc *reminder) Create(ctx context.Context, new *types.Reminder) (res *types.Reminder, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &reminderActionProps{reminder: new, new: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ReminderActionCreate, err)
}

func (svc *reminder) Update(ctx context.Context, upd *types.Reminder) (res *types.Reminder, err error) {
	var (
		aProps = &reminderActionProps{updated: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ReminderActionUpdate, err)
}

func (svc *reminder) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &reminderActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, ReminderActionDelete, err)
}

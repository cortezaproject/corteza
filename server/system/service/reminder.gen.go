package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
	"time"
)

type reminder struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        reminderAccessController
	services  *reminderServices
}

func (svc *reminder) FindByID(ctx context.Context, ID uint64) (res *types.Reminder, err error) {
	var (
		aProps = &reminderActionProps{reminder: &types.Reminder{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
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
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
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
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ReminderActionCreate, err)
}

func (svc *reminder) Update(ctx context.Context, upd *types.Reminder) (res *types.Reminder, err error) {
	var (
		aProps = &reminderActionProps{updated: upd}
		old    *types.Reminder
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadReminder(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setReminder(res)
		aProps.setUpdated(res)
		old = res.Clone()

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ReminderErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Resource = upd.Resource
		res.Payload = upd.Payload
		res.AssignedTo = upd.AssignedTo
		res.AssignedBy = upd.AssignedBy
		res.UpdatedAt = now()

		if err = store.UpdateReminder(ctx, s, res); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, ReminderActionUpdate, err, old, res)
}

func (svc *reminder) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &reminderActionProps{}
		res    *types.Reminder
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadReminder(ctx, s, ID); err != nil {
			return
		}

		aProps.setReminder(res)

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ReminderActionDelete, err)
}

func loadReminder(ctx context.Context, s store.Reminders, ID uint64) (res *types.Reminder, err error) {
	if ID == 0 {
		return nil, ReminderErrInvalidID()
	}

	if res, err = store.LookupReminderByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ReminderErrNotFound()
	}

	return
}
func (svc *reminder) guard(_ context.Context, _ *types.Reminder) error { return nil }

func (svc *reminder) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *reminder) FindByIDs(ctx context.Context, IDs []uint64) (rr types.ReminderSet, err error) {
	var (
		aProps = &reminderActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		rr, err = svc.onFindByIDs(ctx, aProps, IDs)
		return err
	}()

	return rr, svc.recordAction(ctx, aProps, ReminderActionFindByIDs, err)
}

func (svc *reminder) Dismiss(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &reminderActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onDismiss(ctx, aProps, ID)
		return err
	}()

	return svc.recordAction(ctx, aProps, ReminderActionDismiss, err)
}

func (svc *reminder) Undismiss(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &reminderActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onUndismiss(ctx, aProps, ID)
		return err
	}()

	return svc.recordAction(ctx, aProps, ReminderActionDismiss, err)
}

func (svc *reminder) Snooze(ctx context.Context, ID uint64, remindAt *time.Time) (err error) {
	var (
		aProps = &reminderActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onSnooze(ctx, aProps, ID, remindAt)
		return err
	}()

	return svc.recordAction(ctx, aProps, ReminderActionSnooze, err)
}

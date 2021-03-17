package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"time"
)

type (
	reminder struct {
		ac reminderAccessController

		actionlog actionlog.Recorder
		store     store.Reminders
	}

	reminderAccessController interface {
		CanAssignReminder(ctx context.Context) bool
	}

	ReminderService interface {
		Find(context.Context, types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error)
		FindByID(context.Context, uint64) (*types.Reminder, error)

		Create(context.Context, *types.Reminder) (*types.Reminder, error)

		Update(context.Context, *types.Reminder) (*types.Reminder, error)

		Dismiss(context.Context, uint64) error
		Snooze(context.Context, uint64, *time.Time) error

		Delete(context.Context, uint64) error
	}
)

func Reminder(ctx context.Context) ReminderService {
	return &reminder{
		ac:    DefaultAccessControl,
		store: DefaultStore,
	}
}

func (svc reminder) Find(ctx context.Context, filter types.ReminderFilter) (rr types.ReminderSet, f types.ReminderFilter, err error) {
	var (
		raProps = &reminderActionProps{filter: &filter}
	)

	err = func() (err error) {
		rr, f, err = store.SearchReminders(ctx, svc.store, filter)
		if err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, ReminderActionSearch, err)
}

func (svc reminder) FindByID(ctx context.Context, ID uint64) (r *types.Reminder, err error) {
	var (
		raProps = &reminderActionProps{reminder: &types.Reminder{ID: ID}}
	)

	err = func() (err error) {
		if ID == 0 {
			return ReminderErrInvalidID()
		}

		r, err = store.LookupReminderByID(ctx, svc.store, ID)
		if err != nil {
			return err
		}

		raProps.setReminder(r)

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ReminderActionLookup, err)
}

func (svc reminder) FindByIDs(ctx context.Context, IDs ...uint64) (rr types.ReminderSet, err error) {
	if len(IDs) == 0 {
		return nil, nil
	}

	rr, _, err = svc.Find(ctx, types.ReminderFilter{ReminderID: IDs})

	return rr, nil
}

func (svc reminder) checkAssignee(ctx context.Context, rm *types.Reminder) (err error) {
	// Check if user is assigning to someone else
	if rm.AssignedTo != svc.currentUser(ctx) {
		if !svc.ac.CanAssignReminder(ctx) {
			return ReminderErrNotAllowedToAssign()
		}
	}

	return nil
}

func (svc reminder) currentUser(ctx context.Context) uint64 {
	return intAuth.GetIdentityFromContext(ctx).Identity()
}

func (svc reminder) Create(ctx context.Context, new *types.Reminder) (r *types.Reminder, err error) {
	var (
		raProps = &reminderActionProps{new: new}
	)

	err = func() (err error) {
		if err := svc.checkAssignee(ctx, new); err != nil {
			return err
		}

		r = new
		r.ID = nextID()
		r.CreatedAt = *now()

		if err = store.CreateReminder(ctx, svc.store, r); err != nil {
			return err
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ReminderActionUpdate, err)

}

func (svc reminder) Update(ctx context.Context, upd *types.Reminder) (r *types.Reminder, err error) {
	var (
		raProps = &reminderActionProps{updated: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return ReminderErrInvalidID()
		}

		if r, err = store.LookupReminderByID(ctx, svc.store, upd.ID); err != nil {
			return
		}

		if err := svc.checkAssignee(ctx, upd); err != nil {
			return err
		}

		// Assign changed values
		if upd.AssignedTo != r.AssignedTo {
			r.AssignedTo = upd.AssignedTo
			r.AssignedBy = svc.currentUser(ctx)
			r.AssignedAt = time.Now()
		}

		r.Payload = upd.Payload
		r.RemindAt = upd.RemindAt
		r.Resource = upd.Resource
		r.UpdatedAt = now()

		if err = store.UpdateReminder(ctx, svc.store, r); err != nil {
			return err
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ReminderActionUpdate, err)
}

func (svc reminder) Dismiss(ctx context.Context, ID uint64) (err error) {
	var (
		r *types.Reminder

		raProps = &reminderActionProps{reminder: &types.Reminder{ID: ID}}
	)

	err = func() (err error) {
		if ID == 0 {
			return ReminderErrInvalidID()
		}

		if r, err = store.LookupReminderByID(ctx, svc.store, ID); err != nil {
			return ReminderErrNotFound()
		}

		raProps.setReminder(r)

		// Assign changed values
		n := time.Now()
		r.DismissedAt = &n
		r.DismissedBy = svc.currentUser(ctx)

		if err = store.UpdateReminder(ctx, svc.store, r); err != nil {
			return err
		}

		return nil
	}()

	return svc.recordAction(ctx, raProps, ReminderActionDismiss, err)
}

func (svc reminder) Snooze(ctx context.Context, ID uint64, remindAt *time.Time) (err error) {
	var (
		r *types.Reminder

		raProps = &reminderActionProps{reminder: &types.Reminder{ID: ID, RemindAt: remindAt}}
	)

	err = func() (err error) {
		if ID == 0 {
			return ReminderErrInvalidID()
		}

		if r, err = store.LookupReminderByID(ctx, svc.store, ID); err != nil {
			return ReminderErrNotFound()
		}

		raProps.setReminder(r)

		// Assign changed values
		r.SnoozeCount++
		r.RemindAt = remindAt

		if err = store.UpdateReminder(ctx, svc.store, r); err != nil {
			return err
		}

		return nil
	}()

	return svc.recordAction(ctx, raProps, ReminderActionSnooze, err)
}

func (svc reminder) Delete(ctx context.Context, ID uint64) (err error) {
	var (
		r *types.Reminder

		raProps = &reminderActionProps{reminder: &types.Reminder{ID: ID}}
	)

	err = func() (err error) {
		if ID == 0 {
			return ReminderErrInvalidID()
		}

		if r, err = svc.FindByID(ctx, ID); err != nil {
			return ReminderErrNotFound()
		}

		r.DeletedAt = now()

		raProps.setReminder(r)

		return store.UpdateReminder(ctx, svc.store, r)
	}()

	return svc.recordAction(ctx, raProps, ReminderActionDelete, err)
}

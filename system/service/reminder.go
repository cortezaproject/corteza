package service

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/system/repository"

	"github.com/titpetric/factory"

	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	reminder struct {
		db *factory.DB
		ac reminderAccessController

		actionlog actionlog.Recorder
		reminder  repository.ReminderRepository
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
	db := repository.DB(ctx)
	return &reminder{
		db:       db,
		ac:       DefaultAccessControl,
		reminder: repository.Reminder(ctx, db),
	}
}

func (svc reminder) Find(ctx context.Context, filter types.ReminderFilter) (rr types.ReminderSet, f types.ReminderFilter, err error) {
	var (
		raProps = &reminderActionProps{filter: &filter}
	)

	err = svc.db.With(ctx).Transaction(func() (err error) {
		rr, f, err = svc.reminder.Find(filter)
		if err != nil {
			return err
		}

		return nil
	})

	return rr, f, svc.recordAction(ctx, raProps, ReminderActionSearch, err)
}

func (svc reminder) FindByID(ctx context.Context, ID uint64) (r *types.Reminder, err error) {
	var (
		raProps = &reminderActionProps{reminder: &types.Reminder{ID: ID}}
	)

	err = svc.db.With(ctx).Transaction(func() (err error) {
		if ID == 0 {
			return ReminderErrInvalidID()
		}

		r, err = svc.reminder.FindByID(ID)
		if err != nil {
			return err
		}

		raProps.setReminder(r)

		return nil
	})

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

	err = svc.db.With(ctx).Transaction(func() (err error) {
		if err := svc.checkAssignee(ctx, new); err != nil {
			return err
		}

		if r, err = svc.reminder.Create(new); err != nil {
			return err
		}

		return nil
	})

	return r, svc.recordAction(ctx, raProps, ReminderActionUpdate, err)

}

func (svc reminder) Update(ctx context.Context, upd *types.Reminder) (r *types.Reminder, err error) {
	var (
		raProps = &reminderActionProps{updated: upd}
	)

	err = svc.db.With(ctx).Transaction(func() (err error) {
		if upd.ID == 0 {
			return ReminderErrInvalidID()
		}

		if r, err = svc.reminder.FindByID(upd.ID); err != nil {
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

		if r, err = svc.reminder.Update(r); err != nil {
			return err
		}

		return nil
	})

	return r, svc.recordAction(ctx, raProps, ReminderActionUpdate, err)
}

func (svc reminder) Dismiss(ctx context.Context, ID uint64) (err error) {
	var (
		r *types.Reminder

		raProps = &reminderActionProps{reminder: &types.Reminder{ID: ID}}
	)

	err = svc.db.With(ctx).Transaction(func() (err error) {
		if ID == 0 {
			return ReminderErrInvalidID()
		}

		if r, err = svc.reminder.FindByID(ID); err != nil {
			return ReminderErrNotFound()
		}

		raProps.setReminder(r)

		// Assign changed values
		n := time.Now()
		r.DismissedAt = &n
		r.DismissedBy = svc.currentUser(ctx)

		if r, err = svc.reminder.Update(r); err != nil {
			return err
		}

		return nil
	})

	return svc.recordAction(ctx, raProps, ReminderActionDismiss, err)
}

func (svc reminder) Snooze(ctx context.Context, ID uint64, remindAt *time.Time) (err error) {
	var (
		r *types.Reminder

		raProps = &reminderActionProps{reminder: &types.Reminder{ID: ID, RemindAt: remindAt}}
	)

	err = svc.db.With(ctx).Transaction(func() (err error) {
		if ID == 0 {
			return ReminderErrInvalidID()
		}

		if r, err = svc.reminder.FindByID(ID); err != nil {
			return ReminderErrNotFound()
		}

		raProps.setReminder(r)

		// Assign changed values
		r.SnoozeCount++
		r.RemindAt = remindAt

		if r, err = svc.reminder.Update(r); err != nil {
			return err
		}

		return nil
	})

	return svc.recordAction(ctx, raProps, ReminderActionSnooze, err)
}

func (svc reminder) Delete(ctx context.Context, ID uint64) (err error) {
	var (
		r *types.Reminder

		raProps = &reminderActionProps{reminder: &types.Reminder{ID: ID}}
	)

	err = svc.db.With(ctx).Transaction(func() (err error) {
		if ID == 0 {
			return ReminderErrInvalidID()
		}

		if r, err = svc.FindByID(ctx, ID); err != nil {
			return ReminderErrNotFound()
		}

		raProps.setReminder(r)

		return svc.reminder.Delete(ID)
	})

	return svc.recordAction(ctx, raProps, ReminderActionDelete, err)
}

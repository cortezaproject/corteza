package service

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/system/repository"

	intAuth "github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/titpetric/factory"
)

type (
	reminder struct {
		db *factory.DB
		ac reminderAccessController

		reminder repository.ReminderRepository
	}

	reminderAccessController interface {
		CanReadAnyReminder(ctx context.Context) bool
		CanCreateReminder(ctx context.Context) bool
		CanReadReminder(ctx context.Context, rm *types.Reminder) bool
		CanUpdateReminder(ctx context.Context, rm *types.Reminder) bool
		CanDeleteReminder(ctx context.Context, rm *types.Reminder) bool
	}

	ReminderService interface {
		Find(context.Context, types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error)
		FindByID(context.Context, uint64) (*types.Reminder, error)
		FindByIDs(context.Context, ...uint64) (types.ReminderSet, error)

		Create(context.Context, *types.Reminder) (*types.Reminder, error)

		Update(context.Context, *types.Reminder) (*types.Reminder, error)

		Dismiss(context.Context, uint64) error
		Snooze(context.Context, uint64, time.Time) error

		Delete(context.Context, uint64) error
	}
)

func Reminder(ctx context.Context) ReminderService {
	db := repository.DB(ctx)
	return (&reminder{
		db:       db,
		ac:       DefaultAccessControl,
		reminder: repository.Reminder(ctx, db),
	})
}

func (svc reminder) Find(ctx context.Context, f types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error) {
	f.PageFilter.NormalizePerPageNoMax()

	f.AccessCheck = permissions.InitAccessCheckFilter(
		"read",
		intAuth.GetIdentityFromContext(ctx).Roles(),
		svc.ac.CanReadAnyReminder(ctx),
	)

	rr, f, err := svc.reminder.Find(f)
	if err != nil {
		return nil, f, err
	}

	return rr, f, nil
}

func (svc reminder) FindByID(ctx context.Context, ID uint64) (*types.Reminder, error) {
	rm, err := svc.reminder.FindByID(ID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadReminder(ctx, rm) {
		return nil, ErrNoReadPermissions
	}

	return rm, nil
}

func (svc reminder) FindByIDs(ctx context.Context, IDs ...uint64) (types.ReminderSet, error) {
	rr, err := svc.reminder.FindByIDs(IDs)
	if err != nil {
		return nil, err
	}

	ret := types.ReminderSet{}
	for _, rm := range rr {
		if svc.ac.CanReadReminder(ctx, rm) {
			ret = append(ret, rm)
		}
	}

	return ret, nil
}

func (svc reminder) Create(ctx context.Context, rm *types.Reminder) (*types.Reminder, error) {
	if !svc.ac.CanCreateReminder(ctx) {
		return nil, ErrNoCreatePermissions
	}

	return svc.reminder.Create(rm)
}

func (svc reminder) Update(ctx context.Context, rm *types.Reminder) (t *types.Reminder, err error) {
	if !svc.ac.CanUpdateReminder(ctx, rm) {
		return nil, ErrNoUpdatePermissions
	}

	return rm, svc.db.Transaction(func() (err error) {
		if t, err = svc.reminder.FindByID(rm.ID); err != nil {
			return
		}

		// Assign changed values
		if rm.AssignedTo != t.AssignedTo {
			t.AssignedTo = rm.AssignedTo
			t.AssignedBy = intAuth.GetIdentityFromContext(ctx).Identity()
			t.AssignedAt = time.Now()
		}
		t.Payload = rm.Payload
		t.RemindAt = rm.RemindAt
		t.Resource = rm.Resource

		if t, err = svc.reminder.Update(t); err != nil {
			return err
		}

		return nil
	})
}

func (svc reminder) Dismiss(ctx context.Context, ID uint64) (err error) {
	return svc.db.Transaction(func() (err error) {
		var t *types.Reminder
		if t, err = svc.reminder.FindByID(ID); err != nil {
			return err
		}

		if !svc.ac.CanUpdateReminder(ctx, t) {
			return ErrNoUpdatePermissions
		}

		// Assign changed values
		n := time.Now()
		t.DismissedAt = &n
		t.DismissedBy = intAuth.GetIdentityFromContext(ctx).Identity()

		if t, err = svc.reminder.Update(t); err != nil {
			return err
		}

		return nil
	})
}

func (svc reminder) Snooze(ctx context.Context, ID uint64, remindAt time.Time) (err error) {
	return svc.db.Transaction(func() (err error) {
		var t *types.Reminder
		if t, err = svc.reminder.FindByID(ID); err != nil {
			return err
		}

		if !svc.ac.CanUpdateReminder(ctx, t) {
			return ErrNoUpdatePermissions
		}

		// Assign changed values
		t.SnoozeCount++
		t.RemindAt = remindAt

		if t, err = svc.reminder.Update(t); err != nil {
			return err
		}

		return nil
	})
}

func (svc reminder) Delete(ctx context.Context, ID uint64) error {
	rm, err := svc.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeleteReminder(ctx, rm) {
		return ErrNoDeletePermissions
	}

	return svc.reminder.Delete(ID)
}

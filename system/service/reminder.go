package service

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/system/repository"

	"github.com/titpetric/factory"

	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	reminder struct {
		db *factory.DB
		ac reminderAccessController

		reminder repository.ReminderRepository
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

func (svc reminder) Find(ctx context.Context, f types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error) {
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

	return rm, nil
}

func (svc reminder) FindByIDs(ctx context.Context, IDs ...uint64) (types.ReminderSet, error) {
	if len(IDs) == 0 {
		return nil, nil
	}

	rr, _, err := svc.reminder.Find(types.ReminderFilter{ReminderID: IDs})
	if err != nil {
		return nil, err
	}

	return rr, nil
}

func (svc reminder) checkAssignee(ctx context.Context, rm *types.Reminder) error {
	// Check if user is assigning to someone else
	if rm.AssignedTo != svc.currentUser(ctx) {
		if !svc.ac.CanAssignReminder(ctx) {
			return ErrNoReminderAssignPermissions
		}
	}

	return nil
}

func (svc reminder) currentUser(ctx context.Context) uint64 {
	return intAuth.GetIdentityFromContext(ctx).Identity()
}

func (svc reminder) Create(ctx context.Context, rm *types.Reminder) (*types.Reminder, error) {
	if err := svc.checkAssignee(ctx, rm); err != nil {
		return nil, err
	}

	return svc.reminder.Create(rm)
}

func (svc reminder) Update(ctx context.Context, rm *types.Reminder) (t *types.Reminder, err error) {
	return rm, svc.db.With(ctx).Transaction(func() (err error) {
		if t, err = svc.reminder.FindByID(rm.ID); err != nil {
			return
		}

		if err := svc.checkAssignee(ctx, rm); err != nil {
			return err
		}

		// Assign changed values
		if rm.AssignedTo != t.AssignedTo {
			t.AssignedTo = rm.AssignedTo
			t.AssignedBy = svc.currentUser(ctx)
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
	return svc.db.With(ctx).Transaction(func() (err error) {
		var t *types.Reminder
		if t, err = svc.reminder.FindByID(ID); err != nil {
			return err
		}

		// Assign changed values
		n := time.Now()
		t.DismissedAt = &n
		t.DismissedBy = svc.currentUser(ctx)

		if t, err = svc.reminder.Update(t); err != nil {
			return err
		}

		return nil
	})
}

func (svc reminder) Snooze(ctx context.Context, ID uint64, remindAt *time.Time) (err error) {
	return svc.db.With(ctx).Transaction(func() (err error) {
		var t *types.Reminder
		if t, err = svc.reminder.FindByID(ID); err != nil {
			return err
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
	_, err := svc.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	return svc.reminder.Delete(ID)
}

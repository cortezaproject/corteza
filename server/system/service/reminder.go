package service

import (
	"context"
	"time"

	intAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

type (
	reminderSender interface {
		Send(kind string, payload interface{}, userIDs ...uint64) error
	}

	reminderServices struct {
		log            *zap.Logger
		reminderSender reminderSender
	}

	reminderAccessController interface {
		CanAssignReminder(ctx context.Context) bool
	}

	ReminderService interface {
		Search(context.Context, types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error)
		FindByID(context.Context, uint64) (*types.Reminder, error)
		FindByIDs(context.Context, []uint64) (types.ReminderSet, error)

		Create(context.Context, *types.Reminder) (*types.Reminder, error)
		Update(context.Context, *types.Reminder) (*types.Reminder, error)

		Dismiss(context.Context, uint64) error
		Undismiss(context.Context, uint64) error
		Snooze(context.Context, uint64, *time.Time) error

		DeleteByID(context.Context, uint64) error

		Watch(ctx context.Context)
	}
)

func Reminder(ctx context.Context, log *zap.Logger, rs reminderSender) ReminderService {
	return &reminder{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		services: &reminderServices{
			log:            log,
			reminderSender: rs,
		},
	}
}

func (svc *reminder) onLookup(ctx context.Context, ID uint64, aProps *reminderActionProps) (*types.Reminder, error) {
	if ID == 0 {
		return nil, ReminderErrInvalidID()
	}

	r, err := store.LookupReminderByID(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	if svc.checkAssignTo(ctx, r) {
		return nil, ReminderErrNotAllowedToRead()
	}

	aProps.setReminder(r)
	return r, nil
}

func (svc *reminder) onSearch(ctx context.Context, filter types.ReminderFilter, aProps *reminderActionProps) (types.ReminderSet, types.ReminderFilter, error) {
	filter.Check = func(r *types.Reminder) (bool, error) {
		return !svc.checkAssignTo(ctx, r), nil
	}

	return store.SearchReminders(ctx, svc.store, filter)
}

func (svc *reminder) onCreate(ctx context.Context, new *types.Reminder) error {
	if err := svc.checkAssignee(ctx, new); err != nil {
		return err
	}

	new.ID = nextID()
	new.CreatedAt = *now()

	return store.CreateReminder(ctx, svc.store, new)
}

func (svc *reminder) onUpdate(ctx context.Context, s store.Storer, upd *types.Reminder, res *types.Reminder, aProps *reminderActionProps, before func() error, after func() error) error {
	if err := before(); err != nil {
		return err
	}

	if err := svc.checkAssignee(ctx, upd); err != nil {
		return err
	}

	if upd.AssignedTo != res.AssignedTo {
		res.AssignedBy = svc.currentUser(ctx)
		res.AssignedAt = time.Now()
	}

	a := res.RemindAt
	if a == nil {
		a = &time.Time{}
	}
	b := upd.RemindAt
	if b == nil {
		b = &time.Time{}
	}
	if !a.Equal(*b) {
		res.SnoozeCount = 0
	}

	res.RemindAt = upd.RemindAt

	return after()
}

func (svc *reminder) onDelete(ctx context.Context, s store.Storer, res *types.Reminder, aProps *reminderActionProps) error {
	res.DeletedAt = now()
	return store.UpdateReminder(ctx, s, res)
}

func (svc *reminder) onFindByIDs(ctx context.Context, aProps *reminderActionProps, IDs []uint64) (types.ReminderSet, error) {
	if len(IDs) == 0 {
		return nil, nil
	}

	rr, _, err := svc.onSearch(ctx, types.ReminderFilter{
		ReminderID: IDs,
		AssignedTo: svc.currentUser(ctx),
	}, aProps)
	return rr, err
}

func (svc *reminder) onDismiss(ctx context.Context, aProps *reminderActionProps, ID uint64) error {
	if ID == 0 {
		return ReminderErrInvalidID()
	}

	r, err := store.LookupReminderByID(ctx, svc.store, ID)
	if err != nil {
		return ReminderErrNotFound()
	}

	if svc.checkAssignTo(ctx, r) {
		return ReminderErrNotAllowedToDismiss()
	}

	aProps.setReminder(r)

	n := time.Now()
	r.DismissedAt = &n
	r.DismissedBy = svc.currentUser(ctx)

	return store.UpdateReminder(ctx, svc.store, r)
}

func (svc *reminder) onUndismiss(ctx context.Context, aProps *reminderActionProps, ID uint64) error {
	if ID == 0 {
		return ReminderErrInvalidID()
	}

	r, err := store.LookupReminderByID(ctx, svc.store, ID)
	if err != nil {
		return ReminderErrNotFound()
	}

	if svc.checkAssignTo(ctx, r) {
		return ReminderErrNotAllowedToUndismiss()
	}

	aProps.setReminder(r)

	r.DismissedAt = nil
	r.DismissedBy = 0

	return store.UpdateReminder(ctx, svc.store, r)
}

func (svc *reminder) onSnooze(ctx context.Context, aProps *reminderActionProps, ID uint64, remindAt *time.Time) error {
	if ID == 0 {
		return ReminderErrInvalidID()
	}

	r, err := store.LookupReminderByID(ctx, svc.store, ID)
	if err != nil {
		return ReminderErrNotFound()
	}

	aProps.setReminder(r)

	r.SnoozeCount++
	r.RemindAt = remindAt

	return store.UpdateReminder(ctx, svc.store, r)
}

func (svc *reminder) checkAssignee(ctx context.Context, rm *types.Reminder) error {
	if svc.checkAssignTo(ctx, rm) {
		if !svc.ac.CanAssignReminder(ctx) {
			return ReminderErrNotAllowedToAssign()
		}
	}
	return nil
}

func (svc *reminder) checkAssignTo(ctx context.Context, rm *types.Reminder) bool {
	return rm.AssignedTo != svc.currentUser(ctx)
}

func (svc *reminder) currentUser(ctx context.Context) uint64 {
	return intAuth.GetIdentityFromContext(ctx).Identity()
}

func (svc *reminder) Watch(ctx context.Context) {
	if svc.services.reminderSender == nil {
		return
	}

	var (
		interval = time.Minute
		rTicker  = time.NewTicker(interval)
	)

	go func() {
		defer sentry.Recover()
		defer rTicker.Stop()
		defer svc.services.log.Info("stopped")

		for {
			select {
			case <-ctx.Done():
				return
			case <-rTicker.C:
				rr, _, err := svc.onSearch(ctx, types.ReminderFilter{
					ExcludeDismissed: true,
					ScheduledOnly:    true,
				}, &reminderActionProps{})

				if err != nil {
					svc.services.log.Error("failed to get reminders of users", zap.Error(err))
				}

				_ = rr.Walk(func(r *types.Reminder) error {
					if r.RemindAt != nil && r.DismissedAt == nil && now().Add(interval).After(*r.RemindAt) {
						if err := svc.services.reminderSender.Send("reminder", r, r.AssignedTo); err != nil {
							svc.services.log.Error("failed to send reminder to user", zap.Error(err))
						}
					}
					return nil
				})
			}
		}
	}()
}

package service

import (
	"context"
	"time"

	"github.com/crusttech/human/server/pkg/actionlog"
	intAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/event"
	"github.com/crusttech/human/server/system/types"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

// The CRUD skeleton (FindByID, Search, Create, Update, DeleteByID) is generated
// in reminder.gen.go from system/reminder.cue. Because reminder access is
// entirely non-RBAC, every op is delegated to an on<Op> handler in this file
// (extracted verbatim from the previous hand-written CRUD methods).
//
// This file owns the struct, access-controller interface, constructor, the
// public ReminderService contract, the on<Op> bodies the generated methods call
// into, and the resource-specific methods: Dismiss / Undismiss / Snooze / Watch
// / FindByIDs and the assignment helpers.

type (
	reminderSender interface {
		Send(kind string, payload interface{}, userIDs ...uint64) error
	}

	reminder struct {
		ac reminderAccessController

		log            *zap.Logger
		actionlog      actionlog.Recorder
		store          store.Storer
		reminderSender reminderSender
		eventbus       eventDispatcher
	}

	reminderAccessController interface {
		CanAssignReminder(ctx context.Context) bool
	}

	ReminderService interface {
		Search(context.Context, types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error)
		FindByID(context.Context, uint64) (*types.Reminder, error)

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
		ac:             DefaultAccessControl,
		log:            log,
		store:          DefaultStore,
		reminderSender: rs,
		eventbus:       eventbus.Service(),
	}
}

// onSearch is the custom body for the generated Search. The generated method
// owns the action-log scaffold + recordAction; the assign-scoped filter.Check
// and the store search live here. There is no RBAC search check (customAccessOps).
func (svc *reminder) onSearch(ctx context.Context, filter types.ReminderFilter, raProps *reminderActionProps) (rr types.ReminderSet, f types.ReminderFilter, err error) {
	filter.Check = func(r *types.Reminder) (bool, error) {
		return !svc.checkAssignTo(ctx, r), nil
	}

	rr, f, err = store.SearchReminders(ctx, svc.store, filter)
	if err != nil {
		return rr, f, err
	}

	return rr, f, nil
}

// onLookup is the custom body for the generated FindByID. Access is gated by
// the assignment check rather than RBAC.
func (svc *reminder) onLookup(ctx context.Context, ID uint64, raProps *reminderActionProps) (r *types.Reminder, err error) {
	if ID == 0 {
		return nil, ReminderErrInvalidID()
	}

	r, err = store.LookupReminderByID(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	if svc.checkAssignTo(ctx, r) {
		return nil, ReminderErrNotAllowedToRead()
	}

	raProps.setReminder(r)

	return r, nil
}

func (svc *reminder) onFindByIDs(ctx context.Context, _ *reminderActionProps, IDs []uint64) (rr types.ReminderSet, err error) {
	if len(IDs) == 0 {
		return nil, nil
	}

	rr, _, err = svc.Search(ctx, types.ReminderFilter{ReminderID: IDs, AssignedTo: svc.currentUser(ctx)})

	return rr, nil
}

func (svc *reminder) checkAssignee(ctx context.Context, rm *types.Reminder) (err error) {
	// Check if user is assigning to someone else
	if svc.checkAssignTo(ctx, rm) {
		if !svc.ac.CanAssignReminder(ctx) {
			return ReminderErrNotAllowedToAssign()
		}
	}

	return nil
}

// checkAssignTo compares current user with reminder.AssignedTo and return bool
func (svc *reminder) checkAssignTo(ctx context.Context, rm *types.Reminder) (valid bool) {
	return rm.AssignedTo != svc.currentUser(ctx)
}

func (svc *reminder) currentUser(ctx context.Context) uint64 {
	return intAuth.GetIdentityFromContext(ctx).Identity()
}

// onCreate is the custom body for the generated Create. The generated method
// owns the action-log scaffold + recordAction; the assignment check, id /
// timestamp assignment and Before/After events live here.
func (svc *reminder) onCreate(ctx context.Context, new *types.Reminder) (err error) {
	if err := svc.checkAssignee(ctx, new); err != nil {
		return err
	}
	r := new
	r.ID = nextID()
	r.CreatedAt = *now()
	if err = svc.eventbus.WaitFor(ctx, event.ReminderBeforeCreate(new, r)); err != nil {
		return
	}
	if err = store.CreateReminder(ctx, svc.store, r); err != nil {
		return err
	}
	svc.eventbus.Dispatch(ctx, event.ReminderAfterCreate(new, r))
	return nil
}

// onUpdate is the custom body for the generated Update. The generated method
// owns the action-log scaffold + recordAction; the assignment check, the
// conditional field copy (assignment / snooze reset) and Before/After events
// live here.
func (svc *reminder) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Reminder, aProps *reminderActionProps, _ func() error, _ func() error) error {
	if upd.ID == 0 {
		return ReminderErrInvalidID()
	}

	if err := svc.checkAssignee(ctx, upd); err != nil {
		return err
	}

	// Assign changed values
	if upd.AssignedTo != res.AssignedTo {
		res.AssignedTo = upd.AssignedTo
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
	res.Payload = upd.Payload
	res.Resource = upd.Resource
	res.UpdatedAt = now()
	if err := svc.eventbus.WaitFor(ctx, event.ReminderBeforeUpdate(upd, res)); err != nil {
		return err
	}
	if err := store.UpdateReminder(ctx, s, res); err != nil {
		return err
	}
	svc.eventbus.Dispatch(ctx, event.ReminderAfterUpdate(upd, res))
	return nil
}

func (svc *reminder) onDismiss(ctx context.Context, aProps *reminderActionProps, ID uint64) (err error) {
	var r *types.Reminder

	aProps.setReminder(&types.Reminder{ID: ID})

	if ID == 0 {
		return ReminderErrInvalidID()
	}

	if r, err = store.LookupReminderByID(ctx, svc.store, ID); err != nil {
		return ReminderErrNotFound()
	}

	if svc.checkAssignTo(ctx, r) {
		return ReminderErrNotAllowedToDismiss()
	}

	aProps.setReminder(r)

	n := time.Now()
	r.DismissedAt = &n
	r.DismissedBy = svc.currentUser(ctx)
	if err = svc.eventbus.WaitFor(ctx, event.ReminderBeforeDismiss(nil, r)); err != nil {
		return
	}
	if err = store.UpdateReminder(ctx, svc.store, r); err != nil {
		return err
	}
	svc.eventbus.Dispatch(ctx, event.ReminderAfterDismiss(nil, r))
	return nil
}

func (svc *reminder) onUndismiss(ctx context.Context, aProps *reminderActionProps, ID uint64) (err error) {
	var r *types.Reminder

	aProps.setReminder(&types.Reminder{ID: ID})

	if ID == 0 {
		return ReminderErrInvalidID()
	}

	if r, err = store.LookupReminderByID(ctx, svc.store, ID); err != nil {
		return ReminderErrNotFound()
	}

	if svc.checkAssignTo(ctx, r) {
		return ReminderErrNotAllowedToUndismiss()
	}

	aProps.setReminder(r)

	r.DismissedAt = nil
	r.DismissedBy = 0
	//pending eventbus integration
	if err = store.UpdateReminder(ctx, svc.store, r); err != nil {
		return err
	}

	return nil
}

func (svc *reminder) onSnooze(ctx context.Context, aProps *reminderActionProps, ID uint64, remindAt *time.Time) (err error) {
	var r *types.Reminder

	aProps.setReminder(&types.Reminder{ID: ID, RemindAt: remindAt})

	if ID == 0 {
		return ReminderErrInvalidID()
	}

	if r, err = store.LookupReminderByID(ctx, svc.store, ID); err != nil {
		return ReminderErrNotFound()
	}

	aProps.setReminder(r)

	r.SnoozeCount++
	r.RemindAt = remindAt
	if err = svc.eventbus.WaitFor(ctx, event.ReminderBeforeSnooze(nil, r)); err != nil {
		return
	}
	if err = store.UpdateReminder(ctx, svc.store, r); err != nil {
		return err
	}
	svc.eventbus.Dispatch(ctx, event.ReminderAfterSnooze(nil, r))

	return nil
}

// onDelete is the custom body for the generated DeleteByID. It performs a
// soft-delete (UpdatedAt-based) via FindByID (which applies the assignment
// access check) plus Before/After delete events.
func (svc *reminder) onDelete(ctx context.Context, s store.Storer, res *types.Reminder, aProps *reminderActionProps) error {
	res.DeletedAt = now()

	aProps.setReminder(res)
	if err := svc.eventbus.WaitFor(ctx, event.ReminderBeforeDelete(nil, res)); err != nil {
		return err
	}
	if err := store.UpdateReminder(ctx, s, res); err != nil {
		return err
	}
	svc.eventbus.Dispatch(ctx, event.ReminderAfterDelete(nil, res))
	return nil
}

func (svc *reminder) Watch(ctx context.Context) {
	if svc.reminderSender != nil {
		var (
			interval = time.Minute
			rTicker  = time.NewTicker(interval)
		)

		go func() {
			defer sentry.Recover()
			defer rTicker.Stop()
			defer svc.log.Info("stopped")

			for {
				select {
				case <-ctx.Done():
					return
				case <-rTicker.C:
					// Get scheduled reminders of users
					rr, _, err := svc.Search(ctx, types.ReminderFilter{
						ExcludeDismissed: true,
						ScheduledOnly:    true,
					})

					if err != nil {
						svc.log.Error("failed to get reminders of users", zap.Error(err))
					}

					// Send scheduled reminders to users
					_ = rr.Walk(func(r *types.Reminder) error {
						if r.RemindAt != nil && r.DismissedAt == nil && now().Add(interval).After(*r.RemindAt) {
							if err := svc.reminderSender.Send("reminder", r, r.AssignedTo); err != nil {
								svc.log.Error("failed to send reminder to user", zap.Error(err))
							}
						}
						return nil
					})
				}
			}
		}()
	}
}

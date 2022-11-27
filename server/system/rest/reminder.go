package rest

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"time"
)

type (
	Reminder struct {
		reminder service.ReminderService
	}

	reminderSetPayload struct {
		Filter types.ReminderFilter `json:"filter"`
		Set    types.ReminderSet    `json:"set,omitempty"`
	}
)

func (Reminder) New() *Reminder {
	ctrl := &Reminder{}
	ctrl.reminder = service.DefaultReminder
	return ctrl
}

func (ctrl *Reminder) List(ctx context.Context, r *request.ReminderList) (interface{}, error) {
	var (
		err error
		f   = types.ReminderFilter{
			AssignedTo:       r.AssignedTo,
			Resource:         r.Resource,
			ScheduledFrom:    r.ScheduledFrom,
			ScheduledUntil:   r.ScheduledUntil,
			ExcludeDismissed: r.ExcludeDismissed,
			IncludeDeleted:   r.IncludeDeleted,
			ScheduledOnly:    r.ScheduledOnly,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err := ctrl.reminder.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl *Reminder) Create(ctx context.Context, r *request.ReminderCreate) (interface{}, error) {
	ntf := &types.Reminder{
		AssignedAt: time.Now(),
		AssignedBy: auth.GetIdentityFromContext(ctx).Identity(),

		AssignedTo: r.AssignedTo,
		Payload:    r.Payload,
		Resource:   r.Resource,
		RemindAt:   r.RemindAt,
	}

	return ctrl.reminder.Create(ctx, ntf)
}

func (ctrl *Reminder) Update(ctx context.Context, r *request.ReminderUpdate) (interface{}, error) {
	ntf := &types.Reminder{
		ID:         r.ReminderID,
		AssignedAt: time.Now(),
		AssignedBy: auth.GetIdentityFromContext(ctx).Identity(),

		AssignedTo: r.AssignedTo,
		Payload:    r.Payload,
		Resource:   r.Resource,
		RemindAt:   r.RemindAt,
	}

	return ctrl.reminder.Update(ctx, ntf)
}

func (ctrl *Reminder) Read(ctx context.Context, r *request.ReminderRead) (interface{}, error) {
	return ctrl.reminder.FindByID(ctx, r.ReminderID)
}

func (ctrl *Reminder) Delete(ctx context.Context, r *request.ReminderDelete) (interface{}, error) {
	return api.OK(), ctrl.reminder.Delete(ctx, r.ReminderID)
}

func (ctrl *Reminder) Dismiss(ctx context.Context, r *request.ReminderDismiss) (interface{}, error) {
	return api.OK(), ctrl.reminder.Dismiss(ctx, r.ReminderID)
}

func (ctrl *Reminder) Snooze(ctx context.Context, r *request.ReminderSnooze) (interface{}, error) {
	return api.OK(), ctrl.reminder.Snooze(ctx, r.ReminderID, r.RemindAt)
}

func (ctrl *Reminder) makeFilterPayload(ctx context.Context, nn types.ReminderSet, f types.ReminderFilter, err error) (*reminderSetPayload, error) {
	if err != nil {
		return nil, err
	}

	return &reminderSetPayload{Filter: f, Set: nn}, nil
}

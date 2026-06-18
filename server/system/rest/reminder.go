package rest

import (
	"context"
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
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

// makeFilter builds the search filter for the generated List controller.
func (ctrl *Reminder) makeFilter(ctx context.Context, r *request.ReminderList) (types.ReminderFilter, error) {
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
		return f, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl *Reminder) beforeCreate(ctx context.Context, res *types.Reminder, r *request.ReminderCreate) error {
	res.Payload = r.Payload
	res.AssignedAt = time.Now()
	res.AssignedBy = auth.GetIdentityFromContext(ctx).Identity()
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID).
func (ctrl *Reminder) beforeUpdate(ctx context.Context, res *types.Reminder, r *request.ReminderUpdate) error {
	res.Payload = r.Payload
	res.AssignedAt = time.Now()
	res.AssignedBy = auth.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (ctrl *Reminder) Dismiss(ctx context.Context, r *request.ReminderDismiss) (interface{}, error) {
	return api.OK(), ctrl.reminder.Dismiss(ctx, r.ReminderID)
}

func (ctrl *Reminder) Undismiss(ctx context.Context, r *request.ReminderUndismiss) (interface{}, error) {
	return api.OK(), ctrl.reminder.Undismiss(ctx, r.ReminderID)
}

func (ctrl *Reminder) Snooze(ctx context.Context, r *request.ReminderSnooze) (interface{}, error) {
	return api.OK(), ctrl.reminder.Snooze(ctx, r.ReminderID, r.RemindAt)
}

func (ctrl *Reminder) makePayload(ctx context.Context, m *types.Reminder, err error) (*types.Reminder, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return m, nil
}

func (ctrl *Reminder) makeFilterPayload(ctx context.Context, nn types.ReminderSet, f types.ReminderFilter, err error) (*reminderSetPayload, error) {
	if err != nil {
		return nil, err
	}

	return &reminderSetPayload{Filter: f, Set: nn}, nil
}

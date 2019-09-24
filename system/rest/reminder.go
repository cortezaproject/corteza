package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/system/rest/request"
)

var _ = errors.Wrap

type Reminder struct {
	// xxx service.XXXService
}

func (Reminder) New() *Reminder {
	return &Reminder{}
}

func (ctrl *Reminder) List(ctx context.Context, r *request.ReminderList) (interface{}, error) {
	return nil, errors.New("Not implemented: Reminder.list")
}

func (ctrl *Reminder) Create(ctx context.Context, r *request.ReminderCreate) (interface{}, error) {
	return nil, errors.New("Not implemented: Reminder.create")
}

func (ctrl *Reminder) Update(ctx context.Context, r *request.ReminderUpdate) (interface{}, error) {
	return nil, errors.New("Not implemented: Reminder.update")
}

func (ctrl *Reminder) Read(ctx context.Context, r *request.ReminderRead) (interface{}, error) {
	return nil, errors.New("Not implemented: Reminder.read")
}

func (ctrl *Reminder) Delete(ctx context.Context, r *request.ReminderDelete) (interface{}, error) {
	return nil, errors.New("Not implemented: Reminder.delete")
}

func (ctrl *Reminder) Dismiss(ctx context.Context, r *request.ReminderDismiss) (interface{}, error) {
	return nil, errors.New("Not implemented: Reminder.dismiss")
}

func (ctrl *Reminder) Snooze(ctx context.Context, r *request.ReminderSnooze) (interface{}, error) {
	return nil, errors.New("Not implemented: Reminder.snooze")
}

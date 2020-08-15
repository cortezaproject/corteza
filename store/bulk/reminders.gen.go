package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/reminders.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	reminderCreate struct {
		Done chan struct{}
		res  *types.Reminder
		err  error
	}

	reminderUpdate struct {
		Done chan struct{}
		res  *types.Reminder
		err  error
	}

	reminderRemove struct {
		Done chan struct{}
		res  *types.Reminder
		err  error
	}
)

// CreateReminder creates a new Reminder
// create job that can be pushed to store's transaction handler
func CreateReminder(res *types.Reminder) *reminderCreate {
	return &reminderCreate{res: res}
}

// Do Executes reminderCreate job
func (j *reminderCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateReminder(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateReminder creates a new Reminder
// update job that can be pushed to store's transaction handler
func UpdateReminder(res *types.Reminder) *reminderUpdate {
	return &reminderUpdate{res: res}
}

// Do Executes reminderUpdate job
func (j *reminderUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateReminder(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveReminder creates a new Reminder
// remove job that can be pushed to store's transaction handler
func RemoveReminder(res *types.Reminder) *reminderRemove {
	return &reminderRemove{res: res}
}

// Do Executes reminderRemove job
func (j *reminderRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveReminder(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

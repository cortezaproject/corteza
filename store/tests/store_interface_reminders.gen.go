package tests

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/reminders.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	remindersStore interface {
		SearchReminders(ctx context.Context, f types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error)
		LookupReminderByID(ctx context.Context, id uint64) (*types.Reminder, error)
		CreateReminder(ctx context.Context, rr ...*types.Reminder) error
		UpdateReminder(ctx context.Context, rr ...*types.Reminder) error
		PartialUpdateReminder(ctx context.Context, onlyColumns []string, rr ...*types.Reminder) error
		RemoveReminder(ctx context.Context, rr ...*types.Reminder) error
		RemoveReminderByID(ctx context.Context, ID uint64) error

		TruncateReminders(ctx context.Context) error
	}
)

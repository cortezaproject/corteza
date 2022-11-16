package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/reminders.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	Reminders interface {
		SearchReminders(ctx context.Context, f types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error)
		LookupReminderByID(ctx context.Context, id uint64) (*types.Reminder, error)

		CreateReminder(ctx context.Context, rr ...*types.Reminder) error

		UpdateReminder(ctx context.Context, rr ...*types.Reminder) error

		UpsertReminder(ctx context.Context, rr ...*types.Reminder) error

		DeleteReminder(ctx context.Context, rr ...*types.Reminder) error
		DeleteReminderByID(ctx context.Context, ID uint64) error

		TruncateReminders(ctx context.Context) error
	}
)

var _ *types.Reminder
var _ context.Context

// SearchReminders returns all matching Reminders from store
func SearchReminders(ctx context.Context, s Reminders, f types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error) {
	return s.SearchReminders(ctx, f)
}

// LookupReminderByID searches for reminder by its ID
//
// It returns reminder even if deleted or suspended
func LookupReminderByID(ctx context.Context, s Reminders, id uint64) (*types.Reminder, error) {
	return s.LookupReminderByID(ctx, id)
}

// CreateReminder creates one or more Reminders in store
func CreateReminder(ctx context.Context, s Reminders, rr ...*types.Reminder) error {
	return s.CreateReminder(ctx, rr...)
}

// UpdateReminder updates one or more (existing) Reminders in store
func UpdateReminder(ctx context.Context, s Reminders, rr ...*types.Reminder) error {
	return s.UpdateReminder(ctx, rr...)
}

// UpsertReminder creates new or updates existing one or more Reminders in store
func UpsertReminder(ctx context.Context, s Reminders, rr ...*types.Reminder) error {
	return s.UpsertReminder(ctx, rr...)
}

// DeleteReminder Deletes one or more Reminders from store
func DeleteReminder(ctx context.Context, s Reminders, rr ...*types.Reminder) error {
	return s.DeleteReminder(ctx, rr...)
}

// DeleteReminderByID Deletes Reminder from store
func DeleteReminderByID(ctx context.Context, s Reminders, ID uint64) error {
	return s.DeleteReminderByID(ctx, ID)
}

// TruncateReminders Deletes all Reminders from store
func TruncateReminders(ctx context.Context, s Reminders) error {
	return s.TruncateReminders(ctx)
}

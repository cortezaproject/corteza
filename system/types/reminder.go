package types

import (
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type (
	Reminder struct {
		ID          uint64         `json:"reminderID,string"`
		Resource    string         `json:"resource"`
		Payload     types.JSONText `json:"payload"`
		SnoozeCount uint           `json:"snoozeCount"`

		AssignedTo uint64    `json:"assignedTo,string"`
		AssignedBy uint64    `json:"assignedBy,string"`
		AssignedAt time.Time `json:"assignedAt"`

		DismissedBy uint64     `json:"dismissedBy,string"`
		DismissedAt *time.Time `json:"dismissedAt"`

		RemindAt *time.Time `json:"remindAt"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	ReminderFilter struct {
		ReminderID       []uint64   `json:"reminderID"`
		Resource         string     `json:"resource"`
		AssignedTo       uint64     `json:"assignedTo,uint64"`
		ScheduledFrom    *time.Time `json:"scheduledFrom"`
		ScheduledUntil   *time.Time `json:"scheduledUntil"`
		ExcludeDismissed bool       `json:"excludeDismissed"`
		IncludeDeleted   bool       `json:"includeDeleted"`
		ScheduledOnly    bool       `json:"scheduledOnly"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Reminder) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

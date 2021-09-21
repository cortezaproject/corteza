package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/system/types"
	"time"
)

func (s Store) convertReminderFilter(f types.ReminderFilter) (query squirrel.SelectBuilder, err error) {
	query = s.remindersSelectBuilder()

	if len(f.ReminderID) > 0 {
		query = query.Where(squirrel.Eq{"rmd.ID": f.ReminderID})
	}

	if f.ExcludeDismissed {
		query = query.Where("rmd.dismissed_at IS NULL")
	}

	if !f.IncludeDeleted {
		query = query.Where("rmd.deleted_at IS NULL")
	}

	if f.ScheduledOnly {
		query = query.Where("rmd.remind_at IS NOT NULL")
	}

	if f.AssignedTo != 0 {
		query = query.Where("rmd.assigned_to = ?", f.AssignedTo)
	}

	if f.Resource != "" {
		query = query.Where("rmd.resource LIKE ?", f.Resource+"%")
	}

	if f.ScheduledFrom != nil {
		query = query.Where("rmd.remind_at >= ?", f.ScheduledFrom.Format(time.RFC3339))
	}
	if f.ScheduledUntil != nil {
		query = query.Where("rmd.remind_at <= ?", f.ScheduledUntil.Format(time.RFC3339))
	}

	return
}

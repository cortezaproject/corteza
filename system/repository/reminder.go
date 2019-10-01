package repository

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/rh"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/titpetric/factory"
	squirrel "gopkg.in/Masterminds/squirrel.v1"
)

type (
	ReminderRepository interface {
		Find(types.ReminderFilter) (set types.ReminderSet, f types.ReminderFilter, err error)
		FindByID(ID uint64) (*types.Reminder, error)
		FindByIDs(ID []uint64) (types.ReminderSet, error)

		Create(mod *types.Reminder) (*types.Reminder, error)
		Update(mod *types.Reminder) (*types.Reminder, error)

		Delete(ID uint64) error
	}

	reminder struct {
		*repository
	}
)

const (
	ErrReminderNotFound = repositoryError("ReminderNotFound")
)

func Reminder(ctx context.Context, db *factory.DB) ReminderRepository {
	rpo := &repository{}
	return (&reminder{
		repository: rpo.With(ctx, db),
	})
}

func (r reminder) table() string {
	return "sys_reminder"
}

func (r reminder) columns() []string {
	return []string{
		"id",
		"resource",
		"payload",
		"snooze_count",

		"assigned_to",
		"assigned_by",
		"assigned_at",

		"dismissed_by",
		"dismissed_at",

		"remind_at",

		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (r reminder) query() squirrel.SelectBuilder {
	return r.queryNoFilter().Where("r.deleted_at IS NULL")
}

func (r reminder) queryNoFilter() squirrel.SelectBuilder {
	return squirrel.
		Select().
		From(r.table() + " AS r").
		Columns(r.columns()...)
}

func (r reminder) Find(filter types.ReminderFilter) (set types.ReminderSet, f types.ReminderFilter, err error) {
	f = filter
	q := r.query()

	if f.ExcludeDismissed {
		q = q.Where("dismissed_at IS NULL")
	}

	if f.ScheduledOnly {
		q = q.Where("remind_at IS NOT NULL")
	}

	if f.AssignedTo != 0 {
		q = q.Where("r.assigned_to = ?", f.AssignedTo)
	}

	if f.Resource != "" {
		q = q.Where("r.resource LIKE ?", f.Resource+"%")
	}

	if f.ScheduledFrom != nil {
		q = q.Where("r.remind_at >= ?", f.ScheduledFrom.Format(time.RFC3339))
	}
	if f.ScheduledUntil != nil {
		q = q.Where("r.remind_at <= ?", f.ScheduledUntil.Format(time.RFC3339))
	}

	if f.AccessCheck.HasOperation() {
		q = q.Where(f.AccessCheck.BindToEnv(
			types.ReminderPermissionResource,
			"sys",
		))
	}

	if f.Count, err = r.count(q); err != nil || f.Count == 0 {
		return
	}

	// @todo allow sorting at some point
	q = q.OrderBy("r.remind_at")

	return set, f, rh.FetchPaged(r.db(), q, f.Page, f.PerPage, &set)
}

func (r reminder) FindByID(ID uint64) (rm *types.Reminder, err error) {
	rm = &types.Reminder{}

	q := r.query().
		Where("r.id = ?", ID)

	err = r.fetchOne(rm, q)
	if err != nil {
		return nil, err
	} else if rm.ID <= 0 {
		return nil, ErrReminderNotFound
	}

	return rm, nil
}

func (r reminder) FindByIDs(IDs []uint64) (rr types.ReminderSet, err error) {
	if len(IDs) == 0 {
		return nil, nil
	}

	var (
		q = r.query().
			Where("r.id IN (?)", IDs)
	)

	return rr, r.fetchSet(&rr, q)
}

func (r reminder) Create(mod *types.Reminder) (rm *types.Reminder, err error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.table(), mod)
}

func (r reminder) Update(mod *types.Reminder) (*types.Reminder, error) {
	mod.UpdatedAt = timeNowPtr()
	return mod, r.db().Replace(r.table(), mod)
}

func (r reminder) Delete(ID uint64) error {
	return r.updateColumnByID(r.table(), "deleted_at", time.Now(), ID)
}

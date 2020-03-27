package repository

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/rh"

	squirrel "github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/titpetric/factory"
)

type (
	ReminderRepository interface {
		Find(types.ReminderFilter) (set types.ReminderSet, f types.ReminderFilter, err error)
		FindByID(ID uint64) (*types.Reminder, error)

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
	return (&reminder{}).With(ctx, db)
}

func (r reminder) With(ctx context.Context, db *factory.DB) ReminderRepository {
	return &reminder{
		repository: r.repository.With(ctx, db),
	}
}

func (r reminder) table() string {
	return "sys_reminder"
}

func (r reminder) columns() []string {
	return []string{
		"r.id",
		"r.resource",
		"r.payload",
		"r.snooze_count",

		"r.assigned_to",
		"r.assigned_by",
		"r.assigned_at",

		"r.dismissed_by",
		"r.dismissed_at",

		"r.remind_at",

		"r.created_at",
		"r.updated_at",
		"r.deleted_at",
	}
}

func (r reminder) query() squirrel.SelectBuilder {
	return r.queryNoFilter().Where("r.deleted_at IS NULL")
}

func (r reminder) queryNoFilter() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table() + " AS r")
}

func (r reminder) FindByID(ID uint64) (rm *types.Reminder, err error) {
	return r.findOneBy("id", ID)
}

func (r reminder) findOneBy(field string, value interface{}) (*types.Reminder, error) {
	var (
		p = &types.Reminder{}

		q = r.query().
			Where(squirrel.Eq{field: value})

		err = rh.FetchOne(r.db(), q, p)
	)

	if err != nil {
		return nil, err
	} else if p.ID == 0 {
		return nil, ErrReminderNotFound
	}

	return p, nil
}

func (r reminder) Find(filter types.ReminderFilter) (set types.ReminderSet, f types.ReminderFilter, err error) {
	f = filter

	if f.Sort == "" {
		f.Sort = "r.remind_at"
	}

	query := r.query()

	if len(f.ReminderID) > 0 {
		query = query.Where(squirrel.Eq{"r.ID": f.ReminderID})
	}

	if f.ExcludeDismissed {
		query = query.Where("r.dismissed_at IS NULL")
	}

	if f.ScheduledOnly {
		query = query.Where("r.remind_at IS NOT NULL")
	}

	if f.AssignedTo != 0 {
		query = query.Where("r.assigned_to = ?", f.AssignedTo)
	}

	if f.Resource != "" {
		query = query.Where("r.resource LIKE ?", f.Resource+"%")
	}

	if f.ScheduledFrom != nil {
		query = query.Where("r.remind_at >= ?", f.ScheduledFrom.Format(time.RFC3339))
	}
	if f.ScheduledUntil != nil {
		query = query.Where("r.remind_at <= ?", f.ScheduledUntil.Format(time.RFC3339))
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, r.columns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	if f.Count, err = rh.Count(r.db(), query); err != nil || f.Count == 0 {
		return
	}

	return set, f, rh.FetchPaged(r.db(), query, f.PageFilter, &set)
}

func (r reminder) Create(mod *types.Reminder) (rm *types.Reminder, err error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.table(), mod)
}

func (r reminder) Update(mod *types.Reminder) (*types.Reminder, error) {
	rh.SetCurrentTimeRounded(&mod.UpdatedAt)
	return mod, r.db().Replace(r.table(), mod)
}

func (r reminder) Delete(ID uint64) error {
	return r.updateColumnByID(r.table(), "deleted_at", time.Now(), ID)
}

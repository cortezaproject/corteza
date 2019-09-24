package repository

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/titpetric/factory"
	squirrel "gopkg.in/Masterminds/squirrel.v1"
)

type (
	ReminderRepository interface {
		Find(types.ReminderFilter) (set types.ReminderSet, err error)
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

func (r reminder) Find(filter types.ReminderFilter) (set types.ReminderSet, err error) {
	q := r.query()

	if filter.AssignedTo != 0 {
		q = q.Where("r.assigned_to = ?", filter.AssignedTo)
	}

	if filter.Resource != "" {
		q = q.Where("r.resource LIKE ?%", filter.Resource)
	}

	// @todo allow sorting at some point
	q = q.OrderBy("r.remind_at")

	return set, r.fetchPaged(&set, q, filter.Page, filter.PerPage)
}

func (r reminder) FindByID(ID uint64) (rm *types.Reminder, err error) {
	q := r.query().
		Where("r.id = ?", ID)

	return rm, isFound(r.fetchOne(rm, q), rm.ID > 0, ErrReminderNotFound)
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

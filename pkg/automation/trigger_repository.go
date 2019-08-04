package automation

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	// repository servs as a db storage layer for permission rules
	triggerRepository struct {
		dbh *factory.DB

		// sql table reference
		dbTablePrefix string
	}
)

func TriggerRepository(db *factory.DB, dbTablePrefix string) *triggerRepository {
	return &triggerRepository{
		dbTablePrefix: dbTablePrefix,
		dbh:           db,
	}
}

func (r *triggerRepository) With(ctx context.Context) *triggerRepository {
	return &triggerRepository{
		dbTablePrefix: r.dbTablePrefix,
		dbh:           r.db().With(ctx),
	}
}

func (r *triggerRepository) db() *factory.DB {
	return r.dbh
}

func (r triggerRepository) table() string {
	return r.dbTablePrefix + "_automation_trigger"
}

func (r triggerRepository) columns() []string {
	return []string{
		"id",
		"event",
		"resource",
		"`condition`",
		"rel_script",
		"enabled",
		"created_at",
		"created_by",
		"updated_at",
		"updated_by",
		"deleted_at",
		"deleted_by",
	}
}

func (r *triggerRepository) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table())
}

// FindByID finds specific trigger
func (r *triggerRepository) FindByID(ctx context.Context, triggerID uint64) (*Trigger, error) {
	var (
		rval = &Trigger{}

		query = r.query().
			Columns(r.columns()...).
			Where("id = ?", triggerID)
	)

	return rval, rh.IsFound(rh.FetchOne(r.db(), query, rval), rval.ID > 0, errors.New("trigger not found"))
}

// Find - finds triggers using given filter
func (r *triggerRepository) Find(ctx context.Context, filter TriggerFilter) (set TriggerSet, f TriggerFilter, err error) {
	f = filter

	query := r.query()

	if f.ScriptID > 0 {
		query = query.Where("rel_script = ?", f.ScriptID)
	}

	if f.Event != "" {
		query = query.Where("resource = ?", f.Event)
	}

	if f.Resource != "" {
		query = query.Where("resource = ?", f.Resource)
	}

	if !filter.IncDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	if f.Count, err = rh.Count(r.db(), query); err != nil || f.Count == 0 {
		return
	}

	query = query.OrderBy("id ASC")

	return set, f, rh.FetchPaged(r.db(), query, f.Page, f.PerPage, &set)
}

// FindAllRunnable - loads and returns all runnable triggers
func (r *triggerRepository) FindAllRunnable() (TriggerSet, error) {
	rr := make([]*Trigger, 0)

	return rr, errors.Wrap(rh.FetchAll(
		r.db(),
		r.query().Where("enabled AND deleted_at IS NULL"),
		&rr,
	), "could not load runnable triggers")
}

func (r *triggerRepository) Create(s *Trigger) (err error) {
	return r.dbh.Transaction(func() error {
		// Generate ID
		s.ID = factory.Sonyflake.NextID()

		if s.CreatedAt.IsZero() {
			// Make sure time of creation is set
			s.CreatedAt = time.Now()
		}

		// Ensure sanity
		s.UpdatedAt, s.UpdatedBy = nil, 0
		s.DeletedAt, s.DeletedBy = nil, 0

		return r.dbh.Insert(r.table(), s)
	})
}

func (r *triggerRepository) Update(s *Trigger) (err error) {
	return r.dbh.Transaction(func() error {
		s.UpdatedAt = &time.Time{}
		*s.UpdatedAt = time.Now()

		return r.dbh.Update(r.table(), s, "id")
	})
}

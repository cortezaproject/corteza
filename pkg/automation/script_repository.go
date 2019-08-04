package automation

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	// repository servs as a db storage layer for permission rules
	scriptRepository struct {
		dbh *factory.DB

		// sql table reference
		dbTablePrefix string
	}
)

func ScriptRepository(db *factory.DB, dbTablePrefix string) *scriptRepository {
	return &scriptRepository{
		dbTablePrefix: dbTablePrefix,
		dbh:           db,
	}
}

func (r *scriptRepository) With(ctx context.Context) *scriptRepository {
	return &scriptRepository{
		dbTablePrefix: r.dbTablePrefix,
		dbh:           r.db().With(ctx),
	}
}

func (r *scriptRepository) db() *factory.DB {
	return r.dbh
}

func (r scriptRepository) table() string {
	return r.dbTablePrefix + "_automation_script"
}

func (r scriptRepository) columns() []string {
	return []string{
		"id",
		"name",
		"source_ref",
		"source",
		"async",
		"rel_runner",
		"run_in_ua",
		"timeout",
		"critical",
		"enabled",
		"created_at",
		"created_by",
		"updated_at",
		"updated_by",
		"deleted_at",
		"deleted_by",
	}
}

func (r *scriptRepository) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table())
}

// FindByID finds specific script
func (r *scriptRepository) FindByID(ctx context.Context, scriptID uint64) (*Script, error) {
	var (
		rval = &Script{}

		query = r.query().
			Columns(r.columns()...).
			Where("id = ?", scriptID)
	)

	return rval, rh.IsFound(rh.FetchOne(r.db(), query, rval), rval.ID > 0, errors.New("script not found"))
}

// Find - finds scripts using given filter
func (r *scriptRepository) Find(ctx context.Context, filter ScriptFilter) (set ScriptSet, f ScriptFilter, err error) {
	f = filter

	query := r.query()

	if !filter.IncDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	if f.Query != "" {
		q := "%" + f.Query + "%"
		query = query.Where("name like ?", q)
	}

	if f.Resource != "" {
		// Making partial trigger repo struct on the fly to help us calculate the name of the triggers table
		ttable := (triggerRepository{dbTablePrefix: r.dbTablePrefix}).table()
		query = query.Where(
			fmt.Sprintf("id IN (SELECT rel_script FROM `%s` WHERE resource = ?", ttable),
			f.Resource,
		)
	}

	if f.Count, err = rh.Count(r.db(), query); err != nil || f.Count == 0 {
		return
	}

	query = query.OrderBy("id ASC")

	return set, f, rh.FetchPaged(r.db(), query, f.Page, f.PerPage, &set)
}

// FindAllRunnable - loads and returns all runnable scripts
func (r *scriptRepository) FindAllRunnable() (ScriptSet, error) {
	rr := make([]*Script, 0)

	return rr, errors.Wrap(rh.FetchAll(
		r.db(),
		r.query().Where("enabled AND deleted_at IS NULL"),
		&rr,
	), "could not load runnable scripts")
}

func (r *scriptRepository) Create(s *Script) (err error) {
	return r.dbh.Transaction(func() error {
		return r.dbh.Insert(r.table(), s)
	})
}

func (r *scriptRepository) Update(s *Script) (err error) {
	return r.dbh.Transaction(func() error {
		return r.dbh.Update(r.table(), s, "id")
	})
}

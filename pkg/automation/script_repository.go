package automation

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	// repository servs as a db storage layer for permission rules
	scriptRepository struct {
		// sql table reference
		dbTablePrefix string
	}
)

func ScriptRepository(dbTablePrefix string) *scriptRepository {
	return &scriptRepository{
		dbTablePrefix: dbTablePrefix,
	}
}

func (r scriptRepository) table() string {
	return r.dbTablePrefix + "_automation_script"
}

func (r scriptRepository) columns() []string {
	return []string{
		"id",
		"rel_namespace",
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
func (r *scriptRepository) findByID(db *factory.DB, scriptID uint64) (rval *Script, err error) {
	var (
		query = r.query().
			Where("id = ?", scriptID)
	)

	rval = &Script{}

	if err = rh.IsFound(rh.FetchOne(db, query, rval), rval.ID > 0, errors.New("script not found")); err != nil {
		return nil, err
	}

	return rval, nil
}

// Find - finds scripts using given filter
func (r *scriptRepository) find(db *factory.DB, filter ScriptFilter) (set ScriptSet, f ScriptFilter, err error) {
	f = filter

	query := r.query()

	if !f.IncDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	if f.NamespaceID > 0 {
		query = query.Where("rel_namespace = ?", f.NamespaceID)
	}

	if f.Name != "" {
		query = query.Where(squirrel.Eq{"name": f.Name})
	}

	if f.Query != "" {
		q := "%" + f.Query + "%"
		query = query.Where("name like ?", q)
	}

	if f.Resource != "" {
		// Making partial trigger repo struct on the fly to help us calculate the name of the triggers table
		query = query.Where(
			fmt.Sprintf("id IN (SELECT rel_script FROM `%s` WHERE resource = ?", r.table()),
			f.Resource,
		)
	}

	// Limit by access
	if f.AccessCheck.HasOperation() {
		query = query.Where(f.AccessCheck.BindToEnv(
			types.AutomationScriptPermissionResource,
			r.dbTablePrefix,
		))
	}

	if f.Count, err = rh.Count(db, query); err != nil || f.Count == 0 {
		return
	}

	query = query.OrderBy("id ASC")

	return set, f, rh.FetchPaged(db, query, f.Page, f.PerPage, &set)
}

// FindAllRunnable - loads and returns all runnable scripts
func (r *scriptRepository) findRunnable(db *factory.DB) (ScriptSet, error) {
	rr := make([]*Script, 0)

	return rr, errors.Wrap(rh.FetchAll(
		db,
		r.query().Where("enabled AND deleted_at IS NULL"),
		&rr,
	), "could not load runnable scripts")
}

func (r *scriptRepository) create(db *factory.DB, s *Script) (err error) {
	s.ID = factory.Sonyflake.NextID()
	return db.Insert(r.table(), s)
}

func (r *scriptRepository) update(db *factory.DB, s *Script) (err error) {
	return db.Update(r.table(), s, "id")
}

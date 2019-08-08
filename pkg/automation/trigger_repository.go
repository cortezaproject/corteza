package automation

import (
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	// repository servs as a db storage layer for permission rules
	triggerRepository struct {
		// sql table reference
		dbTablePrefix string
	}
)

func TriggerRepository(dbTablePrefix string) *triggerRepository {
	return &triggerRepository{
		dbTablePrefix: dbTablePrefix,
	}
}

func (r triggerRepository) table() string {
	return r.dbTablePrefix + "_automation_trigger"
}

func (r triggerRepository) columns() []string {
	return []string{
		"id",
		"resource",
		"event",
		"event_condition",
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
func (r *triggerRepository) findByID(db *factory.DB, triggerID uint64) (*Trigger, error) {
	var (
		rval = &Trigger{}

		query = r.query().Where("id = ?", triggerID)
	)

	return rval, rh.IsFound(rh.FetchOne(db, query, rval), rval.ID > 0, errors.New("trigger not found"))
}

// Find - finds triggers using given filter
func (r *triggerRepository) find(db *factory.DB, filter TriggerFilter) (set TriggerSet, f TriggerFilter, err error) {
	f = filter

	query := r.query()

	if f.ScriptID > 0 {
		query = query.Where("rel_script = ?", f.ScriptID)
	}

	if f.Event != "" {
		query = query.Where("event = ?", f.Event)
	}

	if f.Resource != "" {
		query = query.Where("resource = ?", f.Resource)
	}

	if f.Condition != "" {
		query = query.Where("event_condition = ?", f.Condition)
	}

	if !filter.IncDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	if f.Count, err = rh.Count(db, query); err != nil || f.Count == 0 {
		return
	}

	query = query.OrderBy("id ASC")

	return set, f, rh.FetchPaged(db, query, f.Page, f.PerPage, &set)
}

// FindAllRunnable - loads and returns all runnable triggers
func (r *triggerRepository) findRunnable(db *factory.DB) (TriggerSet, error) {
	rr := make([]*Trigger, 0)

	return rr, errors.Wrap(rh.FetchAll(
		db,
		r.query().Where("enabled AND deleted_at IS NULL"),
		&rr,
	), "could not load runnable triggers")
}

func (r *triggerRepository) replace(db *factory.DB, t *Trigger) (err error) {
	if dup, err := r.checkDuplicate(db, t); err != nil {
		return err
	} else if dup != nil {
		t.ID = dup.ID
	}

	return db.Replace(r.table(), t)
}

func (r *triggerRepository) deleteByScriptID(db *factory.DB, scriptID uint64) (err error) {
	return db.Delete(r.table(), Trigger{ScriptID: scriptID}, "rel_script")
}

// Check for existing events
func (r *triggerRepository) checkDuplicate(db *factory.DB, t *Trigger) (*Trigger, error) {
	if t.IsDeferred() && t.IsInterval() {
		// deferred & interval triggers
		// can have duplicates
		return nil, nil
	}

	tt, _, err := r.find(db, TriggerFilter{
		Resource:   t.Resource,
		Event:      t.Event,
		ScriptID:   t.ScriptID,
		Condition:  t.Condition,
		IncDeleted: false,
	})

	if err != nil || len(tt) == 0 {
		return nil, err
	}

	return tt[0], nil
}

func (r *triggerRepository) mergeSet(db *factory.DB, tms triggersMergeStrategy, scriptID uint64, tt TriggerSet) error {
	if tms == STMS_IGNORE {
		// do nothing
		return nil
	}

	if tms == STMS_REPLACE {
		// Mark all existing as deleted
		//
		// here, we're assuming we have the entire
		// trigger list present (in s.triggers)
		if err := r.deleteByScriptID(db, scriptID); err != nil {
			return err
		}
	}

	return tt.Walk(func(t *Trigger) error {
		if t.ScriptID != scriptID {
			return nil
		}

		// Replace (upsert) all triggers we have
		return r.replace(db, t)
	})
}

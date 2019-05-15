package permissions

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"
)

type (
	// repository servs as a db storage layer for permission rules
	repository struct {
		dbh *factory.DB

		// sql table reference
		dbTable string
	}
)

func Repository(db *factory.DB, table string) *repository {
	return &repository{
		dbTable: table,
		dbh:     db,
	}
}

func (r *repository) db() *factory.DB {
	return r.dbh
}

func (r repository) columns() []string {
	return []string{
		"rel_role",
		"resource",
		"operation",
		"access",
	}
}

func (r *repository) With(ctx context.Context) *repository {
	return &repository{
		dbTable: r.dbTable,
		dbh:     r.db().With(ctx),
	}
}

func (r *repository) Load() (RuleSet, error) {
	rr := make([]*Rule, 0)

	lookup := squirrel.
		Select(r.columns()...).
		From(r.dbTable)

	if query, args, err := lookup.ToSql(); err != nil {
		return nil, errors.Wrap(err, "could not build lookup query for permission rules")
	} else if err = r.dbh.Select(&rr, query, args...); err != nil {
		return nil, errors.Wrap(err, "could not get permission rules")
	}

	return rr, nil
}

func (r *repository) Store(deleteSet, updateSet RuleSet) (err error) {
	if len(deleteSet) == 0 && len(updateSet) == 0 {
		return
	}

	return r.dbh.Transaction(func() error {
		if len(deleteSet) > 0 {
			err = deleteSet.Walk(func(rule *Rule) error {
				return r.dbh.Delete(r.dbTable, rule, "rel_role", "resource", "operation")
			})

			if err != nil {
				return err
			}
		}

		if len(updateSet) > 0 {
			err = updateSet.Walk(func(rule *Rule) error {
				return r.dbh.Replace(r.dbTable, rule)
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

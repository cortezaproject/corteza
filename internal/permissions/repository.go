package permissions

import (
	"context"

	"github.com/titpetric/factory"
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
		"value",
	}
}

func (r *repository) With(ctx context.Context) *repository {
	return &repository{
		dbTable: r.dbTable,
		dbh:     r.db().With(ctx),
	}
}

func (r *repository) Load() (rr RuleSet, err error) {
	// @todo load and return
	return nil, nil
}

package repository

import (
	"context"

	"github.com/lann/builder"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/internal/auth"
)

type (
	repository struct {
		ctx context.Context
		dbh *factory.DB
	}
)

// DB produces a contextual DB handle
func DB(ctx context.Context) *factory.DB {
	return factory.Database.MustGet("compose").With(ctx)
}

// Identity returns the User ID from context
func Identity(ctx context.Context) uint64 {
	return auth.GetIdentityFromContext(ctx).Identity()
}

// With updates repository and database contexts
func (r *repository) With(ctx context.Context, db *factory.DB) *repository {
	return &repository{
		ctx: ctx,
		dbh: db,
	}
}

// Context returns current active repository context
func (r *repository) Context() context.Context {
	return r.ctx
}

// db returns context-aware db handle
func (r *repository) db() *factory.DB {
	if r.dbh != nil {
		return r.dbh
	}
	return DB(r.ctx)
}

// Fetches single row from table
func (r repository) fetchOne(one interface{}, q squirrel.SelectBuilder) (err error) {
	var (
		sql  string
		args []interface{}
	)

	if sql, args, err = q.ToSql(); err != nil {
		return
	}

	if err = r.db().Get(one, sql, args...); err != nil {
		return
	}

	return
}

// Counts all rows that match conditions from given query builder
func (r repository) count(q squirrel.SelectBuilder) (uint, error) {
	var (
		count uint
		cq    = q.Column("COUNT(*)")
	)

	// Remove order-bys for counting
	cq = builder.Delete(cq, "OrderBys").(squirrel.SelectBuilder)

	if sqlSelect, argsSelect, err := cq.ToSql(); err != nil {
		return 0, err
	} else {
		if err := r.db().Get(&count, sqlSelect, argsSelect...); err != nil {
			return 0, err
		}
	}

	return count, nil
}

// Fetches paged rows
func (r repository) fetchPaged(set interface{}, q squirrel.SelectBuilder, page, perPage uint) error {
	if perPage > 0 {
		q = q.Limit(uint64(perPage))
	}

	if page > 0 {
		q = q.Offset(uint64(page * perPage))
	}

	if sqlSelect, argsSelect, err := q.ToSql(); err != nil {
		return err
	} else {
		return r.db().Select(set, sqlSelect, argsSelect...)
	}
}

func normalizePerPage(val, min, max, def uint) uint {
	if val == 0 {
		return def
	}

	if max > 0 && val > max {
		return max
	}

	if min > 0 && val < min {
		return min
	}

	return val
}

func isFound(err error, valid bool, nerr error) error {
	if err != nil {
		return errors.WithStack(err)
	} else if !valid {
		return errors.WithStack(nerr)
	}

	return nil
}

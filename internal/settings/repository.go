package settings

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"
)

type (
	repository struct {
		dbh *factory.DB

		// sql table reference
		dbTable string
	}

	Repository interface {
		With(ctx context.Context) Repository

		Find(filter Filter) (ss ValueSet, err error)

		Get(name string, ownedBy uint64) (value *Value, err error)
		Set(value *Value) error
		BulkSet(vv ValueSet) error
		Delete(name string, ownedBy uint64) (err error)
	}
)

func NewRepository(db *factory.DB, table string) Repository {
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
		"name",
		"value",
		"rel_owner",
		"updated_at",
		"updated_by",
	}
}

func (r *repository) With(ctx context.Context) Repository {
	return &repository{
		dbTable: r.dbTable,
		dbh:     r.db().With(ctx),
	}
}

func (r *repository) Find(f Filter) (ss ValueSet, err error) {
	f.Normalize()
	lookup := squirrel.
		Select(r.columns()...).
		From(r.dbTable).
		// Always filter by owner
		Where("rel_owner = ?", f.OwnedBy)

	if len(f.Prefix) > 0 {
		lookup = lookup.Where("name LIKE ?", f.Prefix+"%")
	}

	if f.Page > 0 {
		lookup = lookup.Offset(f.PerPage * f.Page)
	}

	if f.PerPage > 0 {
		lookup = lookup.Limit(f.PerPage)
	}

	if query, args, err := lookup.ToSql(); err != nil {
		return nil, errors.Wrap(err, "could not build lookup query for settings")
	} else if err = r.db().Select(&ss, query, args...); err != nil {
		return nil, errors.Wrap(err, "could not find settings")
	} else {
		return ss, nil
	}
}

func (r repository) BulkSet(vv ValueSet) error {
	// Save all inside a db transaction
	return r.db().Transaction(func() (err error) {
		return vv.Walk(func(v *Value) error {
			return r.Set(v)
		})
	})
}

func (r *repository) Set(value *Value) error {
	value.UpdatedAt = time.Now()
	return r.db().Replace(r.dbTable, value)
}

func (r *repository) Delete(name string, ownedBy uint64) error {
	_, err := r.db().Exec(
		fmt.Sprintf("DELETE FROM %s WHERE name = ? AND rel_owner = ? ", r.dbTable),
		name,
		ownedBy,
	)
	return err
}

func (r *repository) Get(name string, ownedBy uint64) (value *Value, err error) {
	lookup := squirrel.
		Select(r.columns()...).
		From(r.dbTable).
		Where("rel_owner = ?", ownedBy).
		Where("name = ?", name)

	value = &Value{}

	if query, args, err := lookup.ToSql(); err != nil {
		return nil, errors.Wrap(err, "could not build lookup query for settings")
	} else if err = r.db().Get(value, query, args...); err != nil {
		return nil, errors.Wrap(err, "could not get settings")
	} else if value.Name == "" {
		return nil, nil
	} else {
		return value, nil
	}

}

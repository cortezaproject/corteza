package rh

import (
	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"
)

func Delete(db *factory.DB, table string, cnd squirrel.Sqlizer) error {
	_, err := squirrel.ExecWith(db, squirrel.Delete(table).Where(cnd))
	return err
}

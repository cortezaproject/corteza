package rh

import (
	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"
)

type (
	Set map[string]interface{}
)

// UpdateColumns constructs and executes an update query
func UpdateColumns(db *factory.DB, table string, set Set, cnd squirrel.Sqlizer) error {
	_, err := squirrel.ExecWith(db, squirrel.Update(table).SetMap(set).Where(cnd))
	return err
}

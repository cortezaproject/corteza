package repository

import (
	"fmt"

	"github.com/pkg/errors"
)

func (r repository) updateColumnByID(tableName, columnName string, value interface{}, id uint64) (err error) {
	return exec(r.db().Exec(
		fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", tableName, columnName),
		value,
		id))
}

func exec(_ interface{}, err error) error {
	return errors.WithStack(err)
}

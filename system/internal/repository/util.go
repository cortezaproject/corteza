package repository

import (
	"fmt"
	"time"
)

func (r repository) updateColumnByID(tableName, columnName string, value interface{}, id uint64) (err error) {
	return exec(r.db().Exec(
		fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", tableName, columnName),
		value,
		id))
}

func exec(_ interface{}, err error) error {
	return err
}

// Returns err if set otherwise it returns nerr if not valid
func isFound(err error, valid bool, nerr error) error {
	if err != nil {
		return err
	} else if !valid {
		return nerr
	}

	return nil
}

func timeNowPtr() *time.Time {
	n := time.Now()
	return &n
}

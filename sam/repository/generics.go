package repository

import (
	"context"
	"fmt"
	"github.com/titpetric/factory"
)

func simpleUpdate(ctx context.Context, tableName, columnName string, value interface{}, id uint64) (err error) {
	db := factory.Database.MustGet()

	sql := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", tableName, columnName)

	_, err = db.With(ctx).Exec(sql, value, id)
	return err
}

func simpleDelete(ctx context.Context, tableName string, id uint64) (err error) {
	db := factory.Database.MustGet()

	sql := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)

	_, err = db.With(ctx).Exec(sql, id)
	return err
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

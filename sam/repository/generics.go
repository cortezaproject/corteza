package repository

import (
	"context"
	"fmt"
	"github.com/titpetric/factory"
)

func simpleUpdate(ctx context.Context, tableName, columnName string, value interface{}, id uint64) (err error) {
	db := factory.Database.MustGet()

	sql := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", tableName, columnName)

	_, err = db.ExecContext(ctx, sql, value, id)
	return err
}

func simpleDelete(ctx context.Context, tableName string, id uint64) (err error) {
	db := factory.Database.MustGet()

	sql := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)

	_, err = db.ExecContext(ctx, sql, id)
	return err
}

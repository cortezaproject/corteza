package repository

import (
	"context"
	"fmt"
	"github.com/titpetric/factory"
)

func simpleUpdate(ctx context.Context, tableName, columnName string, value interface{}, id uint64) error {
	db := factory.Database.MustGet()

	sql := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", tableName, columnName)

	if _, err := db.Exec(sql, value, id); err != nil {
		return ErrDatabaseError
	} else {
		return nil
	}
}

func simpleDelete(ctx context.Context, tableName string, id uint64) error {
	db := factory.Database.MustGet()

	sql := fmt.Sprintf("DELETE %s WHERE id = ?", tableName)

	if _, err := db.Exec(sql, id); err != nil {
		return ErrDatabaseError
	} else {
		return nil
	}
}

package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/titpetric/factory"
	"time"
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

func timeNowPtr() *time.Time {
	n := time.Now()
	return &n
}

func coalesceJson(vals ...json.RawMessage) json.RawMessage {
	for _, val := range vals {
		if val != nil {
			return val
		}
	}

	return nil
}

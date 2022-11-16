package rdbms

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/store"
	"strings"
)

// UpsertBuilder combines insert and update builder into upsert (INSERT ... ON CONFLICT UPDATE ...)
//
// It accepts name of the table we're inserting to, payload with values that should be inserted and list of
// columns that represent unique keys
//
// https://www.postgresql.org/docs/12/sql-insert.html
// https://sqlite.org/lang_UPSERT.html
func UpsertBuilder(cfg *Config, table string, payload store.Payload, columns ...string) (squirrel.InsertBuilder, error) {
	var (
		updateCond    = squirrel.Eq{}
		updatePayload = store.Payload{}

		// where to cutoff the update query
		updateCutoff = len("UPDATE " + table + " ")
	)

	for k, v := range payload {
		updatePayload[k] = v
	}

	for _, c := range columns {
		if _, has := updatePayload[c]; has {
			delete(updatePayload, c)
			updateCond[c] = payload[c]
		}
	}

	sql, args, err := squirrel.
		Update(table).
		SetMap(updatePayload).
		//Where(updateCond).
		ToSql()

	if err != nil {
		return squirrel.InsertBuilder{}, err
	}

	suffix := fmt.Sprintf("ON CONFLICT (%s) DO UPDATE %s", strings.Join(columns, ","), sql[updateCutoff:])

	return squirrel.
		Insert(table).
		PlaceholderFormat(cfg.PlaceholderFormat).
		SetMap(payload).
		Suffix(suffix, args...), nil
}

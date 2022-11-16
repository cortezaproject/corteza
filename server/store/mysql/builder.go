package mysql

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/rdbms"
)

// UpsertBuilder combines insert and update builder into upsert (INSERT ... ON CONFLICT UPDATE ...)
//
// It accepts name of the table we're inserting to, payload with values that should be inserted and list of
// columns that represent unique keys
//
// This is MySQL eddition that differs a bit from the standard SQL
// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
func UpsertBuilder(cfg *rdbms.Config, table string, payload store.Payload, columns ...string) (squirrel.InsertBuilder, error) {
	var (
		updatePayload = store.Payload{}

		// where to cutoff the update query
		updateCutoff = len("UPDATE " + table + " SET ")
	)

	for k, v := range payload {
		updatePayload[k] = v
	}

	for _, c := range columns {
		if _, has := updatePayload[c]; has {
			delete(updatePayload, c)
		}
	}

	sql, args, err := squirrel.
		Update(table).
		PlaceholderFormat(cfg.PlaceholderFormat).
		SetMap(updatePayload).
		ToSql()

	if err != nil {
		return squirrel.InsertBuilder{}, err
	}

	return squirrel.
		Insert(table).
		PlaceholderFormat(cfg.PlaceholderFormat).
		SetMap(payload).
		Suffix("ON DUPLICATE KEY UPDATE "+sql[updateCutoff:], args...), nil
}

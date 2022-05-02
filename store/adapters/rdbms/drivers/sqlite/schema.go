package sqlite

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/jmoiron/sqlx"
)

type (
	schema struct {
		dialect *dialect
	}
)

func (s *schema) TableExists(ctx context.Context, db sqlx.QueryerContext, table string) (bool, error) {
	var (
		exists bool
		sql    = `SELECT COUNT(*) > 0 FROM sqlite_master WHERE type = 'table' AND name = ?`
	)

	if err := sqlx.GetContext(ctx, db, &exists, sql, table); err != nil {
		return false, fmt.Errorf("could not check if table exists: %w", err)
	}

	return exists, nil
}

func (s *schema) CreateTable(ctx context.Context, db sqlx.ExtContext, t *ddl.Table) (err error) {
	tt := append(
		[]any{&ddl.CreateTableTemplate{
			Table:         t,
			TrColumnTypes: columnTypTranslator,
		}},
		ddl.CreateIndexTemplates(&ddl.CreateIndexTemplate{OmitFieldLength: true}, t.Indexes...)...,
	)

	return ddl.Exec(ctx, db, tt...)
}

func columnTypTranslator(ct ddl.ColumnType) string {
	switch ct.Type {
	case ddl.ColumnTypeTimestamp:
		return "TIMESTAMP"
	case ddl.ColumnTypeBinary:
		return "BLOB"
	default:
		return ddl.ColumnTypeTranslator(ct)
	}

}

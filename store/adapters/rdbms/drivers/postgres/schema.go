package postgres

import (
	"context"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// PostgreSQL specific prefixes, sql
// templates, functions and other helpers

type (
	schema struct {
		schemaName string
		dialect    goqu.DialectWrapper
	}
)

func (s *schema) TableExists(ctx context.Context, db sqlx.QueryerContext, table string) (bool, error) {
	return ddl.TableExists(ctx, db, s.dialect, table, "public")
}

func (s *schema) CreateTable(ctx context.Context, db sqlx.ExtContext, t *ddl.Table) (err error) {
	tt := append([]any{
		&ddl.CreateTableTemplate{
			Table:        t,
			SuffixClause: " WITHOUT OIDS",
		}},
		ddl.CreateIndexTemplates(&ddl.CreateIndexTemplate{OmitFieldLength: true}, t.Indexes...)...,
	)

	return ddl.Exec(ctx, db, tt...)
}

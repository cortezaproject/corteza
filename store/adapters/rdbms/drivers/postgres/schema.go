package postgres

import (
	"context"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/jmoiron/sqlx"
)

// PostgreSQL specific prefixes, sql
// templates, functions and other helpers

type (
	schema struct {
		schemaName string
	}
)

var (
	_ = &schema{}
)

func (s *schema) TableExists(ctx context.Context, db sqlx.QueryerContext, table string) (bool, error) {
	return ddl.TableExists(ctx, db, Dialect(), table, "public")
}

// ColumnExists  checks if column exists in the MySQL table
func (s *schema) ColumnExists(ctx context.Context, db sqlx.QueryerContext, column, table string) (bool, error) {
	return ddl.ColumnExists(ctx, db, Dialect(), column, table, "public")
}

func (s *schema) CreateTable(ctx context.Context, db sqlx.ExtContext, t *ddl.Table) (err error) {
	tt := append([]any{
		&ddl.CreateTable{
			Dialect:      Dialect(),
			Table:        t,
			SuffixClause: " WITHOUT OIDS",
		}},
		ddl.CreateIndexTemplates(&ddl.CreateIndex{OmitFieldLength: true}, t.Indexes...)...,
	)

	return ddl.Exec(ctx, db, tt...)
}

func (s *schema) AddColumn(ctx context.Context, db sqlx.ExtContext, t *ddl.Table, cc ...*ddl.Column) (err error) {
	var (
		aux    []any
		exists bool
	)

	for _, c := range cc {
		// check column existence
		if exists, err = s.ColumnExists(ctx, db, c.Name, t.Name); err != nil {
			return
		} else if exists {
			// column exists
			continue
		}

		aux = append(aux, &ddl.AddColumn{
			Dialect: dialect,
			Table:   t,
			Column:  c,
		})
	}

	return ddl.Exec(ctx, db, aux...)
}

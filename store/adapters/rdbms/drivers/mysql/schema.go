package mysql

import (
	"context"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/jmoiron/sqlx"
)

type (
	schema struct {
		dbName string
	}
)

// TableExists  checks if table exists in the MySQL database
func (s *schema) TableExists(ctx context.Context, db sqlx.QueryerContext, table string) (bool, error) {
	return ddl.TableExists(ctx, db, Dialect(), table, s.dbName)
}

// ColumnExists  checks if column exists in the MySQL table
func (s *schema) ColumnExists(ctx context.Context, db sqlx.QueryerContext, column, table string) (bool, error) {
	return ddl.ColumnExists(ctx, db, Dialect(), column, table, s.dbName)
}

// CreateTable
//
// MySQL does not hav CREATE-INDEX-IF-NOT-EXISTS; we need to check index existence manually
func (s *schema) CreateTable(ctx context.Context, db sqlx.ExtContext, t *ddl.Table) (err error) {
	tc := &ddl.CreateTable{
		Dialect:      Dialect(),
		Table:        t,
		SuffixClause: "ENGINE=InnoDB DEFAULT CHARSET=utf8",
	}

	if err = ddl.Exec(ctx, db, tc); err != nil {
		return
	}

	for _, index := range t.Indexes {
		if index.Condition != "" {
			// MySQL, sad little DB does not support
			// conditional indexes
			//
			// We'll solve this on an application level
			continue
		}

		var doesIt bool
		if doesIt, err = ddl.IndexExists(ctx, db, Dialect(), index.Name, index.Table, s.dbName); err != nil {
			return
		} else if doesIt {
			continue
		}

		ic := &ddl.CreateIndex{
			Dialect:               Dialect(),
			Index:                 index,
			OmitIfNotExistsClause: true,
		}

		if err = ddl.Exec(ctx, db, ic); err != nil {
			return
		}

	}

	return
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

		// Sadly, some column types in MySQL can not have default values
		if c.Type.Type == ddl.ColumnTypeJson || c.Type.Type == ddl.ColumnTypeBinary || c.Type.Type == ddl.ColumnTypeText {
			c.DefaultValue = ""
		}

		aux = append(aux, &ddl.AddColumn{
			Dialect: dialect,
			Table:   t,
			Column:  c,
		})
	}

	return ddl.Exec(ctx, db, aux...)
}

func (s *schema) DropColumn(ctx context.Context, db sqlx.ExtContext, t *ddl.Table, cc ...string) (err error) {
	var (
		aux    []any
		exists bool
	)

	for _, c := range cc {
		// check column existence
		if exists, err = s.ColumnExists(ctx, db, c, t.Name); err != nil {
			return
		} else if !exists {
			// column exists
			continue
		}

		aux = append(aux, &ddl.DropColumn{
			Dialect: dialect,
			Table:   t,
			Column:  c,
		})
	}

	return ddl.Exec(ctx, db, aux...)
}

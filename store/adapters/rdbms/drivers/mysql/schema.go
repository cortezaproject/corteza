package mysql

import (
	"context"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/jmoiron/sqlx"
)

type (
	schema struct {
		dbName  string
		dialect *dialect
	}
)

func (s *schema) TableExists(ctx context.Context, db sqlx.QueryerContext, table string) (bool, error) {
	return ddl.TableExists(ctx, db, s.dialect, table, "public")
}

// CreateTable
//
// MySQL does not hav CREATE-INDEX-IF-NOT-EXISTS; we need to check index existance manually
func (s *schema) CreateTable(ctx context.Context, db sqlx.ExtContext, t *ddl.Table) (err error) {
	tc := &ddl.CreateTableTemplate{
		Table:         t,
		TrColumnTypes: columnTypTranslator,
		SuffixClause:  "ENGINE=InnoDB DEFAULT CHARSET=utf8",
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
		if doesIt, err = ddl.IndexExists(ctx, db, s.dialect, index.Name, index.Table, s.dbName); err != nil {
			return
		} else if doesIt {
			continue
		}

		ic := &ddl.CreateIndexTemplate{
			Index:                 index,
			OmitIfNotExistsClause: true,
		}

		if err = ddl.Exec(ctx, db, ic); err != nil {
			return
		}

	}

	return
}

func columnTypTranslator(ct ddl.ColumnType) string {
	switch ct.Type {
	case ddl.ColumnTypeIdentifier:
		return "BIGINT UNSIGNED"
	case ddl.ColumnTypeText:
		// @todo when compose_record_value is removed, this will no longer be needed
		if y, has := ct.Flags["mysqlLongText"].(bool); has && y {
			return "LONGTEXT"
		}

		return "TEXT"
	case ddl.ColumnTypeBinary:
		return "BLOB"
	case ddl.ColumnTypeTimestamp:
		return "DATETIME"
	case ddl.ColumnTypeBoolean:
		return "TINYINT(1)"
	default:
		return ddl.ColumnTypeTranslator(ct)
	}
}

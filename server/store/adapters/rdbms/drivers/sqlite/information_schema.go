package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ddl"
	"github.com/jmoiron/sqlx"
)

type (
	// @todo this is unmodified copy of mysql's information schema struct!
	informationSchema struct {
		conn *sqlx.DB
	}
)

func InformationSchema(conn *sqlx.DB) *informationSchema {
	return &informationSchema{
		conn: conn,
	}
}

func (i *informationSchema) TableLookup(ctx context.Context, table, schema, dbname string) (tbl *ddl.Table, err error) {
	var (
		// https://www.sqlite.org/pragma.html
		query = fmt.Sprintf(`PRAGMA table_info(%s)`, Dialect().QuoteIdent(table))
		aux   = make([]struct {
			CID          any            `db:"cid"`
			Name         string         `db:"name"`
			Type         string         `db:"type"`
			NotNull      bool           `db:"notnull"`
			DefaultValue sql.NullString `db:"dflt_value"`
			PrimaryKey   int            `db:"pk"`
		}, 0)
	)

	if err = sqlx.SelectContext(ctx, i.conn, &aux, query); err != nil {
		return
	}

	if len(aux) == 0 {
		return nil, errors.NotFound("table does not exist")
	}

	tbl = &ddl.Table{
		Ident:   table,
		Columns: make([]*ddl.Column, len(aux)),
	}

	// iterate over results (aux) and populate ddl.Index struct
	for p, a := range aux {
		tbl.Columns[p] = &ddl.Column{
			Ident: a.Name,
			Type: &ddl.ColumnType{
				Name: a.Type,
				Null: !a.NotNull,
			},
			Default: a.DefaultValue.String,
		}
	}

	return
}

func (i *informationSchema) IndexLookup(ctx context.Context, index, table, schema string) (*ddl.Index, error) {
	// for now, there is no need do implement this
	return nil, errors.NotFound("index %q not found", index)
}

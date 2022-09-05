package mysql

import (
	"context"
	"database/sql"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type (
	informationSchema struct {
		conn *sqlx.DB
	}
)

func InformationSchema(conn *sqlx.DB) *informationSchema {
	return &informationSchema{
		conn: conn,
	}
}

func (i *informationSchema) TableLookup(ctx context.Context, table, dbname string) (*ddl.Table, error) {
	var (
		oneTable = i.columnSelect().Where(
			exp.ParseIdentifier("TABLE_CATALOG").Eq("def"),
			// schema == dbname in MySQL world
			exp.ParseIdentifier("TABLE_SCHEMA").Eq(dbname),
			exp.ParseIdentifier("TABLE_NAME").Eq(table),
		)
	)

	if out, err := i.scanColumns(ctx, oneTable); err != nil {
		return nil, err
	} else if len(out) > 0 {
		return out[0], nil
	} else {
		return nil, errors.NotFound("table does not exist")
	}
}

func (i *informationSchema) columnSelect() *goqu.SelectDataset {
	return dialect.GOQU().Select(
		"TABLE_NAME",
		"COLUMN_NAME",
		"IS_NULLABLE",
		"DATA_TYPE",
	).
		From("information_schema.columns").
		Order(
			exp.NewOrderedExpression(exp.ParseIdentifier("TABLE_SCHEMA"), exp.AscDir, exp.NoNullsSortType),
			exp.NewOrderedExpression(exp.ParseIdentifier("ORDINAL_POSITION"), exp.AscDir, exp.NoNullsSortType),
		)
}

func (i *informationSchema) scanColumns(ctx context.Context, sd *goqu.SelectDataset) (out []*ddl.Table, err error) {
	var (
		at  int
		has bool
		n2p = make(map[string]int)

		// https://dev.mysql.com/doc/mysql-infoschema-excerpt/5.7/en/information-schema-statistics-table.html
		aux = make([]struct {
			Table      string `db:"TABLE_NAME"`
			Column     string `db:"COLUMN_NAME"`
			IsNullable string `db:"IS_NULLABLE"`
			Type       string `db:"DATA_TYPE"`
		}, 0)
	)

	if err = ddl.Structs(ctx, i.conn, sd, &aux); err != nil {
		return
	}

	out = make([]*ddl.Table, 0, 10)

	for _, v := range aux {
		if at, has = n2p[v.Table]; !has {
			at = len(out)
			n2p[v.Table] = at
			out = append(out, &ddl.Table{Ident: v.Table})
		}

		out[at].Columns = append(out[at].Columns, &ddl.Column{
			Ident: v.Column,
			Type: &ddl.ColumnType{
				Name: v.Type,
				Null: v.IsNullable == "YES",
			},
		})
	}

	return
}

func (i *informationSchema) IndexLookup(ctx context.Context, index, table, dbname string) (*ddl.Index, error) {
	var (
		oneIndexOnly = i.indexSelect().Where(
			exp.ParseIdentifier("TABLE_SCHEMA").Eq(dbname),
			exp.ParseIdentifier("INDEX_NAME").Eq(index),
			exp.ParseIdentifier("TABLE_NAME").Eq(table),
		)
	)

	if out, err := i.scanIndexes(ctx, oneIndexOnly); err != nil {
		return nil, err
	} else if len(out) > 0 {
		return out[0], nil
	} else {
		return nil, nil
	}
}

func (i *informationSchema) indexSelect() *goqu.SelectDataset {
	return dialect.GOQU().Select(
		"INDEX_NAME",
		"TABLE_NAME",
		"COLUMN_NAME",
		"INDEX_TYPE",
		"EXPRESSION",
		"INDEX_COMMENT",
	).
		From("information_schema.statistics").
		Order(
			exp.NewOrderedExpression(exp.ParseIdentifier("TABLE_SCHEMA"), exp.AscDir, exp.NoNullsSortType),
			exp.NewOrderedExpression(exp.ParseIdentifier("INDEX_NAME"), exp.AscDir, exp.NoNullsSortType),
			exp.NewOrderedExpression(exp.ParseIdentifier("SEQ_IN_INDEX"), exp.AscDir, exp.NoNullsSortType),
		)
}

func (i *informationSchema) scanIndexes(ctx context.Context, sd *goqu.SelectDataset) (out []*ddl.Index, err error) {
	var (
		at  int
		has bool
		n2p = make(map[string]int)

		// https://dev.mysql.com/doc/mysql-infoschema-excerpt/5.7/en/information-schema-statistics-table.html
		aux = make([]struct {
			Name      string `db:"INDEX_NAME"`
			Table     string `db:"TABLE_NAME"`
			NonUnique bool   `db:"NON_UNIQUE"`
			Type      string `db:"INDEX_TYPE"`

			// @todo there's also a "COMMENT" column?
			Comment string `db:"INDEX_COMMENT"`

			Expression      sql.NullString `db:"EXPRESSION"`
			ColumnName      string         `db:"COLUMN_NAME"`
			ColumnSubPart   sql.NullInt32  `db:"SUB_PART"`
			ColumnCollation sql.NullString `db:"COLLATION"`

			// stats
			ColumnStatsCardinality int64 `db:"CARDINALITY"`
		}, 0)
	)

	if err = ddl.Structs(ctx, i.conn, sd, &aux); err != nil {
		return
	}

	out = make([]*ddl.Index, 0, 10)

	// iterate over results (aux) and populate ddl.Index struct
	for p, a := range aux {
		if at, has = n2p[a.Name]; !has {
			out = append(out, &ddl.Index{
				Ident:      a.Name,
				TableIdent: a.Table,
				Type:       a.Type,
				Comment:    a.Comment,
				Fields:     make([]*ddl.IndexField, 0),
				Unique:     !a.NonUnique,
				Predicate:  "",
				Meta:       nil,
			})
			n2p[a.Name] = p
			at = p
		}

		col := &ddl.IndexField{
			Length: int(a.ColumnSubPart.Int32),
			Statistics: &ddl.IndexFieldStatistics{
				Cardinality: a.ColumnStatsCardinality,
			},
		}

		switch a.ColumnCollation.String {
		case "A":
			col.Sorted = ddl.IndexFieldSortAsc
		case "D":
			col.Sorted = ddl.IndexFieldSortDesc
		}

		if a.Expression.Valid {
			col.Expression = a.Expression.String
		} else {
			col.Column = a.ColumnName
		}

		out[at].Fields = append(out[at].Fields, col)
	}

	return
}

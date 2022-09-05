package sqlite

import (
	"context"
	"database/sql"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
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

func (i *informationSchema) IndexLookup(ctx context.Context, index, table, schema string) (*ddl.Index, error) {
	var (
		oneIndex = i.indexSelect(schema).Where(
			exp.ParseIdentifier("INDEX_NAME").Eq(index),
			exp.ParseIdentifier("TABLE_NAME").Eq(table),
		)
	)

	if out, err := i.scanIndexes(ctx, oneIndex); err != nil {
		return nil, err
	} else if len(out) > 0 {
		return out[0], nil
	} else {
		return nil, nil
	}
}

func (i *informationSchema) indexSelect(schema string) *goqu.SelectDataset {
	return dialect.GOQU().Select(
		"INDEX_NAME",
		"TABLE_NAME",
		"COLUMN_NAME",
		"INDEX_TYPE",
		"EXPRESSION",
		"INDEX_COMMENT",
	).
		From("information_schema.statistics").
		Where(
			exp.ParseIdentifier("TABLE_SCHEMA").Eq(schema),
		).
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

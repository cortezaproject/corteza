package postgres

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"strings"
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

func (i *informationSchema) TableLookup(ctx context.Context, table, schema, dbname string) (*ddl.Table, error) {
	var (
		oneTable = i.columnSelect().Where(
			exp.ParseIdentifier("table_name").Eq(table),
			exp.ParseIdentifier("table_catalog").Eq(dbname),
			exp.ParseIdentifier("table_schema").Eq(schema),
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
		"table_name",
		"column_name",
		"is_nullable",
		"data_type",
	).
		From("information_schema.columns").
		Order(
			exp.NewOrderedExpression(exp.ParseIdentifier("table_schema"), exp.AscDir, exp.NoNullsSortType),
			exp.NewOrderedExpression(exp.ParseIdentifier("ordinal_position"), exp.AscDir, exp.NoNullsSortType),
		)
}

func (i *informationSchema) scanColumns(ctx context.Context, sd *goqu.SelectDataset) (out []*ddl.Table, err error) {
	var (
		at  int
		has bool
		n2p = make(map[string]int)

		// https://dev.mysql.com/doc/mysql-infoschema-excerpt/5.7/en/information-schema-statistics-table.html
		aux = make([]struct {
			Table      string `db:"table_name"`
			Column     string `db:"column_name"`
			IsNullable bool   `db:"is_nullable"`
			Type       string `db:"data_type"`
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
				Null: v.IsNullable,
			},
		})
	}

	return
}

func (i *informationSchema) IndexLookup(ctx context.Context, index, table, schema string) (*ddl.Index, error) {
	var (
		oneIndex = i.indexSelect().Where(
			exp.ParseIdentifier("indexname").Eq(index),
			exp.ParseIdentifier("tablename").Eq(table),
			exp.ParseIdentifier("schemaname").Eq(schema),
		)
	)

	if out, err := i.scanIndexes(ctx, oneIndex); err != nil {
		return nil, err
	} else if len(out) > 0 {
		return out[0], nil
	} else {
		return nil, errors.NotFound("index does not exist")
	}
}

func (i *informationSchema) indexSelect() *goqu.SelectDataset {
	return dialect.GOQU().Select(
		"tablename",
		"indexname",
		"indexdef",
	).
		From("pg_indexes").
		Order(
			exp.NewOrderedExpression(exp.ParseIdentifier("tablename"), exp.AscDir, exp.NoNullsSortType),
			exp.NewOrderedExpression(exp.ParseIdentifier("indexname"), exp.AscDir, exp.NoNullsSortType),
		)
}

func (i *informationSchema) scanIndexes(ctx context.Context, sd *goqu.SelectDataset) (out []*ddl.Index, err error) {
	var (
		// https://dev.mysql.com/doc/mysql-infoschema-excerpt/5.7/en/information-schema-statistics-table.html
		aux = make([]struct {
			Table string `db:"tablename"`
			Index string `db:"indexname"`
			Def   string `db:"indexdef"`
		}, 0)
	)

	if err = ddl.Structs(ctx, i.conn, sd, &aux); err != nil {
		return
	}

	out = make([]*ddl.Index, len(aux))

	//// iterate over results (aux) and populate ddl.Index struct
	for i, a := range aux {
		out[i] = &ddl.Index{}
		if err = parseIndexDefinition(a.Def, out[i]); err != nil {
			return
		}
	}

	return
}

// parses postgesql specific index definition
func parseIndexDefinition(def string, index *ddl.Index) error {
	const (
		TOKEN_INDEX  = " INDEX "
		TOKEN_ON     = " ON "
		TOKEN_USING  = " USING "
		TOKEN_UNIQUE = " UNIQUE INDEX "
		TOKEN_WHERE  = " WHERE "
	)

	// if def contains TOKEN_UNIQUE, consider this an unique index
	index.Unique = strings.Contains(def, TOKEN_UNIQUE)

	// parse postgesql index definition
	tpi := strings.Index(def, TOKEN_INDEX)
	tpo := strings.Index(def, TOKEN_ON)
	tpu := strings.Index(def, TOKEN_USING)
	tpw := strings.Index(def, TOKEN_WHERE)

	if tpw == -1 {
		tpw = len(def)
	}

	// extract index name
	index.Ident = def[tpi+len(TOKEN_INDEX) : tpo]

	// extract index table
	index.TableIdent = def[tpo+len(TOKEN_ON) : tpu]
	index.TableIdent = strings.SplitN(index.TableIdent, ".", 2)[1]

	// extract index fields
	fieldsDef := def[tpu+len(TOKEN_USING) : tpw]

	typeAndFields := strings.SplitN(fieldsDef, " ", 2)

	index.Type = typeAndFields[0]

	fields := strings.Split(typeAndFields[1][1:len(typeAndFields[1])-1], ", ")
	index.Fields = make([]*ddl.IndexField, len(fields))
	for i, fd := range fields {
		f := &ddl.IndexField{}

		if strings.Contains(fd, "(") {
			f.Expression = fd
		} else {
			f.Column = fd
		}

		// @todo direction/sorting
		index.Fields[i] = f
	}

	// extract predicate/where
	if tpw != len(def) {
		index.Predicate = def[tpw+len(TOKEN_WHERE)+1 : len(def)-1]
	}

	return nil
}

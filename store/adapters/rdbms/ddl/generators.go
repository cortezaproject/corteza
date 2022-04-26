package ddl

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type (
	trColTypeFn func(ColumnType) string

	CreateTableTemplate struct {
		*Table
		OmitIfNotExistsClause bool
		SuffixClause          string
		TrColumnTypes         trColTypeFn
	}

	CreateIndexTemplate struct {
		*Index
		OmitIfNotExistsClause bool
		OmitFieldLength       bool
	}
)

func CreateIndexTemplates(base *CreateIndexTemplate, ii ...*Index) []any {
	var (
		tt = make([]any, len(ii))
	)

	for i := range ii {
		tt[i] = &CreateIndexTemplate{
			Index:                 ii[i],
			OmitIfNotExistsClause: base.OmitIfNotExistsClause,
			OmitFieldLength:       base.OmitFieldLength,
		}
	}

	return tt
}

// utility for executing series fo commands
func Exec(ctx context.Context, db sqlx.ExecerContext, ss ...any) (err error) {
	for _, s := range ss {
		var (
			sql  string
			args []any
		)

		switch c := s.(type) {
		case string:
			sql = c
		case fmt.Stringer:
			sql = c.String()
		case exp.SQLExpression:
			sql, args, err = c.ToSQL()
		default:
			panic(fmt.Sprintf("unexecutable input (%T)", s))
		}

		spew.Dump(sql, args)
		if _, err = db.ExecContext(ctx, sql, args...); err != nil {
			return
		}
	}

	return
}

func TableExists(ctx context.Context, db sqlx.QueryerContext, d goqu.DialectWrapper, table, schema string) (bool, error) {
	return GetBool(ctx, db, GenTableCheck(d, table, schema))
}

func GenTableCheck(d goqu.DialectWrapper, table, schema string) *goqu.SelectDataset {
	return d.Select(goqu.COUNT(goqu.Star()).Gt(0)).
		From("information_schema.tables").
		Where(
			exp.ParseIdentifier("table_name").Eq(table),
			exp.ParseIdentifier("table_schema").Eq(schema),
		)
}

func IndexExists(ctx context.Context, db sqlx.QueryerContext, d goqu.DialectWrapper, index, table, schema string) (bool, error) {
	return GetBool(ctx, db, GenIndexCheck(d, index, table, schema))
}

func GenIndexCheck(d goqu.DialectWrapper, index, table, schema string) *goqu.SelectDataset {
	return d.Select(goqu.COUNT(goqu.Star()).Gt(0)).
		From("information_schema.statistics").
		Where(
			exp.ParseIdentifier("index_name").Eq(index),
			exp.ParseIdentifier("table_name").Eq(table),
			exp.ParseIdentifier("table_schema").Eq(schema),
		)
}

func (t *CreateTableTemplate) String() string {
	if t.TrColumnTypes == nil {
		t.TrColumnTypes = ColumnTypeTranslator
	}

	sql := "CREATE TABLE "

	if !t.OmitIfNotExistsClause {
		sql += "IF NOT EXISTS "
	}

	sql += "\"" + t.Name + "\" (\n"
	sql += GenCreateTableBody(t.Table, t.TrColumnTypes)
	sql += "\n)"
	sql += t.SuffixClause

	return sql
}

func GenCreateTableBody(t *Table, trColType trColTypeFn) string {
	sql := ""

	for c, col := range t.Columns {
		if c == 0 {
			sql += "  "
		} else {
			sql += ", "
		}

		sql += GenTableColumn(col, trColType)

		sql += "\n"
	}

	if t.PrimaryKey != nil {
		sql += "\n, " + GenPrimaryKey(t.PrimaryKey)
	}

	return sql
}

func GenTableColumn(col *Column, trType trColTypeFn) string {
	sql := "\"" + col.Name + "\" " + trType(col.Type)

	if !col.IsNull {
		sql += " NOT NULL"
	}

	if col.DefaultValue > "" {
		sql += " DEFAULT " + col.DefaultValue
	}

	return sql
}

func GenPrimaryKey(pk *Index) string {
	sql := "PRIMARY KEY ("
	for f, field := range pk.Fields {
		if f > 0 {
			sql += ", "
		}
		sql += field.Field
	}
	sql += ")"

	return sql
}

func (t *CreateIndexTemplate) String() string {
	sql := "CREATE "

	if t.Index.Unique {
		sql += "UNIQUE "
	}

	sql += "INDEX "

	if !t.OmitIfNotExistsClause {
		sql += "IF NOT EXISTS "
	}

	sql += "\"" + t.Index.Name + "\" ON \"" + t.Index.Table + "\" ("

	for f, field := range t.Index.Fields {
		if f > 0 {
			sql += ", "
		}

		if field.Expr {
			sql += "("
		}

		sql += field.Field

		if field.Desc {
			sql += " DESC"
		}

		if field.Length > 0 && !t.OmitFieldLength {
			sql += fmt.Sprintf("(%d)", field.Length)
		}

		if field.Expr {
			sql += ")"
		}
	}
	sql += ")"

	if t.Index.Condition != "" {
		sql += " WHERE " + t.Index.Condition
	}

	return sql
}

func GetBool(ctx context.Context, db sqlx.QueryerContext, query exp.SQLExpression) (bool, error) {
	var (
		exists         bool
		sql, args, err = query.ToSQL()
	)

	if err != nil {
		return false, fmt.Errorf("could not generate SQLk")
	}

	if err = sqlx.GetContext(ctx, db, &exists, sql, args...); err != nil {
		return false, err
	}

	return exists, nil
}

// ColumnTypeTranslator is the most generic translator of "corteza types"
// to db-native column types.
//
// @todo it might be smart to merge this with data.AttributeType (part of CRS feature)
func ColumnTypeTranslator(ct ColumnType) string {
	switch ct.Type {
	case ColumnTypeIdentifier:
		return "BIGINT"
	case ColumnTypeVarchar:
		if ct.Length > 0 {
			// VARCHAR(0) is useless
			return fmt.Sprintf("VARCHAR(%d)", ct.Length)
		}
		return "VARCHAR"
	case ColumnTypeText:
		return "TEXT"
	case ColumnTypeJson:
		return "JSON"
	case ColumnTypeBinary:
		return "BYTEA"
	case ColumnTypeTimestamp:
		if ct.Length > -1 {
			// TIMESTAMPTZ(0) strips out milliseconds
			return fmt.Sprintf("TIMESTAMPTZ(%d)", ct.Length)
		}

		return "TIMESTAMPTZ"
	case ColumnTypeInteger:
		return "INTEGER"
	case ColumnTypeBoolean:
		return "BOOLEAN"
	default:
		panic(fmt.Sprintf("unhandled column type: %d ", ct.Type))
	}
}

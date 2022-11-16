package ddl

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type (
	dialect interface {
		QuoteIdent(i string) string
	}

	CreateTable struct {
		Dialect               dialect
		Table                 *Table
		OmitIfNotExistsClause bool
		SuffixClause          string
	}

	DropTable struct {
		Dialect dialect
		Table   string
	}

	CreateIndex struct {
		Dialect               dialect
		Index                 *Index
		OmitIfNotExistsClause bool
		OmitFieldLength       bool
	}

	DropIndex struct {
		Dialect    dialect
		Ident      string
		TableIdent string
	}

	AddColumn struct {
		Dialect dialect
		Table   string
		Column  *Column
	}

	DropColumn struct {
		Dialect dialect
		Table   string
		Column  string
	}

	RenameColumn struct {
		Dialect dialect
		Table   string
		Old     string
		New     string
	}
)

// Exec is a utility for executing series of commands
//
// Parameters can be string, Stringer interface or goqu's exp.SQLExpression
//
// Any other type will result in panic
func Exec(ctx context.Context, db sqlx.ExtContext, ss ...any) (err error) {
	for _, s := range ss {
		var (
			sql  string
			args []any
		)

		switch c := s.(type) {
		case interface{ ToSQL() (string, []any, error) }:
			sql, args, err = c.ToSQL()
			if err != nil {
				return
			}
		case fmt.Stringer:
			sql = c.String()
		case string:
			sql = c
		default:
			panic(fmt.Sprintf("unexecutable input (%T)", s))
		}

		if _, err = db.ExecContext(ctx, sql, args...); err != nil {
			return
		}
	}

	return
}

func (t *CreateTable) String() string {
	sql := "CREATE "
	if t.Table.Temporary {
		sql += "TEMPORARY "
	}

	sql += "TABLE "
	if !t.OmitIfNotExistsClause {
		sql += "IF NOT EXISTS "
	}

	sql += t.Dialect.QuoteIdent(t.Table.Ident) + " (\n"
	sql += GenCreateTableBody(t.Dialect, t.Table)
	sql += ")"
	sql += t.SuffixClause

	return sql
}

func (c *DropTable) ToSQL() (sql string, aa []interface{}, err error) {
	return fmt.Sprintf(
		`DROP TABLE %s`,
		c.Dialect.QuoteIdent(c.Table),
	), nil, nil
}

func GenCreateTableBody(d dialect, t *Table) string {
	sql := ""

	for c, col := range t.Columns {
		if c == 0 {
			sql += "  "
		} else {
			sql += ", "
		}

		sql += GenTableColumn(d, col)

		sql += "\n"
	}

	// check if any of the indexes is a primary key
	for _, pk := range t.Indexes {
		if pk.Ident != PRIMARY_KEY {
			continue
		}

		sql += "\n"
		sql += ", " + GenPrimaryKey(d, pk) + "\n"
		break
	}

	return sql
}

func GenTableColumn(d dialect, col *Column) string {
	sql := d.QuoteIdent(col.Ident) + " " + col.Type.Name + " "

	if col.Type.Null {
		sql += "    NULL"
	} else {
		sql += "NOT NULL"
	}

	if col.Default > "" {
		sql += " DEFAULT " + col.Default
	}

	return sql
}

func GenPrimaryKey(d dialect, pk *Index) string {
	sql := "PRIMARY KEY ("
	for f, field := range pk.Fields {
		if f > 0 {
			sql += ", "
		}

		sql += d.QuoteIdent(field.Column)
	}
	sql += ")"

	return sql
}

func (t *CreateIndex) String() string {
	sql := "CREATE "

	if t.Index.Unique {
		sql += "UNIQUE "
	}

	sql += "INDEX "

	if !t.OmitIfNotExistsClause {
		sql += "IF NOT EXISTS "
	}

	sql += t.Dialect.QuoteIdent(t.Index.Ident) + " ON " + t.Dialect.QuoteIdent(t.Index.TableIdent) + " ("

	for f, field := range t.Index.Fields {
		isExpr := len(field.Expression) > 0

		if f > 0 {
			sql += ", "
		}

		if isExpr {
			sql += "("
			sql += field.Expression
		} else {
			sql += t.Dialect.QuoteIdent(field.Column)
		}

		if field.Length > 0 && !t.OmitFieldLength {
			sql += fmt.Sprintf("(%d)", field.Length)
		}

		if isExpr {
			sql += ")"
		}

		switch field.Sort {
		case dal.IndexFieldSortDesc:
			sql += " DESC"
		case dal.IndexFieldSortAsc:
			sql += " ASC"
		}

		switch field.Nulls {
		case dal.IndexFieldNullsLast:
			sql += " NULLS LAST"
		case dal.IndexFieldNullsFirst:
			sql += " NULLS FIRST"
		}

	}
	sql += ")"

	if t.Index.Predicate != "" {
		sql += " WHERE " + t.Index.Predicate
	}

	return sql
}

func (c *AddColumn) ToSQL() (sql string, aa []interface{}, err error) {
	sql = fmt.Sprintf(
		`ALTER TABLE %s ADD COLUMN %s %s`,
		c.Dialect.QuoteIdent(c.Table),
		c.Dialect.QuoteIdent(c.Column.Ident),
		c.Column.Type.Name,
	)

	if !c.Column.Type.Null {
		sql += " NOT NULL"
	}

	if len(c.Column.Default) > 0 {
		// @todo right now we can (and need to) trust that default
		//       values are unharmful!
		sql += " DEFAULT " + c.Column.Default
	}

	return
}

func (c *DropColumn) ToSQL() (sql string, aa []interface{}, err error) {
	return fmt.Sprintf(
		`ALTER TABLE %s DROP COLUMN %s`,
		c.Dialect.QuoteIdent(c.Table),
		c.Dialect.QuoteIdent(c.Column),
	), nil, nil
}

func (c *DropIndex) ToSQL() (sql string, aa []interface{}, err error) {
	return fmt.Sprintf(
		`DROP INDEX %s ON %s`,
		c.Dialect.QuoteIdent(c.Ident),
		c.Dialect.QuoteIdent(c.TableIdent),
	), nil, nil
}

func (c *RenameColumn) ToSQL() (sql string, aa []interface{}, err error) {
	return fmt.Sprintf(
		`ALTER TABLE %s RENAME COLUMN %s TO %s`,
		c.Dialect.QuoteIdent(c.Table),
		c.Dialect.QuoteIdent(c.Old),
		c.Dialect.QuoteIdent(c.New),
	), nil, nil
}

// GetBool is a utility function to simplify getting a boolean value from a query result.
func GetBool(ctx context.Context, db sqlx.QueryerContext, query exp.SQLExpression) (bool, error) {
	var (
		exists         bool
		sql, args, err = query.ToSQL()
	)

	if err != nil {
		return false, fmt.Errorf("could not generate SQL: %v", err)
	}

	if err = sqlx.GetContext(ctx, db, &exists, sql, args...); err != nil {
		return false, err
	}

	return exists, nil
}

// Structs is a utility function to simplify selecting data into slice of structs
func Structs(ctx context.Context, db sqlx.QueryerContext, query exp.SQLExpression, t any) error {
	var (
		sql, args, err = query.ToSQL()
	)

	if err != nil {
		return fmt.Errorf("could not generate SQL: %v", err)
	}

	return sqlx.SelectContext(ctx, db, t, sql, args...)
}

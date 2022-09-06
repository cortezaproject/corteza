package ddl

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type (
	dialect interface {
		// GOQU returns goqu's dialect wrapper struct
		GOQU() goqu.DialectWrapper
	}

	CreateTable struct {
		Table *Table

		OmitIfNotExistsClause bool
		SuffixClause          string
	}

	CreateIndex struct {
		Index                 *Index
		OmitIfNotExistsClause bool
		OmitFieldLength       bool
	}

	DropIndex struct {
		Ident      exp.IdentifierExpression
		TableIdent exp.IdentifierExpression
	}

	AddColumn struct {
		Table  exp.IdentifierExpression
		Column *Column
	}

	DropColumn struct {
		Table  exp.IdentifierExpression
		Column exp.IdentifierExpression
	}

	RenameColumn struct {
		Table exp.IdentifierExpression
		Old   exp.IdentifierExpression
		New   exp.IdentifierExpression
	}
)

func CreateIndexTemplates(base *CreateIndex, ii ...*Index) []any {
	var (
		tt = make([]any, len(ii))
	)

	for i := range ii {
		tt[i] = &CreateIndex{
			Index:                 ii[i],
			OmitIfNotExistsClause: base.OmitIfNotExistsClause,
			OmitFieldLength:       base.OmitFieldLength,
		}
	}

	return tt
}

// Exec is an utility for executing series of commands
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
		case string:
			sql = c
		case fmt.Stringer:
			sql = c.String()
		case exp.SQLExpression:
			sql, args, err = c.ToSQL()
			if err != nil {
				return
			}
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
	sql := "CREATE TABLE "

	if !t.OmitIfNotExistsClause {
		sql += "IF NOT EXISTS "
	}

	sql += "\"" + t.Table.Ident + "\" (\n"
	sql += GenCreateTableBody(t.Table)
	sql += "\n)"
	sql += t.SuffixClause

	return sql
}

func GenCreateTableBody(t *Table) string {
	sql := ""

	for c, col := range t.Columns {
		if c == 0 {
			sql += "  "
		} else {
			sql += ", "
		}

		sql += GenTableColumn(col)

		sql += "\n"
	}

	// check if any of the indexes is a primary key
	for _, pk := range t.Indexes {
		if pk.Ident != PRIMARY_KEY {
			continue
		}

		sql += ", " + GenPrimaryKey(pk) + "\n"
		break
	}

	return sql
}

func GenTableColumn(col *Column) string {
	sql := `"` + col.Ident + `"` + col.Type.Name + ` `

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

func GenPrimaryKey(pk *Index) string {
	sql := "PRIMARY KEY ("
	for f, field := range pk.Fields {
		if f > 0 {
			sql += ", "
		}
		sql += field.Column
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

	sql += "\"" + t.Index.Ident + "\" ON \"" + t.Index.TableIdent + "\" ("

	for f, field := range t.Index.Fields {
		isExpr := len(field.Expression) > 0

		if f > 0 {
			sql += ", "
		}

		if isExpr {
			sql += "("
			sql += field.Expression
		} else {
			sql += field.Column

		}

		switch field.Sorted {
		case IndexFieldSortDesc:
			sql += " DESC"
		case IndexFieldSortAsc:
			sql += " ASC"
		}

		if field.Length > 0 && !t.OmitFieldLength {
			sql += fmt.Sprintf("(%d)", field.Length)
		}

		if isExpr {
			sql += ")"
		}
	}
	sql += ")"

	if t.Index.Predicate != "" {
		sql += " WHERE " + t.Index.Predicate
	}

	return sql
}

//func (c *AddColumn) String() string {
//	sql := "ALTER TABLE" + " " + c.Table + " ADD COLUMN " + c.Column.Name + " " + c.Column.Type
//	if !c.Column.IsNull {
//		sql += " NOT NULL"
//	}
//
//	if len(c.Column.DefaultValue) > 0 {
//		sql += " DEFAULT " + c.Column.DefaultValue
//	}
//
//	return sql
//}

//func (c *DropColumn) String() string {
//	return "ALTER TABLE" + " " + c.Table.Name + " DROP COLUMN " + c.Column
//}

func (c *DropIndex) Express() exp.SQLExpression {
	return SQLExpression(exp.NewLiteralExpression("DROP INDEX ? ON ?", c.Ident, c.TableIdent))
}

//func (c *RenameColumn) String() string {
//	return "ALTER TABLE" + " " + c.Table.Name + " RENAME COLUMN " + c.Old + " TO " + c.New
//}

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

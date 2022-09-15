package ddl

import "github.com/doug-martin/goqu/v9/exp"

// goqu's little helpers

type (
	// literalExpression is a wrapper around goqu's LiteralExpression
	sqlExpression struct {
		exp.LiteralExpression
	}
)

func SQLExpression(e exp.LiteralExpression) exp.SQLExpression {
	return &sqlExpression{e}
}

func (e *sqlExpression) ToSQL() (string, []any, error) { return e.Literal(), e.Args(), nil }
func (e *sqlExpression) IsPrepared() bool              { return e.IsPrepared() }

var (
	_ exp.SQLExpression = &sqlExpression{}
)

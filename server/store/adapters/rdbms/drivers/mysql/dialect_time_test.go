package mysql

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/ql"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/require"
)

// TIME() takes its argument as text: MySQL 8 turns a bound '02:00:00' in
// TIME(?) into 00:00:00
func TestExprHandlerTimeAsText(t *testing.T) {
	req := require.New(t)

	e, err := mysqlDialect{}.ExprHandler(&ql.ASTNode{Ref: "time"}, exp.NewLiteralExpression("?", "02:00:00"))
	req.NoError(err)

	sql, args, err := goquDialectWrapper.Select(e).Prepared(true).ToSQL()
	req.NoError(err)
	req.Contains(sql, "TIME(CAST(? AS CHAR))")
	req.Equal([]any{"02:00:00"}, args)
}

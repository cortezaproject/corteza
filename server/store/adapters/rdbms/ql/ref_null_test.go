package ql

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/require"
)

// null and nnull test their argument when given one, as the paging cursor
// expression uses them; bare, they stay the keyword of an IS comparison
func TestRefNullWithArgument(t *testing.T) {
	var (
		req  = require.New(t)
		base = goqu.Dialect("sqlite").Select().Prepared(true)
		foo  = exp.NewIdentifierExpression("", "", "foo")
	)

	render := func(e exp.Expression) string {
		sql, _, err := base.Where(e).ToSQL()
		req.NoError(err)
		return sql
	}

	req.Contains(render(ref2exp["null"].Handler(foo)), `("foo" IS NULL)`)
	req.Contains(render(ref2exp["nnull"].Handler(foo)), `("foo" IS NOT NULL)`)
	req.Contains(render(ref2exp["is"].Handler(foo, ref2exp["nnull"].Handler())), `"foo" IS NOT NULL`)
}

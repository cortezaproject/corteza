package ql

import (
	"testing"

	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/require"
)

func TestRef_betweenNeedsBothBounds(t *testing.T) {
	hh := ExprHandlerMap{}.ExprHandlers()
	two := []exp.Expression{exp.NewLiteralExpression("a"), exp.NewLiteralExpression("1")}

	for _, ref := range []string{"between", "nbetween"} {
		_, err := hh[ref].HandlerE(two...)
		require.Error(t, err, ref)
	}
}

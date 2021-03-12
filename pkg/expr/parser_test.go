package expr

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParser(t *testing.T) {
	var (
		req    = require.New(t)
		ctx    = context.Background()
		p      = Parser()
		e, err = p.NewEvaluable("0 == 0")

		result bool
	)

	req.NoError(err)

	result, err = e.EvalBool(ctx, nil)
	req.NoError(err)
	req.True(result)
}

func Example_simpleExpresion() {
	eval(`40 + 2`, nil)
	// output:
	// 42
}

package expr

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func Example_simpleExpression() {
	eval(`40 + 2`, nil)
	// output:
	// 42
}

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

func TestGvalParser(t *testing.T) {
	var (
		req     = require.New(t)
		ctx     = context.Background()
		p       = NewGvalParser()
		vv, err = NewVars(map[string]interface{}{
			"vars":  &Vars{},
			"key":   "foo",
			"value": Must(NewString("foo")),
		})
		result interface{}
	)
	req.NoError(err)

	pp, err := p.Parse("toJSON(set(vars, key, value))")
	req.NoError(err)

	result, err = pp.Eval(ctx, vv)
	req.NoError(err)
	req.Equal("{\"foo\":{\"@value\":\"foo\",\"@type\":\"String\"}}", result)
}

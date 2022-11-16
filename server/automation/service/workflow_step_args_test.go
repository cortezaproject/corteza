package service

import (
	"context"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessArguments(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var (
			req    = require.New(t)
			p, err = processArguments(context.Background(), []*types.Expr{}, nil)

			tBool bool
			tStr  string
			tUint uint64
			tVars *expr.Vars
		)

		req.NoError(err)
		req.NotNil(p)

		req.False(p.bool("not-there", &tBool))
		req.False(tBool)

		req.False(p.string("not-there", &tStr))
		req.Empty(tStr)

		req.False(p.uint64("not-there", &tUint))
		req.Zero(tUint)

		req.False(p.vars("not-there", tVars))
		req.Zero(tVars)
	})

	t.Run("expr", func(t *testing.T) {
		var (
			req = require.New(t)
			ee  = []*types.Expr{
				types.NewTypedExpr("tBool", "true", &expr.Boolean{}),
				types.NewTypedExpr("tString", "\"foo\"", &expr.String{}),
				types.NewTypedExpr("tUint", "42", &expr.UnsignedInteger{}),
				types.NewTypedExpr("tVars", "{}", &expr.Vars{}),
				types.NewTypedExpr("tVars.foo", "\"bar\"", &expr.String{}),
			}
		)

		req.NoError(expr.NewGvalParser().ParseEvaluators(func() []expr.Evaluator {
			oo := make([]expr.Evaluator, len(ee))
			for i, e := range ee {
				oo[i] = e
			}
			return oo
		}()...))

		var (
			p, err = processArguments(context.Background(), ee, nil)

			tBool bool
			tStr  string
			tUint uint64
			tVars = &expr.Vars{}
		)

		req.NoError(err)
		req.NotNil(p)

		req.True(p.bool("tBool", &tBool))
		req.True(tBool)

		req.True(p.string("tString", &tStr))
		req.Equal("foo", tStr)

		req.True(p.uint64("tUint", &tUint))
		req.Equal(uint64(42), tUint)

		req.True(p.vars("tVars", tVars))
		req.Equal(expr.Must(expr.NewVars(map[string]any{"foo": "bar"})), tVars)
	})
}

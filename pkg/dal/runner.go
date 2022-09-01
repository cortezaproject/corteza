package dal

import (
	"context"
)

type (
	tester interface {
		// Test returns a boolean output based on the expression, mostly used for filters
		// @todo remove the error and rely on validator so make sure everything is valid
		Test(ctx context.Context, params any) (bool, error)
	}

	evaluator interface {
		// Eval returns some value based on the expression, mostly used for attribute eval
		// @todo remove the error and rely on validator so make sure eve
		Eval(ctx context.Context, params any) (any, error)
	}
)

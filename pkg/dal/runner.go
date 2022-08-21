package dal

import (
	"context"
)

type (
	tester interface {
		Test(ctx context.Context, params any) bool
	}

	evaluator interface {
		Eval(ctx context.Context, params any) any
	}
)

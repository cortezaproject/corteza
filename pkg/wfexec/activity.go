package wfexec

import (
	"context"
)

type (
	activity struct {
		parameters *Expressions
		results    *Expressions
		fn         ActivityHandler
	}

	ActivityHandler func(ctx context.Context, in Variables) (Variables, error)
)

// Create new activity
func Activity(fn ActivityHandler, aa, rr *Expressions) *activity {
	return &activity{
		parameters: aa,
		results:    rr,
		fn:         fn,
	}
}

// Exec executes activity by running current scope through parameters Expressions, execute wrapped function and
// collect results with result Expressions.
func (a *activity) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	var (
		params, results Variables
		err             error
	)

	if a.parameters != nil {
		params, err = a.parameters.Eval(ctx, r.Scope)
		if err != nil {
			return nil, err
		}
	}

	results, err = a.fn(ctx, params)
	if err != nil {
		return nil, err
	}

	if a.results != nil {
		results, err = a.results.Eval(ctx, results)
		if err != nil {
			return nil, err
		}

		return results, nil
	}

	return Variables{}, nil
}

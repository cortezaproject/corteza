package types

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"time"
)

type (
	delayStep struct {
		wfexec.StepIdentifier
		args  ExprSet
		until time.Time
		now   func() time.Time
	}
)

// DelayStep creates a step that logs entire contents of the scope
func DelayStep(args ExprSet) *delayStep {
	return &delayStep{
		now:  func() time.Time { return time.Now() },
		args: args,
	}
}

// Executes delay step
func (s *delayStep) Exec(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
	if s.until.IsZero() {
		// wait time not yet calculated
		s.until = s.now()

		var (
			result, err = s.args.Eval(ctx, r.Scope)
		)

		if err != nil {
			return nil, err
		}

		const (
			absArgName = "timestamp"
			relArgName = "offset"
		)

		switch {
		case result.Has(absArgName):
			abs, _ := expr.NewDateTime(expr.Must(result.Select(absArgName)))
			s.until = abs.GetValue().Add(0)
		case result.Has(relArgName):
			rel, _ := expr.NewDuration(expr.Must(result.Select(relArgName)))
			s.until = s.until.Add(rel.GetValue())
		}
	}

	if s.now().After(s.until) {
		// Resume if we've delayed enough
		return wfexec.Resume(), nil
	}

	return wfexec.Delay(s.until), nil
}

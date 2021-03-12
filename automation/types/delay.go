package types

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"time"
)

type (
	delayStep struct {
		wfexec.StepIdentifier
		args ExprSet
		now  func() time.Time
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
func (s delayStep) Exec(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
	// session will set "resumed" on input when step is resumed
	if !r.Input.Has("resumed") {
		// not yet resumed
		var until = s.now()

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
			return wfexec.Delay(abs.GetValue().Add(0)), nil

		case result.Has(relArgName):
			rel, _ := expr.NewDuration(expr.Must(result.Select(relArgName)))
			return wfexec.Delay(until.Add(rel.GetValue())), nil
		}

		return nil, errors.InvalidData("failed to execute delay step")
	}

	return wfexec.Resume(), nil
}

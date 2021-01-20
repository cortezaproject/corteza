package types

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/davecgh/go-spew/spew"
)

type (
	promptStep struct {
		*expressionsStep
	}
)

func PromptStep(kind string, wse *expressionsStep) *promptStep {
	return &promptStep{wse}
}

// Executes prompt step
func (p *promptStep) Exec(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
	testResults, err := p.Set.Validate(ctx, r.Scope.Merge(r.Input))
	if err != nil {
		return nil, err
	}

	if len(testResults) > 0 {
		// @todo extend waitforinput to accept test results that are passable back to caller/client?
		spew.Dump(testResults)
		return wfexec.WaitForInput(), nil
	}

	return p.Exec(ctx, r)
}

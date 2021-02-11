package types

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	errorStep struct {
		wfexec.StepIdentifier
		message string
	}
)

func ErrorStep(msg string) *errorStep {
	return &errorStep{message: msg}
}

// Executes prompt step
func (e *errorStep) Exec(context.Context, *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
	return nil, fmt.Errorf(e.message)
}

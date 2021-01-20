package types

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	errorHandlerStep struct {
		identifiableStep
		handler wfexec.Step
	}
)

func ErrorHandlerStep(h wfexec.Step) *errorHandlerStep {
	return &errorHandlerStep{handler: h}
}

// Executes prompt step
func (h *errorHandlerStep) Exec(_ context.Context, _ *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
	return wfexec.ErrorHandler(h.handler), nil
}

package wfexec

import "context"

type (
	Steps []Step
	Step  interface {
		ID() uint64
		SetID(uint64)
		Exec(context.Context, *ExecRequest) (ExecResponse, error)
	}

	StepIdentifier struct{ id uint64 }

	execFn func(context.Context, *ExecRequest) (ExecResponse, error)

	genericStep struct {
		StepIdentifier
		fn execFn
	}
)

func (i *StepIdentifier) ID() uint64      { return i.id }
func (i *StepIdentifier) SetID(id uint64) { i.id = id }

func NewGenericStep(fn execFn) *genericStep {
	return &genericStep{fn: fn}
}

func (g *genericStep) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	return g.fn(ctx, r)
}

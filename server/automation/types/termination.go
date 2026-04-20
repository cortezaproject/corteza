package types

import (
	"context"

	nx "github.com/crusttech/human/server/pkg/automation_exec/types"
)

type (
	terminationStep struct{}
)

// TerminationStep initializes a new termination step
// This step signals that the current execution branch should terminate
func TerminationStep() (*terminationStep, error) {
	return &terminationStep{}, nil
}

// ExecN executes the termination step
func (t *terminationStep) ExecN(ctx context.Context, r *nx.ExecRequest) (nx.ExecResponse, error) {
	// No-op: actual termination logic is handled by the scheduler and runtime
	return nil, nil
}

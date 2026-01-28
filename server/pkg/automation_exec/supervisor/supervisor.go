package supervisor

import (
	"context"
	"sync/atomic"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

type (
	Ledger interface {
		ListExecutions(ctx context.Context) ([]*types.Execution, error)
		GetExecution(ctx context.Context, executableID, executionID id.ID, rev int) (*types.Execution, error)
	}

	RuntimeManager interface {
		Stop(execID id.ID) error
		SetDraining(draining bool)
		IsDraining() bool
	}

	supervisor struct {
		ledger  Ledger
		manager RuntimeManager

		draining atomic.Bool
	}
)

// The Supervisor is the control and inspection plane.
//
// * Reads execution state from the Ledger
// * Issues operational commands (stop, drain)
// * Exposes admin / ops views
// * Never touches execution logic
func Supervisor(ledger Ledger, manager RuntimeManager) *supervisor {
	return &supervisor{
		ledger:  ledger,
		manager: manager,
	}
}

// read side

func (s *supervisor) ListExecutions(ctx context.Context) ([]*types.Execution, error) {
	return s.ledger.ListExecutions(ctx)
}

func (s *supervisor) GetExecution(ctx context.Context, executableID id.ID, executionID id.ID, rev int) (*types.Execution, error) {
	return s.ledger.GetExecution(ctx, executableID, executionID, rev)
}

// control side

func (s *supervisor) Drain() {
	s.draining.Store(true)
	s.manager.SetDraining(true)
}

func (s *supervisor) Undrain() {
	s.draining.Store(false)
	s.manager.SetDraining(false)
}

func (s *supervisor) IsDraining() bool {
	return s.draining.Load()
}

func (s *supervisor) StopExecution(execID id.ID) error {
	return s.manager.Stop(execID)
}

package governor

import (
	"testing"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
)

// TL;DR: Registering the same execution twice returns ErrExecutionExists.
// Example: a bug causes a workflow to be started twice with the same ID before the first run finishes.
func TestAddExecution_DuplicateReturnsError(t *testing.T) {
	g, _ := newTestGovernor(t)
	execID := nextID()

	if err := g.AddExecution(execID, 1, types.Budget{}, types.RateLimit{}); err != nil {
		t.Fatalf("first AddExecution: %v", err)
	}
	if err := g.AddExecution(execID, 1, types.Budget{}, types.RateLimit{}); err != ErrExecutionExists {
		t.Errorf("expected ErrExecutionExists, got %v", err)
	}
}

// TL;DR: An execution whose single-step cost exceeds the global budget cap is rejected immediately.
// Example: a step that costs 10 ops is registered against a system with a global MaxOps=5 cap.
func TestAddExecution_StepTooExpensive_Global(t *testing.T) {
	g, _ := newTestGovernor(t)
	execID := nextID()

	g.SetGlobalBudget(5, 0)

	err := g.AddExecution(execID, 10, types.Budget{}, types.RateLimit{})
	if err != ErrStepTooExpensive {
		t.Errorf("expected ErrStepTooExpensive, got %v", err)
	}
}

// TL;DR: An execution whose single-step cost exceeds its own execution budget cap is rejected immediately.
// Example: a step costing 10 ops is registered with an execution-level MaxOps=5 limit.
func TestAddExecution_StepTooExpensive_Exec(t *testing.T) {
	g, _ := newTestGovernor(t)
	execID := nextID()

	err := g.AddExecution(execID, 10, types.Budget{MaxOps: 5}, types.RateLimit{})
	if err != ErrExecutableTooExpensive {
		t.Errorf("expected ErrExecutableTooExpensive, got %v", err)
	}
}

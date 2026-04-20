package governor

import (
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
)

// TL;DR: When no limits are set, a Request is granted immediately (returns a pre-closed channel).
// Example: a development/test environment with no rate or budget caps configured.
func TestRequest_GrantedImmediately_NoLimits(t *testing.T) {
	g, _ := newTestGovernor(t)
	execID := nextID()

	if err := g.AddExecution(execID, 10, types.Budget{}, types.RateLimit{}); err != nil {
		t.Fatalf("AddExecution: %v", err)
	}

	ch, err := g.Request(execID, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isClosed(ch) {
		t.Fatal("expected closedChan when no limits configured")
	}
}

// TL;DR: A granted request deducts the requested ops from both the exec and global budget counters.
// Example: a step claiming 7 ops reduces the remaining headroom for both that execution and the whole system.
func TestRequest_GrantedImmediately_ConsumesUsed(t *testing.T) {
	g, _ := newTestGovernor(t)
	execID := nextID()

	g.SetGlobalBudget(100, 0)
	g.SetGlobalRate(100, time.Minute)

	if err := g.AddExecution(execID, 10, types.Budget{MaxOps: 50}, types.RateLimit{MaxOps: 50, Window: time.Minute}); err != nil {
		t.Fatalf("AddExecution: %v", err)
	}

	ch, err := g.Request(execID, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isClosed(ch) {
		t.Fatal("expected closedChan on first granted request")
	}

	g.mux.Lock()
	ep := g.exec[execID]
	globalBudgetUsed := g.global.budget.used
	execBudgetUsed := ep.budget.used
	g.mux.Unlock()

	if globalBudgetUsed != 7 {
		t.Errorf("global budget used = %d, want 7", globalBudgetUsed)
	}
	if execBudgetUsed != 7 {
		t.Errorf("exec budget used = %d, want 7", execBudgetUsed)
	}
}

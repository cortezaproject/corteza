package governor

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
)

// TL;DR: When an execution exhausts its per-execution budget, further requests block until the window resets.
// Example: a workflow with a MaxOps=10 budget per minute hits the limit mid-run and waits for the next window.
func TestRequest_BlockedByExecBudget(t *testing.T) {
	g, nowPtr := newTestGovernor(t)
	execID := nextID()

	window := time.Minute
	if err := g.AddExecution(execID, 5, types.Budget{MaxOps: 10, Window: window}, types.RateLimit{}); err != nil {
		t.Fatalf("AddExecution: %v", err)
	}

	// Consume all budget
	if _, err := g.Request(execID, 10); err != nil {
		t.Fatalf("first request: %v", err)
	}

	ch, err := g.Request(execID, 1)
	if err != nil {
		t.Fatalf("second request: %v", err)
	}
	if isClosed(ch) {
		t.Fatal("expected blocking gate channel, got closedChan")
	}

	// Advance time past window reset — gate should open
	*nowPtr = (*nowPtr).Add(window + time.Second)
	chAfter, err := g.Request(execID, 1)
	if err != nil {
		t.Fatalf("request after window reset: %v", err)
	}
	if !isClosed(chAfter) {
		t.Fatal("expected grant after window reset")
	}
	if !isClosed(ch) {
		t.Fatal("expected original gate to be opened after window reset")
	}
}

// TL;DR: When the system-wide budget is exhausted, all executions block regardless of their own limits.
// Example: a shared API quota is spent by one automation; others must wait until the global window resets.
func TestRequest_BlockedByGlobalBudget(t *testing.T) {
	g, nowPtr := newTestGovernor(t)
	execID := nextID()

	window := time.Minute
	g.SetGlobalBudget(10, window)

	if err := g.AddExecution(execID, 5, types.Budget{}, types.RateLimit{}); err != nil {
		t.Fatalf("AddExecution: %v", err)
	}

	if _, err := g.Request(execID, 10); err != nil {
		t.Fatalf("first request: %v", err)
	}

	ch, err := g.Request(execID, 1)
	if err != nil {
		t.Fatalf("second request: %v", err)
	}
	if isClosed(ch) {
		t.Fatal("expected blocking gate channel for global budget")
	}

	*nowPtr = (*nowPtr).Add(window + time.Second)
	chAfter, err := g.Request(execID, 1)
	if err != nil {
		t.Fatalf("request after window reset: %v", err)
	}
	if !isClosed(chAfter) {
		t.Fatal("expected grant after global budget window reset")
	}
	if !isClosed(ch) {
		t.Fatal("expected original global gate to be opened")
	}
}

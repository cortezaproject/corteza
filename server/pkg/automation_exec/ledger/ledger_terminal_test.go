package ledger

import (
	"errors"
	"testing"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/id"
)


func setupExecution(t *testing.T) (*ledger, id.ID, id.ID) {
	t.Helper()
	l := newLedger()
	executableID, execID := nextID(), nextID()
	if err := l.RegisterExecution(ctx, executableID, execID, 1); err != nil {
		t.Fatal(err)
	}
	return l, executableID, execID
}

// TL;DR: ExecutionCompleted marks the execution as completed, sets EndedAt, and clears Error.
// Example: all steps in a workflow finish successfully — the run is sealed with a completion timestamp.
func TestExecutionCompleted_SetsTerminalState(t *testing.T) {
	l, xID, eID := setupExecution(t)
	if err := l.ExecutionCompleted(ctx, xID, eID, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ex, _ := l.GetExecution(ctx, xID, eID, 1)
	if ex.Status != types.StatusCompleted {
		t.Errorf("want %q, got %q", types.StatusCompleted, ex.Status)
	}
	if ex.EndedAt == nil {
		t.Error("EndedAt should be set")
	}
	if ex.Error != nil {
		t.Errorf("Error should be nil, got %v", ex.Error)
	}
}

// TL;DR: ExecutionFailed marks the execution as failed, sets EndedAt, and stores the error.
// Example: a step returns a fatal error — the run is sealed with the cause attached for diagnostics.
func TestExecutionFailed_SetsTerminalState(t *testing.T) {
	l, xID, eID := setupExecution(t)
	sentinel := errors.New("something went wrong")
	if err := l.ExecutionFailed(ctx, xID, eID, 1, sentinel); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ex, _ := l.GetExecution(ctx, xID, eID, 1)
	if ex.Status != types.StatusFailed {
		t.Errorf("want %q, got %q", types.StatusFailed, ex.Status)
	}
	if ex.EndedAt == nil {
		t.Error("EndedAt should be set")
	}
	if !errors.Is(ex.Error, sentinel) {
		t.Errorf("want sentinel error, got %v", ex.Error)
	}
}

// TL;DR: ExecutionCompleted and ExecutionFailed return errors for unknown executions.
// Example: a race condition causes a completion event to arrive after the execution entry was evicted.
func TestExecutionTerminal_UnknownExecution_Errors(t *testing.T) {
	t.Run("completed", func(t *testing.T) {
		l := newLedger()
		if err := l.ExecutionCompleted(ctx, nextID(), nextID(), 1); err == nil {
			t.Error("expected error for unknown execution")
		}
	})
	t.Run("failed", func(t *testing.T) {
		l := newLedger()
		if err := l.ExecutionFailed(ctx, nextID(), nextID(), 1, errors.New("x")); err == nil {
			t.Error("expected error for unknown execution")
		}
	})
}

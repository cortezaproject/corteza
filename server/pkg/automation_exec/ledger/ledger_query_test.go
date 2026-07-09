package ledger

import (
	"errors"
	"testing"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
)

// TL;DR: IsExecutableInUse returns true while the execution is non-terminal, false once it finishes.
// Example: the runtime manager checks in-use before hot-reloading an updated automation definition.
func TestIsExecutableInUse_StatusTransitions(t *testing.T) {
	t.Run("true for non-terminal status", func(t *testing.T) {
		l, xID, _ := setupExecution(t)
		inUse, err := l.IsExecutableInUse(ctx, xID, 1)
		if err != nil {
			t.Fatal(err)
		}
		if !inUse {
			t.Error("expected in-use for created status")
		}
	})

	t.Run("false after completed", func(t *testing.T) {
		l, xID, eID := setupExecution(t)
		_ = l.ExecutionCompleted(ctx, xID, eID, 1)
		inUse, err := l.IsExecutableInUse(ctx, xID, 1)
		if err != nil {
			t.Fatal(err)
		}
		if inUse {
			t.Error("expected not in-use after completed")
		}
	})

	t.Run("false after failed", func(t *testing.T) {
		l, xID, eID := setupExecution(t)
		_ = l.ExecutionFailed(ctx, xID, eID, 1, errors.New("x"))
		inUse, _ := l.IsExecutableInUse(ctx, xID, 1)
		if inUse {
			t.Error("expected not in-use after failed")
		}
	})
}

// TL;DR: IsExecutableInUse returns false (not error) for an unknown executable.
// Example: querying in-use status for an automation that has never been run.
func TestIsExecutableInUse_UnknownExecutable_ReturnsFalse(t *testing.T) {
	l := newLedger()
	inUse, err := l.IsExecutableInUse(ctx, nextID(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if inUse {
		t.Error("expected false for unknown executable")
	}
}

// TL;DR: IsExecutableInUse is scoped to revision — completing rev2 does not affect rev1's in-use state.
// Example: v2 of an automation finishes while v1 is still running; v1 must remain flagged as in-use.
func TestIsExecutableInUse_ScopedToRevision(t *testing.T) {
	l := newLedger()
	xID := nextID()
	execID1, execID2 := nextID(), nextID()

	_ = l.RegisterExecution(ctx, xID, execID1, 1, types.ExecutionParams{})
	_ = l.RegisterExecution(ctx, xID, execID2, 2, types.ExecutionParams{})
	_ = l.ExecutionCompleted(ctx, xID, execID2, 2)

	inUse, _ := l.IsExecutableInUse(ctx, xID, 1)
	if !inUse {
		t.Error("rev1 should still be in-use")
	}
	inUse2, _ := l.IsExecutableInUse(ctx, xID, 2)
	if inUse2 {
		t.Error("rev2 should not be in-use")
	}
}

// TL;DR: ListExecutionsByExecutable returns all executions for a given executable, across all revisions.
// Example: an admin dashboard lists every run of a particular automation regardless of version.
func TestListExecutionsByExecutable_ReturnsAll(t *testing.T) {
	l := newLedger()
	xID := nextID()
	eID1, eID2 := nextID(), nextID()

	_ = l.RegisterExecution(ctx, xID, eID1, 1, types.ExecutionParams{})
	_ = l.RegisterExecution(ctx, xID, eID2, 2, types.ExecutionParams{})

	execs, err := l.ListExecutionsByExecutable(ctx, xID)
	if err != nil {
		t.Fatal(err)
	}
	if len(execs) != 2 {
		t.Errorf("expected 2 executions, got %d", len(execs))
	}
}

// TL;DR: ListExecutionsByExecutable returns an empty slice (not an error) for an unknown executable.
// Example: fetching the run history of an automation that has never been triggered.
func TestListExecutionsByExecutable_UnknownExecutable_ReturnsEmpty(t *testing.T) {
	l := newLedger()
	execs, err := l.ListExecutionsByExecutable(ctx, nextID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(execs) != 0 {
		t.Errorf("expected empty slice, got %d", len(execs))
	}
}

// TL;DR: ListExecutionsByExecutable never returns executions belonging to a different executable.
// Example: two automations run concurrently — fetching one's history must not include the other's runs.
func TestListExecutionsByExecutable_Isolated(t *testing.T) {
	l := newLedger()
	xID1, xID2 := nextID(), nextID()
	_ = l.RegisterExecution(ctx, xID1, nextID(), 1, types.ExecutionParams{})
	_ = l.RegisterExecution(ctx, xID2, nextID(), 1, types.ExecutionParams{})

	execs, _ := l.ListExecutionsByExecutable(ctx, xID1)
	if len(execs) != 1 {
		t.Errorf("expected 1 execution for xID1, got %d", len(execs))
	}
	for _, ex := range execs {
		if ex.ExecutableID != xID1 {
			t.Errorf("got execution belonging to wrong executable: %v", ex.ExecutableID)
		}
	}
}

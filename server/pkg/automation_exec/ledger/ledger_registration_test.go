package ledger

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
)

// TL;DR: A registered execution is stored with StatusCreated and an initialised events slice.
// Example: a new automation run is kicked off — the ledger records it as pending before any step executes.
func TestRegisterExecution_CreatesEntry(t *testing.T) {
	l := newLedger()
	execID := nextID()
	executableID := nextID()

	if err := l.RegisterExecution(ctx, executableID, execID, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ex, err := l.GetExecution(ctx, executableID, execID, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex.Status != types.StatusCreated {
		t.Errorf("want status %q, got %q", types.StatusCreated, ex.Status)
	}
	if ex.Events == nil {
		t.Error("events slice should be initialised, got nil")
	}
	if ex.ExecutableID != executableID {
		t.Errorf("want executableID %v, got %v", executableID, ex.ExecutableID)
	}
	if ex.ID != execID {
		t.Errorf("want executionID %v, got %v", execID, ex.ID)
	}
	if ex.Revision != 1 {
		t.Errorf("want revision 1, got %d", ex.Revision)
	}
}

// TL;DR: Multiple revisions of the same executable are kept as independent entries.
// Example: an automation is updated while a previous run is still tracked — both records coexist.
func TestRegisterExecution_MultipleRevisions_StoredIndependently(t *testing.T) {
	l := newLedger()
	executableID := nextID()
	execID := nextID()

	if err := l.RegisterExecution(ctx, executableID, execID, 1); err != nil {
		t.Fatal(err)
	}
	if err := l.RegisterExecution(ctx, executableID, execID, 2); err != nil {
		t.Fatal(err)
	}

	ex1, err := l.GetExecution(ctx, executableID, execID, 1)
	if err != nil {
		t.Fatal(err)
	}
	ex2, err := l.GetExecution(ctx, executableID, execID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if ex1 == ex2 {
		t.Error("revisions should be stored as separate entries")
	}
	if ex1.Revision != 1 || ex2.Revision != 2 {
		t.Errorf("unexpected revisions: %d, %d", ex1.Revision, ex2.Revision)
	}
}

// TL;DR: GetExecution returns an error for unknown executables, executions, or revisions.
// Example: querying the ledger for an execution that never existed or was cleaned up.
func TestGetExecution_Errors(t *testing.T) {
	t.Run("unknown executable", func(t *testing.T) {
		l := newLedger()
		_, err := l.GetExecution(ctx, nextID(), nextID(), 1)
		if err == nil {
			t.Error("expected error for unknown executable")
		}
	})

	t.Run("unknown execution", func(t *testing.T) {
		l := newLedger()
		executableID := nextID()
		if err := l.RegisterExecution(ctx, executableID, nextID(), 1); err != nil {
			t.Fatal(err)
		}
		_, err := l.GetExecution(ctx, executableID, nextID(), 1)
		if err == nil {
			t.Error("expected error for unknown execution")
		}
	})

	t.Run("unknown revision", func(t *testing.T) {
		l := newLedger()
		executableID := nextID()
		execID := nextID()
		if err := l.RegisterExecution(ctx, executableID, execID, 1); err != nil {
			t.Fatal(err)
		}
		_, err := l.GetExecution(ctx, executableID, execID, 99)
		if err == nil {
			t.Error("expected error for unknown revision")
		}
	})
}

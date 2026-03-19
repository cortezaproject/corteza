package registry

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

func setupOne(t *testing.T) (*registry, id.ID) {
	t.Helper()
	r := newRegistry(nil)
	execID := nextID()
	if err := r.Add(ctx, makeExec(execID, 1)); err != nil {
		t.Fatal(err)
	}
	return r, execID
}

// TL;DR: Deprecating a revision sets StatusDeprecated and stamps DeprecatedAt.
// Example: a new automation version is published — the old revision is deprecated to prevent new runs.
func TestDeprecate_SetsDeprecatedStatus(t *testing.T) {
	r, execID := setupOne(t)

	if err := r.Deprecate(ctx, execID, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e, _ := r.Get(ctx, execID, 1)
	if e.Status != types.StatusDeprecated {
		t.Errorf("want StatusDeprecated, got %v", e.Status)
	}
	if e.DeprecatedAt == nil {
		t.Error("DeprecatedAt should be set")
	}
}

// TL;DR: Calling Deprecate twice on the same revision is a no-op.
// Example: a sync process tries to deprecate a revision that is already deprecated — should not error.
func TestDeprecate_Idempotent(t *testing.T) {
	r, execID := setupOne(t)

	_ = r.Deprecate(ctx, execID, 1)
	if err := r.Deprecate(ctx, execID, 1); err != nil {
		t.Errorf("second Deprecate should be idempotent, got: %v", err)
	}
}

// TL;DR: Deprecating an unknown executable or revision returns the appropriate sentinel error.
// Example: a stale deployment removes an automation that was already cleaned up — error identifies it clearly.
func TestDeprecate_Errors(t *testing.T) {
	t.Run("unknown executable", func(t *testing.T) {
		r := newRegistry(nil)
		if err := r.Deprecate(ctx, nextID(), 1); err != ErrExecutableNotFound {
			t.Errorf("want ErrExecutableNotFound, got %v", err)
		}
	})

	t.Run("unknown revision", func(t *testing.T) {
		r, execID := setupOne(t)
		if err := r.Deprecate(ctx, execID, 99); err != ErrRevisionNotFound {
			t.Errorf("want ErrRevisionNotFound, got %v", err)
		}
	})
}

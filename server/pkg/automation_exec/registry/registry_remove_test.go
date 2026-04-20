package registry

import (
	"testing"

	"github.com/crusttech/human/server/pkg/id"
)

func setupWithChecker(t *testing.T, checker usageChecker) (*registry, id.ID) {
	t.Helper()
	r := newRegistry(checker)
	execID := nextID()
	if err := r.Add(ctx, makeExec(execID, 1)); err != nil {
		t.Fatal(err)
	}
	return r, execID
}

// TL;DR: A deprecated, not-in-use revision is fully removed and the executable map entry is cleaned up.
// Example: the last revision of an old automation is pruned during a housekeeping run.
func TestRemove_DeprecatedNotInUse_Deleted(t *testing.T) {
	r, execID := setupWithChecker(t, &mockUsageChecker{inUse: false})

	_ = r.Deprecate(ctx, execID, 1)
	if err := r.Remove(ctx, execID, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := r.Get(ctx, execID, 1); err != ErrExecutableNotFound && err != ErrRevisionNotFound {
		t.Errorf("expected not-found error after remove, got %v", err)
	}

	if s := r.Stats(ctx); s.TotalExecutables != 0 {
		t.Errorf("want 0 executables after full remove, got %d", s.TotalExecutables)
	}
}

// TL;DR: Remove requires the revision to be deprecated first, otherwise returns ErrExecutableNotDeprecated.
// Example: an operator tries to force-remove a live automation — the registry rejects it to prevent data loss.
func TestRemove_NotDeprecated_Errors(t *testing.T) {
	r, execID := setupWithChecker(t, &mockUsageChecker{inUse: false})
	if err := r.Remove(ctx, execID, 1); err != ErrExecutableNotDeprecated {
		t.Errorf("want ErrExecutableNotDeprecated, got %v", err)
	}
}

// TL;DR: Remove blocks if the execution is still in use, returning ErrExecutableInUse.
// Example: a housekeeping job tries to remove a deprecated automation that still has an active run.
func TestRemove_StillInUse_Errors(t *testing.T) {
	r, execID := setupWithChecker(t, &mockUsageChecker{inUse: true})
	_ = r.Deprecate(ctx, execID, 1)
	if err := r.Remove(ctx, execID, 1); err != ErrExecutableInUse {
		t.Errorf("want ErrExecutableInUse, got %v", err)
	}
}

// TL;DR: Remove requires a usageChecker to be configured; without one it returns ErrUsageCheckerRequired.
// Example: a stripped-down registry instance without ledger integration tries to remove an entry.
func TestRemove_NoUsageChecker_Errors(t *testing.T) {
	r, execID := setupWithChecker(t, nil)
	_ = r.Deprecate(ctx, execID, 1)
	if err := r.Remove(ctx, execID, 1); err != ErrUsageCheckerRequired {
		t.Errorf("want ErrUsageCheckerRequired, got %v", err)
	}
}

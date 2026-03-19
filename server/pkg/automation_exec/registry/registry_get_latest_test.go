package registry

import (
	"testing"
)

// TL;DR: GetLatest returns the highest active revision when multiple are registered.
// Example: a workflow dispatcher looks up the current version of an automation to start a new run.
func TestGetLatest_ReturnsHighestActiveRevision(t *testing.T) {
	r := newRegistry(nil)
	execID := nextID()

	_ = r.Add(ctx, makeExec(execID, 1))
	_ = r.Add(ctx, makeExec(execID, 3))
	_ = r.Add(ctx, makeExec(execID, 2))

	e, err := r.GetLatest(ctx, execID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Executable.Revision != 3 {
		t.Errorf("want revision 3, got %d", e.Executable.Revision)
	}
}

// TL;DR: GetLatest skips deprecated revisions and returns the highest still-active one.
// Example: v2 was rolled back (deprecated) — GetLatest falls back to v1 for new runs.
func TestGetLatest_SkipsDeprecated(t *testing.T) {
	r := newRegistry(nil)
	execID := nextID()

	_ = r.Add(ctx, makeExec(execID, 1))
	_ = r.Add(ctx, makeExec(execID, 2))
	_ = r.Deprecate(ctx, execID, 2)

	e, err := r.GetLatest(ctx, execID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Executable.Revision != 1 {
		t.Errorf("want revision 1 (non-deprecated), got %d", e.Executable.Revision)
	}
}

// TL;DR: GetLatest returns ErrNoActiveRevisions when all revisions are deprecated.
// Example: every version of an automation is deprecated before a new one is published — no run can start.
func TestGetLatest_AllDeprecated_Errors(t *testing.T) {
	r := newRegistry(nil)
	execID := nextID()

	_ = r.Add(ctx, makeExec(execID, 1))
	_ = r.Deprecate(ctx, execID, 1)

	if _, err := r.GetLatest(ctx, execID); err != ErrNoActiveRevisions {
		t.Errorf("want ErrNoActiveRevisions, got %v", err)
	}
}

// TL;DR: GetLatest returns ErrExecutableNotFound for an unknown executable.
// Example: a trigger fires for an automation ID that was never registered.
func TestGetLatest_UnknownExecutable_Errors(t *testing.T) {
	r := newRegistry(nil)
	if _, err := r.GetLatest(ctx, nextID()); err != ErrExecutableNotFound {
		t.Errorf("want ErrExecutableNotFound, got %v", err)
	}
}

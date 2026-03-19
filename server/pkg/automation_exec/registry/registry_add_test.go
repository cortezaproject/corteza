package registry

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
)

// TL;DR: Adding a new executable sets StatusActive and records RegisteredAt.
// Example: a new automation is deployed — the registry marks it as live and ready to run.
func TestAdd_SetsActiveStatus(t *testing.T) {
	r := newRegistry(nil)
	execID := nextID()

	if err := r.Add(ctx, makeExec(execID, 1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e, err := r.Get(ctx, execID, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Status != types.StatusActive {
		t.Errorf("want StatusActive, got %v", e.Status)
	}
	if e.RegisteredAt.IsZero() {
		t.Error("RegisteredAt should be set")
	}
}

// TL;DR: Adding the same revision again overwrites it without error.
// Example: a hot-reload pushes the same revision twice due to a retry — only one entry exists.
func TestAdd_SameRevision_Overwrites(t *testing.T) {
	r := newRegistry(nil)
	execID := nextID()

	_ = r.Add(ctx, makeExec(execID, 1))
	if err := r.Add(ctx, makeExec(execID, 1)); err != nil {
		t.Fatalf("second Add should not error: %v", err)
	}

	e, err := r.Get(ctx, execID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if e.Status != types.StatusActive {
		t.Errorf("replaced entry should still be active, got %v", e.Status)
	}
}

// TL;DR: TotalExecutables only increments once per unique executable, not per revision.
// Example: publishing v1 and v2 of the same automation counts as one executable, two revisions.
func TestAdd_Stats_ExecutablesCountedOnce(t *testing.T) {
	r := newRegistry(nil)
	execID := nextID()

	_ = r.Add(ctx, makeExec(execID, 1))
	if s := r.Stats(ctx); s.TotalExecutables != 1 {
		t.Errorf("want 1 executable, got %d", s.TotalExecutables)
	}

	_ = r.Add(ctx, makeExec(execID, 2))
	if s := r.Stats(ctx); s.TotalExecutables != 1 {
		t.Errorf("Executables should still be 1 after second revision, got %d", s.TotalExecutables)
	}
}

// TL;DR: TotalRevisions increments for each distinct revision registered.
// Example: deploying v1, v2, and v3 of an automation results in TotalRevisions=3.
func TestAdd_Stats_RevisionsIncrementPerRegistration(t *testing.T) {
	r := newRegistry(nil)
	execID := nextID()

	_ = r.Add(ctx, makeExec(execID, 1))
	_ = r.Add(ctx, makeExec(execID, 2))
	if s := r.Stats(ctx); s.TotalRevisions != 2 {
		t.Errorf("want 2 revisions, got %d", s.TotalRevisions)
	}
}

// TL;DR: Re-adding the same revision does not double-count TotalRevisions.
// Example: an idempotent sync re-registers an already-known revision — counters stay correct.
func TestAdd_Stats_ReplacingRevision_NoExtraCount(t *testing.T) {
	r := newRegistry(nil)
	execID := nextID()

	_ = r.Add(ctx, makeExec(execID, 1))
	_ = r.Add(ctx, makeExec(execID, 1))
	if s := r.Stats(ctx); s.TotalRevisions != 1 {
		t.Errorf("replacing revision should not increment TotalRevisions, got %d", s.TotalRevisions)
	}
}

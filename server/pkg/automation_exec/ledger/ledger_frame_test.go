package ledger

import (
	"testing"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
)

// TL;DR: RecordFrame appends a stack frame to the execution trace.
// Example: the runtime records every scheduler tick so operators can replay the execution path.
func TestRecordFrame_AppendsFrame(t *testing.T) {
	l, xID, eID := setupExecution(t)
	frame := types.StackFrame{ID: nextID(), StepID: nextID()}
	if err := l.RecordFrame(ctx, xID, eID, 1, frame); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ex, _ := l.GetExecution(ctx, xID, eID, 1)
	if len(ex.Trace) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(ex.Trace))
	}
}

// TL;DR: Once a parent has MaxIteratorFrames child frames, the next one triggers truncation keeping the count stable.
// Example: a long-running iterator accumulates thousands of frames — the ledger trims old sibling frames to cap memory.
func TestRecordFrame_SiblingTruncation_AtMaxIteratorFrames(t *testing.T) {
	l, xID, eID := setupExecution(t)
	parentID := nextID()

	for i := 0; i < types.MaxIteratorFrames; i++ {
		f := types.StackFrame{ID: nextID(), StepID: nextID(), ParentID: parentID}
		if err := l.RecordFrame(ctx, xID, eID, 1, f); err != nil {
			t.Fatalf("unexpected error at i=%d: %v", i, err)
		}
	}
	ex, _ := l.GetExecution(ctx, xID, eID, 1)
	if len(ex.Trace) != types.MaxIteratorFrames {
		t.Errorf("expected %d frames before truncation, got %d", types.MaxIteratorFrames, len(ex.Trace))
	}

	// One more → truncation fires, count stays stable
	f := types.StackFrame{ID: nextID(), StepID: nextID(), ParentID: parentID}
	if err := l.RecordFrame(ctx, xID, eID, 1, f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ex, _ = l.GetExecution(ctx, xID, eID, 1)
	if len(ex.Trace) != types.MaxIteratorFrames {
		t.Errorf("after truncation expected %d frames, got %d", types.MaxIteratorFrames, len(ex.Trace))
	}
}

// TL;DR: Root frames (zero ParentID) are never subject to sibling truncation.
// Example: top-level sequential steps accumulate freely — only iterator children are capped.
func TestRecordFrame_RootFrames_NeverTruncated(t *testing.T) {
	l, xID, eID := setupExecution(t)

	for i := 0; i < types.MaxIteratorFrames+5; i++ {
		f := types.StackFrame{ID: nextID(), StepID: nextID()} // zero ParentID
		if err := l.RecordFrame(ctx, xID, eID, 1, f); err != nil {
			t.Fatalf("unexpected error at i=%d: %v", i, err)
		}
	}
	ex, _ := l.GetExecution(ctx, xID, eID, 1)
	if len(ex.Trace) != types.MaxIteratorFrames+5 {
		t.Errorf("root frames should not be truncated: expected %d, got %d",
			types.MaxIteratorFrames+5, len(ex.Trace))
	}
}

// TL;DR: RecordFrame errors when the execution doesn't exist.
// Example: a frame is recorded after the execution entry was already purged from the ledger.
func TestRecordFrame_UnknownExecution_Errors(t *testing.T) {
	l := newLedger()
	if err := l.RecordFrame(ctx, nextID(), nextID(), 1, types.StackFrame{}); err == nil {
		t.Error("expected error for unknown execution")
	}
}

// TL;DR: GetTrace hands back a copy, so a caller working on the result cannot rewrite the stored trace.
// Example: a consumer filtering frames in place (out := frames[:0]) would otherwise compact the
// ledger's own slice, and the next read of the same execution would return something else.
func TestGetTrace_ReturnsCopy(t *testing.T) {
	l, xID, eID := setupExecution(t)
	// The dropped frame sits in the middle on purpose: compacting in place only
	// overwrites what follows it, so a trailing one would leave the array intact
	// and the test would pass against the aliasing bug it is here to catch.
	for _, kind := range []string{"trigger", "termination", "function"} {
		f := types.StackFrame{ID: nextID(), StepID: nextID(), Kind: kind}
		if err := l.RecordFrame(ctx, xID, eID, 1, f); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	first, err := l.GetTrace(ctx, xID, eID, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Exactly what an in-place filter does to the returned slice.
	compacted := first[:0]
	for _, f := range first {
		if f.Kind != "termination" {
			compacted = append(compacted, f)
		}
	}

	second, err := l.GetTrace(ctx, xID, eID, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(second) != 3 {
		t.Fatalf("expected the stored trace to keep 3 frames, got %d", len(second))
	}
	for i, kind := range []string{"trigger", "termination", "function"} {
		if second[i].Kind != kind {
			t.Errorf("frame %d: expected kind %q, got %q", i, kind, second[i].Kind)
		}
	}
}

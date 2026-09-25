package ledger

import (
	"context"
	"errors"
	"testing"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
	"go.uber.org/zap"
)

// The sink hears every terminal transition exactly once, with a snapshot that
// carries the outcome but not the trace, and never hears a pause.
func TestTerminalSink(t *testing.T) {
	var got []types.Execution
	l := Ledger(zap.NewNop(), WithTerminalSink(func(_ context.Context, ex types.Execution) {
		got = append(got, ex)
	}))

	xID, eID, pausedID, failedID := nextID(), nextID(), nextID(), nextID()
	if err := l.RegisterExecution(ctx, xID, eID, 1, types.ExecutionParams{EventType: "onManual"}); err != nil {
		t.Fatal(err)
	}
	if err := l.RegisterExecution(ctx, xID, pausedID, 1, types.ExecutionParams{}); err != nil {
		t.Fatal(err)
	}
	if err := l.RegisterExecution(ctx, xID, failedID, 1, types.ExecutionParams{}); err != nil {
		t.Fatal(err)
	}

	if err := l.RecordFrame(ctx, xID, eID, 1, types.StackFrame{}); err != nil {
		t.Fatal(err)
	}
	if err := l.ExecutionCompleted(ctx, xID, eID, 1); err != nil {
		t.Fatal(err)
	}
	if err := l.ExecutionPaused(ctx, xID, pausedID, nextID(), 0, 1, nil); err != nil {
		t.Fatal(err)
	}
	if err := l.ExecutionFailed(ctx, xID, failedID, 1, errors.New("boom")); err != nil {
		t.Fatal(err)
	}

	// a second report of the same end is not a second run
	if err := l.ExecutionCompleted(ctx, xID, eID, 1); err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("want 2 terminal snapshots, got %d", len(got))
	}

	completed, failed := got[0], got[1]
	if completed.ID != eID || completed.Status != types.StatusCompleted || completed.EndedAt == nil {
		t.Errorf("completed snapshot wrong: %+v", completed)
	}
	if completed.EventType != "onManual" || completed.ExecutableID != xID {
		t.Errorf("snapshot lost its params: %+v", completed)
	}
	if completed.Trace != nil {
		t.Errorf("snapshot should not carry the trace")
	}
	if failed.ID != failedID || failed.Status != types.StatusFailed || failed.Error == nil || failed.Error.Error() != "boom" {
		t.Errorf("failed snapshot wrong: %+v", failed)
	}
}

// Without a sink the ledger behaves exactly as before.
func TestTerminalSink_Absent(t *testing.T) {
	l, xID, eID := setupExecution(t)
	if err := l.ExecutionCompleted(ctx, xID, eID, 1); err != nil {
		t.Fatal(err)
	}
}

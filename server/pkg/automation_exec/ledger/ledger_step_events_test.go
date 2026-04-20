package ledger

import (
	"errors"
	"testing"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
)

// TL;DR: Step lifecycle events are appended in order with the correct type.
// Example: a step starts, finishes, then fails retry — all three events appear in the execution trace.
func TestStepEvents_AppendedInOrder(t *testing.T) {
	l, xID, eID := setupExecution(t)
	stepID := nextID()

	_ = l.StepStarted(ctx, xID, eID, stepID, 1)
	_ = l.StepCompleted(ctx, xID, eID, stepID, 1, "result")
	_ = l.StepFailed(ctx, xID, eID, stepID, 1, errors.New("fail"))

	ex, _ := l.GetExecution(ctx, xID, eID, 1)
	if len(ex.Events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(ex.Events))
	}
	if ex.Events[0].Type != types.EventStepStarted {
		t.Errorf("event[0] type: want %q, got %q", types.EventStepStarted, ex.Events[0].Type)
	}
	if ex.Events[1].Type != types.EventStepCompleted {
		t.Errorf("event[1] type: want %q, got %q", types.EventStepCompleted, ex.Events[1].Type)
	}
	if ex.Events[2].Type != types.EventStepFailed {
		t.Errorf("event[2] type: want %q, got %q", types.EventStepFailed, ex.Events[2].Type)
	}
}

// TL;DR: StepCompleted carries an arbitrary payload that can be retrieved from the event.
// Example: a function step returns the ID of a newly created record — stored in the event for downstream use.
func TestStepCompleted_CarriesPayload(t *testing.T) {
	l, xID, eID := setupExecution(t)
	stepID := nextID()
	_ = l.StepCompleted(ctx, xID, eID, stepID, 1, "my-payload")
	ex, _ := l.GetExecution(ctx, xID, eID, 1)
	if ex.Events[0].Payload != "my-payload" {
		t.Errorf("unexpected payload: %v", ex.Events[0].Payload)
	}
}

// TL;DR: StepFailed stores the originating error on the event for later inspection.
// Example: an HTTP call fails — the error is attached to the event so operators can diagnose it.
func TestStepFailed_CarriesError(t *testing.T) {
	l, xID, eID := setupExecution(t)
	stepID := nextID()
	sentinel := errors.New("step-err")
	_ = l.StepFailed(ctx, xID, eID, stepID, 1, sentinel)
	ex, _ := l.GetExecution(ctx, xID, eID, 1)
	if !errors.Is(ex.Events[0].Error, sentinel) {
		t.Errorf("unexpected error: %v", ex.Events[0].Error)
	}
}

// TL;DR: UpdatedAt advances (never goes backwards) after each step event.
// Example: a UI poll relies on UpdatedAt to detect new activity on a running execution.
func TestStepEvents_UpdatedAtAdvances(t *testing.T) {
	l, xID, eID := setupExecution(t)
	stepID := nextID()

	ex, _ := l.GetExecution(ctx, xID, eID, 1)
	_ = l.StepStarted(ctx, xID, eID, stepID, 1)
	_ = l.StepCompleted(ctx, xID, eID, stepID, 1, nil)
	ex2, _ := l.GetExecution(ctx, xID, eID, 1)
	if ex2.UpdatedAt.Before(ex.UpdatedAt) {
		t.Error("UpdatedAt went backwards")
	}
}

// TL;DR: StepStarted, StepCompleted, and StepFailed all error on an unknown execution.
// Example: a step completion event is routed to the wrong execution ID due to a mapping bug.
func TestStepEvents_UnknownExecution_Errors(t *testing.T) {
	t.Run("started", func(t *testing.T) {
		l := newLedger()
		if err := l.StepStarted(ctx, nextID(), nextID(), nextID(), 1); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("completed", func(t *testing.T) {
		l := newLedger()
		if err := l.StepCompleted(ctx, nextID(), nextID(), nextID(), 1, nil); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("failed", func(t *testing.T) {
		l := newLedger()
		if err := l.StepFailed(ctx, nextID(), nextID(), nextID(), 1, errors.New("x")); err == nil {
			t.Error("expected error")
		}
	})
}

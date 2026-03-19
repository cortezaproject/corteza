package manager

import (
	"context"
	"errors"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

// TL;DR: A successful execution calls ExecutionCompleted on the ledger and RemoveExecution on the governor.
// Example: all steps finish without error — the run is sealed, resources are released, and Done is closed.
func TestOnExit_Success(t *testing.T) {
	exec := singleStepExec()
	led := &mockLedger{}
	gov := &mockGovernor{}
	rm := newRM(t, &mockRegistry{exec: exec}, led, gov, defaultConfig())

	eid, err := rm.Start(context.Background(), exec.ID, 1, types.ExecutionParams{})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	e, err := rm.Get(eid)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	waitDone(t, e.Done)

	led.mu.Lock()
	completed, failed := led.completedCalled, led.failedCalled
	led.mu.Unlock()
	gov.mu.Lock()
	removed := gov.removeCalled
	gov.mu.Unlock()

	if completed < 1 {
		t.Errorf("want ExecutionCompleted called ≥1, got %d", completed)
	}
	if failed != 0 {
		t.Errorf("want ExecutionFailed not called, got %d", failed)
	}
	if removed < 1 {
		t.Errorf("want RemoveExecution called ≥1, got %d", removed)
	}
}

// TL;DR: A failed execution calls ExecutionFailed on the ledger; the governor always releases the slot.
// Example: a step returns a fatal error — the run is recorded as failed and its governor slot freed.
func TestOnExit_Failure(t *testing.T) {
	stepErr := errors.New("step boom")
	exec := types.Executable{
		ID:       id.MustNumID(id.Next()),
		Revision: 1,
		Steps: []types.Step{
			{ID: id.MustNumID(id.Next()), Handle: "boom", Kind: "action", Handler: &alwaysFailHandler{err: stepErr}},
		},
	}
	led := &mockLedger{}
	gov := &mockGovernor{}
	rm := newRM(t, &mockRegistry{exec: exec}, led, gov, defaultConfig())

	eid, err := rm.Start(context.Background(), exec.ID, 1, types.ExecutionParams{})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	e, err := rm.Get(eid)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	waitDone(t, e.Done)

	led.mu.Lock()
	failedCalled := led.failedCalled
	led.mu.Unlock()
	gov.mu.Lock()
	removed := gov.removeCalled
	gov.mu.Unlock()

	if failedCalled < 1 {
		t.Errorf("want ExecutionFailed called ≥1, got %d", failedCalled)
	}
	if removed < 1 {
		t.Errorf("want RemoveExecution called even on failure, got %d", removed)
	}
}

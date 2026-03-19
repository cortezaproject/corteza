package manager

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

// TL;DR: Start returns ErrSystemDraining when the manager is shutting down.
// Example: a trigger fires while the server is gracefully stopping — new automations are refused.
func TestStart_Draining_Errors(t *testing.T) {
	rm := newRM(t, &mockRegistry{exec: singleStepExec()}, &mockLedger{}, &mockGovernor{}, defaultConfig())
	rm.SetDraining(true)

	_, err := rm.Start(context.Background(), id.MustNumID(id.Next()), 1, types.ExecutionParams{})
	if !errors.Is(err, ErrSystemDraining) {
		t.Errorf("want ErrSystemDraining, got %v", err)
	}
}

// TL;DR: Start returns ErrQueueFull when MaxQueued pending executions are already waiting.
// Example: a burst of triggers arrives during a slow run — the queue cap protects against runaway memory growth.
func TestStart_QueueFull_Errors(t *testing.T) {
	release := make(chan struct{})
	exec := types.Executable{
		ID:       id.MustNumID(id.Next()),
		Revision: 1,
		Steps: []types.Step{
			{ID: id.MustNumID(id.Next()), Handle: "block", Kind: "action", Handler: &blockHandler{release: release}},
		},
	}
	rm := newRM(t, &mockRegistry{exec: exec}, &mockLedger{}, &mockGovernor{}, Config{
		MaxConcurrent: 1,
		MaxQueued:     1,
		SlotTimeout:   5 * time.Second,
	})
	ctx := context.Background()

	eid1, err := rm.Start(ctx, exec.ID, 1, types.ExecutionParams{})
	if err != nil {
		t.Fatalf("first Start: %v", err)
	}
	if e1, _ := rm.Get(eid1); e1 != nil {
		select {
		case <-e1.Started:
		case <-time.After(3 * time.Second):
			t.Fatal("first execution never started")
		}
	}

	if _, err = rm.Start(ctx, exec.ID, 1, types.ExecutionParams{}); err != nil {
		t.Fatalf("second Start: %v", err)
	}

	_, err = rm.Start(ctx, exec.ID, 1, types.ExecutionParams{})
	if !errors.Is(err, ErrQueueFull) {
		t.Errorf("want ErrQueueFull, got %v", err)
	}
	close(release)
}

// TL;DR: Start returns a non-zero execution ID and synchronously calls RegisterExecution on the ledger.
// Example: a scheduler kicks off a workflow — the ledger immediately records the new run before any step fires.
func TestStart_RegistersExecution(t *testing.T) {
	exec := singleStepExec()
	led := &mockLedger{}
	rm := newRM(t, &mockRegistry{exec: exec}, led, &mockGovernor{}, defaultConfig())

	eid, err := rm.Start(context.Background(), exec.ID, 1, types.ExecutionParams{})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if eid == id.Zero() {
		t.Error("returned execution ID should not be zero")
	}

	led.mu.Lock()
	rc := led.registerCalled
	led.mu.Unlock()
	if rc < 1 {
		t.Errorf("want RegisterExecution called at least once, got %d", rc)
	}
}

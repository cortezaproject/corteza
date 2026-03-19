package manager

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

// TL;DR: With MaxConcurrent=1, a second Start queues and only runs after the first finishes.
// Example: a limited deployment allows only one automation run at a time — concurrent triggers wait their turn.
func TestSlotEnforcement_SecondStartWaitsForFirst(t *testing.T) {
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
		SlotTimeout:   5 * time.Second,
	})
	ctx := context.Background()

	eid1, err := rm.Start(ctx, exec.ID, 1, types.ExecutionParams{})
	if err != nil {
		t.Fatalf("first Start: %v", err)
	}

	e1, err := rm.Get(eid1)
	if err != nil {
		t.Fatalf("Get first: %v", err)
	}

	select {
	case <-e1.Started:
	case <-time.After(3 * time.Second):
		t.Fatal("first execution never started")
	}

	eid2, err := rm.Start(ctx, exec.ID, 1, types.ExecutionParams{})
	if err != nil {
		t.Fatalf("second Start: %v", err)
	}

	// Second must not have started while first holds the slot
	if e2, _ := rm.Get(eid2); e2 != nil {
		select {
		case <-e2.Started:
			t.Error("second execution started while first still holds the slot")
		default:
		}
	}

	close(release)
	waitDone(t, e1.Done)

	// Second should eventually start and complete
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if e2, err := rm.Get(eid2); err == nil {
			waitDone(t, e2.Done)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

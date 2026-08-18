package automation_exec

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/id"
)

// A ledger that answers GetTrace and nothing else; the embedded nil interface
// satisfies the rest of ledgerAPI, which this test never reaches.
type traceOnlyLedger struct {
	ledgerAPI
	trace []types.StackFrame
}

func (l traceOnlyLedger) GetTrace(context.Context, id.ID, id.ID, int) ([]types.StackFrame, error) {
	return l.trace, nil
}

// TL;DR: the execution trace carries every recorded frame, terminations included.
// Example: a branch arm ending in End is only distinguishable from the arm that was
// not taken by the presence of that arm's termination frame — the builder canvas
// colours the taken path from it.
func TestGetExecutionTrace_KeepsTerminationFrames(t *testing.T) {
	svc := &automationService{led: traceOnlyLedger{trace: []types.StackFrame{
		{Handle: "trigger_1", Kind: "trigger"},
		{Handle: "step_2", Kind: "gatewayExclusive"},
		{Handle: "step_3", Kind: "function"},
		{Handle: "step_5", Kind: "termination"},
	}}}

	out, err := svc.GetExecutionTrace(context.Background(), id.ID{}, 0, id.ID{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) != 4 {
		t.Fatalf("expected all 4 frames, got %d", len(out))
	}

	var terminations int
	for _, f := range out {
		if f.Kind == "termination" {
			terminations++
		}
	}
	if terminations != 1 {
		t.Errorf("expected the termination frame to survive, got %d", terminations)
	}
}

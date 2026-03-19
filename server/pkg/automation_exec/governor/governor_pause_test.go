package governor

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
)

// TL;DR: While paused, all requests block; ResumeAll opens the gate for all waiting requests.
// Example: an operator pauses all automations during a maintenance window then resumes them afterwards.
func TestPause_BlocksRequests_ResumeOpensGate(t *testing.T) {
	g, _ := newTestGovernor(t)
	execID := nextID()

	if err := g.AddExecution(execID, 10, types.Budget{}, types.RateLimit{}); err != nil {
		t.Fatalf("AddExecution: %v", err)
	}

	g.PauseAll()

	ch, err := g.Request(execID, 1)
	if err != nil {
		t.Fatalf("Request while paused: %v", err)
	}
	if isClosed(ch) {
		t.Fatal("expected blocking gate while paused")
	}

	g.ResumeAll()
	if !isClosed(ch) {
		t.Fatal("expected gate to be opened after ResumeAll")
	}
}

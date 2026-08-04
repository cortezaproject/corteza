package agentic

import (
	"strings"
	"testing"

	"github.com/crusttech/human/server/system/agentic/toolkit"

	autoTypes "github.com/crusttech/human/server/automation/types"
)

// TestRefuseMultipleEntryPoints covers the one rule this handler enforces that
// the service does not. The server checks it when a session starts, so a
// multi-entry workflow stores clean and fails on its first run; everything here
// is about catching it before the write.
func TestRefuseMultipleEntryPoints(t *testing.T) {
	step := func(id uint64, kind autoTypes.WorkflowStepKind) *autoTypes.WorkflowStep {
		return &autoTypes.WorkflowStep{ID: id, Kind: kind}
	}
	path := func(parent, child uint64) *autoTypes.WorkflowPath {
		return &autoTypes.WorkflowPath{ParentID: parent, ChildID: child}
	}

	tests := []struct {
		name    string
		steps   autoTypes.WorkflowStepSet
		paths   autoTypes.WorkflowPathSet
		wantErr bool
	}{
		{
			name: "empty workflow",
		},
		{
			name:  "single step, no paths",
			steps: autoTypes.WorkflowStepSet{step(1, autoTypes.WorkflowStepKindExpressions)},
		},
		{
			name:  "chain has one entry",
			steps: autoTypes.WorkflowStepSet{step(1, autoTypes.WorkflowStepKindExpressions), step(2, autoTypes.WorkflowStepKindTermination)},
			paths: autoTypes.WorkflowPathSet{path(1, 2)},
		},
		{
			name: "branches rejoin, still one entry",
			steps: autoTypes.WorkflowStepSet{
				step(1, autoTypes.WorkflowStepKindGateway),
				step(2, autoTypes.WorkflowStepKindExpressions),
				step(3, autoTypes.WorkflowStepKindExpressions),
				step(4, autoTypes.WorkflowStepKindTermination),
			},
			paths: autoTypes.WorkflowPathSet{path(1, 2), path(1, 3), path(2, 4), path(3, 4)},
		},
		{
			name: "two disconnected steps",
			steps: autoTypes.WorkflowStepSet{
				step(1, autoTypes.WorkflowStepKindExpressions),
				step(2, autoTypes.WorkflowStepKindExpressions),
			},
			wantErr: true,
		},
		{
			// A visual step is dropped by the converter before the graph is
			// built, so it is never an entry point however it is connected.
			name: "visual step is not an entry point",
			steps: autoTypes.WorkflowStepSet{
				step(1, autoTypes.WorkflowStepKindExpressions),
				step(2, autoTypes.WorkflowStepKindVisual),
			},
		},
		{
			// A cycle leaves no entry at all. The server rejects that at run
			// time too ("could not find starting step"), but a workflow being
			// built up over several updates passes through states like this,
			// so it is not refused here.
			name: "no entry point at all is allowed through",
			steps: autoTypes.WorkflowStepSet{
				step(1, autoTypes.WorkflowStepKindExpressions),
				step(2, autoTypes.WorkflowStepKindExpressions),
			},
			paths: autoTypes.WorkflowPathSet{path(1, 2), path(2, 1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := refuseMultipleEntryPoints(tt.steps, tt.paths)
			if (err != nil) != tt.wantErr {
				t.Fatalf("refuseMultipleEntryPoints() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "stepIDs") {
				t.Errorf("error does not name the offending steps: %v", err)
			}
		})
	}
}

// TestWorkflowGraphJSON covers both wire forms a client may use for 'steps' and
// 'paths', and the ,string ID tags that make a bare number fail.
//
// The decoding lives in toolkit.JSONArg, shared with the TAQ tools; what is
// worth testing here is that it lands correctly in the workflow types, since the
// ,string tags are what turns a bare number into an error rather than a
// truncated ID.
func TestWorkflowGraphJSON(t *testing.T) {
	steps := func(raw any) (autoTypes.WorkflowStepSet, error) {
		var out autoTypes.WorkflowStepSet
		_, err := toolkit.JSONArg(map[string]any{"steps": raw}, "steps", "steps", &out)
		return out, err
	}

	t.Run("JSON string", func(t *testing.T) {
		got, err := steps(`[{"stepID":"7","kind":"expressions"}]`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 || got[0].ID != 7 {
			t.Fatalf("got %#v, want one step with ID 7", got)
		}
	})

	t.Run("decoded array", func(t *testing.T) {
		var got autoTypes.WorkflowPathSet
		_, err := toolkit.JSONArg(
			map[string]any{"paths": []any{map[string]any{"parentID": "1", "childID": "2"}}},
			"paths", "paths", &got,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 || got[0].ParentID != 1 || got[0].ChildID != 2 {
			t.Fatalf("got %#v, want one path 1 -> 2", got)
		}
	})

	t.Run("absent and empty are not a change", func(t *testing.T) {
		for _, raw := range []any{nil, "", "   "} {
			got, err := steps(raw)
			if err != nil {
				t.Fatalf("unexpected error for %#v: %v", raw, err)
			}
			if got != nil {
				t.Fatalf("got %#v for %#v, want nil so the service reads it as unchanged", got, raw)
			}
		}
	})

	t.Run("explicit empty array is a wipe, not absence", func(t *testing.T) {
		got, err := steps("[]")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatal("got nil, want an empty non-nil set so the service replaces rather than preserves")
		}
	})

	t.Run("unquoted stepID is refused with a message that says why", func(t *testing.T) {
		_, err := steps(`[{"stepID":7,"kind":"expressions"}]`)
		if err == nil {
			t.Fatal("expected an error for a bare numeric stepID")
		}
		if !strings.Contains(err.Error(), "quoted string") {
			t.Errorf("error does not tell the caller to quote the ID: %v", err)
		}
	})
}

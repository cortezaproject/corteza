package envoy

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/automation/types"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
)

func sTestAutomationWorkflow(ctx context.Context, t *testing.T, s store.Storer, pfx string) *types.Workflow {
	wf := &types.Workflow{
		ID:     su.NextID(),
		Handle: pfx + "_handle",
		Meta: &types.WorkflowMeta{
			Name:        pfx + "_name",
			Description: pfx + "_description",
		},
		Enabled:      true,
		Trace:        true,
		KeepSessions: 10,
		Steps: types.WorkflowStepSet{
			&types.WorkflowStep{
				ID:   11,
				Kind: "function",
			},
			&types.WorkflowStep{
				ID:   12,
				Kind: "function",
			},
		},
		Paths: types.WorkflowPathSet{
			&types.WorkflowPath{
				ParentID: 11,
				ChildID:  12,
				Expr:     "qwerty",
			},
		},

		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}

	err := store.CreateAutomationWorkflow(ctx, s, wf)
	if err != nil {
		t.Fatal(err)
	}

	return wf
}

func sTestAutomationTrigger(ctx context.Context, t *testing.T, s store.Storer, wfID uint64, pfx string) *types.Trigger {
	wf := &types.Trigger{
		ID: su.NextID(),

		Enabled:    true,
		WorkflowID: wfID,
		StepID:     11,

		ResourceType: "testko:test:",

		Constraints: types.TriggerConstraintSet{
			&types.TriggerConstraint{
				Name: "qwerty",
				Op:   "=",
				Values: []string{
					"a",
					"b",
				},
			},
		},

		Meta: &types.TriggerMeta{
			Description: pfx + "_description",
		},

		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}

	err := store.CreateAutomationTrigger(ctx, s, wf)
	if err != nil {
		t.Fatal(err)
	}

	return wf
}

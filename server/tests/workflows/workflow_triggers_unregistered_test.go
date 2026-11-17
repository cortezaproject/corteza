package workflows

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/service"
	"github.com/cortezaproject/corteza/server/automation/types"
	cmpEvent "github.com/cortezaproject/corteza/server/compose/service/event"
	cmpTypes "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/stretchr/testify/require"
)

// A workflow stops firing once it is disabled, deleted or saved with issues;
// its record triggers then leave saves on the module alone.
func Test_workflow_triggers_unregistered(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())

		brokenStep = &types.WorkflowStep{
			ID:        2,
			Kind:      types.WorkflowStepKindExpressions,
			Arguments: []*types.Expr{{Target: "b", Expr: "(((", Type: "Integer"}},
		}

		// a before-create workflow on the module with the given handle
		// that refuses the save, so a save that passes ran nothing
		setup = func(t *testing.T, handle string) (*types.Workflow, func() error) {
			req := require.New(t)

			wf, err := service.DefaultWorkflow.Create(ctx, &types.Workflow{
				Handle:  handle,
				Enabled: true,
				Meta:    &types.WorkflowMeta{Name: handle},
				Steps: types.WorkflowStepSet{{
					ID:        1,
					Kind:      types.WorkflowStepKindError,
					Arguments: []*types.Expr{{Target: "message", Expr: `"fired"`, Type: "String"}},
				}},
			})
			req.NoError(err)

			_, err = service.DefaultTrigger.Create(ctx, &types.Trigger{
				WorkflowID:   wf.ID,
				Enabled:      true,
				ResourceType: "compose:record",
				EventType:    "beforeCreate",
				StepID:       1,
				Constraints:  types.TriggerConstraintSet{{Name: "module", Op: "=", Values: []string{handle}}},
			})
			req.NoError(err)

			save := func() error {
				ns := &cmpTypes.Namespace{ID: id.Next(), Slug: "crm"}
				m := &cmpTypes.Module{ID: id.Next(), NamespaceID: ns.ID, Handle: handle}
				r := &cmpTypes.Record{ModuleID: m.ID, NamespaceID: ns.ID}
				return eventBus.WaitFor(ctx, cmpEvent.RecordBeforeCreate(r, nil, m, ns, nil, nil))
			}

			req.ErrorContains(save(), "fired")
			return wf, save
		}
	)

	cleanup(t)

	t.Run("disabled with an issue", func(t *testing.T) {
		req := require.New(t)
		wf, save := setup(t, "disabled_with_issue")

		wf.Enabled = false
		wf.Steps = append(wf.Steps, brokenStep)
		wf, err := service.DefaultWorkflow.Update(ctx, wf)
		req.NoError(err)
		req.NotEmpty(wf.Issues)

		req.NoError(save())
	})

	t.Run("enabled with an issue", func(t *testing.T) {
		req := require.New(t)
		wf, save := setup(t, "enabled_with_issue")

		wf.Steps = append(wf.Steps, brokenStep)
		wf, err := service.DefaultWorkflow.Update(ctx, wf)
		req.NoError(err)
		req.NotEmpty(wf.Issues)

		req.NoError(save())
	})

	t.Run("deleted", func(t *testing.T) {
		req := require.New(t)
		wf, save := setup(t, "deleted")

		req.NoError(service.DefaultWorkflow.DeleteByID(ctx, wf.ID))
		req.NoError(save())
	})

	t.Run("deleted with an issue", func(t *testing.T) {
		req := require.New(t)
		wf, save := setup(t, "deleted_with_issue")

		wf.Steps = append(wf.Steps, brokenStep)
		wf, err := service.DefaultWorkflow.Update(ctx, wf)
		req.NoError(err)
		req.NotEmpty(wf.Issues)

		req.NoError(service.DefaultWorkflow.DeleteByID(ctx, wf.ID))
		req.NoError(save())
	})
}

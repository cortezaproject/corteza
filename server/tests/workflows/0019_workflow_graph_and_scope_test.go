package workflows

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
)

func Test0019_workflow_graph_and_scope(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
	)

	loadScenario(ctx, t)

	t.Run("iterator inside fork/join resolves", func(t *testing.T) {
		req := require.New(t)
		_, _, _, err := execWorkflow(ctx, "iterator_in_fork", types.WorkflowExecParams{})
		if issues, is := err.(types.WorkflowIssueSet); is {
			for _, i := range issues {
				t.Logf("issue: %s", i.Description)
			}
		}
		req.NoError(err)
	})

	t.Run("mid-branch variable survives join", func(t *testing.T) {
		req := require.New(t)
		vars, _ := mustExecWorkflow(ctx, t, "join_new_var", types.WorkflowExecParams{})
		req.Equal(expr.Must(expr.NewString("fromB")), vars.GetValue()["foo"])
	})

	t.Run("mid-branch change to pre-existing variable survives join", func(t *testing.T) {
		req := require.New(t)
		vars, _ := mustExecWorkflow(ctx, t, "join_existing_var", types.WorkflowExecParams{})
		req.Equal(expr.Must(expr.NewString("fromB")), vars.GetValue()["foo"])
	})

	t.Run("nested parallel section keeps its own changes", func(t *testing.T) {
		req := require.New(t)
		vars, _ := mustExecWorkflow(ctx, t, "join_nested", types.WorkflowExecParams{})
		req.Equal(expr.Must(expr.NewString("fromA")), vars.GetValue()["foo"])
		req.Equal(expr.Must(expr.NewString("b1")), vars.GetValue()["b1"])
		req.Equal(expr.Must(expr.NewString("b2")), vars.GetValue()["b2"])
	})

	t.Run("error handler vars survive a normal step", func(t *testing.T) {
		req := require.New(t)
		vars, _ := mustExecWorkflow(ctx, t, "errhandler_after_step", types.WorkflowExecParams{})
		req.Equal(expr.Must(expr.NewString("Boom")), vars.GetValue()["caught"])
	})
}

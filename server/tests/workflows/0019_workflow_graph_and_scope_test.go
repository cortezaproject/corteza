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

	t.Run("error handler vars survive a normal step", func(t *testing.T) {
		req := require.New(t)
		vars, _ := mustExecWorkflow(ctx, t, "errhandler_after_step", types.WorkflowExecParams{})
		req.Equal(expr.Must(expr.NewString("Boom")), vars.GetValue()["caught"])
	})
}

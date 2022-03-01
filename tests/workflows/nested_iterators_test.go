package workflows

import (
	"context"
	_ "embed"
	"testing"

	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
)

// make test.integration.workflows TEST_FLAGS="-run Test_nested_iterators -v"
func Test_nested_iterators(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)

		aux = struct {
			Token string
		}{}

		input, err = expr.NewVars(map[string]interface{}{
			"L1": 3,
			"L2": 3,
			"L3": 3,
		})
	)

	req.NoError(err)

	req.NoError(defStore.TruncateAttachments(ctx))

	loadNewScenario(ctx, t)

	vars, _ := mustExecWorkflow(ctx, t, "nested_iterators", types.WorkflowExecParams{Input: input})

	req.NoError(vars.Decode(&aux))
}

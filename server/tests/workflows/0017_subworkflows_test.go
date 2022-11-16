package workflows

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/types"
)

func Test0017_subworkflows(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	loadScenario(ctx, t)

	type (
		testInput struct {
			Out string
		}
	)
	var (
		aux      = testInput{}
		vars, _  = mustExecWorkflow(ctx, t, "main", types.WorkflowExecParams{})
		expected = testInput{Out: "main + sub"}
	)

	req.NoError(vars.Decode(&aux))
	req.Equal(expected, aux)

}

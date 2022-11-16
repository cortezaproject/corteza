package workflows

import (
	"context"
	"testing"

	autTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/stretchr/testify/require"
)

func Test0012_termination_implicit(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	loadScenario(ctx, t)

	var (
		_, trace = mustExecWorkflow(ctx, t, "testing", autTypes.WorkflowExecParams{})
	)

	req.Len(trace, 2)
	req.Equal(uint64(10), trace[0].StepID)
	req.Equal(uint64(0), trace[1].StepID)
}

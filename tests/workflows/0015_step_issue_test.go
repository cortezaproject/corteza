package workflows

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/automation/types"
)

func Test0015_step_issue(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	loadScenario(ctx, t)

	t.Run("exclusive gateway step issue", func(t *testing.T) {
		t.Skipf("workflow step resolution & validation need to be to be fixed")
		_, _, err := execWorkflow(ctx, "case1", types.WorkflowExecParams{})

		issues, is := err.(types.WorkflowIssueSet)
		req.True(is)
		// It should return only 3 issues
		// 1. gateway step expects at least 1 outbound path(s)
		// 2. expecting at least two paths for exclusive gateway
		// 3. failed to resolve workflow step dependencies
		req.Len(issues, 3)
	})

	t.Run("inclusive gateway step issue", func(t *testing.T) {
		t.Skipf("workflow step resolution & validation need to be to be fixed")
		_, _, err := execWorkflow(ctx, "case2", types.WorkflowExecParams{})

		issues, is := err.(types.WorkflowIssueSet)
		req.True(is)
		// It should return only 3 issues
		// 1. gateway step expects at least 1 outbound path(s)
		// 2. expecting at least two paths for inclusive gateway
		// 3. failed to resolve workflow step dependencies
		req.Len(issues, 3)
	})

	t.Run("function step issue", func(t *testing.T) {
		t.Skipf("workflow step resolution & validation need to be to be fixed")
		_, _, err := execWorkflow(ctx, "case3", types.WorkflowExecParams{})

		issues, is := err.(types.WorkflowIssueSet)
		req.True(is)
		// It should return only 2 issues
		// 1. failed to verify argument expressions for function logInfo: parameter message is required
		// 2. failed to resolve workflow step dependencies
		req.Len(issues, 2)
	})

	t.Run("iterator step issue", func(t *testing.T) {
		t.Skipf("workflow step resolution & validation need to be to be fixed")
		_, _, err := execWorkflow(ctx, "case4", types.WorkflowExecParams{})

		issues, is := err.(types.WorkflowIssueSet)
		req.True(is)
		// It should return only 4 issues
		// 1. iterator step expects reference
		// 2. iterator step expects exactly 2 outbound path(s)
		// 3. unknown function ""
		// 4. failed to resolve workflow step dependencies
		req.Len(issues, 4)
	})

}

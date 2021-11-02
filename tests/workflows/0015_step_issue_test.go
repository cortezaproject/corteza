package workflows

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cortezaproject/corteza-server/automation/types"
)

func Test0015_step_issue(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	loadScenario(ctx, t)

	t.Run("exclusive gateway step issue", func(t *testing.T) {
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
		_, _, err := execWorkflow(ctx, "case3", types.WorkflowExecParams{})

		issues, is := err.(types.WorkflowIssueSet)
		req.True(is)
		// It should return only 2 issues
		// 1. failed to verify argument expressions for function logInfo: parameter message is required
		// 2. failed to resolve workflow step dependencies
		req.Len(issues, 2)
	})

	t.Run("iterator step issue", func(t *testing.T) {
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

	t.Run("multiple gateway with single function step issues", func(t *testing.T) {
		_, _, err := execWorkflow(ctx, "case5", types.WorkflowExecParams{})

		issues, is := err.(types.WorkflowIssueSet)
		req.True(is)
		// It should return only 12 issues
		// 1. gateway step expects at least 1 outbound path(s)
		// 2. expecting at least two paths for inclusive gateway
		// 3. gateway step expects at least 1 outbound path(s)
		// 4. expecting at least two paths for exclusive gateway
		// 5. gateway step expects at least 1 outbound path(s)
		// 6. expecting at least two paths for inclusive gateway
		// 7. gateway step expects at least 1 outbound path(s)
		// 8. expecting at least two paths for exclusive gateway
		// 9. gateway step expects at least 1 outbound path(s)
		// 10. expecting at least two paths for inclusive gateway
		// 11. failed to verify argument expressions for function logInfo: parameter message is required
		// 12. failed to resolve workflow step dependencies
		req.Len(issues, 12)
	})

	t.Run("two gateway and a function step issues", func(t *testing.T) {
		_, _, err := execWorkflow(ctx, "case6", types.WorkflowExecParams{})

		issues, is := err.(types.WorkflowIssueSet)
		req.True(is)
		// It should return only 12 issues
		// 1. gateway step expects at least 1 outbound path(s)
		// 2. expecting at least two paths for exclusive gateway
		// 3. function step expects reference
		// 4. unknown function ""
		// 5. failed to resolve workflow step dependencies
		// 6. failed to resolve step with ID 2
		// 7. failed to resolve step with ID 3
		req.Len(issues, 7)
	})

	t.Run("error handler step issue", func(t *testing.T) {
		_, _, err := execWorkflow(ctx, "case7", types.WorkflowExecParams{})

		issues, is := err.(types.WorkflowIssueSet)
		req.True(is)
		// It should return only 3 issues
		// 1. error-handler step expects at least 1 outbound path(s)
		// 2. expecting at least one path out of error handling step
		// 3. failed to resolve workflow step dependencies
		req.Len(issues, 3)
	})

	t.Run("multiple error handler and gateway step with single function step issues", func(t *testing.T) {
		_, _, err := execWorkflow(ctx, "case8", types.WorkflowExecParams{})

		issues, is := err.(types.WorkflowIssueSet)
		req.True(is)
		// It should return only 12 issues
		// 1. gateway step expects at least 1 outbound path(s)
		// 2. expecting at least two paths for inclusive gateway
		// 3. error-handler step expects at least 1 outbound path(s)
		// 4. expecting at least one path out of error handling step
		// 5. gateway step expects at least 1 outbound path(s)
		// 6. expecting at least two paths for exclusive gateway
		// 7. error-handler step expects at least 1 outbound path(s)
		// 8. expecting at least one path out of error handling step
		// 9. failed to verify argument expressions for function logInfo: parameter message is required
		// 10. failed to resolve workflow step dependencies
		req.Len(issues, 10)
	})
}

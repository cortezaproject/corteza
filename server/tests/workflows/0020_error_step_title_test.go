package workflows

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/automation/service"
	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/stretchr/testify/require"
)

// The error raised by an error step is what the workflow author configured:
// no session prefix in front of the message, and an optional title for the UI
func Test0020_error_step_title(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	loadScenario(ctx, t)

	t.Run("message and title", func(t *testing.T) {
		_, _, _, err := execWorkflow(ctx, "error_step_with_title", types.WorkflowExecParams{})
		req.Error(err)
		req.Equal("Order is not ready", err.Error())

		var e *errors.Error
		req.True(errors.As(err, &e))
		req.True(errors.IsAutomation(e))
		req.Equal("Order check", e.Meta()[service.ErrorStepTitleMetaKey])
	})

	t.Run("message only", func(t *testing.T) {
		_, _, _, err := execWorkflow(ctx, "error_step_without_title", types.WorkflowExecParams{})
		req.Error(err)
		req.Equal("Order is not ready", err.Error())

		var e *errors.Error
		req.True(errors.As(err, &e))
		req.Nil(e.Meta()[service.ErrorStepTitleMetaKey])
	})

	t.Run("failing argument names the function and the parameter", func(t *testing.T) {
		_, _, _, err := execWorkflow(ctx, "failing_argument", types.WorkflowExecParams{})
		req.Error(err)
		req.Contains(err.Error(), `function "logInfo", arguments`)
		req.Contains(err.Error(), `could not evaluate "message"`)
		// other failures keep the session context
		req.Contains(err.Error(), "execution failed")
	})
}

package workflows

import (
	"context"
	"testing"
	"time"

	autTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/id"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

// Test_results_without_trace asserts on what a workflow returns rather than on
// its stacktrace, so it holds whether or not frame collection is turned on
// (WORKFLOW_STACK_TRACE_ENABLED)
func Test_results_without_trace(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
		aux = struct {
			Total uint64
			Found int
		}{}
	)

	req.NoError(defStore.TruncateReminders(ctx))

	rr := make([]*sysTypes.Reminder, 5)
	for i := range rr {
		rr[i] = &sysTypes.Reminder{
			ID:         id.Next(),
			Resource:   "test:resource",
			AssignedTo: auth.GetIdentityFromContext(ctx).Identity(),
			CreatedAt:  time.Now(),
		}
	}
	req.NoError(defStore.CreateReminder(ctx, rr...))

	loadNewScenario(ctx, t)

	vars, _ := mustExecWorkflow(ctx, t, "testing", autTypes.WorkflowExecParams{})

	req.NoError(vars.Decode(&aux))
	req.Equal(uint64(5), aux.Total)
	req.Equal(5, aux.Found)
}

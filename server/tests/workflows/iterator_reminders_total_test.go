package workflows

import (
	"context"
	"testing"
	"time"

	autTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

// Test_iterator_reminders_total iterates with incTotal over more reminders than
// the buffer holds; total is fetched with the first page and reported on every
// iteration
func Test_iterator_reminders_total(t *testing.T) {
	wfexec.MaxIteratorBufferSize = 2
	defer func() {
		wfexec.MaxIteratorBufferSize = wfexec.DefaultMaxIteratorBufferSize
	}()

	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	req.NoError(defStore.TruncateReminders(ctx))

	loadNewScenario(ctx, t)

	for j := 0; j < 5; j++ {
		req.NoError(defStore.CreateReminder(ctx, &sysTypes.Reminder{
			ID:         id.Next(),
			Resource:   "test:resource",
			AssignedTo: auth.GetIdentityFromContext(ctx).Identity(),
			CreatedAt:  time.Now(),
		}))
	}

	var (
		_, trace = mustExecWorkflow(ctx, t, "testing", autTypes.WorkflowExecParams{})
	)

	// 6x iterator, 5x continue, 1x terminator, 1x completed
	req.Len(trace, 13)

	for j := 0; j <= 4; j++ {
		total, err := expr.Integer{}.Cast(trace[j*2].Results.GetValue()["t"])
		req.NoError(err)
		req.Equal(int64(5), total.Get().(int64))
	}
}

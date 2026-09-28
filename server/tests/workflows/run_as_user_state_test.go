package workflows

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/automation/types"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func Test_run_as_user_state(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)

		runner = func() (string, error) {
			aux := struct{ Runner *sysTypes.User }{}
			vars, _, _, err := execWorkflow(ctx, "runner", types.WorkflowExecParams{Wait: true})
			if err != nil {
				return "", err
			}

			req.NoError(vars.Decode(&aux))
			return aux.Runner.Handle, nil
		}
	)

	req.NoError(defStore.TruncateUsers(ctx))
	loadNewScenario(ctx, t)

	u, err := defStore.LookupUserByHandle(ctx, "wf-runner")
	req.NoError(err)

	handle, err := runner()
	req.NoError(err)
	req.Equal("wf-runner", handle)

	now := time.Now()
	u.SuspendedAt = &now
	req.NoError(defStore.UpdateUser(ctx, u))

	_, err = runner()
	req.ErrorContains(err, "suspended")

	u.SuspendedAt = nil
	req.NoError(defStore.UpdateUser(ctx, u))

	handle, err = runner()
	req.NoError(err)
	req.Equal("wf-runner", handle)
}

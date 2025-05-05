package workflows

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/service"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	"github.com/stretchr/testify/require"
)

func Test_exec_permissions(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	req.NoError(defStore.TruncateUsers(ctx))
	req.NoError(defStore.TruncateRoles(ctx))
	req.NoError(defStore.TruncateRoleMembers(ctx))
	req.NoError(defStore.TruncateRbacRules(ctx))

	loadNewScenario(ctx, t)

	// user that the workflow is configured to use for run-as
	execAllowed, err := defStore.LookupUserByHandle(ctx, "exec-allowed")
	req.NoError(err)

	// user that the workflow is configured to use for run-as
	execDenied, err := defStore.LookupUserByHandle(ctx, "exec-denied")
	req.NoError(err)

	// invokers group with permissions to execute workflow
	executors, err := defStore.LookupRoleByHandle(ctx, "executors")
	req.NoError(err)

	//err = defStore.CreateRoleMember(ctx, &sysTypes.RoleMember{UserID: wfInvoker.ID, RoleID: wfInvokers.ID})
	//req.NoError(err)

	execAllowed.SetRoles(executors.ID)

	helpers.UpdateRBAC(
		executors.ID,
	)

	t.Run("exec allowed", func(t *testing.T) {
		ctx = auth.SetIdentityToContext(ctx, execAllowed)
		_, _ = mustExecWorkflow(ctx, t, "wf", types.WorkflowExecParams{})
	})

	t.Run("exec denied", func(t *testing.T) {
		req = require.New(t)
		ctx = auth.SetIdentityToContext(ctx, execDenied)
		_, _, _, err = execWorkflow(ctx, "wf", types.WorkflowExecParams{})
		req.ErrorIs(err, service.WorkflowErrNotAllowedToExecute())
	})
}

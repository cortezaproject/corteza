package workflows

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	"github.com/stretchr/testify/require"
)

func Test_invoker_and_runner_in_scope(t *testing.T) {
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
	wfRunner, err := defStore.LookupUserByHandle(ctx, "wf-runner")
	req.NoError(err)

	// user invoking the workflow
	wfInvoker, err := defStore.LookupUserByHandle(ctx, "wf-invoker")
	req.NoError(err)

	// invokers group with permissions to execute workflow
	wfInvokers, err := defStore.LookupRoleByHandle(ctx, "wf-invokers")
	req.NoError(err)
	err = defStore.CreateRoleMember(ctx, &sysTypes.RoleMember{UserID: wfInvoker.ID, RoleID: wfInvokers.ID})
	req.NoError(err)

	wfInvoker.SetRoles(wfInvokers.ID)
	ctx = auth.SetIdentityToContext(ctx, wfInvoker)

	helpers.UpdateRBAC(
		wfInvokers.ID,
	)

	t.Run("invoker set in scope", func(t *testing.T) {
		var (
			req = require.New(t)
			aux = struct {
				Invoker *sysTypes.User
				Runner  *sysTypes.User
			}{}
		)

		vars, _ := mustExecWorkflow(ctx, t, "invoker", types.WorkflowExecParams{})
		req.NoError(vars.Decode(&aux))

		// Expecting both, invoker & runner to be same as invoker
		req.NotNil(aux.Runner)
		req.NotNil(aux.Invoker)
		req.Equal(aux.Runner.Handle, wfInvoker.Handle)
		req.Equal(aux.Invoker.Handle, wfInvoker.Handle)
	})

	t.Run("runner set in scope", func(t *testing.T) {
		var (
			req = require.New(t)
			aux = struct {
				Invoker *sysTypes.User
				Runner  *sysTypes.User
			}{}
		)

		vars, _ := mustExecWorkflow(ctx, t, "runner", types.WorkflowExecParams{})
		req.NoError(vars.Decode(&aux))

		// Expecting runner and invoker to be different.
		req.NotNil(aux.Runner)
		req.NotNil(aux.Invoker)
		req.Equal(aux.Runner.Handle, wfRunner.Handle)
		req.Equal(aux.Invoker.Handle, wfInvoker.Handle)
	})
}

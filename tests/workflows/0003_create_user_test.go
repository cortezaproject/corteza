package workflows

import (
	"context"
	"testing"

	autTypes "github.com/cortezaproject/corteza-server/automation/types"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func Test0003_create_user(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	loadScenario(ctx, t)

	var (
		aux = struct {
			User *sysTypes.User
		}{}
		vars, _ = mustExecWorkflow(ctx, t, "create-user", autTypes.WorkflowExecParams{})
	)

	req.NoError(vars.Decode(&aux))
}

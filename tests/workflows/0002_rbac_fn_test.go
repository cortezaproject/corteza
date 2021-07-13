package workflows

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/stretchr/testify/require"
)

func Test0002_rbac_fn(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	loadScenario(ctx, t)

	req.Len(rbac.Global().Rules(), 0)

	var (
		aux = struct {
			CanCurrentRead string
			CanOtherRead   string
		}{}
		vars, _ = mustExecWorkflow(ctx, t, "check-and-grant", types.WorkflowExecParams{})
	)

	req.NoError(vars.Decode(&aux))
	req.Equal("y", aux.CanCurrentRead)
	req.Equal("n", aux.CanOtherRead)
	req.Len(rbac.Global().Rules(), 1)
}

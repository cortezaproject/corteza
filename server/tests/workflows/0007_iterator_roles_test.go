package workflows

import (
	"context"
	"fmt"
	"testing"

	autTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/cortezaproject/corteza/server/system/automation"
	"github.com/stretchr/testify/require"
)

func Test0007_iterator_roles(t *testing.T) {
	wfexec.MaxIteratorBufferSize = wfexec.DefaultMaxIteratorBufferSize
	defer func() {
		wfexec.MaxIteratorBufferSize = wfexec.DefaultMaxIteratorBufferSize
	}()

	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	req.NoError(defStore.TruncateRoles(ctx))

	loadScenario(ctx, t)

	var (
		_, trace = mustExecWorkflow(ctx, t, "testing", autTypes.WorkflowExecParams{})
	)

	// 6x iterator, 5x continue, 1x terminator, 1x completed
	req.Len(trace, 13)

	// there are 4 iterator calls; each on the *2 index
	ctr := int64(-1)
	for j := 0; j <= 4; j++ {
		ix := j * 2
		ctr++

		frame := trace[ix]
		req.Equal(uint64(10), frame.StepID)

		i, err := expr.Integer{}.Cast(frame.Results.GetValue()["i"])
		req.NoError(err)
		req.Equal(ctr, i.Get().(int64))

		usr, err := automation.NewRole(frame.Results.GetValue()["r"])
		req.NoError(err)
		req.Equal(fmt.Sprintf("r%d", ctr+1), usr.GetValue().Handle)
	}
}

func Test0008_iterator_roles_chunked(t *testing.T) {
	wfexec.MaxIteratorBufferSize = 2
	defer func() {
		wfexec.MaxIteratorBufferSize = wfexec.DefaultMaxIteratorBufferSize
	}()

	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	req.NoError(defStore.TruncateRoles(ctx))

	loadScenarioWithName(ctx, t, "S0007_iterator_roles")

	var (
		_, trace = mustExecWorkflow(ctx, t, "testing", autTypes.WorkflowExecParams{})
	)

	// 6x iterator, 5x continue, 1x terminator, 1x completed
	req.Len(trace, 13)

	// there are 4 iterator calls; each on the *2 index
	ctr := int64(-1)
	for j := 0; j <= 4; j++ {
		ix := j * 2
		ctr++

		frame := trace[ix]
		req.Equal(uint64(10), frame.StepID)

		i, err := expr.Integer{}.Cast(frame.Results.GetValue()["i"])
		req.NoError(err)
		req.Equal(ctr, i.Get().(int64))

		usr, err := automation.NewRole(frame.Results.GetValue()["r"])
		req.NoError(err)
		req.Equal(fmt.Sprintf("r%d", ctr+1), usr.GetValue().Handle)
	}
}

func Test0008_iterator_roles_limited(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	req.NoError(defStore.TruncateRoles(ctx))

	loadScenarioWithName(ctx, t, "iterator_roles_limit")

	var (
		_, trace = mustExecWorkflow(ctx, t, "testing", autTypes.WorkflowExecParams{})
	)

	// 3x iterator, 2x continue, 1x terminator, 1x completed
	req.Len(trace, 7)

	// there are 4 iterator calls; each on the *2 index
	ctr := int64(-1)
	for j := 0; j <= 1; j++ {
		ix := j * 2
		ctr++

		frame := trace[ix]
		req.Equal(uint64(10), frame.StepID)

		i, err := expr.Integer{}.Cast(frame.Results.GetValue()["i"])
		req.NoError(err)
		req.Equal(ctr, i.Get().(int64))

		usr, err := automation.NewRole(frame.Results.GetValue()["r"])
		req.NoError(err)
		req.Equal(fmt.Sprintf("r%d", ctr+1), usr.GetValue().Handle)
	}
}

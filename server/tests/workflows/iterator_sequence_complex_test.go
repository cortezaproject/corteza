package workflows

import (
	"context"
	"testing"

	autTypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func Test_iterator_sequence_complex(t *testing.T) {
	wfexec.MaxIteratorBufferSize = wfexec.DefaultMaxIteratorBufferSize
	defer func() {
		wfexec.MaxIteratorBufferSize = wfexec.DefaultMaxIteratorBufferSize
	}()

	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	req.NoError(truncateRecords(ctx))
	req.NoError(defStore.TruncateComposeModules(ctx))
	req.NoError(defStore.TruncateComposeNamespaces(ctx))

	loadNewScenario(ctx, t)

	var (
		_, trace = mustExecWorkflow(ctx, t, "testing", autTypes.WorkflowExecParams{})
	)

	req.Len(trace, 7)

	frame := trace[1]
	req.True(frame.Scope.GetValue()["first"].Get().(bool))
	req.False(frame.Scope.GetValue()["last"].Get().(bool))
	aux := frame.Scope.GetValue()["i"].Get()
	i, err := cast.ToIntE(aux)
	req.NoError(err)
	req.Equal(2, i)

	frame = trace[len(trace)-1]
	req.False(frame.Scope.GetValue()["first"].Get().(bool))
	req.True(frame.Scope.GetValue()["last"].Get().(bool))
	aux = frame.Scope.GetValue()["i"].Get()
	i, err = cast.ToIntE(aux)
	req.NoError(err)
	req.Equal(4, i)
}

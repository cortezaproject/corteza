package workflows

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/stretchr/testify/require"
)

func Test0001_basics(t *testing.T) {
	ctx := bypassRBAC(context.Background())
	loadScenario(ctx, t)

	var (
		req     = require.New(t)
		aux     = struct{ Foo int64 }{}
		vars, _ = mustExecWorkflow(ctx, t, "basic", types.WorkflowExecParams{})
	)

	req.NoError(vars.Decode(&aux))
	req.Equal(int64(42), aux.Foo)
}

func Test0001_basics_detect_datarace_issues(t *testing.T) {
	ctx := bypassRBAC(context.Background())
	loadScenarioWithName(ctx, t, "S0001_basics")

	var (
		wg = &sync.WaitGroup{}
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			defer wg.Done()
			var (
				aux     = struct{ Foo int64 }{}
				req     = require.New(t)
				vars, _ = mustExecWorkflow(ctx, t, "basic", types.WorkflowExecParams{})
			)
			req.NoError(vars.Decode(&aux))
			req.Equal(int64(42), aux.Foo)
		})
	}

	wg.Wait()
}

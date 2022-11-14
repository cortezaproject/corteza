package workflows

import (
	"testing"
)

func Test0010_stacktrace(t *testing.T) {
	t.Skip("@todo envoy not yet refactored")

	// var (
	// 	ctx = bypassRBAC(context.Background())
	// 	req = require.New(t)
	// )

	// req.NoError(truncateRecords(ctx))
	// req.NoError(defStore.TruncateComposeModules(ctx))
	// req.NoError(defStore.TruncateComposeNamespaces(ctx))

	// loadScenario(ctx, t)

	// for rep := 0; rep < 11; rep++ {
	// 	t.Run(fmt.Sprintf("iteration %d", rep), func(t *testing.T) {
	// 		var (
	// 			_, trace = mustExecWorkflow(ctx, t, "testing", autTypes.WorkflowExecParams{})
	// 		)

	// 		// 6x iterator, 5x continue, 1x terminator, 1x completed
	// 		req.Len(trace, 13)

	// 		steps := []uint64{
	// 			10,
	// 			11,
	// 			10,
	// 			11,
	// 			10,
	// 			11,
	// 			10,
	// 			11,
	// 			10,
	// 			11,

	// 			10,
	// 			12,
	// 			0,
	// 		}

	// 		for i := 0; i < 13; i++ {
	// 			req.Equal(steps[i], trace[i].StepID)
	// 		}
	// 	})
	// }
}

package workflows

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/types"
)

func Test0018_error_handler_step(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	loadScenario(ctx, t)

	t.Run("error handler with result", func(t *testing.T) {
		type (
			testInput struct {
				E    expr.TypedValue
				EM   expr.TypedValue
				EsID expr.TypedValue
			}
		)

		var (
			aux      = testInput{}
			vars, _  = mustExecWorkflow(ctx, t, "error_handler_results", types.WorkflowExecParams{})
			expected = testInput{
				E:    expr.Must(expr.NewAny(errors.New(errors.KindAutomation, "TestingError"))),
				EM:   expr.Must(expr.NewString("TestingError")),
				EsID: expr.Must(expr.NewInteger(2)),
			}
		)

		req.NoError(vars.Decode(&aux))

		errA := expr.Must(expr.NewAny(aux.E.Get()))
		req.Error(errors.Unwrap(errA.Get().(*errors.Error)))
		req.Equal(errors.Unwrap(errA.Get().(*errors.Error)).Error(), "TestingError")
		req.True(errors.IsKind(errA.Get().(*errors.Error), errors.KindAutomation))

		req.Equal(expected.EM, aux.EM)
		req.Equal(expected.EsID, aux.EsID)
	})
}

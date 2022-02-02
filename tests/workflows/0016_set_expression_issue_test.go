package workflows

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cortezaproject/corteza-server/automation/types"
)

func Test0016_set_expression_issue(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	loadScenario(ctx, t)

	t.Run("set expressions", func(t *testing.T) {
		var (
			aux     = struct{ Out map[string]expr.TypedValue }{}
			vars, _ = mustExecWorkflow(ctx, t, "set_expression", types.WorkflowExecParams{})

			testString = expr.Must(expr.NewString("testing string"))
			testInt    = expr.Must(expr.NewInteger(40))
			testVar    = expr.Must(expr.NewVars(map[string]interface{}{
				"testString": testString,
				"testFloat":  expr.Must(expr.NewFloat(50)),
			}))
			expected = map[string]expr.TypedValue{
				"testString": testString,
				"testInt":    testInt,
				"testVar":    testVar,
			}
		)

		req.NoError(vars.Decode(&aux))
		req.Equal(expected, aux.Out)
	})
}

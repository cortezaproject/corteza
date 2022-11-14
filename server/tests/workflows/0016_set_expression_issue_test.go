package workflows

import (
	"context"
	cmpTypes "github.com/cortezaproject/corteza-server/compose/types"
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
		type (
			testInput struct {
				Out         map[string]expr.TypedValue
				OutConstStr map[string]expr.TypedValue
				OutConstInt map[string]expr.TypedValue
				OutRv       *cmpTypes.Record
			}
		)
		var (
			aux     = testInput{}
			vars, _ = mustExecWorkflow(ctx, t, "set_expression", types.WorkflowExecParams{})

			testString = expr.Must(expr.NewString("testing string"))
			testInt    = expr.Must(expr.NewInteger(40))
			testVar    = expr.Must(expr.NewVars(map[string]interface{}{
				"testString": testString,
				"testFloat":  expr.Must(expr.NewFloat(50)),
			}))

			expectedVars = map[string]expr.TypedValue{
				"testString": testString,
				"testInt":    testInt,
				"testVar":    testVar,
			}
			expected = testInput{
				Out: expectedVars,
				OutConstStr: map[string]expr.TypedValue{
					"testConstKey": expr.Must(expr.NewString("testConstValue")),
				},
				OutConstInt: map[string]expr.TypedValue{
					"testInt": testInt,
				},

				OutRv: &cmpTypes.Record{
					Values: []*cmpTypes.RecordValue{
						{
							Name:  "testRv",
							Value: "testing string",
						},
						{
							Name:  "testFloat",
							Value: "50",
						},
					},
				},
			}
		)

		req.NoError(vars.Decode(&aux))
		req.Equal(expected, aux)
	})

	t.Run("omit expressions", func(t *testing.T) {
		type (
			testInput struct {
				Out         map[string]expr.TypedValue
				OutConstInt map[string]expr.TypedValue
				OutRv       *cmpTypes.Record
			}
		)
		var (
			aux     = testInput{}
			vars, _ = mustExecWorkflow(ctx, t, "omit_expression", types.WorkflowExecParams{})

			expected = testInput{
				Out: map[string]expr.TypedValue{},
				OutConstInt: map[string]expr.TypedValue{},
				OutRv: &cmpTypes.Record{
					Values: []*cmpTypes.RecordValue{
						{
							Name:  "testFloat",
							Value: "50",
						},
					},
				},
			}
		)

		req.NoError(vars.Decode(&aux))
		req.Equal(expected, aux)
	})
}

package service

import (
	"testing"

	"github.com/crusttech/human/server/automation/types"
	"github.com/stretchr/testify/require"
)

func TestVerifyStepScopeTargets_ShadowedBuiltin(t *testing.T) {
	tcs := []struct {
		name  string
		step  *types.WorkflowStep
		wants bool
	}{
		{
			name: "expressions arg target shadows sum",
			step: &types.WorkflowStep{
				Kind:      types.WorkflowStepKindExpressions,
				Arguments: []*types.Expr{{Target: "sum", Expr: "2"}},
			},
			wants: true,
		},
		{
			name: "expressions arg target safe name",
			step: &types.WorkflowStep{
				Kind:      types.WorkflowStepKindExpressions,
				Arguments: []*types.Expr{{Target: "total", Expr: "2"}},
			},
			wants: false,
		},
		{
			name: "result target shadows count",
			step: &types.WorkflowStep{
				Kind:    types.WorkflowStepKindIterator,
				Ref:     "composeRecordsEach",
				Results: []*types.Expr{{Target: "count", Expr: "record"}},
			},
			wants: true,
		},
		{
			name: "function/prompt arg target is a param name, not a scope var",
			step: &types.WorkflowStep{
				// prompt arg targets ("message") are parameter names — must not flag
				Kind:      types.WorkflowStepKindPrompt,
				Ref:       "notification",
				Arguments: []*types.Expr{{Target: "message", Expr: "x"}, {Target: "sort", Value: "primary"}},
			},
			wants: false,
		},
		{
			name: "pathed target uses base segment",
			step: &types.WorkflowStep{
				Kind:      types.WorkflowStepKindExpressions,
				Arguments: []*types.Expr{{Target: "min.foo", Expr: "2"}},
			},
			wants: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ii := verifyStepScopeTargets(tc.step)
			require.Equal(t, tc.wants, len(ii) > 0, "issues: %v", ii)
		})
	}
}

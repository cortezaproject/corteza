package service

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/stretchr/testify/require"
)

// A disabled construct is listed but never converted into a runnable step
func TestStepConvFunction_disabledConstruct(t *testing.T) {
	req := require.New(t)

	ConstructLibrary().AddFunctions(types.ConstructFunction{
		Ref:      "testDisabledFunction",
		Kind:     "function",
		Disabled: true,
		Meta:     &types.ConstructFunctionMeta{Short: "test"},
		Handler: func(ctx context.Context, in *expr.Vars) (*expr.Vars, error) {
			return in, nil
		},
	})

	_, err := stepConvFunction(&types.NgAutomationStep{
		Ref:  "testDisabledFunction",
		Kind: "function",
	})

	req.Error(err)
	req.Contains(err.Error(), "corredor is disabled")
}

// A TAQ using a disabled construct is stored with one issue naming why, so it
// is never registered as runnable
func TestBuildExecSteps_disabledConstructIsOneIssue(t *testing.T) {
	req := require.New(t)

	ConstructLibrary().AddFunctions(types.ConstructFunction{
		Ref:      "testDisabledFunctionForIssue",
		Kind:     "function",
		Disabled: true,
		Meta:     &types.ConstructFunctionMeta{Short: "test"},
		Handler: func(ctx context.Context, in *expr.Vars) (*expr.Vars, error) {
			return in, nil
		},
	})

	_, _, issues := buildExecSteps(nil, map[uint64]*types.NgAutomationStep{
		7: {Ref: "testDisabledFunctionForIssue", Kind: "function"},
	})

	req.Len(issues, 1)
	req.Contains(issues[0].Message, "corredor is disabled")
}

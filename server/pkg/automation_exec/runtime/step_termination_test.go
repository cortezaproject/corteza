package runtime

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/stretchr/testify/require"
)

// TL;DR: TerminateBranch returns done=true only when the last active branch is terminated.
// Example: a parallel workflow with two paths — the automation ends only after both paths finish.
func TestTerminateBranch_DoneWhenLastBranchTerminated(t *testing.T) {
	term1 := mkStep(3, "termination")
	term2 := mkStep(4, "termination")
	term1.Parents = []types.Step{{ID: stepID(1)}}
	term2.Parents = []types.Step{{ID: stepID(2)}}

	ep1 := types.Step{ID: stepID(1), Handle: "ep1", Kind: "function", Handler: noopHandler{}, Children: []types.Step{term1}}
	ep2 := types.Step{ID: stepID(2), Handle: "ep2", Kind: "function", Handler: noopHandler{}, Children: []types.Step{term2}}

	exe := mkExe(ep1, ep2, term1, term2)
	ss := newScheduler(exe)
	ctx := context.Background()
	require.Equal(t, 2, len(ss.activeBranches), "expect 2 entry-point branches")

	s1, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.NoError(t, ss.StoreOutputs(s1.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, s1.ID))

	t1, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)

	done, err := ss.TerminateBranch(t1.ID)
	require.NoError(t, err)
	require.False(t, done)

	s2, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.NoError(t, ss.StoreOutputs(s2.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, s2.ID))

	t2, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)

	done, err = ss.TerminateBranch(t2.ID)
	require.NoError(t, err)
	require.True(t, done)
}

// TL;DR: TerminateBranch removes only the frames belonging to the terminated branch, leaving others intact.
// Example: one path in a parallel workflow hits an error terminator while the other path continues normally.
func TestTerminateBranch_OnlyTargetBranchCleared(t *testing.T) {
	step3 := mkStep(3, "function")
	step4 := mkStep(4, "function")
	step3.Parents = []types.Step{{ID: stepID(1)}}
	step4.Parents = []types.Step{{ID: stepID(2)}}

	ep1 := types.Step{ID: stepID(1), Handle: "ep1", Kind: "function", Handler: noopHandler{}, Children: []types.Step{step3}}
	ep2 := types.Step{ID: stepID(2), Handle: "ep2", Kind: "function", Handler: noopHandler{}, Children: []types.Step{step4}}

	exe := mkExe(ep1, ep2, step3, step4)
	ss := newScheduler(exe)
	ctx := context.Background()

	s1, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)

	stackBefore := len(ss.stack)

	_, err = ss.TerminateBranch(s1.ID)
	require.NoError(t, err)
	require.Less(t, len(ss.stack), stackBefore)

	s2, _, _, hasNext, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext, "second branch still alive")
	require.NotEqual(t, s1.ID, s2.ID)
}

// TL;DR: FindOutput returns variables stored by StoreOutputs after a step completes.
// Example: a downstream step reads the "result" variable produced by an upstream function step.
func TestFindOutput_FromCompletedOutputs(t *testing.T) {
	step1 := mkStep(1, "function")
	step1.Handle = "myStep"

	exe := mkExe(step1)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.Equal(t, step1.ID, s.ID)

	out := map[string]expr.TypedValue{
		"result": expr.Must(expr.NewString("hello")),
	}
	require.NoError(t, ss.StoreOutputs(step1.ID, out))

	vars, err := ss.FindOutput("myStep")
	require.NoError(t, err)
	require.NotNil(t, vars)
}

// TL;DR: FindOutput can read outputs that are still on the stack (step in progress, not yet completed).
// Example: an expression step reading variables from the currently-executing parent iterator frame.
func TestFindOutput_FromStack(t *testing.T) {
	step1 := mkStep(1, "function")
	step1.Handle = "myStep"

	exe := mkExe(step1)
	ss := newScheduler(exe)
	ctx := context.Background()

	_, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)

	ss.stack[len(ss.stack)-1].outputs = map[string]expr.TypedValue{
		"x": expr.Must(expr.NewString("val")),
	}

	vars, err := ss.FindOutput("myStep")
	require.NoError(t, err)
	require.NotNil(t, vars)
}

// TL;DR: FindOutput returns ErrOutputNotFound for an unknown step handle.
// Example: a step referencing a handle that was never executed or was mistyped.
func TestFindOutput_NotFound(t *testing.T) {
	exe := mkExe(mkStep(1, "function"))
	ss := newScheduler(exe)

	_, err := ss.FindOutput("nonexistent")
	require.ErrorIs(t, err, ErrOutputNotFound)
}

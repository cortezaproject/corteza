package runtime

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/stretchr/testify/require"
)

// TL;DR: Steps execute in order and the scheduler signals completion when all are done.
// Example: a workflow with three sequential function steps (fetch → transform → save).
func TestLinearChain(t *testing.T) {
	step3 := mkStep(3, "function")
	step2 := mkStep(2, "function", step3)
	step1 := mkStep(1, "function", step2)

	step2.Parents = []types.Step{step1}
	step3.Parents = []types.Step{step2}

	exe := mkExe(step1, step2, step3)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, hasNext, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, step1.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(step1.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, step1.ID))

	s, _, _, hasNext, err = ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, step2.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(step2.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, step2.ID))

	s, _, _, hasNext, err = ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, step3.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(step3.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, step3.ID))

	// Stack empty → hasNext=false
	_, _, _, hasNext, err = ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.False(t, hasNext)
}

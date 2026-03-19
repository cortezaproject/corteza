package runtime

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/stretchr/testify/require"
)

// TL;DR: When the iterator has items, the first Next() after entry yields the body step.
// Example: a "for each record" loop that begins processing the first record.
func TestIterator_FirstTick(t *testing.T) {
	bodyStep := mkStep(2, "function")

	iterH := &mockIteratorHandler{maxIter: 2}
	iterStep := types.Step{
		ID:       stepID(1),
		Handle:   "iter",
		Kind:     "iterator",
		Handler:  iterH,
		Children: []types.Step{bodyStep},
	}
	exe := mkExe(iterStep)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, hasNext, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, iterStep.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(iterStep.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, iterStep.ID))

	s, _, _, hasNext, err = ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, bodyStep.ID, s.ID)
	require.Equal(t, 1, iterH.calls, "Next called once per iteration")
}

// TL;DR: When the collection is empty the iterator skips the body and jumps to the exit step.
// Example: a "for each record" loop where the query returns zero results → post-loop step runs immediately.
func TestIterator_EmptyCollection_JumpsToExit(t *testing.T) {
	exitStep := mkStep(3, "function")
	bodyStep := mkStep(2, "function")

	iterH := &mockIteratorHandler{maxIter: 0}
	iterStep := types.Step{
		ID:       stepID(1),
		Handle:   "iter",
		Kind:     "iterator",
		Handler:  iterH,
		Children: []types.Step{bodyStep, exitStep},
	}
	exe := mkExe(iterStep)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.Equal(t, iterStep.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(iterStep.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, iterStep.ID))

	s, _, _, hasNext, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, exitStep.ID, s.ID)
}

// TL;DR: The body step runs exactly N times then execution continues with the exit step.
// Example: a loop over 3 records processes each one then proceeds to a "send summary" step.
func TestIterator_RepeatsBodyThenExits(t *testing.T) {
	const iterations = 3
	exitStep := mkStep(3, "function")
	bodyStep := mkStep(2, "function")
	bodyStep.Parents = []types.Step{{ID: stepID(1)}}
	exitStep.Parents = []types.Step{{ID: stepID(1)}}

	iterH := &mockIteratorHandler{maxIter: iterations}
	iterStep := types.Step{
		ID:       stepID(1),
		Handle:   "iter",
		Kind:     "iterator",
		Handler:  iterH,
		Children: []types.Step{bodyStep, exitStep},
	}
	exe := mkExe(iterStep, bodyStep, exitStep)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, hasNext, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, iterStep.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(iterStep.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, iterStep.ID))

	for i := 0; i < iterations; i++ {
		s, _, _, hasNext, err = ss.Next(ctx, nil, "")
		require.NoError(t, err)
		require.True(t, hasNext, "iteration %d: expected body step", i)
		require.Equal(t, bodyStep.ID, s.ID, "iteration %d: expected body step ID", i)
		require.Equal(t, i+1, iterH.calls, "iteration %d: handler call count", i)
		require.NoError(t, ss.StoreOutputs(bodyStep.ID, nil))
		require.NoError(t, ss.OnStepComplete(ctx, bodyStep.ID))
	}

	s, _, _, hasNext, err = ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, exitStep.ID, s.ID)
}

// TL;DR: If the iterator step has no body child wired, the scheduler returns an error.
// Example: a misconfigured workflow where the loop body edge was never drawn.
func TestIterator_MissingBodyEdge_Errors(t *testing.T) {
	iterH := &mockIteratorHandler{maxIter: 1}
	iterStep := types.Step{
		ID:       stepID(1),
		Handle:   "iter",
		Kind:     "iterator",
		Handler:  iterH,
		Children: nil,
	}
	exe := mkExe(iterStep)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.Equal(t, iterStep.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(iterStep.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, iterStep.ID))

	_, _, _, _, err = ss.Next(ctx, nil, "")
	require.Error(t, err)
}

// TL;DR: After exhausting an iterator with no exit edge, the stack empties cleanly.
// Example: a "fire and forget" loop with no downstream step — automation simply ends.
func TestIterator_NoExitNode(t *testing.T) {
	const iterations = 2
	bodyStep := mkStep(2, "function")
	bodyStep.Parents = []types.Step{{ID: stepID(1)}}

	iterH := &mockIteratorHandler{maxIter: iterations}
	iterStep := types.Step{
		ID:       stepID(1),
		Handle:   "iter",
		Kind:     "iterator",
		Handler:  iterH,
		Children: []types.Step{bodyStep},
	}
	exe := mkExe(iterStep, bodyStep)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, hasNext, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, iterStep.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(iterStep.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, iterStep.ID))

	for i := 0; i < iterations; i++ {
		s, _, _, hasNext, err = ss.Next(ctx, nil, "")
		require.NoError(t, err)
		require.True(t, hasNext, "iteration %d: expected body step", i)
		require.Equal(t, bodyStep.ID, s.ID)
		require.NoError(t, ss.StoreOutputs(bodyStep.ID, nil))
		require.NoError(t, ss.OnStepComplete(ctx, bodyStep.ID))
	}

	_, _, _, hasNext, err = ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.False(t, hasNext)
}

// TL;DR: An iterator whose body is itself an iterator runs the inner body outerN×innerN times.
// Example: "for each order" → "for each line item" → process line item.
func TestIterator_Nested(t *testing.T) {
	const outerIter, innerIter = 2, 3

	bodyStep := mkStep(4, "function")
	bodyStep.Parents = []types.Step{{ID: stepID(3)}}

	innerH := &mockIteratorHandler{maxIter: innerIter}
	innerStep := types.Step{
		ID:       stepID(3),
		Handle:   "inner",
		Kind:     "iterator",
		Handler:  innerH,
		Children: []types.Step{bodyStep},
	}
	innerStep.Parents = []types.Step{{ID: stepID(1)}}

	outerH := &mockIteratorHandler{maxIter: outerIter}
	outerStep := types.Step{
		ID:       stepID(1),
		Handle:   "outer",
		Kind:     "iterator",
		Handler:  outerH,
		Children: []types.Step{innerStep},
	}

	exe := mkExe(outerStep, innerStep, bodyStep)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, hasNext, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, outerStep.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(outerStep.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, outerStep.ID))

	bodyRuns := 0
	for o := 0; o < outerIter; o++ {
		s, _, _, hasNext, err = ss.Next(ctx, nil, "")
		require.NoError(t, err)
		require.True(t, hasNext, "outer %d: expected inner iterator step", o)
		require.Equal(t, innerStep.ID, s.ID, "outer %d: expected inner iterator ID", o)
		require.NoError(t, ss.StoreOutputs(innerStep.ID, nil))
		require.NoError(t, ss.OnStepComplete(ctx, innerStep.ID))

		innerH.calls = 0

		for i := 0; i < innerIter; i++ {
			s, _, _, hasNext, err = ss.Next(ctx, nil, "")
			require.NoError(t, err)
			require.True(t, hasNext, "outer %d inner %d: expected body step", o, i)
			require.Equal(t, bodyStep.ID, s.ID, "outer %d inner %d: expected body step ID", o, i)
			require.NoError(t, ss.StoreOutputs(bodyStep.ID, nil))
			require.NoError(t, ss.OnStepComplete(ctx, bodyStep.ID))
			bodyRuns++
		}
	}

	require.Equal(t, outerIter*innerIter, bodyRuns)

	_, _, _, hasNext, err = ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.False(t, hasNext)
}

package runtime

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/stretchr/testify/require"
)

// TL;DR: Exclusive gateway routes to exactly one child; the other is never scheduled.
// Example: an "if approved / if rejected" split where only one path runs.
func TestExclusiveGateway_SingleBranch(t *testing.T) {
	child0 := mkStep(2, "function")
	child1 := mkStep(3, "function")
	child0.Parents = []types.Step{{ID: stepID(1)}}
	child1.Parents = []types.Step{{ID: stepID(1)}}

	gwH := &mockGatewayHandler{indices: []int{0}}
	gwStep := types.Step{
		ID:       stepID(1),
		Handle:   "gw",
		Kind:     "gatewayExclusive",
		Handler:  gwH,
		Children: []types.Step{child0, child1},
	}
	exe := mkExe(gwStep, child0, child1)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.Equal(t, gwStep.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(gwStep.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, gwStep.ID))

	s, _, _, hasNext, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext)
	require.Equal(t, child0.ID, s.ID)

	require.NoError(t, ss.StoreOutputs(child0.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, child0.ID))
	_, _, _, hasNext, err = ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.False(t, hasNext)
}

// TL;DR: An out-of-range index from the gateway handler is silently ignored, leaving the stack empty.
// Example: a gateway whose condition evaluates to a branch that no longer exists after a schema change.
func TestExclusiveGateway_OutOfRangeIndex_EmptyStack(t *testing.T) {
	child0 := mkStep(2, "function")

	gwH := &mockGatewayHandler{indices: []int{5}}
	gwStep := types.Step{
		ID:       stepID(1),
		Handle:   "gw",
		Kind:     "gatewayExclusive",
		Handler:  gwH,
		Children: []types.Step{child0},
	}
	exe := mkExe(gwStep)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.Equal(t, gwStep.ID, s.ID)
	require.NoError(t, ss.StoreOutputs(gwStep.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, gwStep.ID))

	_, _, _, hasNext, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.False(t, hasNext)
}

// TL;DR: Inclusive gateway schedules all selected children and registers a distinct branch ID for each.
// Example: a "notify all relevant teams" step that fans out to HR, Finance, and Legal simultaneously.
func TestInclusiveGateway_AllBranchesScheduled(t *testing.T) {
	child0 := mkStep(2, "function")
	child1 := mkStep(3, "function")
	child0.Parents = []types.Step{{ID: stepID(1)}}
	child1.Parents = []types.Step{{ID: stepID(1)}}

	gwH := &mockGatewayHandler{indices: []int{0, 1}}
	gwStep := types.Step{
		ID:       stepID(1),
		Handle:   "gw",
		Kind:     "gatewayInclusive",
		Handler:  gwH,
		Children: []types.Step{child0, child1},
	}
	exe := mkExe(gwStep, child0, child1)
	ss := newScheduler(exe)
	ctx := context.Background()

	s, _, _, _, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.Equal(t, gwStep.ID, s.ID)

	branchIDsBefore := make(map[id.ID]bool)
	for bid := range ss.activeBranches {
		branchIDsBefore[bid] = true
	}

	require.NoError(t, ss.StoreOutputs(gwStep.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, gwStep.ID))

	s1, _, _, hasNext1, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext1)

	newBranches := 0
	for bid := range ss.activeBranches {
		if !branchIDsBefore[bid] {
			newBranches++
		}
	}
	require.Equal(t, 2, newBranches, "gateway should register 2 new branch IDs")

	seen := make(map[id.ID]bool)
	seen[s1.ID] = true
	require.NoError(t, ss.StoreOutputs(s1.ID, nil))
	require.NoError(t, ss.OnStepComplete(ctx, s1.ID))

	s2, _, _, hasNext2, err := ss.Next(ctx, nil, "")
	require.NoError(t, err)
	require.True(t, hasNext2)
	seen[s2.ID] = true

	require.True(t, seen[child0.ID], "child0 should have been scheduled")
	require.True(t, seen[child1.ID], "child1 should have been scheduled")
}

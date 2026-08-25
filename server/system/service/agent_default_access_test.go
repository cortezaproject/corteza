package service

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func agentSvcWithNS(id uint64, err error) *agent {
	svc := &agent{services: &agentServices{}}
	svc.WithNamespaceResolver(func(context.Context, string) (uint64, error) { return id, err })
	return svc
}

// TL;DR: an agent scoped to a namespace but granted nothing gets read access to
// that namespace.
// Example: access is deny-by-default, so a freshly created agent refused every
// question — which reads as broken rather than as unconfigured.
func TestDefaultAccess(t *testing.T) {
	ctx := context.Background()
	scoped := func() *types.Agent {
		return &types.Agent{Access: types.AgentAccess{
			Context: types.AgentAccessContext{Namespace: "mtg_collection"},
		}}
	}

	t.Run("fills in a read-only grant for the namespace", func(t *testing.T) {
		a := scoped()
		require.NoError(t, agentSvcWithNS(100, nil).defaultAccess(ctx, a))
		require.Len(t, a.Access.Tools, 1)

		g := a.Access.Tools[0]
		require.Equal(t, "usage", g.Group)
		require.Equal(t, "read", g.MaxRisk)
		require.Equal(t, uint64(100), g.Allow[0].NamespaceID)
		require.Empty(t, g.Allow[0].ModuleIDs, "an empty module list is the whole namespace")

		// A grant with nothing telling the agent what it reaches leaves it
		// guessing handles and reporting they do not exist.
		require.True(t, a.Behavior.InjectSystemContext)
	})

	t.Run("an already-configured agent keeps its own context setting", func(t *testing.T) {
		a := scoped()
		a.Access.Tools = []types.AgentAccessTool{{Name: "compose_record_lookup"}}
		require.NoError(t, agentSvcWithNS(100, nil).defaultAccess(ctx, a))
		require.False(t, a.Behavior.InjectSystemContext, "nothing was defaulted, so nothing was turned on")
	})

	// Anything the author actually asked for is left exactly as written.
	t.Run("an existing grant is never touched", func(t *testing.T) {
		a := scoped()
		a.Access.Tools = []types.AgentAccessTool{{Name: "compose_record_create"}}
		require.NoError(t, agentSvcWithNS(100, nil).defaultAccess(ctx, a))
		require.Equal(t, []types.AgentAccessTool{{Name: "compose_record_create"}}, a.Access.Tools)
	})

	t.Run("a TAQ or workflow grant counts as configured", func(t *testing.T) {
		for _, a := range []*types.Agent{
			{Access: types.AgentAccess{
				Context: types.AgentAccessContext{Namespace: "x"},
				TAQs:    []types.AgentAccessTAQ{{ID: 1}},
			}},
			{Access: types.AgentAccess{
				Context:   types.AgentAccessContext{Namespace: "x"},
				Workflows: []types.AgentAccessWorkflow{{ID: 1}},
			}},
		} {
			require.NoError(t, agentSvcWithNS(100, nil).defaultAccess(ctx, a))
			require.Empty(t, a.Access.Tools)
		}
	})

	// Nothing to scope the grant to means no grant: one with no scope is denied
	// anyway, so inventing it would only look like access.
	t.Run("no namespace means no default", func(t *testing.T) {
		a := &types.Agent{}
		require.NoError(t, agentSvcWithNS(100, nil).defaultAccess(ctx, a))
		require.Empty(t, a.Access.Tools)
	})

	t.Run("an unresolvable namespace is left to the caller", func(t *testing.T) {
		a := scoped()
		require.NoError(t, agentSvcWithNS(0, context.Canceled).defaultAccess(ctx, a))
		require.Empty(t, a.Access.Tools)
	})

	t.Run("no resolver wired means no default", func(t *testing.T) {
		a := scoped()
		svc := &agent{services: &agentServices{}}
		require.NoError(t, svc.defaultAccess(ctx, a))
		require.Empty(t, a.Access.Tools)
	})
}

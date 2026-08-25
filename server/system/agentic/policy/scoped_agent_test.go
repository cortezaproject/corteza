package policy

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/assert"
)

// The shape a read-only assistant is given: a handful of lookup tools, each
// pinned to one namespace and its modules. Everything this asserts is a
// property someone relies on when they scope an agent to a namespace — that it
// cannot read a neighbouring one, and cannot write at all.
func TestReadOnlyNamespaceScopedAgent(t *testing.T) {
	const (
		ns      = "100"
		other   = "200"
		modCard = "11"
		modDeck = "12"
		modAway = "99" // in the same namespace, deliberately not granted
	)

	agent := &types.Agent{
		Access: types.AgentAccess{
			Tools: []types.AgentAccessTool{
				{Name: "compose_namespace_lookup", Allow: []types.AgentAccessAllow{{NamespaceID: 100}}},
				{Name: "compose_record_lookup", Allow: []types.AgentAccessAllow{
					{NamespaceID: 100, ModuleIDs: []uint64{11, 12}},
				}},
			},
		},
	}

	ctx := context.Background()

	t.Run("reads a granted module", func(t *testing.T) {
		d := Evaluate(ctx, agent, "compose_record_lookup",
			MapValues{"namespaceID": ns, "moduleID": modCard}, nil)
		assert.True(t, d.Allowed, d.Reason)
	})

	t.Run("cannot read another namespace", func(t *testing.T) {
		d := Evaluate(ctx, agent, "compose_record_lookup",
			MapValues{"namespaceID": other, "moduleID": modCard}, nil)
		assert.False(t, d.Allowed)
	})

	t.Run("cannot read a module it was not granted", func(t *testing.T) {
		d := Evaluate(ctx, agent, "compose_record_lookup",
			MapValues{"namespaceID": ns, "moduleID": modAway}, nil)
		assert.False(t, d.Allowed)
	})

	// The whole point of a read-only agent: naming the lookup tool must not
	// carry its write siblings in with it.
	for _, tool := range []string{
		"compose_record_create", "compose_record_update", "compose_record_delete",
		"compose_module_create", "compose_namespace_delete",
	} {
		t.Run("cannot "+tool, func(t *testing.T) {
			d := Evaluate(ctx, agent, tool,
				MapValues{"namespaceID": ns, "moduleID": modCard}, nil)
			assert.False(t, d.Allowed, "%s must not be reachable", tool)
		})
	}

	t.Run("namespace lookup is allowed by its own entry", func(t *testing.T) {
		d := Evaluate(ctx, agent, "compose_namespace_lookup", MapValues{"namespaceID": ns}, nil)
		assert.True(t, d.Allowed, d.Reason)
	})

	// Listing a resource-mapped tool with no allow entries reads as
	// "unrestricted" and is the opposite: checkAllow denies it outright. Worth
	// pinning, because the mistake is silent — the agent simply cannot work.
	t.Run("an empty allow denies rather than widens", func(t *testing.T) {
		bare := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{{Name: "compose_namespace_lookup"}},
			},
		}
		d := Evaluate(ctx, bare, "compose_namespace_lookup", MapValues{"namespaceID": ns}, nil)
		assert.False(t, d.Allowed)
		assert.Contains(t, d.Reason, "allow entries")
	})
}

package runtime

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

type fakeResolver struct {
	ns  map[uint64]string
	mod map[uint64]string
}

func (f fakeResolver) Resolve(context.Context, string, string) (uint64, uint64, error) {
	return 0, 0, nil
}

func (f fakeResolver) LookupNamespace(_ context.Context, id uint64) (NsHandle, error) {
	h, ok := f.ns[id]
	if !ok {
		return NsHandle{}, fmt.Errorf("no namespace %d", id)
	}
	return NsHandle{ID: id, Handle: h}, nil
}

func (f fakeResolver) LookupModule(_ context.Context, _, modID uint64) (ModHandle, error) {
	h, ok := f.mod[modID]
	if !ok {
		return ModHandle{}, fmt.Errorf("no module %d", modID)
	}
	return ModHandle{ID: modID, Handle: h}, nil
}

func scopedAgent() *types.Agent {
	return &types.Agent{Access: types.AgentAccess{
		Tools: []types.AgentAccessTool{{
			Name:  "compose_record_lookup",
			Allow: []types.AgentAccessAllow{{NamespaceID: 100, ModuleIDs: []uint64{11}}},
		}},
	}}
}

// TL;DR: the accessible list is framed as a grant, not an inventory.
// Example: a namespace-scoped agent asked what namespaces exist answered "the
// only namespace on this instance is X" — its own access, stated as a fact
// about the world, to someone with no way to tell the difference.
func TestBuildComposeContext_GrantIsNotAnInventory(t *testing.T) {
	r := fakeResolver{ns: map[uint64]string{100: "mtg_collection"}, mod: map[uint64]string{11: "card"}}

	out := buildComposeContext(context.Background(), scopedAgent(), r)
	require.NotEmpty(t, out)

	t.Run("still lists what it can reach", func(t *testing.T) {
		require.Contains(t, out, "mtg_collection")
		require.Contains(t, out, "card")
	})

	t.Run("says the list is a grant", func(t *testing.T) {
		require.Contains(t, out, "NOT an inventory")
	})

	t.Run("warns others may exist", func(t *testing.T) {
		require.True(t,
			strings.Contains(out, "others may exist") || strings.Contains(out, "cannot see"),
			"the model has to be told its view is partial: %s", out)
	})

	// An agent granted nothing gets no section at all rather than an empty one
	// that reads as "there is nothing".
	t.Run("no grant produces no section", func(t *testing.T) {
		require.Empty(t, buildComposeContext(context.Background(), &types.Agent{}, r))
	})
}

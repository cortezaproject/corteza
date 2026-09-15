package service

import (
	"testing"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

// An agent is stored with a name and nothing else; provider and model are
// required by the editor, not here.
func TestAgentCreate_RequiresOnlyAName(t *testing.T) {
	ctx, s := newScenarioTestStore(t)
	svc := &agent{store: s}

	t.Run("a name alone is stored, with no provider or model", func(t *testing.T) {
		res := &types.Agent{Meta: types.AgentMeta{Short: "Label agent"}}
		require.NoError(t, svc.onCreate(ctx, res))

		stored, err := store.LookupAgentByID(ctx, s, res.ID)
		require.NoError(t, err)
		require.Equal(t, "Label agent", stored.Meta.Short)
		require.Zero(t, stored.Execution.Model.LLMProviderID)
		require.Empty(t, stored.Execution.Model.Model)
	})

	for name, short := range map[string]string{"no name": "", "a blank name": "   "} {
		t.Run(name+" is refused", func(t *testing.T) {
			err := svc.onCreate(ctx, &types.Agent{Meta: types.AgentMeta{Short: short}})
			require.Error(t, err)
			require.Contains(t, err.Error(), "missing name")
		})
	}
}

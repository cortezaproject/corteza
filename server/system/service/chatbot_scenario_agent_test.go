package service

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newScenarioTestStore(t *testing.T) (context.Context, store.Storer) {
	t.Helper()
	ctx := context.Background()
	s, err := sqlite.ConnectInMemory(ctx)
	require.NoError(t, err)
	require.NoError(t, store.Upgrade(ctx, zap.NewNop(), s))
	return ctx, s
}

// TL;DR: a conversation scenario must name an agent that is actually there.
// Example: a chatbot saved with a mistyped or since-deleted agentID looked
// healthy in the admin list and failed the conversation in front of a visitor —
// the agent is resolved when someone first writes, not when the chatbot is
// saved.
func TestValidateChatbotScenarios_AgentMustExist(t *testing.T) {
	ctx, s := newScenarioTestStore(t)

	live := &types.Agent{ID: nextID(), Handle: "live_agent", CreatedAt: *now()}
	require.NoError(t, store.CreateAgent(ctx, s, live))

	deletedAt := time.Now()
	gone := &types.Agent{ID: nextID(), Handle: "gone_agent", CreatedAt: *now(), DeletedAt: &deletedAt}
	require.NoError(t, store.CreateAgent(ctx, s, gone))

	scenario := func(agentID uint64) types.ChatbotScenarios {
		return types.ChatbotScenarios{{ID: "qa", Name: "QA", Type: "conversation", AgentID: agentID}}
	}

	t.Run("an existing agent passes", func(t *testing.T) {
		require.NoError(t, validateChatbotScenarios(ctx, s, scenario(live.ID)))
	})

	t.Run("an agent that never existed is refused", func(t *testing.T) {
		err := validateChatbotScenarios(ctx, s, scenario(999999999999999999))
		require.Error(t, err)
		require.Contains(t, err.Error(), "does not exist")
	})

	t.Run("a deleted agent is refused", func(t *testing.T) {
		err := validateChatbotScenarios(ctx, s, scenario(gone.ID))
		require.Error(t, err)
		require.Contains(t, err.Error(), "deleted")
	})

	// The rule this replaced, still enforced.
	t.Run("no agent at all is refused", func(t *testing.T) {
		require.Error(t, validateChatbotScenarios(ctx, s, scenario(0)))
	})

	// Only a conversation hands off to an agent; other scenario types carry no
	// agentID and must not be judged as though they did.
	t.Run("a non-conversation scenario is left alone", func(t *testing.T) {
		other := types.ChatbotScenarios{{ID: "form", Name: "Form", Type: "form"}}
		require.NoError(t, validateChatbotScenarios(ctx, s, other))
	})

	t.Run("every scenario is checked, not just the first", func(t *testing.T) {
		ss := types.ChatbotScenarios{
			{ID: "ok", Name: "OK", Type: "conversation", AgentID: live.ID},
			{ID: "bad", Name: "Bad", Type: "conversation", AgentID: 999999999999999999},
		}
		err := validateChatbotScenarios(ctx, s, ss)
		require.Error(t, err)
		require.Contains(t, err.Error(), `"bad"`)
	})
}

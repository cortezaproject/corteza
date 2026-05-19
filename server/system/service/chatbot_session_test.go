package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/types"
)

func newTestChatbotSession(t *testing.T) (*chatbotSession, store.Storer, uint64) {
	t.Helper()
	req := require.New(t)
	ctx := context.Background()

	s, err := sqlite.ConnectInMemory(ctx)
	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))

	sessionID := nextID()
	req.NoError(store.CreateChatbotSession(ctx, s, &types.ChatbotSession{
		ID:        sessionID,
		ChatbotID: nextID(),
		Status:    "active",
	}))

	svc := &chatbotSession{
		actionlog: actionlog.NewService(s, zap.NewNop(), zap.NewNop(), actionlog.MakeDisabledPolicy()),
		store:     s,
		ac:        DefaultAccessControl,
	}
	return svc, s, sessionID
}

func TestStartStep(t *testing.T) {
	ctx := context.Background()
	svc, s, sessionID := newTestChatbotSession(t)
	convID := nextID()

	t.Run("creates a step in active state", func(t *testing.T) {
		scenario := &types.ChatbotScenario{ID: "intro", Type: "form"}
		step, err := svc.StartStep(ctx, sessionID, scenario, convID, 0, nil)
		assert.NoError(t, err)
		assert.NotNil(t, step)
		assert.Equal(t, "active", step.Status)
		assert.Equal(t, 0, step.ScenarioIndex)
		assert.Equal(t, convID, step.ConversationID)

		stored, err := store.LookupChatbotSessionStepByID(ctx, s, step.ID)
		assert.NoError(t, err)
		assert.Equal(t, "active", stored.Status)
	})

	t.Run("reuses an existing step row", func(t *testing.T) {
		scenario := &types.ChatbotScenario{ID: "intro2", Type: "form"}
		s1, err := svc.StartStep(ctx, sessionID, scenario, convID, 1, nil)
		assert.NoError(t, err)
		s2, err := svc.StartStep(ctx, sessionID, scenario, convID, 1, nil)
		assert.NoError(t, err)
		assert.Equal(t, s1.ID, s2.ID)
	})

	t.Run("before-automation failure marks step failed", func(t *testing.T) {
		scenario := &types.ChatbotScenario{
			ID:   "broken",
			Type: "conversation",
			Automation: types.ChatbotScenarioAutomation{
				Before: types.ChatbotAutomationHook{
					Automation: fmt.Sprintf("corteza::automation:ng-automation/%d", nextID()),
				},
			},
		}
		step, err := svc.StartStep(ctx, sessionID, scenario, convID, 2, nil)
		assert.Error(t, err)
		assert.NotNil(t, step)
		assert.Equal(t, "failed", step.Status)
	})
}

func TestFinalizeStep(t *testing.T) {
	ctx := context.Background()
	svc, s, sessionID := newTestChatbotSession(t)
	convID := nextID()

	scenario := &types.ChatbotScenario{ID: "x", Type: "conversation"}
	step, err := svc.StartStep(ctx, sessionID, scenario, convID, 0, nil)
	require.NoError(t, err)

	require.NoError(t, svc.FinalizeStep(ctx, step.ID, scenario, nil))

	stored, err := store.LookupChatbotSessionStepByID(ctx, s, step.ID)
	require.NoError(t, err)
	assert.Equal(t, "complete", stored.Status)
}

func TestHandoffLifecycle(t *testing.T) {
	ctx := context.Background()
	svc, s, sessionID := newTestChatbotSession(t)
	convID := nextID()

	scenario := &types.ChatbotScenario{ID: "conv", Type: "conversation"}
	step, err := svc.StartStep(ctx, sessionID, scenario, convID, 0, nil)
	require.NoError(t, err)

	h, err := svc.RequestHandoff(ctx, sessionID, step.ID)
	require.NoError(t, err)
	assert.Equal(t, "requested", h.Status)

	sess, err := svc.FindByID(ctx, sessionID)
	require.NoError(t, err)
	assert.Equal(t, "handoff_requested", sess.Status)

	require.NoError(t, svc.ActivateHandoff(ctx, h.ID))
	sess, _ = svc.FindByID(ctx, sessionID)
	assert.Equal(t, "handoff_active", sess.Status)

	// Closing the handoff MUST resume the session at "active" (not "closed")
	// so the AI loop can pick up further user messages on the same step.
	require.NoError(t, svc.CloseHandoff(ctx, h.ID))
	sess, _ = svc.FindByID(ctx, sessionID)
	assert.Equal(t, "active", sess.Status)

	closed, err := store.LookupChatbotSessionHandoffByID(ctx, s, h.ID)
	require.NoError(t, err)
	assert.Equal(t, "closed", closed.Status)
	assert.NotNil(t, closed.ClosedAt)
}

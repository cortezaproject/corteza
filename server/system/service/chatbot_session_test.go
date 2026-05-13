package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/types"
)

func TestExecuteStep(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()

		s store.Storer
		err error

		sessionID      = nextID()
		conversationID = nextID()
		scenarioIdx    = 0
	)

	// Setup in-memory SQLite
	if s, err = sqlite.ConnectInMemory(ctx); err != nil {
		req.NoError(err)
	} else if err = store.Upgrade(ctx, zap.NewNop(), s); err != nil {
		req.NoError(err)
	}

	// Create session
	sess := &types.ChatbotSession{
		ID:        sessionID,
		ChatbotID: nextID(),
		Status:    "active",
	}
	req.NoError(store.CreateChatbotSession(ctx, s, sess))

	// Create service with mock access control
	svc := &chatbotSession{
		actionlog: actionlog.NewService(s, zap.NewNop(), zap.NewNop(), actionlog.MakeDisabledPolicy()),
		store:     s,
		ac:        DefaultAccessControl,
	}

	t.Run("creates step and marks complete", func(t *testing.T) {
		scenario := &types.ChatbotScenario{
			ID:   "test-scenario",
			Type: "conversation",
		}

		step, err := svc.ExecuteStep(ctx, sessionID, scenario, conversationID, scenarioIdx, "")
		assert.NoError(t, err)
		assert.NotNil(t, step)
		assert.Equal(t, "complete", step.Status)
		assert.Equal(t, scenarioIdx, step.ScenarioIndex)
		assert.Equal(t, conversationID, step.ConversationID)

		// Verify step was persisted
		stored, err := store.LookupChatbotSessionStepByID(ctx, s, step.ID)
		assert.NoError(t, err)
		assert.NotNil(t, stored)
		assert.Equal(t, "complete", stored.Status)
	})

	t.Run("before automation failure marks step failed", func(t *testing.T) {
		automationID := nextID()
		scenario := &types.ChatbotScenario{
			ID:                 "test-scenario-before-fail",
			Type:               "conversation",
			BeforeAutomationID: &automationID,
		}

		// Without DefaultNgAutomation set, invokeAutomation will return error
		step, err := svc.ExecuteStep(ctx, sessionID, scenario, conversationID, scenarioIdx+1, "")

		// Execution should fail
		assert.Error(t, err)
		assert.NotNil(t, step)
		assert.Equal(t, "failed", step.Status)
	})

	t.Run("reuses existing step on retry", func(t *testing.T) {
		// Create first step
		scenario := &types.ChatbotScenario{
			ID:   "reuse-scenario",
			Type: "conversation",
		}

		step1, err := svc.ExecuteStep(ctx, sessionID, scenario, conversationID, scenarioIdx+2, "")
		assert.NoError(t, err)
		stepID := step1.ID

		// Call ExecuteStep again without creating new step
		step2, err := svc.ExecuteStep(ctx, sessionID, scenario, conversationID, scenarioIdx+2, "")
		assert.NoError(t, err)
		assert.Equal(t, stepID, step2.ID)
		assert.Equal(t, "complete", step2.Status)
	})
}

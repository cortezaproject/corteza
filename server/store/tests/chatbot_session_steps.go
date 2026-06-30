package tests

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func testChatbotSessionSteps(t *testing.T, s store.ChatbotSessionSteps) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func() *types.ChatbotSessionStep {
			return &types.ChatbotSessionStep{
				ID:        id.Next(),
				SessionID: 42,
				Status:    "active",
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ChatbotSessionStep) {
			req := require.New(t)
			req.NoError(s.TruncateChatbotSessionSteps(ctx))
			res := makeNew()
			req.NoError(s.CreateChatbotSessionStep(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateChatbotSessionSteps(ctx))
		req.NoError(s.CreateChatbotSessionStep(ctx, makeNew()))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, step := truncAndCreate(t)
		fetched, err := s.LookupChatbotSessionStepByID(ctx, step.ID)
		req.NoError(err)
		req.Equal(step.ID, fetched.ID)
		req.Equal(step.Status, fetched.Status)
	})

	t.Run("update", func(t *testing.T) {
		req, step := truncAndCreate(t)
		step.Status = "done"
		req.NoError(s.UpdateChatbotSessionStep(ctx, step))
		fetched, err := s.LookupChatbotSessionStepByID(ctx, step.ID)
		req.NoError(err)
		req.Equal("done", fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		req, step := truncAndCreate(t)
		req.NoError(s.DeleteChatbotSessionStepByID(ctx, step.ID))
	})
}

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

func testChatbotSessionHandoffs(t *testing.T, s store.ChatbotSessionHandoffs) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func() *types.ChatbotSessionHandoff {
			return &types.ChatbotSessionHandoff{
				ID:          id.Next(),
				SessionID:   42,
				Status:      "pending",
				InitiatedAt: time.Now(),
				CreatedAt:   time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ChatbotSessionHandoff) {
			req := require.New(t)
			req.NoError(s.TruncateChatbotSessionHandoffs(ctx))
			res := makeNew()
			req.NoError(s.CreateChatbotSessionHandoff(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateChatbotSessionHandoffs(ctx))
		req.NoError(s.CreateChatbotSessionHandoff(ctx, makeNew()))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, h := truncAndCreate(t)
		fetched, err := s.LookupChatbotSessionHandoffByID(ctx, h.ID)
		req.NoError(err)
		req.Equal(h.ID, fetched.ID)
		req.Equal(h.Status, fetched.Status)
	})

	t.Run("update", func(t *testing.T) {
		req, h := truncAndCreate(t)
		h.Status = "closed"
		req.NoError(s.UpdateChatbotSessionHandoff(ctx, h))
		fetched, err := s.LookupChatbotSessionHandoffByID(ctx, h.ID)
		req.NoError(err)
		req.Equal("closed", fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		req, h := truncAndCreate(t)
		req.NoError(s.DeleteChatbotSessionHandoffByID(ctx, h.ID))
	})
}

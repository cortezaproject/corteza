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

func testChatbotSessions(t *testing.T, s store.ChatbotSessions) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func() *types.ChatbotSession {
			return &types.ChatbotSession{
				ID:        id.Next(),
				ChatbotID: 42,
				Status:    "active",
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ChatbotSession) {
			req := require.New(t)
			req.NoError(s.TruncateChatbotSessions(ctx))
			res := makeNew()
			req.NoError(s.CreateChatbotSession(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateChatbotSessions(ctx))
		cs := makeNew()
		req.NoError(s.CreateChatbotSession(ctx, cs))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, cs := truncAndCreate(t)
		fetched, err := s.LookupChatbotSessionByID(ctx, cs.ID)
		req.NoError(err)
		req.Equal(cs.ID, fetched.ID)
		req.Equal(cs.Status, fetched.Status)
	})

	t.Run("update", func(t *testing.T) {
		req, cs := truncAndCreate(t)
		cs.Status = "closed"
		req.NoError(s.UpdateChatbotSession(ctx, cs))
		fetched, err := s.LookupChatbotSessionByID(ctx, cs.ID)
		req.NoError(err)
		req.Equal("closed", fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		req, cs := truncAndCreate(t)
		req.NoError(s.DeleteChatbotSessionByID(ctx, cs.ID))
	})
}

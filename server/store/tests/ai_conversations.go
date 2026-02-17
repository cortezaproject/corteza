package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func testAiConversations(t *testing.T, s store.AiConversations) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func() *types.AiConversation {
			return &types.AiConversation{
				ID:        id.Next(),
				AgentID:   id.Next(),
				CreatedAt: time.Now(),
				Messages:  types.AiConversationMessages{},
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.AiConversation) {
			req := require.New(t)
			req.NoError(s.TruncateAiConversations(ctx))
			res := makeNew()
			req.NoError(s.CreateAiConversation(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateAiConversations(ctx))
		conv := makeNew()
		req.NoError(s.CreateAiConversation(ctx, conv))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, conv := truncAndCreate(t)
		fetched, err := s.LookupAiConversationByID(ctx, conv.ID)
		req.NoError(err)
		req.Equal(conv.ID, fetched.ID)
		req.Equal(conv.AgentID, fetched.AgentID)
	})

	t.Run("delete", func(t *testing.T) {
		req, conv := truncAndCreate(t)
		req.NoError(s.DeleteAiConversationByID(ctx, conv.ID))
	})
}

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/rand"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func testChatbots(t *testing.T, s store.Chatbots) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.Chatbot {
			return &types.Chatbot{
				ID:        id.Next(),
				Handle:    handle,
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Chatbot) {
			req := require.New(t)
			req.NoError(s.TruncateChatbots(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateChatbot(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateChatbots(ctx))
		c := makeNew("ChatbotCRUD")
		req.NoError(s.CreateChatbot(ctx, c))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, c := truncAndCreate(t)
		fetched, err := s.LookupChatbotByID(ctx, c.ID)
		req.NoError(err)
		req.Equal(c.Handle, fetched.Handle)
		req.Equal(c.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, c := truncAndCreate(t)
		c.Handle = string(rand.Bytes(10))
		req.NoError(s.UpdateChatbot(ctx, c))
		fetched, err := s.LookupChatbotByID(ctx, c.ID)
		req.NoError(err)
		req.Equal(c.Handle, fetched.Handle)
	})

	t.Run("delete", func(t *testing.T) {
		req, c := truncAndCreate(t)
		req.NoError(s.DeleteChatbotByID(ctx, c.ID))
	})
}

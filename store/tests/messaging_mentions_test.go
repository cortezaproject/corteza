package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func testMessagingMentions(t *testing.T, s store.MessagingMentions) {
	var (
		ctx = context.Background()
		req = require.New(t)

		channelID = id.Next()
		messageID = id.Next()
		userID    = id.Next()

		makeNew = func(mention string) *types.Mention {
			// minimum data set for new messagingMention
			return &types.Mention{
				ID:        id.Next(),
				ChannelID: channelID,
				MessageID: messageID,
				UserID:    userID,
				CreatedAt: *now(),
			}
		}
	)

	t.Run("create", func(t *testing.T) {
		messagingMention := makeNew("f")
		req.NoError(s.CreateMessagingMention(ctx, messagingMention))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup by ID", func(t *testing.T) {
		messagingMention := makeNew("look-up-by-id")
		req.NoError(s.CreateMessagingMention(ctx, messagingMention))
		fetched, err := s.LookupMessagingMentionByID(ctx, messagingMention.ID)
		req.NoError(err)
		req.Equal(messagingMention.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
	})

	t.Run("Delete", func(t *testing.T) {
		messagingMention := makeNew("f")
		req.NoError(s.CreateMessagingMention(ctx, messagingMention))
		req.NoError(s.DeleteMessagingMention(ctx))
	})

	t.Run("Delete by ID", func(t *testing.T) {
		messagingMention := makeNew("Delete-by-id")
		req.NoError(s.CreateMessagingMention(ctx, messagingMention))
		req.NoError(s.DeleteMessagingMention(ctx))
	})
}

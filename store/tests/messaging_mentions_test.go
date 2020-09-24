package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
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

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Mention) {
			req := require.New(t)
			req.NoError(s.TruncateMessagingMentions(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateMessagingMention(ctx, res))
			return req, res
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
		req, messagingMention := truncAndCreate(t)
		fetched, err := s.LookupMessagingMentionByID(ctx, messagingMention.ID)
		req.NoError(err)
		req.Equal(messagingMention.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Mention", func(t *testing.T) {
			req, messagingMention := truncAndCreate(t)
			req.NoError(s.DeleteMessagingMention(ctx, messagingMention))
			_, err := s.LookupMessagingMentionByID(ctx, messagingMention.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, messagingMention := truncAndCreate(t)
			req.NoError(s.DeleteMessagingMentionByID(ctx, messagingMention.ID))
			_, err := s.LookupMessagingMentionByID(ctx, messagingMention.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})
}

package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func testMessagingUnread(t *testing.T, s store.MessagingUnreads) {
	var (
		ctx = context.Background()

		channelID = id.Next()
		userID    = id.Next()
		threadID  = id.Next()

		makeNew = func(channelID, threadID, userID uint64) *types.Unread {
			// minimum data set for new messagingUnread
			return &types.Unread{
				ChannelID: channelID,
				ReplyTo:   threadID,
				UserID:    userID,
			}
		}
	)

	t.Run("create", func(t *testing.T) {
		var req = require.New(t)

		messagingUnread := makeNew(channelID, threadID, userID)
		req.NoError(s.CreateMessagingUnread(ctx, messagingUnread))
	})

	t.Run("count", func(t *testing.T) {
		var (
			err error
			req = require.New(t)
		)

		_, err = s.CountMessagingUnread(ctx, 1, 2, 3, 4)
		req.NoError(err)

		_, err = s.CountMessagingUnread(ctx, 1, 2, 3)
		req.NoError(err)

		_, err = s.CountMessagingUnread(ctx, 1, 2)
		req.NoError(err)

		_, err = s.CountMessagingUnread(ctx, 1, 0)
		req.NoError(err)

		_, err = s.CountMessagingUnread(ctx, 0, 0)
		req.NoError(err)

	})
}

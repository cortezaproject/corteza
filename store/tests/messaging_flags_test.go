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

func testMessagingFlags(t *testing.T, s store.MessagingFlags) {
	var (
		ctx = context.Background()
		req = require.New(t)

		channelID = id.Next()
		messageID = id.Next()
		userID    = id.Next()

		makeNew = func(flag string) *types.MessageFlag {
			// minimum data set for new messagingFlag
			return &types.MessageFlag{
				ID:        id.Next(),
				ChannelID: channelID,
				MessageID: messageID,
				UserID:    userID,
				Flag:      flag,
				CreatedAt: *now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.MessageFlag) {
			req := require.New(t)
			req.NoError(s.TruncateMessagingFlags(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateMessagingFlag(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		messagingFlag := makeNew("f")
		req.NoError(s.CreateMessagingFlag(ctx, messagingFlag))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup by ID", func(t *testing.T) {
		messagingFlag := makeNew("look-up-by-id")
		req.NoError(s.CreateMessagingFlag(ctx, messagingFlag))
		fetched, err := s.LookupMessagingFlagByID(ctx, messagingFlag.ID)
		req.NoError(err)
		req.Equal(messagingFlag.Flag, fetched.Flag)
		req.Equal(messagingFlag.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		messagingFlag := makeNew("update-me")
		req.NoError(s.CreateMessagingFlag(ctx, messagingFlag))

		messagingFlag = &types.MessageFlag{
			ID:        messagingFlag.ID,
			CreatedAt: messagingFlag.CreatedAt,
			Flag:      "flg",
		}
		req.NoError(s.UpdateMessagingFlag(ctx, messagingFlag))

		updated, err := s.LookupMessagingFlagByID(ctx, messagingFlag.ID)
		req.NoError(err)
		req.Equal(messagingFlag.Flag, updated.Flag)
	})

	t.Run("update with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Message Flag", func(t *testing.T) {
			req, messagingFlag := truncAndCreate(t)
			req.NoError(s.DeleteMessagingFlag(ctx, messagingFlag))
			_, err := s.LookupMessagingFlagByID(ctx, messagingFlag.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, messagingFlag := truncAndCreate(t)
			req.NoError(s.DeleteMessagingFlagByID(ctx, messagingFlag.ID))
			_, err := s.LookupMessagingFlagByID(ctx, messagingFlag.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.MessageFlag{
			makeNew("flag-1-1"),
			makeNew("flag-1-2"),
			makeNew("flag-2-1"),
			makeNew("flag-2-2"),
		}

		count := len(prefill)

		req.NoError(s.TruncateMessagingFlags(ctx))
		req.NoError(s.CreateMessagingFlag(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchMessagingFlags(ctx, types.MessageFlagFilter{})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		set, f, err = s.SearchMessagingFlags(ctx, types.MessageFlagFilter{Flag: "flag-2-1"})
		req.NoError(err)
		req.Len(set, 1)

		_ = f // dummy
	})
}

package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testMessagingChannels(t *testing.T, s store.MessagingChannels) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(name string) *types.Channel {
			// minimum data set for new messagingChannel
			return &types.Channel{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Name:      name,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Channel) {
			req := require.New(t)
			req.NoError(s.TruncateMessagingChannels(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateMessagingChannel(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		messagingChannel := makeNew("MessagingChannelCRUD")
		req.NoError(s.CreateMessagingChannel(ctx, messagingChannel))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup by ID", func(t *testing.T) {
		messagingChannel := makeNew("look up by id")
		req.NoError(s.CreateMessagingChannel(ctx, messagingChannel))
		fetched, err := s.LookupMessagingChannelByID(ctx, messagingChannel.ID)
		req.NoError(err)
		req.Equal(messagingChannel.Name, fetched.Name)
		req.Equal(messagingChannel.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		messagingChannel := makeNew("update me")
		req.NoError(s.CreateMessagingChannel(ctx, messagingChannel))

		messagingChannel = &types.Channel{
			ID:        messagingChannel.ID,
			CreatedAt: messagingChannel.CreatedAt,
			Name:      "MessagingChannelCRUD+2",
		}
		req.NoError(s.UpdateMessagingChannel(ctx, messagingChannel))

		updated, err := s.LookupMessagingChannelByID(ctx, messagingChannel.ID)
		req.NoError(err)
		req.Equal(messagingChannel.Name, updated.Name)
	})

	t.Run("update with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Channel", func(t *testing.T) {
			req, messagingChannel := truncAndCreate(t)
			req.NoError(s.DeleteMessagingChannel(ctx, messagingChannel))
			_, err := s.LookupMessagingChannelByID(ctx, messagingChannel.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, messagingChannel := truncAndCreate(t)
			req.NoError(s.DeleteMessagingChannelByID(ctx, messagingChannel.ID))
			_, err := s.LookupMessagingChannelByID(ctx, messagingChannel.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Channel{
			makeNew("one-one"),
			makeNew("one-two"),
			makeNew("two-one"),
			makeNew("two-two"),
			makeNew("two-deleted"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateMessagingChannels(ctx))
		req.NoError(s.CreateMessagingChannel(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchMessagingChannels(ctx, types.ChannelFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchMessagingChannels(ctx, types.ChannelFilter{IncludeDeleted: true})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// find all prefixed
		set, f, err = s.SearchMessagingChannels(ctx, types.ChannelFilter{Query: "two-"})
		req.NoError(err)
		req.Len(set, 2)

		_ = f // dummy
	})
}

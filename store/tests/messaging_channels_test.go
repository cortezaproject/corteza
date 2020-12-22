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

func testMessagingChannels(t *testing.T, s store.Storer) {
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
		req, messagingChannel := truncAndCreate(t)
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

	t.Run("lookup by member set", func(t *testing.T) {
		var (
			req = require.New(t)

			ch1, ch2, ch3, ch4 = makeNew("one"), makeNew("two"), makeNew("three"), makeNew("four")
		)

		ch1.Type = types.ChannelTypeGroup
		ch2.Type = types.ChannelTypeGroup
		ch3.Type = types.ChannelTypeGroup
		ch4.Type = types.ChannelTypeGroup

		req.NoError(store.TruncateMessagingChannels(ctx, s))
		req.NoError(store.TruncateMessagingChannelMembers(ctx, s))
		req.NoError(store.CreateMessagingChannel(ctx, s, ch1, ch2, ch3))
		req.NoError(store.CreateMessagingChannelMember(ctx, s,
			// fits
			&types.ChannelMember{CreatedAt: time.Time{}, ChannelID: ch1.ID, UserID: 1000},
			&types.ChannelMember{CreatedAt: time.Time{}, ChannelID: ch1.ID, UserID: 2000},

			// one to many
			&types.ChannelMember{CreatedAt: time.Time{}, ChannelID: ch2.ID, UserID: 1000},
			&types.ChannelMember{CreatedAt: time.Time{}, ChannelID: ch2.ID, UserID: 2000},
			&types.ChannelMember{CreatedAt: time.Time{}, ChannelID: ch2.ID, UserID: 3000},

			// no diff
			&types.ChannelMember{CreatedAt: time.Time{}, ChannelID: ch3.ID, UserID: 1000},
			&types.ChannelMember{CreatedAt: time.Time{}, ChannelID: ch3.ID, UserID: 5000},

			// one only
			&types.ChannelMember{CreatedAt: time.Time{}, ChannelID: ch4.ID, UserID: 1000},
			&types.ChannelMember{CreatedAt: time.Time{}, ChannelID: ch4.ID, UserID: 4000},
			&types.ChannelMember{CreatedAt: time.Time{}, ChannelID: ch4.ID, UserID: 5000},
		))

		ch, err := s.LookupMessagingChannelByMemberSet(ctx, 1000, 2000)
		req.NoError(err)
		req.Equal(ch.ID, ch1.ID)
	})
}

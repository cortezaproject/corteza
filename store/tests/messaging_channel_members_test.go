package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testMessagingChannelMembers(t *testing.T, s store.MessagingChannelMembers) {
	var (
		ctx = context.Background()

		channelID = id.Next()
		userID    = id.Next()

		makeNew = func(channelID, userID uint64) *types.ChannelMember {
			// minimum data set for new messagingChannelMember
			return &types.ChannelMember{
				ChannelID: channelID,
				UserID:    userID,

				Type: types.ChannelMembershipType("owner"),
				Flag: types.ChannelMembershipFlag(""),

				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ChannelMember) {
			req := require.New(t)
			req.NoError(s.TruncateMessagingChannelMembers(ctx))
			res := makeNew(channelID, userID)
			req.NoError(s.CreateMessagingChannelMember(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateMessagingChannelMembers(ctx))
		messagingChannelMember := makeNew(channelID, userID)
		req.NoError(s.CreateMessagingChannelMember(ctx, messagingChannelMember))
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by ChannelMember", func(t *testing.T) {
			req, messagingChannelMember := truncAndCreate(t)
			req.NoError(s.DeleteMessagingChannelMember(ctx, messagingChannelMember))
			set, _, err := s.SearchMessagingChannelMembers(ctx, types.ChannelMemberFilter{ChannelID: []uint64{messagingChannelMember.ChannelID}, MemberID: []uint64{messagingChannelMember.UserID}})
			req.NoError(err)
			req.Len(set, 0)
		})

		t.Run("by ID", func(t *testing.T) {
			req, messagingChannelMember := truncAndCreate(t)
			req.NoError(s.DeleteMessagingChannelMemberByChannelIDUserID(ctx, messagingChannelMember.ChannelID, messagingChannelMember.UserID))
			set, _, err := s.SearchMessagingChannelMembers(ctx, types.ChannelMemberFilter{ChannelID: []uint64{messagingChannelMember.ChannelID}, MemberID: []uint64{messagingChannelMember.UserID}})
			req.NoError(err)
			req.Len(set, 0)
		})
	})
}

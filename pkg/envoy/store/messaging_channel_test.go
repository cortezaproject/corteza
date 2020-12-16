package store

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/messaging/types"
	st "github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"
)

func TestMessagingChannel_Merger(t *testing.T) {
	req := require.New(t)

	now := time.Time{}
	nowP := &time.Time{}

	empty := &types.Channel{}
	full := &types.Channel{
		Name:             "name",
		Topic:            "topic",
		Type:             types.ChannelTypeGroup,
		Meta:             st.JSONText{},
		MembershipPolicy: types.ChannelMembershipPolicyFeatured,
		CreatorID:        1,

		CreatedAt:  now,
		UpdatedAt:  nowP,
		ArchivedAt: nowP,
		DeletedAt:  nowP,
	}

	t.Run("merge on empty", func(t *testing.T) {
		c := mergeMessagingChannels(empty, full)
		req.Equal("name", c.Name)
		req.Equal("topic", c.Topic)
		req.Equal(types.ChannelTypeGroup, c.Type)
		req.NotNil(c.Meta)
		req.Equal(types.ChannelMembershipPolicyFeatured, c.MembershipPolicy)
		req.Equal(uint64(1), c.CreatorID)

		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.ArchivedAt)
		req.Equal(nowP, c.DeletedAt)
	})

	t.Run("merge with empty", func(t *testing.T) {
		c := mergeMessagingChannels(full, empty)
		req.Equal("name", c.Name)
		req.Equal("topic", c.Topic)
		req.Equal(types.ChannelTypeGroup, c.Type)
		req.NotNil(c.Meta)
		req.Equal(types.ChannelMembershipPolicyDefault, c.MembershipPolicy)
		req.Equal(uint64(0), c.CreatorID)

		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.ArchivedAt)
		req.Equal(nowP, c.DeletedAt)
	})
}

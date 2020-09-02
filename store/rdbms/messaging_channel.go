package rdbms

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

// CountReplies counts unread thread info
func (s Store) LookupMessagingChannelByMemberSet(ctx context.Context, memberIDs ...uint64) (ch *types.Channel, err error) {
	return s.execLookupMessagingChannel(ctx, squirrel.And{
		squirrel.Eq{"type": types.ChannelTypeGroup},
		squirrel.Eq{"id": s.getMessagingChannelMembersQuery(memberIDs...)},
	})
}

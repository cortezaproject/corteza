package rdbms

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"strings"
)

// CountReplies counts unread thread info
func (s Store) LookupMessagingChannelByMemberSet(ctx context.Context, memberIDs ...uint64) (ch *types.Channel, err error) {
	return s.execLookupMessagingChannel(ctx, squirrel.And{
		squirrel.Eq{"type": types.ChannelTypeGroup},
		squirrel.Eq{"id": s.getMessagingChannelMembersQuery(memberIDs...)},
	})
}

func (s Store) convertMessagingChannelFilter(f types.ChannelFilter) (query squirrel.SelectBuilder, err error) {
	query = s.messagingChannelsSelectBuilder()

	query = query.Where(squirrel.Eq{"mch.archived_at": nil})

	if !f.IncludeDeleted {
		query = query.Where(squirrel.Eq{"mch.deleted_at": nil})
	}

	if len(f.ChannelID) > 0 {
		query = query.Where(squirrel.Eq{"mch.id": f.ChannelID})
	}

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"LOWER(mch.name)": q},
			squirrel.Like{"LOWER(mch.topic)": q},
		})
	}

	if f.CurrentUserID > 0 {
		qcm := s.SelectBuilder(s.messagingChannelMemberTable("mcm"), "mcm.rel_channel").
			Where(squirrel.Eq{"cmc.rel_user:": f.CurrentUserID})

		query = query.Where(squirrel.Or{
			squirrel.Eq{"c.type": types.ChannelTypePublic},
			squirrel.Eq{"c.id": qcm},
		})
	}

	return
}

package rdbms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"strings"
)

// CountReplies counts unread thread info
func (s Store) LookupMessagingChannelByMemberSet(ctx context.Context, memberIDs ...uint64) (ch *types.Channel, err error) {
	// prepare subquery that merges
	mcmq := s.getMessagingChannelMembersQuery(squirrel.Expr("mch.id = mcm.rel_channel"), memberIDs...)

	if sql, args, err := mcmq.ToSql(); err != nil {
		return nil, err
	} else {
		return s.execLookupMessagingChannel(ctx, squirrel.And{
			squirrel.Eq{"type": types.ChannelTypeGroup},
			squirrel.Expr(fmt.Sprintf("EXISTS (%s)", sql), args...),
		})
	}
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
			// minor squirrel quirk: change placeholder format to question
			// on sub-queries. This will be reset to the right format internally
			// when merged with main query
			PlaceholderFormat(squirrel.Question).
			Where(squirrel.Eq{"mcm.rel_user": f.CurrentUserID})

		qcmSql, qcmArgs, err := qcm.ToSql()
		if err != nil {
			return query, err
		}

		query = query.Where(squirrel.Or{
			squirrel.Eq{"mch.type": types.ChannelTypePublic},
			squirrel.Expr("mch.id IN ("+qcmSql+")", qcmArgs...),
		})
	}

	return
}

package rdbms

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/store"
)

// CountReplies counts unread thread info
func (s Store) CountMessagingUnreadThreads(ctx context.Context, userID, channelID uint64) (uu types.UnreadSet, err error) {
	var (
		rows *sql.Rows

		q = squirrel.Select(
			"rel_channel",
			"rel_user",
			"sum(count) AS count",
			"sum(CASE WHEN count > 0 THEN 1 ELSE 0 END) AS total",
		).
			From(s.messagingUnreadTable()).
			Where("rel_reply_to > 0 AND count > 0").
			GroupBy("rel_channel", "rel_user")
	)

	if userID > 0 {
		q = q.Where(squirrel.Eq{"rel_user": userID})
	}

	if channelID > 0 {
		q = q.Where(squirrel.Eq{"rel_channel": channelID})
	}

	if rows, err = s.Query(ctx, q); err != nil {
		return nil, err
	}

	defer rows.Close()

	uu = make([]*types.Unread, 0, 512)
	for rows.Next() {
		u := &types.Unread{}
		if err = rows.Scan(&u.ChannelID, &u.UserID, &u.ThreadCount, &u.ThreadTotal); err != nil {
			return
		}

		uu = append(uu, u)
	}

	return uu, nil
}

// Count returns counts unread channel info
func (s Store) CountMessagingUnread(ctx context.Context, userID, channelID uint64, threadIDs ...uint64) (uu types.UnreadSet, err error) {
	var (
		q = squirrel.
			Select(
				"rel_channel",
				"rel_last_message",
				"rel_user",
				"rel_reply_to",
				"count",
			).
			From(s.messagingUnreadTable())
	)

	if userID > 0 {
		q = q.Where(squirrel.Eq{"rel_user": userID})
	}

	if channelID > 0 {
		q = q.Where(squirrel.Eq{"rel_channel": channelID})
	}

	if len(threadIDs) == 0 {
		q = q.Where(squirrel.Eq{"rel_reply_to": 0})
	} else {
		q = q.Where(squirrel.Eq{"rel_reply_to": threadIDs})
	}

	uu, _, _, err = s.QueryMessagingUnreads(ctx, q, nil)
	return
}

func (s Store) ResetMessagingUnreadThreads(ctx context.Context, userID, channelID uint64) error {
	var (
		cnd = squirrel.And{
			squirrel.Eq{"rel_user": userID, "rel_channel": channelID},
			squirrel.Gt{"rel_reply_to": 0},
		}
		set = store.Payload{"count": 0}
	)

	return s.execUpdateComposeRecordValues(ctx, cnd, set)

}

func (s Store) IncMessagingUnreadCount(ctx context.Context, channelID uint64, threadID uint64, userID uint64) error {
	var (
		cnd = squirrel.And{
			squirrel.Eq{"rel_reply_to": threadID, "rel_channel": channelID},
			squirrel.NotEq{"rel_user": userID},
		}

		upd = s.UpdateBuilder(s.messagingUnreadTable()).Where(cnd).Set("count", squirrel.Expr("count + 1"))
	)

	return s.config.ErrorHandler(s.Exec(ctx, upd))

}

func (s Store) DecMessagingUnreadCount(ctx context.Context, channelID uint64, threadID uint64, userID uint64) (err error) {
	var (
		cnd = squirrel.And{
			squirrel.Eq{"rel_reply_to": threadID, "rel_channel": channelID},
			squirrel.Gt{"count": 0},
		}

		upd = s.UpdateBuilder(s.messagingUnreadTable()).Where(cnd).Set("count", squirrel.Expr("count - 1"))
	)

	if err = s.Exec(ctx, upd); err != nil {
		return s.config.ErrorHandler(err)
	}

	err = s.UpsertMessagingUnread(ctx, &types.Unread{
		ChannelID: channelID,
		ReplyTo:   threadID,
		UserID:    userID,
		Count:     0,
	})

	if err != nil {
		return s.config.ErrorHandler(err)
	}

	return nil
}

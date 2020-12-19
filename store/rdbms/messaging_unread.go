package rdbms

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/slice"
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
			PlaceholderFormat(s.config.PlaceholderFormat).
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
			PlaceholderFormat(s.config.PlaceholderFormat).
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

	uu, err = s.QueryMessagingUnreads(ctx, q, nil)
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

	return s.execUpdateMessagingUnreads(ctx, cnd, set)
}

// PresetMessagingUnread presets channel unread records for all users (and threads in that channel)
//
// Whenever channel member is added or a new thread is created
// we generate records
func (s Store) PresetMessagingUnread(ctx context.Context, channelID, threadID uint64, userIDs ...uint64) (err error) {
	if channelID == 0 {
		return
	}

	var (
		set types.UnreadSet
	)

	{
		var (
			// map all user-ids to uint64-bool map
			// when we get the actual list of existing entries,
			// we'll remove them from the map and use the existing
			// entries to re-insert them back
			uuIDs = slice.ToUint64BoolMap(userIDs)

			q = s.messagingUnreadsSelectBuilder().Where(squirrel.Eq{
				"rel_channel":  channelID,
				"rel_reply_to": threadID,
				"rel_user":     userIDs,
			})
		)

		if set, err = s.QueryMessagingUnreads(ctx, q, nil); err != nil {
			return
		}

		_ = set.Walk(func(u *types.Unread) error {
			delete(uuIDs, u.UserID)
			return nil
		})

		for userID := range uuIDs {
			if userID == 0 {
				continue
			}

			err = s.CreateMessagingUnread(ctx, &types.Unread{
				ChannelID: channelID,
				ReplyTo:   threadID,
				UserID:    userID,
			})

			if err != nil {
				return
			}
		}
	}

	if threadID == 0 {
		// @todo can we find a way around this presets?!
		//const (
		//	sqlUnreadPresetThreads = `INSERT IGNORE INTO messaging_unread (rel_channel, rel_reply_to, rel_user)
		//            SELECT rel_channel, id, ?
		//              FROM messaging_message
		//             WHERE rel_channel = ?
		//               AND replies > 0`
		//)
		//
		//var (
		//	// fresh map used for creating presets for
		//	// unread messages for threads
		//	uuIDs = slice.ToUint64BoolMap(userIDs)
		//
		//	q = s.messagingMessagesSelectBuilder().
		//		LeftJoin(s.messagingUnreadTable("mur")).
		//		Where(squirrel.And{
		//			// Select 1st messages in thread in a spec. channel
		//			squirrel.Eq{"msg.rel_channel": channelID},
		//			squirrel.Gt{"msg.replies": 0},
		//
		//			// Select unreads in the channel that
		//			squirrel.Eq{"mur.rel_channel": channelID},
		//		})
		//
		//)

	}

	return
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

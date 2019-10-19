package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	// UnreadRepository interface to channel member repository
	UnreadRepository interface {
		With(ctx context.Context, db *factory.DB) UnreadRepository

		Count(userID, channelID uint64, threadIDs ...uint64) (types.UnreadSet, error)
		CountThreads(userID, channelID uint64) (types.UnreadSet, error)

		Preset(channelID, threadID uint64, userIDs ...uint64) (err error)
		Record(userID, channelID, threadID, lastReadMessageID uint64, count uint32) error
		Inc(channelID, replyTo, userID uint64) error
		Dec(channelID, replyTo, userID uint64) error
		ClearThreads(channelID, userID uint64) (err error)
	}

	unread struct {
		*repository
	}
)

const (
	sqlResetThreads = `UPDATE messaging_unread
                          SET count = 0
                        WHERE rel_reply_to > 0 AND rel_channel = ? AND rel_user = ?`

	sqlUnreadIncCount = `UPDATE messaging_unread 
                                  SET count = count + 1
                                WHERE rel_channel = ? AND rel_reply_to = ? AND rel_user <> ?`

	sqlUnreadDecCount = `UPDATE messaging_unread 
                                  SET count = count - 1
                                WHERE rel_channel = ? AND rel_reply_to = ? AND count > 0`

	sqlResetCount = `REPLACE INTO messaging_unread (rel_channel, rel_reply_to, rel_user, count) VALUES (?, ?, ?, 0)`

	sqlUnreadPresetChannel = `INSERT IGNORE INTO messaging_unread (rel_channel, rel_reply_to, rel_user) VALUES (?, ?, ?)`
	sqlUnreadPresetThreads = `INSERT IGNORE INTO messaging_unread (rel_channel, rel_reply_to, rel_user) 
                     SELECT rel_channel, id, ? 
                       FROM messaging_message 
                      WHERE rel_channel = ?
                        AND replies > 0`
)

// Unread creates new instance of channel member repository
func Unread(ctx context.Context, db *factory.DB) UnreadRepository {
	return (&unread{}).With(ctx, db)
}

func (r unread) table() string {
	return "messaging_unread"
}

// With context...
func (r *unread) With(ctx context.Context, db *factory.DB) UnreadRepository {
	return &unread{
		repository: r.repository.With(ctx, db),
	}
}

// Count returns counts unread channel info
func (r *unread) Count(userID, channelID uint64, threadIDs ...uint64) (types.UnreadSet, error) {
	var (
		uu = types.UnreadSet{}
		q  = squirrel.
			Select(
				"rel_channel",
				"rel_last_message",
				"rel_user",
				"rel_reply_to",
				"count",
			).
			From(r.table())
	)

	if userID > 0 {
		q = q.Where("rel_user = ?", userID)
	}

	if channelID > 0 {
		q = q.Where("rel_channel = ?", channelID)
	}

	if len(threadIDs) == 0 {
		q = q.Where("rel_reply_to = 0")
	} else {
		q = q.Where(squirrel.Eq{"rel_reply_to": threadIDs})
	}

	return uu, rh.FetchAll(r.db(), q, &uu)
}

// CountReplies counts unread thread info
func (r unread) CountThreads(userID, channelID uint64) (types.UnreadSet, error) {
	type (
		u struct {
			Rel_channel, Rel_user uint64
			Total, Count          uint32
		}
	)
	var (
		err error

		uu = types.UnreadSet{}

		temp = []*u{}

		q = squirrel.
			Select(
				"rel_channel",
				"rel_user",
				"sum(count) AS count",
				"sum(CASE WHEN count > 0 THEN 1 ELSE 0 END) AS total",
			).
			From(r.table()).
			Where("rel_reply_to > 0 AND count > 0").
			GroupBy("rel_channel", "rel_user")
	)

	if userID > 0 {
		q = q.Where("rel_user = ?", userID)
	}

	if channelID > 0 {
		q = q.Where("rel_channel = ?", channelID)
	}

	err = rh.FetchAll(r.db(), q, &temp)
	if err != nil {
		return nil, err
	}

	for _, t := range temp {
		uu = append(uu, &types.Unread{
			ChannelID:   t.Rel_channel,
			UserID:      t.Rel_user,
			ThreadCount: t.Count,
			ThreadTotal: t.Total,
		})
	}

	return uu, nil
}

func (r unread) ClearThreads(channelID, userID uint64) (err error) {
	_, err = r.db().Exec(sqlResetThreads, channelID, userID)
	return
}

// Preset channel unread records for all users (and threads in that channel)
//
// Whenever channel member is added or a new thread is created
// we generate records
func (r unread) Preset(channelID, threadID uint64, userIDs ...uint64) (err error) {
	if channelID == 0 {
		return
	}

	for _, userID := range userIDs {
		if userID == 0 {
			continue
		}

		_, err = r.db().Exec(sqlUnreadPresetChannel, channelID, threadID, userID)

		if err != nil {
			return
		}

		if threadID == 0 {
			// Preset for all threads in the channel
			_, err = r.db().Exec(sqlUnreadPresetThreads, userID, channelID)

			if err != nil {
				return
			}
		}
	}

	return
}

// Record channel/thread view
func (r *unread) Record(userID, channelID, threadID, lastReadMessageID uint64, count uint32) error {
	mod := &types.Unread{
		ChannelID:     channelID,
		UserID:        userID,
		ReplyTo:       threadID,
		LastMessageID: lastReadMessageID,
		Count:         count,
	}

	return r.db().Replace("messaging_unread", mod)
}

// Inc increments unread message count on a channel/thread for all but one user
func (r *unread) Inc(channelID, threadID, userID uint64) (err error) {
	_, err = r.db().Exec(sqlUnreadIncCount, channelID, threadID, userID)
	if err != nil {
		return err
	}

	return nil
}

// Dec decrements unread message count on a channel/thread for all but one user
func (r *unread) Dec(channelID, threadID, userID uint64) (err error) {
	_, err = r.db().Exec(sqlUnreadDecCount, channelID, threadID)
	if err != nil {
		return err
	}
	_, err = r.db().Exec(sqlResetCount, channelID, threadID, userID)
	if err != nil {
		return err
	}

	return nil
}

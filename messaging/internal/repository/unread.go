package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	// UnreadRepository interface to channel member repository
	UnreadRepository interface {
		With(ctx context.Context, db *factory.DB) UnreadRepository

		Find(filter *types.UnreadFilter) (types.UnreadSet, error)
		Preset(channelID, threadID uint64, userIDs ...uint64) (err error)
		Record(userID, channelID, threadID, lastReadMessageID uint64, count uint32) error
		Inc(channelID, replyTo, userID uint64) error
		Dec(channelID, replyTo, userID uint64) error

		CountOwned(userID uint64) (c int, err error)
		ChangeOwner(userID, target uint64) error
	}

	unread struct {
		*repository
	}
)

const (
	// Fetching channel members of all channels a specific user has access to
	sqlUnreadSelect = `SELECT rel_channel, rel_reply_to, rel_user, count, rel_last_message 
                         FROM messaging_unread
                        WHERE count > 0 && rel_last_message > 0 `

	// Fetching channel members of all channels a specific user has access to
	sqlThreadUnreadSelect = `SELECT rel_channel, sum(count) as count 
                               FROM messaging_unread
                              WHERE rel_user = ? AND rel_reply_to > 0
                           GROUP BY rel_channel`

	sqlUnreadIncCount = `UPDATE messaging_unread 
                                  SET count = count + 1
                                WHERE rel_channel = ? AND rel_reply_to = ? AND rel_user <> ?`

	sqlUnreadDecCount = `UPDATE messaging_unread 
                                  SET count = count - 1
                                WHERE rel_channel = ? AND rel_reply_to = ? AND count > 0`

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

// With context...
func (r *unread) With(ctx context.Context, db *factory.DB) UnreadRepository {
	return &unread{
		repository: r.repository.With(ctx, db),
	}
}

// Find unread info
func (r *unread) Find(filter *types.UnreadFilter) (uu types.UnreadSet, err error) {
	params := make([]interface{}, 0)
	sql := sqlUnreadSelect

	if filter != nil {
		if filter.UserID > 0 {
			// scope: only channel we have access to
			sql += ` AND rel_user = ?`
			params = append(params, filter.UserID)
		}

		if filter.ChannelID > 0 {
			// scope: only channel we have access to
			sql += ` AND rel_channel = ?`
			params = append(params, filter.ChannelID)
		}

		if len(filter.ThreadIDs) > 0 {
			sql += ` AND rel_reply_to IN (?)`
			params = append(params, filter.ThreadIDs)
		} else {
			sql += ` AND rel_reply_to = 0`
		}
	}

	if sql, params, err = sqlx.In(sql, params...); err != nil {
		return nil, err
	} else if err = r.db().Select(&uu, sql, params...); err != nil {
		return nil, err
	} else if len(filter.ThreadIDs) == 0 && filter.UserID > 0 {
		// Check for unread thread messages

		// We'll abuse Unread/UnreadSet
		tt := types.UnreadSet{}

		err = r.db().Select(&tt, sqlThreadUnreadSelect, filter.UserID)

		_ = tt.Walk(func(t *types.Unread) error {
			c := uu.FindByChannelId(t.ChannelID)
			if c != nil {
				c.InThreadCount = t.Count
			} else {
				// No un-reads in channel but we have them in threads (of that channel)
				// swap values and append
				t.InThreadCount, t.Count = t.Count, 0
				uu = append(uu, t)
			}
			return nil
		})
	}

	return uu, nil
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
func (r *unread) Inc(channelID, threadID, userID uint64) error {
	_, err := r.db().Exec(sqlUnreadIncCount, channelID, threadID, userID)
	return err
}

// Dec decrements unread message count on a channel/thread for all but one user
func (r *unread) Dec(channelID, threadID, userID uint64) error {
	_, err := r.db().Exec(sqlUnreadDecCount, channelID, threadID)
	return err
}

func (r *unread) CountOwned(userID uint64) (c int, err error) {
	return c, r.db().Get(&c,
		"SELECT COUNT(*) FROM messaging_unread WHERE rel_user = ?",
		userID)
}

func (r *unread) ChangeOwner(userID, target uint64) (err error) {
	// Remove dups
	// with an ugly mysql workaround
	_, err = r.db().Exec(
		"DELETE FROM messaging_unread WHERE rel_user = ? "+
			"AND rel_channel IN (SELECT rel_channel FROM (SELECT * FROM messaging_unread) AS workaround WHERE rel_user = ?)",
		userID,
		target)

	if err != nil {
		return err
	}

	_, err = r.db().Exec(
		"UPDATE messaging_unread SET rel_user = ? WHERE rel_user = ?",
		target,
		userID)

	return err
}

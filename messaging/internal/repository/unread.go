package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	// UnreadRepository interface to channel member repository
	UnreadRepository interface {
		With(ctx context.Context, db *factory.DB) UnreadRepository

		Find(filter *types.UnreadFilter) (types.UnreadSet, error)
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
                              WHERE true `

	sqlUnreadIncCount = `UPDATE messaging_unread 
                                  SET count = count + 1
                                WHERE rel_channel = ? AND rel_reply_to = ? AND rel_user <> ?`

	sqlUnreadDecCount = `UPDATE messaging_unread 
                                  SET count = count - 1
                                WHERE rel_channel = ? AND rel_reply_to = ? AND rel_user <> ? AND count > 0`
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

// FindMembers fetches membership info
//
// If channelID > 0 it returns members of a specific channel
// If userID    > 0 it returns members of all channels this user is member of
func (r *unread) Find(filter *types.UnreadFilter) (types.UnreadSet, error) {
	params := make([]interface{}, 0)
	vv := types.UnreadSet{}
	sql := sqlUnreadSelect

	if filter != nil {
		if filter.UserID > 0 {
			// scope: only channel we have access to
			sql += ` AND rel_user = ?`
			params = append(params, filter.UserID)
		}
	}

	return vv, r.db().Select(&vv, sql, params...)
}

// Records channel view
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

// Increments unread message count on a channel/thread for all but one user
func (r *unread) Inc(channelID, threadID, userID uint64) error {
	_, err := r.db().Exec(sqlUnreadIncCount, channelID, threadID, userID)
	return err
}

// Decrements unread message count on a channel/thread for all but one user
func (r *unread) Dec(channelID, threadID, userID uint64) error {
	_, err := r.db().Exec(sqlUnreadDecCount, channelID, threadID, userID)
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

package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/types"
)

type (
	// UnreadRepository interface to channel member repository
	UnreadRepository interface {
		With(ctx context.Context, db *factory.DB) UnreadRepository

		Find(filter *types.UnreadFilter) (types.UnreadSet, error)
		Record(userID, channelID, threadID, lastReadMessageID uint64, count uint32) error
		Inc(channelID, replyTo, userID uint64) error
		Dec(channelID, replyTo, userID uint64) error
	}

	unread struct {
		*repository
	}
)

const (
	// Fetching channel members of all channels a specific user has access to
	sqlUnreadSelect = `SELECT rel_channel, rel_reply_to, rel_user, count, rel_last_message 
                               FROM unreads
                              WHERE true `

	sqlUnreadIncCount = `UPDATE unreads 
                                  SET count = count + 1
                                WHERE rel_channel = ? AND rel_reply_to = ? AND rel_user <> ?`

	sqlUnreadDecCount = `UPDATE unreads 
                                  SET count = count - 1
                                WHERE rel_channel = ? AND rel_reply_to = ? AND rel_user <> ? AND count > 0`
)

// ChannelView creates new instance of channel member repository
func ChannelView(ctx context.Context, db *factory.DB) UnreadRepository {
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

	return r.db().Replace("unreads", mod)
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

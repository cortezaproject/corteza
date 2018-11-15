package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/types"
)

type (
	// ChannelViewRepository interface to channel member repository
	ChannelViewRepository interface {
		With(ctx context.Context, db *factory.DB) ChannelViewRepository

		Find(filter *types.UnreadFilter) (types.UnreadSet, error)
		Record(userID, channelID, replyTo, lastMessageID uint64, count uint32) error
		Inc(channelID, userID uint64) error
		Dec(channelID, userID uint64) error
	}

	channelViews struct {
		*repository
	}
)

const (
	// Fetching channel members of all channels a specific user has access to
	sqlChannelViewsSelect = `SELECT rel_channel, rel_user, count, rel_last_message 
                               FROM unreads
                              WHERE true `

	sqlChannelViewsIncCount = `UPDATE unreads 
                                  SET count = count + 1
                                WHERE rel_channel = ? AND rel_user <> ?`

	sqlChannelViewsDecCount = `UPDATE unreads 
                                  SET count = count - 1
                                WHERE rel_channel = ? AND rel_user <> ? AND count > 0`
)

// ChannelView creates new instance of channel member repository
func ChannelView(ctx context.Context, db *factory.DB) ChannelViewRepository {
	return (&channelViews{}).With(ctx, db)
}

// With context...
func (r *channelViews) With(ctx context.Context, db *factory.DB) ChannelViewRepository {
	return &channelViews{
		repository: r.repository.With(ctx, db),
	}
}

// FindMembers fetches membership info
//
// If channelID > 0 it returns members of a specific channel
// If userID    > 0 it returns members of all channels this user is member of
func (r *channelViews) Find(filter *types.UnreadFilter) (types.UnreadSet, error) {
	params := make([]interface{}, 0)
	vv := types.UnreadSet{}
	sql := sqlChannelViewsSelect

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
func (r *channelViews) Record(userID, channelID, replyTo, lastMessageID uint64, count uint32) error {
	mod := &types.Unread{
		ChannelID:     channelID,
		UserID:        userID,
		ReplyTo:       replyTo,
		LastMessageID: lastMessageID,
		Count:         count,
	}

	return r.db().Replace("unreads", mod)
}

// Increments unread (new) message count on a channel for all but one user
func (r *channelViews) Inc(channelID, userID uint64) error {
	_, err := r.db().Exec(sqlChannelViewsIncCount, channelID, userID)
	return err
}

// Increments unread (new) message count on a channel for all but one user
func (r *channelViews) Dec(channelID, userID uint64) error {
	_, err := r.db().Exec(sqlChannelViewsDecCount, channelID, userID)
	return err
}

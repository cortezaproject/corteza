package repository

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	MessageRepository interface {
		With(ctx context.Context, db *factory.DB) MessageRepository

		FindByID(id uint64) (*types.Message, error)
		Find(filter *types.MessageFilter) (types.MessageSet, error)
		FindThreads(filter *types.MessageFilter) (types.MessageSet, error)
		CountFromMessageID(channelID, threadID, messageID uint64) (uint32, error)
		LastMessageID(channelID, threadID uint64) (uint64, error)
		PrefillThreadParticipants(mm types.MessageSet) error

		Create(mod *types.Message) (*types.Message, error)
		Update(mod *types.Message) (*types.Message, error)
		DeleteByID(ID uint64) error

		BindAvatar(message *types.Message, avatar io.Reader) (*types.Message, error)

		IncReplyCount(ID uint64) error
		DecReplyCount(ID uint64) error
	}

	message struct {
		*repository
	}
)

const (
	MESSAGES_MAX_LIMIT = 100

	// subquery that filters out all channels that current user has access to as a member
	// or via channel type (public channels)
	sqlChannelAccess = ` (
				SELECT id
                  FROM messaging_channel c
                       LEFT OUTER JOIN messaging_channel_member AS m ON (c.id = m.rel_channel)
                 WHERE rel_user = ?
              UNION
                SELECT id
                  FROM messaging_channel c
                 WHERE c.type = ?
			)`

	sqlMessageColumns = "id, " +
		"COALESCE(type,'') AS type, " +
		"message, " +
		"rel_user, " +
		"rel_channel, " +
		"reply_to, " +
		"replies, " +
		"created_at, " +
		"updated_at, " +
		"deleted_at"
	sqlMessageScope = "deleted_at IS NULL"

	sqlMessagesSelect = `SELECT ` + sqlMessageColumns + `
        FROM messaging_message
       WHERE ` + sqlMessageScope

	sqlMessagesThreads = "WITH originals AS (" +
		" SELECT id AS original_id " +
		"   FROM messaging_message " +
		"  WHERE " + sqlMessageScope +
		"    AND rel_channel IN " + sqlChannelAccess +
		"    AND reply_to = 0 " +
		"    AND replies > 0 " +
		// for finding only threads we've created or replied to
		"    AND (rel_user = ? OR id IN (SELECT DISTINCT reply_to FROM messaging_message WHERE rel_user = ?))" +
		"  ORDER BY id DESC " +
		"  LIMIT ? " +
		")" +
		" SELECT " + sqlMessageColumns +
		"   FROM messaging_message, originals " +
		"  WHERE " + sqlMessageScope +
		"    AND original_id IN (id, reply_to)"

	sqlThreadParticipantsByMessageID = "SELECT DISTINCT reply_to, rel_user FROM messaging_message WHERE reply_to IN (?)"

	sqlCountFromMessageID = "SELECT COUNT(*) AS count " +
		"FROM messaging_message " +
		"WHERE rel_channel = ? " +
		"AND reply_to = ? " +
		"AND COALESCE(type, '') NOT IN (?) " +
		"AND id > ? AND deleted_at IS NULL"

	sqlLastMessageID = "SELECT COALESCE(MAX(id), 0) AS last " +
		"FROM messaging_message " +
		"WHERE rel_channel = ? " +
		"AND reply_to = ? " +
		"AND COALESCE(type, '') NOT IN (?) " +
		"AND deleted_at IS NULL"

	sqlMessageRepliesIncCount = `UPDATE messaging_message SET replies = replies + 1 WHERE id = ? AND reply_to = 0`
	sqlMessageRepliesDecCount = `UPDATE messaging_message SET replies = replies - 1 WHERE id = ? AND reply_to = 0`

	ErrMessageNotFound = repositoryError("MessageNotFound")
)

func Message(ctx context.Context, db *factory.DB) MessageRepository {
	return (&message{}).With(ctx, db)
}

func (r *message) With(ctx context.Context, db *factory.DB) MessageRepository {
	return &message{
		repository: r.repository.With(ctx, db),
	}
}

func (r message) table() string {
	return "messaging_message"
}

func (r *message) FindByID(id uint64) (*types.Message, error) {
	mod := &types.Message{}
	sql := sqlMessagesSelect + " AND id = ?"

	return mod, rh.IsFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrMessageNotFound)
}

func (r *message) Find(filter *types.MessageFilter) (types.MessageSet, error) {
	r.sanitizeFilter(filter)

	params := make([]interface{}, 0)
	rval := make(types.MessageSet, 0)

	sql := sqlMessagesSelect

	if filter.Query != "" {
		sql += " AND message LIKE ?"
		params = append(params, "%"+filter.Query+"%")
	}

	if len(filter.ChannelID) > 0 {
		sql += " AND rel_channel IN (" + strings.Repeat(",?", len(filter.ChannelID))[1:] + ")"
		for _, id := range filter.ChannelID {
			params = append(params, id)
		}
	}

	if len(filter.UserID) > 0 {
		sql += " AND rel_user IN (" + strings.Repeat(",?", len(filter.UserID))[1:] + ")"
		for _, id := range filter.UserID {
			params = append(params, id)
		}
	}

	if len(filter.ThreadID) > 0 {
		sql += " AND reply_to IN (" + strings.Repeat(",?", len(filter.ThreadID))[1:] + ")"
		for _, id := range filter.ThreadID {
			params = append(params, id)
		}
	} else {
		sql += " AND reply_to = 0 "
	}

	if len(filter.Type) > 0 {
		sql += " AND type IN (" + strings.Repeat(",?", len(filter.Type))[1:] + ")"
		for _, id := range filter.Type {
			params = append(params, id)
		}
	}

	// first, exclusive
	if filter.AfterID > 0 {
		sql += " AND id > ? "
		params = append(params, filter.AfterID)
	}

	// from, inclusive
	if filter.FromID > 0 {
		sql += " AND id >= ? "
		params = append(params, filter.FromID)
	}

	// last, exclusive
	if filter.BeforeID > 0 {
		sql += " AND id < ? "
		params = append(params, filter.BeforeID)
	}

	// to, inclusive
	if filter.ToID > 0 {
		sql += " AND id <= ? "
		params = append(params, filter.ToID)
	}

	if filter.BookmarkedOnly || filter.PinnedOnly {
		sql += " AND id IN (SELECT rel_message FROM messaging_message_flag WHERE flag = ?) "

		if filter.PinnedOnly {
			params = append(params, types.MessageFlagPinnedToChannel)
		} else {
			params = append(params, types.MessageFlagBookmarkedMessage)
		}
	}

	if filter.AttachmentsOnly {
		sql += " AND type IN (?, ?) "
		params = append(
			params,
			types.MessageTypeAttachment,
			types.MessageTypeInlineImage,
		)
	}

	sql += " AND rel_channel IN " + sqlChannelAccess
	params = append(params, filter.CurrentUserID, types.ChannelTypePublic)

	sql += " ORDER BY id DESC"

	sql += " LIMIT ? "
	params = append(params, filter.Limit)

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *message) FindThreads(filter *types.MessageFilter) (types.MessageSet, error) {
	r.sanitizeFilter(filter)

	params := make([]interface{}, 0)
	rval := make(types.MessageSet, 0)

	// for sqlChannelAccess
	params = append(params, filter.CurrentUserID, types.ChannelTypePublic)

	// for finding only threads we've created or replied to
	params = append(params, filter.CurrentUserID, filter.CurrentUserID)

	// for sqlMessagesThreads
	params = append(params, filter.Limit)

	sql := sqlMessagesThreads

	if len(filter.ChannelID) > 0 {
		sql += " AND rel_channel IN (" + strings.Repeat(",?", len(filter.ChannelID))[1:] + ")"
		for _, id := range filter.ChannelID {
			params = append(params, id)
		}
	}

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *message) CountFromMessageID(channelID, threadID, lastReadMessageID uint64) (uint32, error) {
	if lastReadMessageID == 0 {
		// No need for counting, zero unread messages...
		return 0, nil
	}
	rval := struct{ Count uint32 }{}
	return rval.Count, r.db().Get(&rval,
		sqlCountFromMessageID,
		channelID,
		threadID,
		types.MessageTypeChannelEvent,
		lastReadMessageID,
	)
}

func (r *message) LastMessageID(channelID, threadID uint64) (uint64, error) {
	rval := struct{ Last uint64 }{}
	return rval.Last, r.db().Get(&rval,
		sqlLastMessageID,
		channelID,
		threadID,
		types.MessageTypeChannelEvent,
	)
}

func (r *message) PrefillThreadParticipants(mm types.MessageSet) error {
	var rval = []struct {
		ReplyTo uint64 `db:"reply_to"`
		UserID  uint64 `db:"rel_user"`
	}{}

	// Filter out only relevant messages -- ones with replies
	mm, _ = mm.Filter(func(m *types.Message) (b bool, e error) {
		return m.Replies > 0, nil
	})

	if len(mm) == 0 {
		return nil
	}

	if sql, args, err := sqlx.In(sqlThreadParticipantsByMessageID, mm.IDs()); err != nil {
		return err
	} else if err = r.db().Select(&rval, sql, args...); err != nil {
		return err
	} else {
		for _, p := range rval {
			mm.FindByID(p.ReplyTo).RepliesFrom = append(mm.FindByID(p.ReplyTo).RepliesFrom, p.UserID)
		}
	}

	return nil
}

func (r *message) sanitizeFilter(filter *types.MessageFilter) {
	if filter == nil {
		filter = &types.MessageFilter{}
	}

	if filter.Limit == 0 || filter.Limit > MESSAGES_MAX_LIMIT {
		filter.Limit = MESSAGES_MAX_LIMIT
	}
}

func (r *message) Create(mod *types.Message) (*types.Message, error) {
	mod.ID = factory.Sonyflake.NextID()
	rh.SetCurrentTimeRounded(&mod.CreatedAt)

	return mod, r.db().Insert("messaging_message", mod)
}

func (r *message) Update(mod *types.Message) (*types.Message, error) {
	rh.SetCurrentTimeRounded(&mod.UpdatedAt)

	return mod, r.db().Replace("messaging_message", mod)
}

func (svc *message) BindAvatar(in *types.Message, avatar io.Reader) (*types.Message, error) {
	// @todo: implement setting avatar on a message
	in.Meta.Avatar = ""
	return in, nil
}

func (r *message) DeleteByID(ID uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"deleted_at": time.Now()}, squirrel.Eq{"id": ID})
}

func (r *message) IncReplyCount(ID uint64) error {
	_, err := r.db().Exec(sqlMessageRepliesIncCount, ID)
	return err
}

func (r *message) DecReplyCount(ID uint64) error {
	_, err := r.db().Exec(sqlMessageRepliesDecCount, ID)
	return err
}

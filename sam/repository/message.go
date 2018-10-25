package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/types"
)

type (
	MessageRepository interface {
		With(ctx context.Context, db *factory.DB) MessageRepository

		FindMessageByID(id uint64) (*types.Message, error)
		FindMessages(filter *types.MessageFilter) (types.MessageSet, error)
		CreateMessage(mod *types.Message) (*types.Message, error)
		UpdateMessage(mod *types.Message) (*types.Message, error)
		DeleteMessageByID(ID uint64) error
		IncReplyCount(ID uint64) error
		DecReplyCount(ID uint64) error
	}

	message struct {
		*repository
	}
)

const (
	MESSAGES_MAX_LIMIT = 100

	sqlMessageScope = "deleted_at IS NULL"

	sqlMessagesSelect = `SELECT id,
             COALESCE(type,'') AS type,
             message,
             rel_user,
             rel_channel,
             reply_to,
             replies,
             created_at,
             updated_at,
             deleted_at
        FROM messages
       WHERE ` + sqlMessageScope

	sqlMessageRepliesIncCount = `UPDATE messages SET replies = replies + 1 WHERE id = ? AND reply_to = 0`
	sqlMessageRepliesDecCount = `UPDATE messages SET replies = replies - 1 WHERE id = ? AND reply_to = 0`

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

func (r *message) FindMessageByID(id uint64) (*types.Message, error) {
	mod := &types.Message{}
	sql := sqlMessagesSelect + " AND id = ?"

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrMessageNotFound)
}

func (r *message) FindMessages(filter *types.MessageFilter) (types.MessageSet, error) {
	params := make([]interface{}, 0)
	rval := make(types.MessageSet, 0)

	sql := sqlMessagesSelect

	if filter != nil {
		if filter.Query != "" {
			sql += " AND message LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	if filter.ChannelID == 0 && filter.RepliesTo == 0 {
		// Channel history or replies to a message...
		// nothing more.
		return nil, nil
	}

	if filter.ChannelID > 0 {
		sql += " AND rel_channel = ? "
		params = append(params, filter.ChannelID)
	}

	if filter.RepliesTo > 0 {
		sql += " AND reply_to = ? "
		params = append(params, filter.RepliesTo)
	} else {
		sql += " AND reply_to = 0 "
	}

	if filter.FirstID > 0 || filter.LastID > 0 {
		// Fetching (exclusively) range of messages, without reply
		if filter.FirstID > 0 {
			sql += " AND id > ? "
			params = append(params, filter.FirstID)
		}

		if filter.LastID > 0 {
			sql += " AND id < ? "
			params = append(params, filter.LastID)
		}
	}

	sql += " ORDER BY id DESC"

	if filter.Limit == 0 || filter.Limit > MESSAGES_MAX_LIMIT {
		filter.Limit = MESSAGES_MAX_LIMIT
	}

	sql += " LIMIT ? "
	params = append(params, filter.Limit)

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *message) CreateMessage(mod *types.Message) (*types.Message, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert("messages", mod)
}

func (r *message) UpdateMessage(mod *types.Message) (*types.Message, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, r.db().Replace("messages", mod)
}

func (r *message) DeleteMessageByID(ID uint64) error {
	return r.updateColumnByID("messages", "deleted_at", time.Now(), ID)
}

func (r *message) IncReplyCount(ID uint64) error {
	_, err := r.db().Exec(sqlMessageRepliesIncCount, ID)
	return err
}

func (r *message) DecReplyCount(ID uint64) error {
	_, err := r.db().Exec(sqlMessageRepliesDecCount, ID)
	return err
}

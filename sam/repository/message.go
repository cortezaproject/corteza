package repository

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

type (
	Message interface {
		FindMessageByID(id uint64) (*types.Message, error)
		FindMessages(filter *types.MessageFilter) (types.MessageSet, error)
		CreateMessage(mod *types.Message) (*types.Message, error)
		UpdateMessage(mod *types.Message) (*types.Message, error)
		DeleteMessageByID(id uint64) error
	}
)

const (
	sqlMessageScope = "deleted_at IS NULL"

	sqlMessagesSelect = `SELECT id,
             COALESCE(type,'') AS type,
             message,
             rel_user,
             rel_channel,
             COALESCE(reply_to, 0) AS reply_to,
             created_at,
             updated_at,
             deleted_at
        FROM messages
       WHERE ` + sqlMessageScope

	ErrMessageNotFound = repositoryError("MessageNotFound")
)

func (r *repository) FindMessageByID(id uint64) (*types.Message, error) {
	mod := &types.Message{}
	sql := sqlMessagesSelect + " AND id = ?"

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrMessageNotFound)
}

func (r *repository) FindMessages(filter *types.MessageFilter) (types.MessageSet, error) {
	params := make([]interface{}, 0)
	rval := make(types.MessageSet, 0)

	sql := sqlMessagesSelect

	if filter != nil {
		if filter.Query != "" {
			sql += " AND message LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	if filter.ChannelID > 0 {
		sql += " AND rel_channel = ? "
		params = append(params, filter.ChannelID)
	}

	if filter.FromMessageID > 0 {
		sql += " AND id > ? "
		params = append(params, filter.FromMessageID)
	}

	if filter.UntilMessageID > 0 {
		sql += " AND id < ? "
		params = append(params, filter.UntilMessageID)
	}

	sql += " ORDER BY id ASC"

	if filter.Limit > 0 {
		// @todo implement some kind of protection
		sql += " LIMIT ? "
		params = append(params, filter.Limit)
	}
	return rval, r.db().Select(&rval, sql, params...)
}

func (r *repository) CreateMessage(mod *types.Message) (*types.Message, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert("messages", mod)
}

func (r *repository) UpdateMessage(mod *types.Message) (*types.Message, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, r.db().Replace("messages", mod)
}

func (r *repository) DeleteMessageByID(id uint64) error {
	return r.updateColumnByID("messages", "deleted_at", nil, id)
}

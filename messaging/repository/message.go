package repository

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	MessageRepository interface {
		With(ctx context.Context, db *factory.DB) MessageRepository

		FindByID(id uint64) (*types.Message, error)
		Find(types.MessageFilter) (types.MessageSet, types.MessageFilter, error)
		FindThreads(types.MessageFilter) (types.MessageSet, types.MessageFilter, error)
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

func (r message) columns() []string {
	return []string{
		"m.id",
		"COALESCE(m.type,'') AS type",
		"m.message",
		"m.rel_user",
		"m.rel_channel",
		"m.reply_to",
		"m.replies",
		"m.created_at",
		"m.updated_at",
		"m.deleted_at",
	}
}

func (r message) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table() + " AS m").
		Where(squirrel.Eq{"m.deleted_at": nil})
}

func (r message) FindByID(id uint64) (*types.Message, error) {
	return r.findOneBy(squirrel.Eq{"m.id": id})
}

func (r message) findOneBy(cnd squirrel.Sqlizer) (*types.Message, error) {
	var (
		ch = &types.Message{}

		q = r.query().
			Where(cnd)

		err = rh.FetchOne(r.db(), q, ch)
	)

	if err != nil {
		return nil, err
	} else if ch.ID == 0 {
		return nil, ErrMessageNotFound
	}

	return ch, nil
}

func (r message) Find(filter types.MessageFilter) (set types.MessageSet, f types.MessageFilter, err error) {
	f = r.sanitizeFilter(filter)

	query := r.query()

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Like{"LOWER(m.message)": q})
	}

	if len(f.ChannelID) > 0 {
		query = query.Where(squirrel.Eq{"m.rel_channel": f.ChannelID})
	}

	if len(f.UserID) > 0 {
		query = query.Where(squirrel.Eq{"m.rel_user": f.UserID})
	}

	if len(f.ThreadID) > 0 {
		query = query.Where(squirrel.Eq{"m.reply_to": f.ThreadID})
	} else {
		query = query.Where(squirrel.Eq{"m.reply_to": 0})
	}

	if f.AttachmentsOnly {
		// Override Type filter
		f.Type = []string{
			types.MessageTypeAttachment.String(),
			types.MessageTypeInlineImage.String(),
		}
	}

	if len(f.Type) > 0 {
		query = query.Where(squirrel.Eq{"m.type": f.Type})
	}

	// first, exclusive
	if f.AfterID > 0 {
		query = query.OrderBy("m.id ASC")
		query = query.Where(squirrel.Gt{"m.id": f.AfterID})
	}

	// from, inclusive
	if f.FromID > 0 {
		query = query.OrderBy("m.id ASC")
		query = query.Where(squirrel.GtOrEq{"m.id": f.FromID})
	}

	// last, exclusive
	if f.BeforeID > 0 {
		query = query.OrderBy("m.id DESC")
		query = query.Where(squirrel.Lt{"m.id": f.BeforeID})
	}

	// to, inclusive
	if f.ToID > 0 {
		query = query.OrderBy("m.id DESC")
		query = query.Where(squirrel.LtOrEq{"m.id": f.ToID})
	}

	if f.BookmarkedOnly || f.PinnedOnly {
		flag := types.MessageFlagBookmarkedMessage
		if f.PinnedOnly {
			flag = types.MessageFlagPinnedToChannel
		}

		query = query.
			Where(squirrel.ConcatExpr("m.id IN(", (messageFlag{}).queryMessagesWithFlags(flag), ")"))
	}

	query = query.
		OrderBy("id DESC").
		Limit(uint64(f.Limit))

	return set, f, rh.FetchAll(r.db(), query, &set)
}

func (r *message) FindThreads(filter types.MessageFilter) (set types.MessageSet, f types.MessageFilter, err error) {
	f = r.sanitizeFilter(filter)

	// Selecting first valid (deleted_at IS NULL) messages in threads (replies > 0 && reply_to = 0)
	//   that belong to filtered channels and we've contributed to (or stated it)
	originals := squirrel.
		Select("id AS original_id").
		From(r.table()).
		Where(squirrel.And{
			squirrel.Eq{
				"deleted_at":  nil,
				"rel_channel": f.ChannelID,
				"reply_to":    0,
			},
			squirrel.Gt{"replies": 0},
			squirrel.Or{
				squirrel.Eq{"rel_user": filter.CurrentUserID},
				squirrel.Expr(
					"id IN (SELECT DISTINCT reply_to FROM messaging_message WHERE rel_user = ?)",
					filter.CurrentUserID),
			},
		}).
		OrderBy("id DESC").
		Limit(uint64(f.Limit))

	// Prepare the actual message selector
	query := r.query().Join("originals ON (original_id IN (id, reply_to))")

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Like{"LOWER(m.message)": q})
	}

	// And create CTE
	cte := squirrel.ConcatExpr("WITH originals AS (", originals, ") ", query)

	return set, f, rh.FetchAll(r.db(), cte, &set)
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

func (r *message) PrefillThreadParticipants(mm types.MessageSet) (err error) {
	var rval []struct {
		ReplyTo uint64 `db:"reply_to"`
		UserID  uint64 `db:"rel_user"`
	}

	// Filter out only relevant messages -- ones with replies
	mm, _ = mm.Filter(func(m *types.Message) (b bool, e error) {
		return m.Replies > 0, nil
	})

	if len(mm) == 0 {
		return nil
	}

	query := squirrel.
		Select("reply_to", "rel_user").
		From(r.table()).
		Where(squirrel.Eq{"reply_to": mm.IDs()})

	err = rh.FetchAll(r.db(), query, &rval)
	if err != nil {
		return
	}

	for _, p := range rval {
		mm.FindByID(p.ReplyTo).RepliesFrom = append(mm.FindByID(p.ReplyTo).RepliesFrom, p.UserID)
	}

	return nil
}

func (r *message) sanitizeFilter(f types.MessageFilter) types.MessageFilter {
	if f.Limit == 0 || f.Limit > MESSAGES_MAX_LIMIT {
		f.Limit = MESSAGES_MAX_LIMIT
	}

	return f
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

func (r *message) BindAvatar(in *types.Message, avatar io.Reader) (*types.Message, error) {
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

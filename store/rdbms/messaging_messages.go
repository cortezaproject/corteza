package rdbms

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"strings"
)

const (
	messagingMessagesMaxLimit = 100
)

func (s Store) convertMessagingMessageFilter(f types.MessageFilter) (query squirrel.SelectBuilder, err error) {
	if f.Limit == 0 || f.Limit > messagingMessagesMaxLimit {
		f.Limit = messagingMessagesMaxLimit
	}

	query = s.messagingMessagesSelectBuilder()

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
		var (
			flagQuery squirrel.SelectBuilder
			flag      = types.MessageFlagBookmarkedMessage
		)
		if f.PinnedOnly {
			flag = types.MessageFlagPinnedToChannel
		}

		if flagQuery, err = s.convertMessagingFlagFilter(types.MessageFlagFilter{Flag: flag}); err != nil {
			return
		} else {
			query = query.
				Where(squirrel.Eq{"m.id": flagQuery})
		}
	}

	query = query.
		OrderBy("id DESC").
		Limit(uint64(f.Limit))
	return
}

func (s Store) SearchMessagingThreads(ctx context.Context, filter types.MessageFilter) (set types.MessageSet, f types.MessageFilter, err error) {
	if f.Limit == 0 || f.Limit > messagingMessagesMaxLimit {
		f.Limit = messagingMessagesMaxLimit
	}

	// Selecting first valid (deleted_at IS NULL) messages in threads (replies > 0 && reply_to = 0)
	//   that belong to filtered channels and we've contributed to (or stated it)
	originals := squirrel.
		Select("id AS original_id").
		From(s.messagingMessageTable()).
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
	base := s.messagingMessagesSelectBuilder().
		Where(squirrel.Eq{"m.deleted_at": nil}).
		Join("originals ON (original_id IN (id, reply_to))")

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		base = base.Where(squirrel.Like{"LOWER(m.message)": q})
	}

	// Create CTE with originals & base
	cte := squirrel.ConcatExpr("WITH originals AS (", originals, ") ", base)

	if set, _, _, err = s.QueryMessagingMessages(ctx, cte, nil); err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// CountMessagingMessagesFromID returns number of messages from last message on
func (s Store) CountMessagingMessagesFromID(ctx context.Context, channelID, threadID, messageID uint64) (c uint32, err error) {
	if messageID == 0 {
		// No need for counting, zero unread messages...
		return
	}

	var (
		row *sql.Row
		//rval := struct{ Count uint32 }{}
		q = s.SelectBuilder(s.messagingMessageTable("msg"), "COUNT(*)").
			Where(squirrel.NotEq{"type": types.MessageTypeChannelEvent}).
			Where(squirrel.Gt{"id": messageID}).
			Where(squirrel.Eq{
				"rel_channel": channelID,
				"reply_to":    threadID,
				"deleted_at":  nil,
			})
	)

	if row, err = s.QueryRow(ctx, q); err != nil {
		return
	}

	if err = row.Scan(&c); err != nil {
		return
	}

	return
}

// LastMessagingMessageID  returns last message in the channel or thread
func (s Store) LastMessagingMessageID(ctx context.Context, channelID, threadID uint64) (ID uint64, err error) {
	var (
		row *sql.Row
		//rval := struct{ Count uint32 }{}
		q = s.SelectBuilder(s.messagingMessageTable("msg"), "COALESCE(MAX(id), 0)").
			Where(squirrel.NotEq{"type": types.MessageTypeChannelEvent}).
			Where(squirrel.Eq{
				"rel_channel": channelID,
				"reply_to":    threadID,
				"deleted_at":  nil,
			})
	)

	if row, err = s.QueryRow(ctx, q); err != nil {
		return
	}

	if err = row.Scan(&ID); err != nil {
		return
	}

	return
}

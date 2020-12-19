package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/store"
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
	query = query.Where(squirrel.Eq{"msg.deleted_at": nil})

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Like{"LOWER(msg.message)": q})
	}

	if len(f.ChannelID) > 0 {
		query = query.Where(squirrel.Eq{"msg.rel_channel": f.ChannelID})
	}

	if len(f.UserID) > 0 {
		query = query.Where(squirrel.Eq{"msg.rel_user": f.UserID})
	}

	if len(f.ThreadID) > 0 {
		query = query.Where(squirrel.Eq{"msg.reply_to": f.ThreadID})
	} else {
		query = query.Where(squirrel.Eq{"msg.reply_to": 0})
	}

	if f.AttachmentsOnly {
		// Override Type filter
		f.Type = []string{
			types.MessageTypeAttachment.String(),
			types.MessageTypeInlineImage.String(),
		}
	}

	if len(f.Type) > 0 {
		query = query.Where(squirrel.Eq{"msg.type": f.Type})
	}

	// first, exclusive
	if f.AfterID > 0 {
		query = query.OrderBy("msg.id ASC")
		query = query.Where(squirrel.Gt{"msg.id": f.AfterID})
	}

	// from, inclusive
	if f.FromID > 0 {
		query = query.OrderBy("msg.id ASC")
		query = query.Where(squirrel.GtOrEq{"msg.id": f.FromID})
	}

	// last, exclusive
	if f.BeforeID > 0 {
		query = query.OrderBy("msg.id DESC")
		query = query.Where(squirrel.Lt{"msg.id": f.BeforeID})
	}

	// to, inclusive
	if f.ToID > 0 {
		query = query.OrderBy("msg.id DESC")
		query = query.Where(squirrel.LtOrEq{"msg.id": f.ToID})
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

			if fqSql, fqArgs, fqErr := flagQuery.Where("msg.id = mmf.rel_message").ToSql(); err != nil {
				return query, fqErr
			} else {
				query = query.
					Where(fmt.Sprintf("EXISTS (%s)", fqSql), fqArgs...)
			}
		}
	}

	// Manually sorting & limiting for BC
	query = query.
		PlaceholderFormat(s.config.PlaceholderFormat).
		OrderBy("id DESC").
		Limit(uint64(f.Limit))
	return
}

func (s Store) SearchMessagingThreads(ctx context.Context, filter types.MessageFilter) (set types.MessageSet, f types.MessageFilter, err error) {
	if filter.Limit == 0 || filter.Limit > messagingMessagesMaxLimit {
		filter.Limit = messagingMessagesMaxLimit
	}

	// Selecting first valid (deleted_at IS NULL) messages in threads (replies > 0 && reply_to = 0)
	//   that belong to filtered channels and we've contributed to (or stated it)
	originals := squirrel.
		Select("id AS original_id").
		// reset placeholder to question;
		// this will help us a bit lower with the CTE on
		// postgresql (uses $<number> placeholder)
		PlaceholderFormat(squirrel.Question).
		From(s.messagingMessageTable()).
		Where(squirrel.And{
			squirrel.Eq{
				"deleted_at":  nil,
				"rel_channel": filter.ChannelID,
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
		Limit(uint64(filter.Limit))

	// Prepare the actual message selector
	base := s.messagingMessagesSelectBuilder().
		// reset placeholder to question;
		// this will help us a bit lower with the CTE on
		// postgresql (uses $<number> placeholder)
		PlaceholderFormat(squirrel.Question).
		Where(squirrel.Eq{"msg.deleted_at": nil}).
		Join("originals ON (original_id IN (id, reply_to))")

	if filter.Query != "" {
		q := "%" + strings.ToLower(filter.Query) + "%"
		base = base.Where(squirrel.Like{"LOWER(msg.message)": q})
	}

	// Create CTE with originals & base
	cte := SquirrelConcatExpr("WITH originals AS (", originals, ") ", base)
	cte = cte.PlaceholderFormat(s.config.PlaceholderFormat)
	if set, err = s.QueryMessagingMessages(ctx, cte, nil); err != nil {
		return nil, filter, err
	}

	return set, filter, nil
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

func (s Store) UpdateMessagingMessageReplyCount(ctx context.Context, messageID uint64, replies uint) error {
	return s.partialMessagingMessageUpdate(ctx, []string{"replies"}, &types.Message{ID: messageID, Replies: replies})
}

func (s Store) encodeMessagingMessage(res *types.Message) store.Payload {
	// Custom message encoder to make sure that type is not converted to
	// NULL when empty string is used (MessageTypeSimpleMessage)
	return store.Payload{
		"id":          res.ID,
		"type":        sql.NullString{String: string(res.Type), Valid: true},
		"message":     res.Message,
		"meta":        res.Meta,
		"rel_user":    res.UserID,
		"rel_channel": res.ChannelID,
		"reply_to":    res.ReplyTo,
		"replies":     res.Replies,
		"created_at":  res.CreatedAt,
		"updated_at":  res.UpdatedAt,
		"deleted_at":  res.DeletedAt,
	}
}

package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/messaging_mentions.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchMessagingMentions returns all matching rows
//
// This function calls convertMessagingMentionFilter with the given
// types.MentionFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchMessagingMentions(ctx context.Context, f types.MentionFilter) (types.MentionSet, types.MentionFilter, error) {
	var (
		err error
		set []*types.Mention
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertMessagingMentionFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryMessagingMentions(ctx, q, nil)
		return err
	}()
}

// QueryMessagingMentions queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryMessagingMentions(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Mention) (bool, error),
) ([]*types.Mention, error) {
	var (
		set = make([]*types.Mention, 0, DefaultSliceCapacity)
		res *types.Mention

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagingMentionRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupMessagingMentionByID searches for attachment by its ID
//
// It returns attachment even if deleted
func (s Store) LookupMessagingMentionByID(ctx context.Context, id uint64) (*types.Mention, error) {
	return s.execLookupMessagingMention(ctx, squirrel.Eq{
		s.preprocessColumn("msg.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateMessagingMention creates one or more rows in messaging_mention table
func (s Store) CreateMessagingMention(ctx context.Context, rr ...*types.Mention) (err error) {
	for _, res := range rr {
		err = s.checkMessagingMentionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateMessagingMentions(ctx, s.internalMessagingMentionEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateMessagingMention updates one or more existing rows in messaging_mention
func (s Store) UpdateMessagingMention(ctx context.Context, rr ...*types.Mention) error {
	return s.partialMessagingMentionUpdate(ctx, nil, rr...)
}

// partialMessagingMentionUpdate updates one or more existing rows in messaging_mention
func (s Store) partialMessagingMentionUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Mention) (err error) {
	for _, res := range rr {
		err = s.checkMessagingMentionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagingMentions(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("msg.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalMessagingMentionEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertMessagingMention updates one or more existing rows in messaging_mention
func (s Store) UpsertMessagingMention(ctx context.Context, rr ...*types.Mention) (err error) {
	for _, res := range rr {
		err = s.checkMessagingMentionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertMessagingMentions(ctx, s.internalMessagingMentionEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingMention Deletes one or more rows from messaging_mention table
func (s Store) DeleteMessagingMention(ctx context.Context, rr ...*types.Mention) (err error) {
	for _, res := range rr {

		err = s.execDeleteMessagingMentions(ctx, squirrel.Eq{
			s.preprocessColumn("msg.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingMentionByID Deletes row from the messaging_mention table
func (s Store) DeleteMessagingMentionByID(ctx context.Context, ID uint64) error {
	return s.execDeleteMessagingMentions(ctx, squirrel.Eq{
		s.preprocessColumn("msg.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateMessagingMentions Deletes all rows from the messaging_mention table
func (s Store) TruncateMessagingMentions(ctx context.Context) error {
	return s.Truncate(ctx, s.messagingMentionTable())
}

// execLookupMessagingMention prepares MessagingMention query and executes it,
// returning types.Mention (or error)
func (s Store) execLookupMessagingMention(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Mention, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.messagingMentionsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalMessagingMentionRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateMessagingMentions updates all matched (by cnd) rows in messaging_mention with given data
func (s Store) execCreateMessagingMentions(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.messagingMentionTable()).SetMap(payload))
}

// execUpdateMessagingMentions updates all matched (by cnd) rows in messaging_mention with given data
func (s Store) execUpdateMessagingMentions(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.messagingMentionTable("msg")).Where(cnd).SetMap(set))
}

// execUpsertMessagingMentions inserts new or updates matching (by-primary-key) rows in messaging_mention with given data
func (s Store) execUpsertMessagingMentions(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.messagingMentionTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteMessagingMentions Deletes all matched (by cnd) rows in messaging_mention with given data
func (s Store) execDeleteMessagingMentions(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.messagingMentionTable("msg")).Where(cnd))
}

func (s Store) internalMessagingMentionRowScanner(row rowScanner) (res *types.Mention, err error) {
	res = &types.Mention{}

	if _, has := s.config.RowScanners["messagingMention"]; has {
		scanner := s.config.RowScanners["messagingMention"].(func(_ rowScanner, _ *types.Mention) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.MessageID,
			&res.ChannelID,
			&res.UserID,
			&res.MentionedByID,
			&res.CreatedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan messagingMention db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryMessagingMentions returns squirrel.SelectBuilder with set table and all columns
func (s Store) messagingMentionsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.messagingMentionTable("msg"), s.messagingMentionColumns("msg")...)
}

// messagingMentionTable name of the db table
func (Store) messagingMentionTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "messaging_mention" + alias
}

// MessagingMentionColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) messagingMentionColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_message",
		alias + "rel_channel",
		alias + "rel_user",
		alias + "rel_mentioned_by",
		alias + "created_at",
	}
}

// {true true false false false false}

// internalMessagingMentionEncoder encodes fields from types.Mention to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeMessagingMention
// func when rdbms.customEncoder=true
func (s Store) internalMessagingMentionEncoder(res *types.Mention) store.Payload {
	return store.Payload{
		"id":               res.ID,
		"rel_message":      res.MessageID,
		"rel_channel":      res.ChannelID,
		"rel_user":         res.UserID,
		"rel_mentioned_by": res.MentionedByID,
		"created_at":       res.CreatedAt,
	}
}

// checkMessagingMentionConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkMessagingMentionConstraints(ctx context.Context, res *types.Mention) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	if !valid {
		return nil
	}

	return nil
}

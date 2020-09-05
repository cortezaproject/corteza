package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/messaging_channels.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchMessagingChannels returns all matching rows
//
// This function calls convertMessagingChannelFilter with the given
// types.ChannelFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchMessagingChannels(ctx context.Context, f types.ChannelFilter) (types.ChannelSet, types.ChannelFilter, error) {
	var (
		err error
		set []*types.Channel
		q   squirrel.SelectBuilder
	)
	q, err = s.convertMessagingChannelFilter(f)
	if err != nil {
		return nil, f, err
	}

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reversedCursor := f.PageCursor != nil && f.PageCursor.Reverse

	// If paging with reverse cursor, change the sorting
	// direction for all columns we're sorting by
	curSort := f.Sort.Clone()
	if reversedCursor {
		curSort.Reverse()
	}

	return set, f, s.config.ErrorHandler(func() error {
		set, err = s.fetchFullPageOfMessagingChannels(ctx, q, curSort, f.PageCursor, f.Limit, f.Check)

		if err != nil {
			return err
		}

		if f.Limit > 0 && len(set) > 0 {
			if f.PageCursor != nil && (!f.PageCursor.Reverse || uint(len(set)) == f.Limit) {
				f.PrevPage = s.collectMessagingChannelCursorValues(set[0], curSort.Columns()...)
				f.PrevPage.Reverse = true
			}

			// Less items fetched then requested by page-limit
			// not very likely there's another page
			f.NextPage = s.collectMessagingChannelCursorValues(set[len(set)-1], curSort.Columns()...)
		}

		f.PageCursor = nil
		return nil
	}())
}

// fetchFullPageOfMessagingChannels collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - sorting rules (order by ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn). Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfMessagingChannels(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.Channel) (bool, error),
) ([]*types.Channel, error) {
	var (
		set  = make([]*types.Channel, 0, DefaultSliceCapacity)
		aux  []*types.Channel
		last *types.Channel

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedCursor = cursor != nil && cursor.Reverse

		// copy of the select builder
		tryQuery squirrel.SelectBuilder

		fetched uint
		err     error
	)

	// Make sure we always end our sort by primary keys
	if sort.Get("id") == nil {
		sort = append(sort, &filter.SortExpr{Column: "id"})
	}

	// Apply sorting expr from filter to query
	if q, err = setOrderBy(q, sort, s.sortableMessagingChannelColumns()...); err != nil {
		return nil, err
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit))
		}

		if aux, fetched, last, err = s.QueryMessagingChannels(ctx, tryQuery, check); err != nil {
			return nil, err
		}

		if limit > 0 && uint(len(aux)) >= limit {
			// we should use only as much as requested
			set = append(set, aux[0:limit]...)
			break
		} else {
			set = append(set, aux...)
		}

		// if limit is not set or we've already collected enough items
		// we can break the loop right away
		if limit == 0 || fetched == 0 || fetched < limit {
			break
		}

		// In case limit is set very low and we've missed records in the first fetch,
		// make sure next fetch limit is a bit higher
		if limit < MinEnsureFetchLimit {
			limit = MinEnsureFetchLimit
		}

		// @todo improve strategy for collecting next page with lower limit

		// Point cursor to the last fetched element
		if cursor = s.collectMessagingChannelCursorValues(last, sort.Columns()...); cursor == nil {
			break
		}
	}

	if reversedCursor {
		// Cursor for previous page was used
		// Fetched set needs to be reverseCursor because we've forced a descending order to
		// get the previous page
		for i, j := 0, len(set)-1; i < j; i, j = i+1, j-1 {
			set[i], set[j] = set[j], set[i]
		}
	}

	return set, nil
}

// QueryMessagingChannels queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryMessagingChannels(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Channel) (bool, error),
) ([]*types.Channel, uint, *types.Channel, error) {
	var (
		set = make([]*types.Channel, 0, DefaultSliceCapacity)
		res *types.Channel

		// Query rows with
		rows, err = s.Query(ctx, q)

		fetched uint
	)

	if err != nil {
		return nil, 0, nil, err
	}

	defer rows.Close()
	for rows.Next() {
		fetched++
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagingChannelRowScanner(rows)
		}

		if err != nil {
			return nil, 0, nil, err
		}

		// If check function is set, call it and act accordingly
		if check != nil {
			if chk, err := check(res); err != nil {
				return nil, 0, nil, err
			} else if !chk {
				// did not pass the check
				// go with the next row
				continue
			}
		}

		set = append(set, res)
	}

	return set, fetched, res, rows.Err()
}

// LookupMessagingChannelByID searches for attachment by its ID
//
// It returns attachment even if deleted
func (s Store) LookupMessagingChannelByID(ctx context.Context, id uint64) (*types.Channel, error) {
	return s.execLookupMessagingChannel(ctx, squirrel.Eq{
		s.preprocessColumn("mch.id", ""): s.preprocessValue(id, ""),
	})
}

// CreateMessagingChannel creates one or more rows in messaging_channel table
func (s Store) CreateMessagingChannel(ctx context.Context, rr ...*types.Channel) (err error) {
	for _, res := range rr {
		err = s.checkMessagingChannelConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateMessagingChannels(ctx, s.internalMessagingChannelEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateMessagingChannel updates one or more existing rows in messaging_channel
func (s Store) UpdateMessagingChannel(ctx context.Context, rr ...*types.Channel) error {
	return s.config.ErrorHandler(s.PartialMessagingChannelUpdate(ctx, nil, rr...))
}

// PartialMessagingChannelUpdate updates one or more existing rows in messaging_channel
func (s Store) PartialMessagingChannelUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Channel) (err error) {
	for _, res := range rr {
		err = s.checkMessagingChannelConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagingChannels(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("mch.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalMessagingChannelEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertMessagingChannel updates one or more existing rows in messaging_channel
func (s Store) UpsertMessagingChannel(ctx context.Context, rr ...*types.Channel) (err error) {
	for _, res := range rr {
		err = s.checkMessagingChannelConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.config.ErrorHandler(s.execUpsertMessagingChannels(ctx, s.internalMessagingChannelEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingChannel Deletes one or more rows from messaging_channel table
func (s Store) DeleteMessagingChannel(ctx context.Context, rr ...*types.Channel) (err error) {
	for _, res := range rr {

		err = s.execDeleteMessagingChannels(ctx, squirrel.Eq{
			s.preprocessColumn("mch.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteMessagingChannelByID Deletes row from the messaging_channel table
func (s Store) DeleteMessagingChannelByID(ctx context.Context, ID uint64) error {
	return s.execDeleteMessagingChannels(ctx, squirrel.Eq{
		s.preprocessColumn("mch.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateMessagingChannels Deletes all rows from the messaging_channel table
func (s Store) TruncateMessagingChannels(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.messagingChannelTable()))
}

// execLookupMessagingChannel prepares MessagingChannel query and executes it,
// returning types.Channel (or error)
func (s Store) execLookupMessagingChannel(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Channel, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.messagingChannelsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalMessagingChannelRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateMessagingChannels updates all matched (by cnd) rows in messaging_channel with given data
func (s Store) execCreateMessagingChannels(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.messagingChannelTable()).SetMap(payload)))
}

// execUpdateMessagingChannels updates all matched (by cnd) rows in messaging_channel with given data
func (s Store) execUpdateMessagingChannels(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.messagingChannelTable("mch")).Where(cnd).SetMap(set)))
}

// execUpsertMessagingChannels inserts new or updates matching (by-primary-key) rows in messaging_channel with given data
func (s Store) execUpsertMessagingChannels(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.messagingChannelTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteMessagingChannels Deletes all matched (by cnd) rows in messaging_channel with given data
func (s Store) execDeleteMessagingChannels(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.messagingChannelTable("mch")).Where(cnd)))
}

func (s Store) internalMessagingChannelRowScanner(row rowScanner) (res *types.Channel, err error) {
	res = &types.Channel{}

	if _, has := s.config.RowScanners["messagingChannel"]; has {
		scanner := s.config.RowScanners["messagingChannel"].(func(_ rowScanner, _ *types.Channel) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Name,
			&res.Topic,
			&res.Type,
			&res.Meta,
			&res.MembershipPolicy,
			&res.CreatorID,
			&res.OrganisationID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.ArchivedAt,
			&res.DeletedAt,
			&res.LastMessageID,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for MessagingChannel: %w", err)
	} else {
		return res, nil
	}
}

// QueryMessagingChannels returns squirrel.SelectBuilder with set table and all columns
func (s Store) messagingChannelsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.messagingChannelTable("mch"), s.messagingChannelColumns("mch")...)
}

// messagingChannelTable name of the db table
func (Store) messagingChannelTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "messaging_channel" + alias
}

// MessagingChannelColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) messagingChannelColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "name",
		alias + "topic",
		alias + "type",
		alias + "meta",
		alias + "membership_policy",
		alias + "rel_creator",
		alias + "rel_organisation",
		alias + "created_at",
		alias + "updated_at",
		alias + "archived_at",
		alias + "deleted_at",
		alias + "rel_last_message",
	}
}

// {true true true true true}

// sortableMessagingChannelColumns returns all MessagingChannel columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableMessagingChannelColumns() []string {
	return []string{
		"id",
	}
}

// internalMessagingChannelEncoder encodes fields from types.Channel to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeMessagingChannel
// func when rdbms.customEncoder=true
func (s Store) internalMessagingChannelEncoder(res *types.Channel) store.Payload {
	return store.Payload{
		"id":                res.ID,
		"name":              res.Name,
		"topic":             res.Topic,
		"type":              res.Type,
		"meta":              res.Meta,
		"membership_policy": res.MembershipPolicy,
		"rel_creator":       res.CreatorID,
		"rel_organisation":  res.OrganisationID,
		"created_at":        res.CreatedAt,
		"updated_at":        res.UpdatedAt,
		"archived_at":       res.ArchivedAt,
		"deleted_at":        res.DeletedAt,
		"rel_last_message":  res.LastMessageID,
	}
}

// collectMessagingChannelCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectMessagingChannelCursorValues(res *types.Channel, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		// All known primary key columns

		pkId bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "id":
					cursor.Set(c, res.ID, false)

					pkId = true

				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !(pkId && true) {
		collect("id")
	}

	return cursor
}

func (s *Store) checkMessagingChannelConstraints(ctx context.Context, res *types.Channel) error {

	return nil
}

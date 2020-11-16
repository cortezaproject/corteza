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
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
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

	return set, f, func() error {
		q, err = s.convertMessagingChannelFilter(f)
		if err != nil {
			return err
		}

		// Paging enabled
		// {search: {enablePaging:true}}
		// Cleanup unwanted cursors (only relevant is f.PageCursor, next&prev are reset and returned)
		f.PrevPage, f.NextPage = nil, nil

		if f.PageCursor != nil {
			// Page cursor exists so we need to validate it against used sort
			if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
				return err
			}
		}

		if len(f.Sort) == 0 {
			f.Sort = filter.SortExprSet{}
		}

		// Make sure results are always sorted at least by primary keys
		if f.Sort.Get("id") == nil {
			f.Sort = append(f.Sort, &filter.SortExpr{Column: "id"})
		}

		sort := f.Sort.Clone()

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		if f.PageCursor != nil && f.PageCursor.Reverse {
			sort.Reverse()
		}

		// Apply sorting expr from filter to query
		if q, err = setOrderBy(q, sort, s.sortableMessagingChannelColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfMessagingChannels(ctx, q, sort.Columns(), sort.Reversed(), f.PageCursor, f.Limit, f.Check)
		if err != nil {
			return err
		}

		f.PageCursor = nil
		return nil
	}()
}

// fetchFullPageOfMessagingChannels collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfMessagingChannels(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sortColumns []string,
	sortDesc bool,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.Channel) (bool, error),
) (set []*types.Channel, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.Channel

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = cursor != nil && cursor.Reverse

		// copy of the select builder
		tryQuery squirrel.SelectBuilder

		fetched uint
	)

	set = make([]*types.Channel, 0, DefaultSliceCapacity)

	if cursor != nil {
		cursor.Reverse = sortDesc
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit + 1))
		}

		if aux, err = s.QueryMessagingChannels(ctx, tryQuery, check); err != nil {
			return nil, nil, nil, err
		}

		fetched = uint(len(aux))
		if cursor != nil && prev == nil && fetched > 0 {
			// Cursor for previous page is calculated only when cursor is used (so, not on first page)
			prev = s.collectMessagingChannelCursorValues(aux[0], sortColumns...)
		}

		// Point cursor to the last fetched element
		if fetched > limit && limit > 0 {
			next = s.collectMessagingChannelCursorValues(aux[limit-1], sortColumns...)

			// we should use only as much as requested
			set = append(set, aux[:limit]...)
			break
		} else {
			set = append(set, aux...)
		}

		// if limit is not set or we've already collected enough items
		// we can break the loop right away
		if limit == 0 || fetched == 0 || fetched <= limit {
			break
		}

		// In case limit is set very low and we've missed records in the first fetch,
		// make sure next fetch limit is a bit higher
		if limit < MinEnsureFetchLimit {
			limit = MinEnsureFetchLimit
		}

		// @todo improve strategy for collecting next page with lower limit
	}

	if reversedOrder {
		// Fetched set needs to be reversed because we've forced a descending order to get the previous page
		for i, j := 0, len(set)-1; i < j; i, j = i+1, j-1 {
			set[i], set[j] = set[j], set[i]
		}

		// and flip prev/next cursors too
		prev, next = next, prev
	}

	if prev != nil {
		prev.Reverse = true
	}

	return set, prev, next, nil
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
) ([]*types.Channel, error) {
	var (
		set = make([]*types.Channel, 0, DefaultSliceCapacity)
		res *types.Channel

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagingChannelRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if check != nil {
			if chk, err := check(res); err != nil {
				return nil, err
			} else if !chk {
				continue
			}
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupMessagingChannelByID searches for attachment by its ID
//
// It returns attachment even if deleted
func (s Store) LookupMessagingChannelByID(ctx context.Context, id uint64) (*types.Channel, error) {
	return s.execLookupMessagingChannel(ctx, squirrel.Eq{
		s.preprocessColumn("mch.id", ""): store.PreprocessValue(id, ""),
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
	return s.partialMessagingChannelUpdate(ctx, nil, rr...)
}

// partialMessagingChannelUpdate updates one or more existing rows in messaging_channel
func (s Store) partialMessagingChannelUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Channel) (err error) {
	for _, res := range rr {
		err = s.checkMessagingChannelConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagingChannels(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("mch.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalMessagingChannelEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
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

		err = s.execUpsertMessagingChannels(ctx, s.internalMessagingChannelEncoder(res))
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
			s.preprocessColumn("mch.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingChannelByID Deletes row from the messaging_channel table
func (s Store) DeleteMessagingChannelByID(ctx context.Context, ID uint64) error {
	return s.execDeleteMessagingChannels(ctx, squirrel.Eq{
		s.preprocessColumn("mch.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateMessagingChannels Deletes all rows from the messaging_channel table
func (s Store) TruncateMessagingChannels(ctx context.Context) error {
	return s.Truncate(ctx, s.messagingChannelTable())
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
	return s.Exec(ctx, s.InsertBuilder(s.messagingChannelTable()).SetMap(payload))
}

// execUpdateMessagingChannels updates all matched (by cnd) rows in messaging_channel with given data
func (s Store) execUpdateMessagingChannels(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.messagingChannelTable("mch")).Where(cnd).SetMap(set))
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

	return s.Exec(ctx, upsert)
}

// execDeleteMessagingChannels Deletes all matched (by cnd) rows in messaging_channel with given data
func (s Store) execDeleteMessagingChannels(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.messagingChannelTable("mch")).Where(cnd))
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
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.ArchivedAt,
			&res.DeletedAt,
			&res.LastMessageID,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan messagingChannel db row").Wrap(err)
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
		alias + "created_at",
		alias + "updated_at",
		alias + "archived_at",
		alias + "deleted_at",
		alias + "rel_last_message",
	}
}

// {true true false true true true}

// sortableMessagingChannelColumns returns all MessagingChannel columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableMessagingChannelColumns() map[string]string {
	return map[string]string{
		"id": "id",
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

// checkMessagingChannelConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkMessagingChannelConstraints(ctx context.Context, res *types.Channel) error {
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

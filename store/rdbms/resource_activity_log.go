package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/discovery/types"
	"github.com/cortezaproject/corteza-server/store"
)

func (s Store) convertResourceActivityLogFilter(f types.ResourceActivityFilter) (query squirrel.SelectBuilder, err error) {
	query = s.actionlogsSelectBuilder()

	// Always sort by ID descending
	query = query.OrderBy("id DESC")

	if f.FromTimestamp != nil {
		query = query.Where(squirrel.GtOrEq{"ts": f.FromTimestamp})
	}

	if f.ToTimestamp != nil {
		query = query.Where(squirrel.LtOrEq{"ts": f.ToTimestamp})
	}

	if f.Limit == 0 || f.Limit > MaxLimit {
		f.Limit = MaxLimit
	}

	query = query.Limit(uint64(f.Limit))

	return
}

func (s Store) scanResourceActivityLogRow(row rowScanner, res *types.ResourceActivity) (err error) {
	err = row.Scan(
		&res.ID,
		&res.ResourceID,
		&res.ResourceType,
		&res.ResourceAction,
		&res.Timestamp,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s Store) encodeResourceActivityLog(res *types.ResourceActivity) store.Payload {
	// ActivityLogEnc encodes fields from discovery.Action to store.Payload (map)
	out := store.Payload{
		"id":              res.ID,
		"rel_resource":    res.ResourceID,
		"resource_type":   res.ResourceType,
		"resource_action": res.ResourceAction,
		"ts":              res.Timestamp,
	}

	return out
}

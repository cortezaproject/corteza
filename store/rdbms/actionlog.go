package rdbms

import (
	"encoding/json"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/store"
)

func (s Store) convertActionlogFilter(f actionlog.Filter) (query squirrel.SelectBuilder, err error) {
	query = s.actionlogsSelectBuilder()

	if f.From != nil {
		query = query.Where(squirrel.GtOrEq{"ts": f.From})
	}

	if f.To != nil {
		query = query.Where(squirrel.LtOrEq{"ts": f.To})
	}

	if len(f.ActorID) > 0 {
		query = query.Where(squirrel.Eq{"actor_id": f.ActorID})
	}

	if f.Resource != "" {
		query = query.Where(squirrel.Eq{"resource": f.Resource})
	}

	if f.Action != "" {
		query = query.Where(squirrel.Eq{"action": f.Action})
	}

	return
}

func (s Store) scanActionlogRow(row rowScanner, res *actionlog.Action) (err error) {
	var metaBuf json.RawMessage

	err = row.Scan(
		&res.ID,
		&res.Timestamp,
		&res.RequestOrigin,
		&res.RequestID,
		&res.ActorIPAddr,
		&res.ActorID,
		&res.Resource,
		&res.Action,
		&res.Error,
		&res.Severity,
		&res.Description,
		&metaBuf,
	)

	if err != nil {
		return err
	}

	// Ignoring unmarshal errors
	_ = json.Unmarshal(metaBuf, &res.Meta)

	return nil
}

func (s Store) encodeActionlog(res *actionlog.Action) store.Payload {
	// ActionlogEnc encodes fields from actionlog.Action to store.Payload (map)
	out := store.Payload{
		"id":             res.ID,
		"ts":             res.Timestamp,
		"request_origin": res.RequestOrigin,
		"request_id":     res.RequestID,
		"actor_ip_addr":  res.ActorIPAddr,
		"actor_id":       res.ActorID,
		"resource":       res.Resource,
		"action":         res.Action,
		"error":          res.Error,
		"severity":       res.Severity,
		"description":    res.Description,
		"meta":           []byte("{}"),
	}

	if res.Meta != nil {
		out["meta"], _ = json.Marshal(res.Meta)
	}

	return out
}

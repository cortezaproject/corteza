package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/titpetric/factory"
	"time"
)

type (
	// Basic mysql storage backend for audit log events
	//
	// this does not follow the usual (one) repository pattern
	// but tries to move towards multi-flavoured repository support
	mysql struct {
		dbh *factory.DB
		tbl string
	}

	event struct {
		Timestamp     time.Time       `db:"ts"`
		RequestOrigin string          `db:"request_origin"`
		RequestID     string          `db:"request_id"`
		ActorIPAddr   string          `db:"actor_ip_addr"`
		ActorID       uint64          `db:"actor_id"`
		Resource      string          `db:"resource"`
		Action        string          `db:"action"`
		Error         string          `db:"error"`
		Severity      int             `db:"severity"`
		Description   string          `db:"description"`
		Meta          json.RawMessage `db:"meta"`
	}
)

func Mysql(db *factory.DB, tbl string) *mysql {
	return &mysql{
		// connection
		dbh: db,

		// table to store the data
		tbl: tbl,
	}
}

func (r *mysql) db() *factory.DB {
	return r.dbh
}

func (r mysql) columns() []string {
	return []string{
		"ts",
		"request_origin",
		"request_id",
		"actor_ip_addr",
		"actor_id",
		"resource",
		"action",
		"error",
		"severity",
		"description",
		"meta",
	}
}

func (r mysql) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.tbl)
}

func (r *mysql) Find(ctx context.Context, flt actionlog.Filter) (set actionlog.ActionSet, f actionlog.Filter, err error) {
	f = flt

	query := r.query()

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

	// @todo implement filtering with query (via pkg/ql)

	query = query.OrderBy("ts DESC")

	results := make([]*event, 0)
	if err = rh.FetchPaged(r.db().With(ctx), query, f.PageFilter, &results); err != nil {
		return nil, f, err
	}

	set = make(actionlog.ActionSet, len(results))
	for i, r := range results {
		set[i] = &actionlog.Action{
			Timestamp:     r.Timestamp,
			RequestOrigin: r.RequestOrigin,
			RequestID:     r.RequestID,
			ActorIPAddr:   r.ActorIPAddr,
			ActorID:       r.ActorID,
			Resource:      r.Resource,
			Action:        r.Action,
			Error:         r.Error,
			Severity:      actionlog.Severity(r.Severity),
			Description:   r.Description,
		}

		// ignore all unmarshaling issues
		_ = json.Unmarshal(r.Meta, &set[i].Meta)
	}

	return set, f, nil
}

// Record stores audit event
func (r *mysql) Record(ctx context.Context, e *actionlog.Action) error {
	m, err := json.Marshal(e.Meta)
	if err != nil {
		return fmt.Errorf("could not format auditlog event: %w", err)
	}

	return r.dbh.With(ctx).InsertIgnore(r.tbl, event{
		Timestamp:     e.Timestamp,
		RequestOrigin: e.RequestOrigin,
		RequestID:     e.RequestID,
		ActorIPAddr:   e.ActorIPAddr,
		ActorID:       e.ActorID,
		Resource:      e.Resource,
		Action:        e.Action,
		Error:         e.Error,
		Severity:      int(e.Severity),
		Description:   e.Description,
		Meta:          m,
	})
}

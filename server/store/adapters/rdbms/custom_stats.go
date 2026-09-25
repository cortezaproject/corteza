package rdbms

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	automationType "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/store/adapters/rdbms/ql"
	systemType "github.com/crusttech/human/server/system/types"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

// The admin dashboard's aggregate. Every series is grouped per day in SQL,
// through the dialect's own date truncation, so rolling days into weeks or
// months stays a Go concern and the SQL stays portable.

type (
	// statsResource describes one inventory table: how a row is classified
	// into a status (first matching class wins) and which rows count at all.
	statsResource struct {
		table    exp.IdentifierExpression
		where    []exp.Expression
		classes  []statsClass
		fallback exp.Expression
	}

	statsClass struct {
		name string
		cond exp.Expression
	}

	// statsLabel says which columns name a row of an inventory table: the
	// label is the first non-empty of label, metaName (a key of the meta
	// JSON) and handle.
	statsLabel struct {
		label    string
		metaName string
		handle   string
	}
)

const (
	statsRecentLimit = 8
	statsTopLimit    = 8

	// action-log rows that count as an error on the activity chart
	statsErrorSeverity = actionlog.Error

	// action-log rows a finished TAQ run leaves behind
	statsTaqRunResource = "automation:ng-automation"
	statsTaqRunAction   = "run"
)

func notNull(col string) exp.Expression {
	return exp.NewLiteralExpression("? IS NOT NULL", goqu.C(col))
}

func isFalse(col string) exp.Expression {
	return goqu.C(col).Eq(false)
}

var (
	statsDeleted = statsClass{"deleted", notNull("deleted_at")}

	// inventory resources in the order the payload lists them
	statsResources = map[string]statsResource{
		systemType.SystemStatsUsers: {
			table:    userTable,
			where:    []exp.Expression{goqu.C("kind").Eq("")},
			classes:  []statsClass{statsDeleted, {"suspended", notNull("suspended_at")}},
			fallback: goqu.V("active"),
		},
		systemType.SystemStatsRoles: {
			table:    roleTable,
			classes:  []statsClass{statsDeleted, {"archived", notNull("archived_at")}},
			fallback: goqu.V("active"),
		},
		systemType.SystemStatsApplications: {
			table:    applicationTable,
			classes:  []statsClass{statsDeleted, {"disabled", isFalse("enabled")}},
			fallback: goqu.V("enabled"),
		},
		systemType.SystemStatsAuthClients: {
			table:    authClientTable,
			classes:  []statsClass{statsDeleted, {"disabled", isFalse("enabled")}},
			fallback: goqu.V("enabled"),
		},
		systemType.SystemStatsAgents: {
			table:    agentTable,
			classes:  []statsClass{statsDeleted},
			fallback: goqu.V("active"),
		},
		systemType.SystemStatsChatbots: {
			table:    chatbotTable,
			classes:  []statsClass{statsDeleted, {"disabled", isFalse("enabled")}},
			fallback: goqu.V("enabled"),
		},
		systemType.SystemStatsProjects: {
			table:    projectTable,
			classes:  []statsClass{statsDeleted},
			fallback: goqu.C("status"),
		},
		systemType.SystemStatsWorkflows: {
			table:    automationWorkflowTable,
			classes:  []statsClass{statsDeleted, {"disabled", isFalse("enabled")}},
			fallback: goqu.V("enabled"),
		},
		systemType.SystemStatsTaqs: {
			table:    automationNgAutomationTable,
			classes:  []statsClass{statsDeleted, {"disabled", isFalse("enabled")}},
			fallback: goqu.V("enabled"),
		},
		systemType.SystemStatsNamespaces: {
			table:    composeNamespaceTable,
			classes:  []statsClass{statsDeleted, {"disabled", isFalse("enabled")}},
			fallback: goqu.V("active"),
		},
		systemType.SystemStatsModules: {
			table:    composeModuleTable,
			classes:  []statsClass{statsDeleted},
			fallback: goqu.V("active"),
		},
		systemType.SystemStatsConnections: {
			table:    connectionTable,
			classes:  []statsClass{statsDeleted},
			fallback: goqu.C("status"),
		},
		systemType.SystemStatsDataSources: {
			table:    dalConnectionTable,
			classes:  []statsClass{statsDeleted},
			fallback: goqu.V("active"),
		},
	}
)

var statsLabels = map[string]statsLabel{
	systemType.SystemStatsUsers:        {label: "name", handle: "email"},
	systemType.SystemStatsRoles:        {label: "name", handle: "handle"},
	systemType.SystemStatsApplications: {label: "name"},
	systemType.SystemStatsAuthClients:  {handle: "handle"},
	systemType.SystemStatsAgents:       {metaName: "short", handle: "handle"},
	systemType.SystemStatsChatbots:     {label: "name", handle: "handle"},
	systemType.SystemStatsProjects:     {metaName: "name", handle: "handle"},
	systemType.SystemStatsWorkflows:    {metaName: "name", handle: "handle"},
	systemType.SystemStatsTaqs:         {metaName: "short", handle: "handle"},
	systemType.SystemStatsNamespaces:   {label: "name", handle: "slug"},
	systemType.SystemStatsModules:      {label: "name", handle: "handle"},
	systemType.SystemStatsConnections:  {metaName: "name", handle: "handle"},
	systemType.SystemStatsDataSources:  {metaName: "name", handle: "handle"},
}

func (s *Store) SystemStats(ctx context.Context, r systemType.SystemStatsRange) (raw *systemType.SystemStatsRaw, err error) {
	raw = &systemType.SystemStatsRaw{
		Resources: make(map[string]*systemType.SystemStatsRawResource, len(statsResources)),
	}

	for name, res := range statsResources {
		rr := &systemType.SystemStatsRawResource{}
		if rr.Status, err = s.statsStatusCounts(ctx, res); err != nil {
			return nil, fmt.Errorf("%s status: %w", name, err)
		}

		if rr.Created, err = s.statsDaily(ctx, res.table, "created_at", r, res.where, nil); err != nil {
			return nil, fmt.Errorf("%s created: %w", name, err)
		}

		raw.Resources[name] = rr
	}

	if err = s.statsWorkflowRuns(ctx, r, raw); err != nil {
		return nil, fmt.Errorf("workflow runs: %w", err)
	}

	if err = s.statsTaqRuns(ctx, r, raw); err != nil {
		return nil, fmt.Errorf("taq runs: %w", err)
	}

	if err = s.statsActivity(ctx, r, raw); err != nil {
		return nil, fmt.Errorf("activity: %w", err)
	}

	if err = s.statsSignins(ctx, r, raw); err != nil {
		return nil, fmt.Errorf("signins: %w", err)
	}

	return raw, nil
}

// statsStatusCounts counts every row of the table once, under the first
// class whose condition holds, or the fallback.
func (s *Store) statsStatusCounts(ctx context.Context, res statsResource) (map[string]uint, error) {
	c := goqu.Case()
	for _, cls := range res.classes {
		c = c.When(cls.cond, goqu.V(cls.name))
	}

	// classified in a subquery so the outer GROUP BY names a plain column;
	// grouping by the CASE itself fails on postgres once it holds placeholders
	classified := s.Dialect.GOQU().
		From(res.table).
		Select(c.Else(res.fallback).As("status")).
		Where(res.where...)

	query := s.Dialect.GOQU().
		From(classified.As("classified")).
		Select(goqu.C("status"), goqu.COUNT(goqu.Star()).As("count")).
		GroupBy(goqu.C("status"))

	rows, err := s.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := make(map[string]uint)
	for rows.Next() {
		var (
			key   any
			count uint
		)

		if err = rows.Scan(&key, &count); err != nil {
			return nil, err
		}

		out[statsKey(key)] = count
	}

	return out, rows.Err()
}

// statsDaily counts rows per day of tsCol within the range, split by keyExpr
// when one is given.
func (s *Store) statsDaily(ctx context.Context, tbl exp.IdentifierExpression, tsCol string, r systemType.SystemStatsRange, where []exp.Expression, keyExpr exp.Expression) (out []systemType.SystemStatsDaily, err error) {
	day, err := s.Dialect.ExprHandler(&ql.ASTNode{Ref: "date"}, goqu.C(tsCol))
	if err != nil {
		return nil, err
	}

	// the day and key are computed in a subquery so the outer GROUP BY names
	// plain columns; grouping by an expression that holds placeholders fails
	// on postgres
	inner := []any{exp.NewAliasExpression(day, "day")}
	outer := []any{goqu.C("day")}
	if keyExpr != nil {
		inner = append(inner, exp.NewAliasExpression(keyExpr, "key"))
		outer = append(outer, goqu.C("key"))
	}

	bucketed := s.Dialect.GOQU().
		From(tbl).
		Select(inner...).
		Where(append(where, statsRangeWhere(goqu.C(tsCol), r)...)...)

	query := s.Dialect.GOQU().
		From(bucketed.As("bucketed")).
		Select(append(outer, goqu.COUNT(goqu.Star()).As("count"))...).
		GroupBy(outer...)

	rows, err := s.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out = make([]systemType.SystemStatsDaily, 0)
	for rows.Next() {
		var (
			dayVal, keyVal any
			count          uint
			dest           = []any{&dayVal, &count}
		)

		if keyExpr != nil {
			dest = []any{&dayVal, &keyVal, &count}
		}

		if err = rows.Scan(dest...); err != nil {
			return nil, err
		}

		out = append(out, systemType.SystemStatsDaily{
			Day:   statsDay(dayVal),
			Key:   statsKey(keyVal),
			Count: count,
		})
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].Day < out[j].Day })

	return out, rows.Err()
}

func statsRangeWhere(tsCol exp.IdentifierExpression, r systemType.SystemStatsRange) []exp.Expression {
	return []exp.Expression{
		tsCol.Gte(r.From),
		tsCol.Lt(r.To),
	}
}

func (s *Store) statsWorkflowRuns(ctx context.Context, r systemType.SystemStatsRange, raw *systemType.SystemStatsRaw) (err error) {
	raw.WorkflowRuns, err = s.statsDaily(ctx, automationSessionTable, "created_at", r, nil, goqu.C("status"))
	if err != nil {
		return
	}

	// the status column is the numeric enum; the payload speaks its names
	raw.WorkflowRunTotals = make(map[string]uint)
	for i, d := range raw.WorkflowRuns {
		name := statsSessionStatus(d.Key)
		raw.WorkflowRuns[i].Key = name
		raw.WorkflowRunTotals[name] += d.Count
	}

	wfName, err := s.Dialect.JsonExtractUnquote(goqu.I("w.meta"), "name")
	if err != nil {
		return err
	}

	query := s.Dialect.GOQU().
		From(automationSessionTable.As("s")).
		LeftJoin(automationWorkflowTable.As("w"), goqu.On(goqu.I("w.id").Eq(goqu.I("s.rel_workflow")))).
		Select(goqu.I("s.id"), goqu.I("s.rel_workflow"), goqu.COALESCE(goqu.I("w.handle"), goqu.V("")), goqu.COALESCE(wfName, goqu.V("")), goqu.I("s.event_type"), goqu.I("s.error"), goqu.I("s.created_at")).
		Where(append(
			[]exp.Expression{goqu.I("s.status").Eq(int(automationType.SessionFailed))},
			statsRangeWhere(goqu.I("s.created_at"), r)...,
		)...).
		Order(goqu.I("s.created_at").Desc()).
		Limit(statsRecentLimit)

	rows, err := s.Query(ctx, query)
	if err != nil {
		return err
	}

	defer rows.Close()

	raw.WorkflowFailures = make([]*systemType.SystemStatsWorkflowFailure, 0)
	for rows.Next() {
		f := &systemType.SystemStatsWorkflowFailure{}
		var errText sql.NullString
		if err = rows.Scan(&f.SessionID, &f.WorkflowID, &f.WorkflowHandle, &f.WorkflowName, &f.EventType, &errText, &f.CreatedAt); err != nil {
			return err
		}
		f.Error = errText.String
		raw.WorkflowFailures = append(raw.WorkflowFailures, f)
	}

	return rows.Err()
}

func (s *Store) statsTaqRuns(ctx context.Context, r systemType.SystemStatsRange, raw *systemType.SystemStatsRaw) (err error) {
	outcome, err := s.Dialect.JsonExtractUnquote(goqu.C("meta"), "status")
	if err != nil {
		return err
	}

	where := []exp.Expression{
		goqu.C("resource").Eq(statsTaqRunResource),
		goqu.C("action").Eq(statsTaqRunAction),
	}

	raw.TaqRuns, err = s.statsDaily(ctx, actionlogTable, "ts", r, where, outcome)
	if err != nil {
		return
	}

	raw.TaqRunTotals = make(map[string]uint)
	for _, d := range raw.TaqRuns {
		raw.TaqRunTotals[d.Key] += d.Count
	}

	raw.TaqFailures, err = s.statsLogEntries(ctx, r, append(where, goqu.C("error").Neq("")))
	return
}

func (s *Store) statsActivity(ctx context.Context, r systemType.SystemStatsRange, raw *systemType.SystemStatsRaw) (err error) {
	isErr := goqu.C("severity").Lte(int(statsErrorSeverity))
	outcome := goqu.Case().When(isErr, goqu.V("error")).Else(goqu.V("ok"))

	raw.Activity, err = s.statsDaily(ctx, actionlogTable, "ts", r, nil, outcome)
	if err != nil {
		return
	}

	raw.ActivityTotals = make(map[string]uint)
	for _, d := range raw.Activity {
		raw.ActivityTotals[d.Key] += d.Count
	}

	if raw.RecentErrors, err = s.statsLogEntries(ctx, r, []exp.Expression{isErr}); err != nil {
		return
	}

	raw.TopResources, err = s.statsTopResources(ctx, r)
	return
}

// statsLogEntries returns the newest action-log rows matching where, within the range.
func (s *Store) statsLogEntries(ctx context.Context, r systemType.SystemStatsRange, where []exp.Expression) ([]*systemType.SystemStatsLogEntry, error) {
	query := s.Dialect.GOQU().
		From(actionlogTable).
		Select("id", "ts", "resource", "action", "description", "error", "actor_id", "meta").
		Where(append(where, statsRangeWhere(goqu.C("ts"), r)...)...).
		Order(goqu.C("ts").Desc()).
		Limit(statsRecentLimit)

	rows, err := s.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := make([]*systemType.SystemStatsLogEntry, 0)
	for rows.Next() {
		var (
			e    = &systemType.SystemStatsLogEntry{}
			meta rawJson
		)

		if err = rows.Scan(&e.ActionID, &e.Timestamp, &e.Resource, &e.Action, &e.Description, &e.Error, &e.ActorID, &meta); err != nil {
			return nil, err
		}

		if len(meta) > 0 {
			// meta that fails to decode is dropped, not fatal: the entry still lists
			_ = json.Unmarshal(meta, &e.Meta)
		}

		out = append(out, e)
	}

	return out, rows.Err()
}

// statsTopResources ranks action-log activity by resource type within the
// range. Rows are grouped per stored resource string and folded in Go, since
// the column mixes types ("system:user") with instances
// ("corteza::compose:namespace/123").
func (s *Store) statsTopResources(ctx context.Context, r systemType.SystemStatsRange) ([]systemType.SystemStatsKeyCount, error) {
	query := s.Dialect.GOQU().
		From(actionlogTable).
		Select(goqu.C("resource"), goqu.COUNT(goqu.Star()).As("count")).
		Where(statsRangeWhere(goqu.C("ts"), r)...).
		GroupBy(goqu.C("resource"))

	rows, err := s.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	folded := make(map[string]uint)
	for rows.Next() {
		var (
			resource string
			count    uint
		)

		if err = rows.Scan(&resource, &count); err != nil {
			return nil, err
		}

		folded[statsResourceType(resource)] += count
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	out := make([]systemType.SystemStatsKeyCount, 0, len(folded))
	for k, c := range folded {
		out = append(out, systemType.SystemStatsKeyCount{Key: k, Count: c})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Key < out[j].Key
	})

	if len(out) > statsTopLimit {
		out = out[:statsTopLimit]
	}

	return out, nil
}

func (s *Store) statsSignins(ctx context.Context, r systemType.SystemStatsRange, raw *systemType.SystemStatsRaw) (err error) {
	if raw.Signins, err = s.statsDaily(ctx, authSessionTable, "created_at", r, nil, nil); err != nil {
		return
	}

	for _, d := range raw.Signins {
		raw.SigninTotal += d.Count
	}

	aux := struct {
		Users uint `db:"users"`
	}{}

	query := s.Dialect.GOQU().
		From(authSessionTable).
		Select(goqu.COUNT(goqu.DISTINCT("rel_user")).As("users")).
		Where(statsRangeWhere(goqu.C("created_at"), r)...)

	if err = s.QueryOne(ctx, query, &aux); err != nil {
		return
	}
	raw.SigninUsers = aux.Users

	live := struct {
		Live uint `db:"live"`
	}{}

	query = s.Dialect.GOQU().
		From(authSessionTable).
		Select(goqu.COUNT(goqu.Star()).As("live")).
		Where(goqu.C("expires_at").Gt(time.Now()))

	if err = s.QueryOne(ctx, query, &live); err != nil {
		return
	}
	raw.LiveSessions = live.Live

	return nil
}

// statsDay normalises a driver's date value (time.Time, string or bytes) to YYYY-MM-DD.
func statsDay(v any) string {
	const format = "2006-01-02"

	switch d := v.(type) {
	case time.Time:
		return d.Format(format)
	case []byte:
		v = string(d)
	}

	if str, ok := v.(string); ok && len(str) >= len(format) {
		return str[:len(format)]
	}

	return ""
}

// statsKey normalises a driver's grouping value to a string.
func statsKey(v any) string {
	switch k := v.(type) {
	case nil:
		return ""
	case string:
		return k
	case []byte:
		return string(k)
	default:
		return fmt.Sprint(k)
	}
}

func statsSessionStatus(key string) string {
	var n int
	if _, err := fmt.Sscanf(key, "%d", &n); err != nil {
		return key
	}

	return automationType.SessionStatus(n).String()
}

// statsResourceType reduces a stored resource identifier to its type:
// "corteza::compose:namespace/123" and "compose:namespace/*" both become
// "compose:namespace".
func statsResourceType(resource string) string {
	resource = strings.TrimPrefix(resource, "corteza::")
	if i := strings.IndexByte(resource, '/'); i >= 0 {
		resource = resource[:i]
	}

	return resource
}

// SystemStatsResourceDetail answers the drill-down for one inventory resource.
func (s *Store) SystemStatsResourceDetail(ctx context.Context, resource string, r systemType.SystemStatsRange) (raw *systemType.SystemStatsDetailRaw, err error) {
	res, ok := statsResources[resource]
	if !ok {
		return nil, fmt.Errorf("unknown inventory resource %q", resource)
	}

	raw = &systemType.SystemStatsDetailRaw{}

	if raw.Status, err = s.statsStatusCounts(ctx, res); err != nil {
		return nil, fmt.Errorf("%s status: %w", resource, err)
	}

	raw.Movement = make([]systemType.SystemStatsDaily, 0)
	for _, mv := range []struct{ key, col string }{
		{systemType.SystemStatsCreated, "created_at"},
		{systemType.SystemStatsUpdated, "updated_at"},
		{systemType.SystemStatsDeleted, "deleted_at"},
	} {
		daily, err := s.statsDaily(ctx, res.table, mv.col, r, res.where, nil)
		if err != nil {
			return nil, fmt.Errorf("%s %s: %w", resource, mv.key, err)
		}

		for _, d := range daily {
			d.Key = mv.key
			raw.Movement = append(raw.Movement, d)
		}
	}

	if raw.Recent, err = s.statsRecent(ctx, resource, res); err != nil {
		return nil, fmt.Errorf("%s recent: %w", resource, err)
	}

	return raw, nil
}

// statsRecent lists the newest rows of an inventory table, deleted ones included.
func (s *Store) statsRecent(ctx context.Context, resource string, res statsResource) ([]*systemType.SystemStatsItem, error) {
	var (
		lbl                   = statsLabels[resource]
		empty                 = goqu.V("")
		label                 = []any{}
		handle exp.Expression = empty
	)

	if lbl.label != "" {
		label = append(label, goqu.Func("NULLIF", goqu.C(lbl.label), empty))
	}
	if lbl.metaName != "" {
		name, err := s.Dialect.JsonExtractUnquote(goqu.C("meta"), lbl.metaName)
		if err != nil {
			return nil, err
		}
		label = append(label, goqu.Func("NULLIF", name, empty))
	}
	if lbl.handle != "" {
		handle = goqu.COALESCE(goqu.C(lbl.handle), empty)
		label = append(label, goqu.C(lbl.handle))
	}
	label = append(label, empty)

	c := goqu.Case()
	for _, cls := range res.classes {
		c = c.When(cls.cond, goqu.V(cls.name))
	}

	query := s.Dialect.GOQU().
		From(res.table).
		Select(goqu.C("id"), goqu.COALESCE(label...), handle, c.Else(res.fallback), goqu.C("created_at"), goqu.C("updated_at"), goqu.C("deleted_at")).
		Where(res.where...).
		Order(goqu.C("created_at").Desc()).
		Limit(statsRecentLimit + 2)

	rows, err := s.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := make([]*systemType.SystemStatsItem, 0)
	for rows.Next() {
		var (
			it     = &systemType.SystemStatsItem{}
			status any
		)

		if err = rows.Scan(&it.ID, &it.Label, &it.Handle, &status, &it.CreatedAt, &it.UpdatedAt, &it.DeletedAt); err != nil {
			return nil, err
		}

		it.Status = statsKey(status)
		out = append(out, it)
	}

	return out, rows.Err()
}

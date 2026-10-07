package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	atypes "github.com/crusttech/human/server/automation/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/api/drivers/google"
	restDriver "github.com/crusttech/human/server/store/adapters/api/drivers/rest"
	sysEvent "github.com/crusttech/human/server/system/service/event"
	"github.com/crusttech/human/server/system/types"
)

const pollCursorPrefix = "pollCursor:"

// errPollCursorStale marks a list call the provider rejected because the cursor
// is too old (e.g. Gmail history.list 404); the engine reseeds and carries on.
var errPollCursorStale = errors.New("poll cursor stale")

// StartConnectorPollLoop runs the declarative poll engine on an interval. It
// delivers webhook events for connectors that declare a poll spec instead of
// pushing — no per-connector Go.
func (svc *configuredConnection) StartConnectorPollLoop(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				svc.pollAllConnections(ctx)
			}
		}
	}()
}

func (svc *configuredConnection) pollAllConnections(ctx context.Context) {
	// Service identity so eventbus handlers and store writes are authorized.
	ctx = a.SetIdentityToContext(ctx, a.ServiceUser())

	set, _, err := store.SearchConfiguredConnections(ctx, svc.store, types.ConfiguredConnectionFilter{
		Status: []string{"active"},
	})
	if err != nil {
		svc.services.logger.Warn("connector poll: failed to load connections", zap.Error(err))
		return
	}

	for _, cc := range set {
		// The poll spec lives on the parent connection; the configured
		// connection's embedded copy can be stale after a catalog update, so
		// load the parent fresh.
		conn, err := loadConnection(ctx, svc.store, cc.ConnectionID)
		if err != nil || !hasPollSpec(conn) {
			continue
		}
		runner, err := svc.runnerFor(ctx, cc)
		if err != nil {
			svc.services.logger.Warn("connector poll: runner failed", zap.Uint64("ccID", cc.ID), zap.Error(err))
			continue
		}
		for _, res := range conn.Resources {
			for _, wh := range res.Webhooks {
				if wh.Poll == nil || wh.Poll.List == nil {
					continue
				}
				// A poll whose List references placeholders beyond {{cursor}}
				// (e.g. {{spreadsheetId}}) can't run connection-wide — those
				// values come from each trigger's scope, so poll once per scope.
				scopeParams := actionPlaceholders(wh.Poll.List)
				delete(scopeParams, "cursor")
				if len(scopeParams) == 0 {
					if err := svc.pollWebhook(ctx, cc, runner, wh, nil, wh.Event); err != nil {
						svc.services.logger.Warn("connector poll: webhook failed",
							zap.Uint64("ccID", cc.ID), zap.String("event", wh.Event), zap.Error(err))
					}
					continue
				}
				for _, scope := range svc.triggerScopes(ctx, cc.ConnectionID, cc.ID, wh.Event, scopeParams) {
					key := wh.Event + "|" + scopeSig(scope, scopeParams)
					if err := svc.pollWebhook(ctx, cc, runner, wh, scope, key); err != nil {
						svc.services.logger.Warn("connector poll: scoped webhook failed",
							zap.Uint64("ccID", cc.ID), zap.String("event", wh.Event),
							zap.String("scope", key), zap.Error(err))
					}
				}
			}
		}
	}
}

func hasPollSpec(conn *types.Connection) bool {
	for _, res := range conn.Resources {
		for _, wh := range res.Webhooks {
			if wh.Poll != nil && wh.Poll.List != nil {
				return true
			}
		}
	}
	return false
}

// runnerFor builds the authenticated HTTP runner for a connection — Google
// wrapper for OAuth/google connectors, generic REST driver otherwise.
func (svc *configuredConnection) runnerFor(ctx context.Context, cc *types.ConfiguredConnection) (connectionRunner, error) {
	resolved := svc.resolveTemplates(&cc.Connection, cc.Config.Params)
	if activeAuth(&cc.Connection, cc).Method == "oauth2_authorization_code" ||
		strings.Contains(resolved.Service.BaseURL.Value, "googleapis.com") {
		ensureCredential(ctx, svc.store, cc, &cc.Connection)
		return google.NewWrapper(resolved.Service.BaseURL.Value, cc.ID), nil
	}
	return restDriver.RunnerFromConnection(resolved, cc.ID)
}

// pollWebhook runs one poll pass. scopeVars holds a trigger scope's constraint
// values (nil for connection-level polls); they are injected into the poll
// action's templates and merged into every dispatched event so the trigger's
// constraints match. cursorKey isolates this scope's cursor from others.
func (svc *configuredConnection) pollWebhook(ctx context.Context, cc *types.ConfiguredConnection, runner connectionRunner, wh types.ConnectionWebhook, scopeVars map[string]string, cursorKey string) error {
	poll := wh.Poll
	cursor := readPollCursor(cc, cursorKey)

	switch poll.CursorMode {
	case "count":
		return svc.pollWebhookCount(ctx, cc, runner, wh, scopeVars, cursorKey, cursor)
	case "snapshot":
		return svc.pollWebhookSnapshot(ctx, cc, runner, wh, scopeVars, cursorKey)
	case "notfound":
		return svc.pollWebhookNotFound(ctx, cc, runner, wh, scopeVars, cursorKey, cursor)
	}

	if cursor == "" {
		// First poll: record a baseline without firing on pre-existing items.
		nc, err := runPollBaseline(ctx, runner, poll, scopeVars)
		if err != nil {
			return err
		}
		return svc.savePollCursor(ctx, cc.ID, cursorKey, nc)
	}

	items, newCursor, err := collectPollItems(ctx, runner, poll, scopeVars, cursor)
	if err != nil {
		if errors.Is(err, errPollCursorStale) {
			if nc, berr := runPollBaseline(ctx, runner, poll, scopeVars); berr == nil {
				return svc.savePollCursor(ctx, cc.ID, cursorKey, nc)
			}
		}
		return err
	}

	for _, item := range items {
		if !pollItemPasses(item, poll.Filter) {
			continue
		}
		item = enrichPollItem(ctx, runner, poll, item)
		svc.dispatchPollItem(ctx, cc, wh, item, scopeVars)
	}

	if newCursor != "" && newCursor != cursor {
		return svc.savePollCursor(ctx, cc.ID, cursorKey, newCursor)
	}
	return nil
}

// pollWebhookCount handles append-only lists whose cursor is the item count
// (no provider-side cursor field) — e.g. spreadsheet rows. New items are those
// past the stored count.
func (svc *configuredConnection) pollWebhookCount(ctx context.Context, cc *types.ConfiguredConnection, runner connectionRunner, wh types.ConnectionWebhook, scopeVars map[string]string, cursorKey, cursor string) error {
	poll := wh.Poll
	resp, _, err := runPollAction(ctx, runner, *poll.List, pollVars(scopeVars, ""))
	if err != nil {
		return err
	}
	count := len(collectByPath(resp, poll.Items))

	// Baseline (first poll) or a shrunk list (rows removed): record and move on.
	if cursor == "" {
		return svc.savePollCursor(ctx, cc.ID, cursorKey, strconv.Itoa(count))
	}
	prev, _ := strconv.Atoi(cursor)
	if count <= prev {
		if count < prev {
			return svc.savePollCursor(ctx, cc.ID, cursorKey, strconv.Itoa(count))
		}
		return nil
	}

	items := collectByPath(resp, poll.Items)
	for _, item := range items[prev:] {
		if !pollItemPasses(item, poll.Filter) {
			continue
		}
		svc.dispatchPollItem(ctx, cc, wh, item, scopeVars)
	}
	return svc.savePollCursor(ctx, cc.ID, cursorKey, strconv.Itoa(count))
}

// pollWebhookSnapshot diffs the listed items against the previous poll's
// snapshot (keyed per scope) and fires for the diff kind the webhook declares:
// "added" for new keys, "removed" for vanished ones, "changed" for keys whose
// body differs. Covers create/update/delete style events with no provider feed.
func (svc *configuredConnection) pollWebhookSnapshot(ctx context.Context, cc *types.ConfiguredConnection, runner connectionRunner, wh types.ConnectionWebhook, scopeVars map[string]string, cursorKey string) error {
	poll := wh.Poll
	resp, _, err := runPollAction(ctx, runner, *poll.List, pollVars(scopeVars, ""))
	if err != nil {
		return err
	}
	newSnap, bodies := buildSnapshot(collectByPath(resp, poll.Items), poll.Key)

	old := readSnapshot(cc, cursorKey)
	if old == nil {
		// First poll: record a baseline without firing on pre-existing items.
		return svc.saveSnapshot(ctx, cc.ID, cursorKey, newSnap)
	}

	added, removed, changed := diffSnapshots(old, newSnap)
	switch poll.Emit {
	case "added":
		for _, k := range added {
			svc.dispatchPollItem(ctx, cc, wh, bodies[k], scopeVars)
		}
	case "changed":
		for _, k := range changed {
			svc.dispatchPollItem(ctx, cc, wh, bodies[k], scopeVars)
		}
	case "removed":
		for _, k := range removed {
			var item any
			_ = json.Unmarshal([]byte(old[k]), &item)
			svc.dispatchPollItem(ctx, cc, wh, item, scopeVars)
		}
	}
	return svc.saveSnapshot(ctx, cc.ID, cursorKey, newSnap)
}

// pollWebhookNotFound fires once when a resource that previously existed starts
// returning 404 — e.g. a deleted spreadsheet. The cursor is a presence sentinel.
func (svc *configuredConnection) pollWebhookNotFound(ctx context.Context, cc *types.ConfiguredConnection, runner connectionRunner, wh types.ConnectionWebhook, scopeVars map[string]string, cursorKey, cursor string) error {
	_, status, err := runPollAction(ctx, runner, *wh.Poll.List, pollVars(scopeVars, ""))
	if status == 0 {
		return err // transport error — retry next cycle, do not change state
	}
	exists := status < 400

	if cursor == "" {
		return svc.savePollCursor(ctx, cc.ID, cursorKey, presence(exists))
	}
	if cursor == "exists" && !exists {
		svc.dispatchPollItem(ctx, cc, wh, map[string]any{}, scopeVars)
		return svc.savePollCursor(ctx, cc.ID, cursorKey, "gone")
	}
	if cursor == "gone" && exists {
		return svc.savePollCursor(ctx, cc.ID, cursorKey, "exists")
	}
	return nil
}

func presence(exists bool) string {
	if exists {
		return "exists"
	}
	return "gone"
}

// buildSnapshot indexes items by their key path (falling back to position when
// the key is absent), returning both the canonical-JSON snapshot for diffing and
// the parsed item bodies for dispatch.
func buildSnapshot(items []any, keyPath []string) (map[string]string, map[string]any) {
	snap := make(map[string]string, len(items))
	bodies := make(map[string]any, len(items))
	for i, item := range items {
		k := strconv.Itoa(i)
		if len(keyPath) > 0 {
			if ks := extractString(item, keyPath); ks != "" {
				k = ks
			}
		}
		b, _ := json.Marshal(item)
		snap[k] = string(b)
		bodies[k] = item
	}
	return snap, bodies
}

// diffSnapshots classifies keys between the previous and current snapshot:
// present only now (added), only before (removed), or present in both with a
// differing body (changed).
func diffSnapshots(old, cur map[string]string) (added, removed, changed []string) {
	for k, v := range cur {
		if prev, ok := old[k]; !ok {
			added = append(added, k)
		} else if prev != v {
			changed = append(changed, k)
		}
	}
	for k := range old {
		if _, ok := cur[k]; !ok {
			removed = append(removed, k)
		}
	}
	return
}

const pollSnapshotPrefix = "pollSnapshot:"

// readSnapshot returns the stored snapshot for a scope, or nil if none (baseline).
func readSnapshot(cc *types.ConfiguredConnection, key string) map[string]string {
	raw, ok := cc.Config.Discovery[pollSnapshotPrefix+key]
	if !ok {
		return nil
	}
	var m map[string]string
	if json.Unmarshal(raw, &m) != nil {
		return nil
	}
	return m
}

// saveSnapshot persists a scope's snapshot, reloading first so a concurrent
// config write is not clobbered.
func (svc *configuredConnection) saveSnapshot(ctx context.Context, ccID uint64, key string, snap map[string]string) error {
	cc, err := loadConfiguredConnection(ctx, svc.store, ccID)
	if err != nil {
		return err
	}
	if cc.Config.Discovery == nil {
		cc.Config.Discovery = map[string]json.RawMessage{}
	}
	raw, _ := json.Marshal(snap)
	cc.Config.Discovery[pollSnapshotPrefix+key] = raw
	return store.UpdateConfiguredConnection(ctx, svc.store, cc)
}

// dispatchPollItem builds the event vars (payload selectors + scope values) and
// fires the webhook event.
func (svc *configuredConnection) dispatchPollItem(ctx context.Context, cc *types.ConfiguredConnection, wh types.ConnectionWebhook, item any, scopeVars map[string]string) {
	vars := pollItemVars(item, cc.ID, wh.Payload)
	// Scope values (e.g. spreadsheetId/sheetName) are not in the item itself;
	// merge them so the trigger's constraints match and downstream steps see them.
	for k, v := range scopeVars {
		if _, ok := vars[k]; !ok {
			vars[k] = v
		}
	}
	ev := sysEvent.ConnectionWebhookEvent(cc.ConnectionID, cc.ID, wh.Event, vars)
	if derr := eventbus.Service().WaitFor(ctx, ev); derr != nil {
		svc.services.logger.Warn("connector poll: dispatch failed",
			zap.Uint64("ccID", cc.ID), zap.String("event", wh.Event), zap.Error(derr))
	}
}

// pollVars merges a trigger scope with the current cursor into one var map.
func pollVars(scopeVars map[string]string, cursor string) map[string]any {
	vars := map[string]any{"cursor": cursor}
	for k, v := range scopeVars {
		vars[k] = v
	}
	return vars
}

// runPollBaseline establishes the initial cursor without emitting events.
func runPollBaseline(ctx context.Context, runner connectionRunner, poll *types.ConnectionWebhookPoll, scopeVars map[string]string) (string, error) {
	if poll.Seed != nil {
		resp, _, err := runPollAction(ctx, runner, *poll.Seed, pollVars(scopeVars, ""))
		if err != nil {
			return "", err
		}
		return extractString(resp, poll.SeedCursor), nil
	}
	resp, _, err := runPollAction(ctx, runner, *poll.List, pollVars(scopeVars, ""))
	if err != nil {
		return "", err
	}
	return extractString(resp, poll.Cursor), nil
}

// collectPollItems runs the list call (following pagination) and returns the new
// items plus the next cursor.
func collectPollItems(ctx context.Context, runner connectionRunner, poll *types.ConnectionWebhookPoll, scopeVars map[string]string, cursor string) ([]any, string, error) {
	resp, status, err := runPollAction(ctx, runner, *poll.List, pollVars(scopeVars, cursor))
	if err != nil {
		if status == 404 {
			return nil, "", errPollCursorStale
		}
		return nil, "", err
	}

	items := collectByPath(resp, poll.Items)
	newCursor := extractString(resp, poll.Cursor)

	if len(poll.PageToken) > 0 && poll.PageParam != "" {
		token := extractString(resp, poll.PageToken)
		for guard := 0; token != "" && guard < 50; guard++ {
			page := *poll.List
			qp := map[string]types.ConnectionTemplate{}
			for k, v := range poll.List.QueryParams {
				qp[k] = v
			}
			qp[poll.PageParam] = types.ConnectionTemplate{Value: token}
			page.QueryParams = qp

			pr, _, perr := runPollAction(ctx, runner, page, pollVars(scopeVars, cursor))
			if perr != nil {
				break
			}
			items = append(items, collectByPath(pr, poll.Items)...)
			token = extractString(pr, poll.PageToken)
		}
	}

	return items, newCursor, nil
}

// actionPlaceholders returns the set of {{name}} placeholders referenced in an
// HTTP action's path, query, headers and body.
func actionPlaceholders(action *types.ConnectionHTTPAction) map[string]bool {
	set := map[string]bool{}
	add := func(s string) {
		for _, m := range placeholderRe.FindAllStringSubmatch(s, -1) {
			set[m[1]] = true
		}
	}
	add(action.Path.Value)
	for _, t := range action.QueryParams {
		add(t.Value)
	}
	for _, t := range action.Headers {
		add(t.Value)
	}
	add(action.BodyTemplate.Value)
	return set
}

// scopeSig builds a stable cursor-key suffix from a scope's values for the
// given parameters, so each distinct spreadsheet/tab keeps its own cursor.
func scopeSig(scope map[string]string, params map[string]bool) string {
	keys := make([]string, 0, len(params))
	for p := range params {
		keys = append(keys, p)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+scope[k])
	}
	return strings.Join(parts, "|")
}

// triggerScopes returns the distinct constraint-value sets of the enabled
// triggers bound to this configured connection for the given event, limited to
// scopes that supply every poll parameter. Each set drives one scoped poll.
func (svc *configuredConnection) triggerScopes(ctx context.Context, connID, ccID uint64, event string, params map[string]bool) []map[string]string {
	aa, _, err := store.SearchAutomationNgAutomations(ctx, svc.store, atypes.NgAutomationFilter{})
	if err != nil {
		return nil
	}
	rt := connectionWebhookResourceType(connID)
	ccIDStr := strconv.FormatUint(ccID, 10)
	seen := map[string]bool{}
	var out []map[string]string
	for _, au := range aa {
		if !au.Enabled {
			continue
		}
		for _, tr := range au.Triggers {
			if !tr.Enabled || tr.ResourceType != rt || tr.EventType != event {
				continue
			}
			scope := map[string]string{}
			for _, c := range tr.Constraints {
				if len(c.Values) > 0 {
					scope[c.Name] = c.Values[0].Value
				}
			}
			if scope["configurationID"] != ccIDStr {
				continue
			}
			complete := true
			for p := range params {
				if scope[p] == "" {
					complete = false
					break
				}
			}
			if !complete {
				continue
			}
			sig := scopeSig(scope, params)
			if seen[sig] {
				continue
			}
			seen[sig] = true
			out = append(out, scope)
		}
	}
	return out
}

func enrichPollItem(ctx context.Context, runner connectionRunner, poll *types.ConnectionWebhookPoll, item any) any {
	if poll.Enrich == nil {
		return item
	}
	key := poll.EnrichKey
	if key == "" {
		key = "id"
	}
	id := extractString(item, []string{key})
	if id == "" {
		return item
	}
	if enriched, _, err := runPollAction(ctx, runner, *poll.Enrich, map[string]any{"id": id}); err == nil && enriched != nil {
		return enriched
	}
	return item
}

// runPollAction resolves and executes one HTTP action, returning its parsed JSON
// body and status. Numbers are decoded as json.Number to avoid precision loss on
// large cursors (e.g. Gmail historyId).
func runPollAction(ctx context.Context, runner connectionRunner, action types.ConnectionHTTPAction, vars map[string]any) (any, int, error) {
	path, headers, payload, err := buildHTTPRequest(action, nil, vars)
	if err != nil {
		return nil, 0, err
	}
	method := action.Method
	if method == "" {
		method = "GET"
	}

	status, _, body, err := runner.Run(ctx, method, path, payload, headers)
	if err != nil {
		return nil, status, err
	}
	if status >= 400 {
		return nil, status, fmt.Errorf("poll action %s %s failed: HTTP %d", method, path, status)
	}
	if len(body) == 0 {
		return nil, status, nil
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var data any
	if err := dec.Decode(&data); err != nil {
		return nil, status, err
	}
	return data, status, nil
}

func pollItemVars(item any, ccID uint64, payload []types.ConnectionWebhookField) map[string]any {
	vars := map[string]any{"configuredConnectionID": ccID}
	for _, f := range payload {
		if len(f.Selector) == 0 {
			continue
		}
		if v := extractByPath(item, f.Selector); v != nil {
			vars[f.Name] = v
		}
	}
	return vars
}

func pollItemPasses(item any, filter *types.ConnectionPollFilter) bool {
	if filter == nil || len(filter.Field) == 0 {
		return true
	}
	switch v := extractByPath(item, filter.Field).(type) {
	case string:
		return v == filter.Contains
	case []any:
		for _, e := range v {
			if s, ok := e.(string); ok && s == filter.Contains {
				return true
			}
		}
	}
	return false
}

// collectByPath walks path, flattening on "*" segments (iterate every array
// element). Non-* segments index maps or a single array element via
// extractByPath. Returns a flat slice of matched values.
func collectByPath(data any, path []string) []any {
	cur := []any{data}
	for _, seg := range path {
		var next []any
		for _, node := range cur {
			if seg == "*" {
				if arr, ok := node.([]any); ok {
					next = append(next, arr...)
				}
				continue
			}
			if v := extractByPath(node, []string{seg}); v != nil {
				next = append(next, v)
			}
		}
		cur = next
	}
	return cur
}

func extractString(data any, path []string) string {
	if len(path) == 0 {
		return ""
	}
	switch v := extractByPath(data, path).(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatInt(int64(v), 10)
	}
	return ""
}

func readPollCursor(cc *types.ConfiguredConnection, event string) string {
	raw, ok := cc.Config.Discovery[pollCursorPrefix+event]
	if !ok {
		return ""
	}
	var s string
	_ = json.Unmarshal(raw, &s)
	return s
}

// savePollCursor persists the cursor, reloading first so a concurrent config
// write (e.g. discovery refresh) is not clobbered.
func (svc *configuredConnection) savePollCursor(ctx context.Context, ccID uint64, event, cursor string) error {
	cc, err := loadConfiguredConnection(ctx, svc.store, ccID)
	if err != nil {
		return err
	}
	if cc.Config.Discovery == nil {
		cc.Config.Discovery = map[string]json.RawMessage{}
	}
	raw, _ := json.Marshal(cursor)
	cc.Config.Discovery[pollCursorPrefix+event] = raw
	return store.UpdateConfiguredConnection(ctx, svc.store, cc)
}

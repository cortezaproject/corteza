package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

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
				if err := svc.pollWebhook(ctx, cc, runner, wh); err != nil {
					svc.services.logger.Warn("connector poll: webhook failed",
						zap.Uint64("ccID", cc.ID), zap.String("event", wh.Event), zap.Error(err))
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

func (svc *configuredConnection) pollWebhook(ctx context.Context, cc *types.ConfiguredConnection, runner connectionRunner, wh types.ConnectionWebhook) error {
	poll := wh.Poll
	cursor := readPollCursor(cc, wh.Event)

	if cursor == "" {
		// First poll: record a baseline without firing on pre-existing items.
		nc, err := runPollBaseline(ctx, runner, poll)
		if err != nil {
			return err
		}
		return svc.savePollCursor(ctx, cc.ID, wh.Event, nc)
	}

	items, newCursor, err := collectPollItems(ctx, runner, poll, cursor)
	if err != nil {
		if errors.Is(err, errPollCursorStale) {
			if nc, berr := runPollBaseline(ctx, runner, poll); berr == nil {
				return svc.savePollCursor(ctx, cc.ID, wh.Event, nc)
			}
		}
		return err
	}

	for _, item := range items {
		if !pollItemPasses(item, poll.Filter) {
			continue
		}
		item = enrichPollItem(ctx, runner, poll, item)
		vars := pollItemVars(item, cc.ID, wh.Payload)
		ev := sysEvent.ConnectionWebhookEvent(cc.ConnectionID, cc.ID, wh.Event, vars)
		if derr := eventbus.Service().WaitFor(ctx, ev); derr != nil {
			svc.services.logger.Warn("connector poll: dispatch failed",
				zap.Uint64("ccID", cc.ID), zap.String("event", wh.Event), zap.Error(derr))
		}
	}

	if newCursor != "" && newCursor != cursor {
		return svc.savePollCursor(ctx, cc.ID, wh.Event, newCursor)
	}
	return nil
}

// runPollBaseline establishes the initial cursor without emitting events.
func runPollBaseline(ctx context.Context, runner connectionRunner, poll *types.ConnectionWebhookPoll) (string, error) {
	if poll.Seed != nil {
		resp, _, err := runPollAction(ctx, runner, *poll.Seed, nil)
		if err != nil {
			return "", err
		}
		return extractString(resp, poll.SeedCursor), nil
	}
	resp, _, err := runPollAction(ctx, runner, *poll.List, map[string]any{"cursor": ""})
	if err != nil {
		return "", err
	}
	return extractString(resp, poll.Cursor), nil
}

// collectPollItems runs the list call (following pagination) and returns the new
// items plus the next cursor.
func collectPollItems(ctx context.Context, runner connectionRunner, poll *types.ConnectionWebhookPoll, cursor string) ([]any, string, error) {
	resp, status, err := runPollAction(ctx, runner, *poll.List, map[string]any{"cursor": cursor})
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

			pr, _, perr := runPollAction(ctx, runner, page, map[string]any{"cursor": cursor})
			if perr != nil {
				break
			}
			items = append(items, collectByPath(pr, poll.Items)...)
			token = extractString(pr, poll.PageToken)
		}
	}

	return items, newCursor, nil
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

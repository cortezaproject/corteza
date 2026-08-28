package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/jmoiron/sqlx/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for these handlers are in reminder_tools.go, in the same order.
type reminderHandler struct {
	reg toolRegistrar
}

func ReminderHandler(reg toolRegistrar) *reminderHandler {
	h := &reminderHandler{reg: reg}
	h.register()
	return h
}

// reminderItem is the compact projection returned by list mode. Single-item
// lookups return the service type unchanged; lists must not, because a payload
// per row would flood the caller's context.
type reminderItem struct {
	ReminderID  string     `json:"reminderID"`
	Resource    string     `json:"resource"`
	RemindAt    *time.Time `json:"remindAt"`
	AssignedTo  string     `json:"assignedTo"`
	DismissedAt *time.Time `json:"dismissedAt"`
}

func (h *reminderHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// Reference param is never Required, so an empty ref means "list".
	if id, err := toolkit.ID(args, "reminderID"); err != nil {
		return nil, err
	} else if id > 0 {
		rm, err := sysService.DefaultReminder.FindByID(ctx, id)
		if err != nil {
			return nil, toolkit.Errf("reminder lookup", err)
		}
		return toolkit.JSONResult(rm)
	}

	f := sysTypes.ReminderFilter{
		Resource:         toolkit.Str(args, "resource"),
		ExcludeDismissed: toolkit.Bool(args, "excludeDismissed"),
		IncludeDeleted:   toolkit.Bool(args, "includeDeleted"),
		ScheduledOnly:    toolkit.Bool(args, "scheduledOnly"),
	}

	if f.AssignedTo, err = toolkit.ID(args, "assignedTo"); err != nil {
		return nil, err
	}
	if f.ScheduledFrom, err = optTime(args, "scheduledFrom"); err != nil {
		return nil, err
	}
	if f.ScheduledUntil, err = optTime(args, "scheduledUntil"); err != nil {
		return nil, err
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := sysService.DefaultReminder.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("reminder list", err)
	}

	items := make([]reminderItem, 0, len(set))
	for _, rm := range set {
		items = append(items, reminderItem{
			ReminderID:  fmt.Sprintf("%d", rm.ID),
			Resource:    rm.Resource,
			RemindAt:    rm.RemindAt,
			AssignedTo:  fmt.Sprintf("%d", rm.AssignedTo),
			DismissedAt: rm.DismissedAt,
		})
	}

	// Pass the cursor itself, not its String(): String is a human-readable debug
	// rendering ("<id: 123, [FWD]>") while parseCursor expects base64 of the
	// cursor JSON, which is what MarshalJSON emits. json.Marshal renders a nil
	// pointer as null, so the last page is safe.
	return toolkit.JSONResult(map[string]any{
		"reminders":      items,
		"nextPageCursor": out.NextPage,
	})
}

func (h *reminderHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	resource, err := toolkit.ReqStr(args, "resource")
	if err != nil {
		return nil, err
	}

	rm := &sysTypes.Reminder{Resource: resource}

	if rm.AssignedTo, err = toolkit.ID(args, "assignedTo"); err != nil {
		return nil, err
	}
	if rm.RemindAt, err = optTime(args, "remindAt"); err != nil {
		return nil, err
	}
	payload, _, err := reminderPayload(args)
	if err != nil {
		return nil, err
	}
	if len(payload) > 0 {
		rm.Payload = payload
	}

	rm, err = sysService.DefaultReminder.Create(ctx, rm)
	if err != nil {
		return nil, toolkit.Errf("reminder creation", err)
	}
	return toolkit.JSONResult(rm)
}

func (h *reminderHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	id, err := toolkit.ReqID(args, "reminderID")
	if err != nil {
		return nil, err
	}

	rm, err := sysService.DefaultReminder.FindByID(ctx, id)
	if err != nil {
		return nil, toolkit.Errf("reminder lookup", err)
	}

	// Absent leaves the field alone; present-and-empty clears it.
	if v, ok := args["resource"]; ok {
		rm.Resource, _ = v.(string)
	}
	if payload, present, pErr := reminderPayload(args); pErr != nil {
		return nil, pErr
	} else if present {
		rm.Payload = payload
	}
	if v, ok := args["remindAt"]; ok {
		if s, _ := v.(string); s == "" {
			rm.RemindAt = nil
		} else if rm.RemindAt, err = optTime(args, "remindAt"); err != nil {
			return nil, err
		}
	}

	rm, err = sysService.DefaultReminder.Update(ctx, rm)
	if err != nil {
		return nil, toolkit.Errf("reminder update", err)
	}
	return toolkit.JSONResult(rm)
}

func (h *reminderHandler) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := reminderID(req)
	if err != nil {
		return nil, err
	}
	if err := sysService.DefaultReminder.DeleteByID(ctx, id); err != nil {
		return nil, toolkit.Errf("reminder deletion", err)
	}
	return toolkit.TextResult("reminder %d deleted", id), nil
}

func (h *reminderHandler) dismiss(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := reminderID(req)
	if err != nil {
		return nil, err
	}
	if err := sysService.DefaultReminder.Dismiss(ctx, id); err != nil {
		return nil, toolkit.Errf("reminder dismissal", err)
	}
	return toolkit.TextResult("reminder %d dismissed", id), nil
}

func (h *reminderHandler) undismiss(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := reminderID(req)
	if err != nil {
		return nil, err
	}
	if err := sysService.DefaultReminder.Undismiss(ctx, id); err != nil {
		return nil, toolkit.Errf("reminder undismissal", err)
	}
	return toolkit.TextResult("reminder %d undismissed", id), nil
}

func (h *reminderHandler) snooze(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	id, err := toolkit.ReqID(args, "reminderID")
	if err != nil {
		return nil, err
	}

	remindAt, err := optTime(args, "remindAt")
	if err != nil {
		return nil, err
	}
	if remindAt == nil {
		return nil, fmt.Errorf("remindAt is required")
	}

	if err := sysService.DefaultReminder.Snooze(ctx, id, remindAt); err != nil {
		return nil, toolkit.Errf("reminder snooze", err)
	}
	return toolkit.TextResult("reminder %d snoozed until %s", id, remindAt.Format(time.RFC3339)), nil
}

// reminderID unwraps the single required argument the domain ops share.
func reminderID(req mcp.CallToolRequest) (uint64, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return 0, err
	}
	return toolkit.ReqID(args, "reminderID")
}

// reminderPayload reads the reminder's contents, whether the caller sent them
// as a JSON object or as the JSON text of one.
//
// The parameter is declared a string and described as "JSON object", so an
// object is what a caller reaches for — and reading it with Str answered "" for
// anything that was not already a string. The reminder then stored an empty
// payload, the create reported success, and a reminder whose whole content was
// that object came back holding nothing.
//
// Invalid JSON is refused rather than stored: the column is JSON, and text that
// is not fails on the way back out, where nothing can say what put it there.
func reminderPayload(args map[string]any) (out types.JSONText, present bool, err error) {
	raw, ok := args["payload"]
	if !ok || raw == nil {
		return nil, false, nil
	}

	var buf []byte

	switch v := raw.(type) {
	case string:
		// Present and empty clears it, which is how every other optional
		// argument here spells "remove this".
		if strings.TrimSpace(v) == "" {
			return types.JSONText(""), true, nil
		}
		buf = []byte(v)
	case map[string]any, []any:
		if buf, err = json.Marshal(v); err != nil {
			return nil, false, fmt.Errorf("payload cannot be encoded: %w", err)
		}
	default:
		return nil, false, fmt.Errorf(
			"payload is a %T; send a JSON object with the reminder's contents, or the JSON text of one", raw)
	}

	if !json.Valid(buf) {
		return nil, false, fmt.Errorf(
			"payload is not valid JSON: %s. Send a JSON object with the reminder's contents, or the JSON text of one",
			string(buf))
	}

	return types.JSONText(buf), true, nil
}

// optTime parses an optional RFC3339 argument. Absent or empty yields nil
// rather than an error, so callers can clear a time by passing "".
func optTime(args map[string]any, key string) (*time.Time, error) {
	s := toolkit.Str(args, key)
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: expected RFC3339 (e.g. 2026-08-03T09:00:00Z), got %q", key, s)
	}
	return &t, nil
}

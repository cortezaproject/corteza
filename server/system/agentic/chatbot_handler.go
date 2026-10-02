package agentic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for these handlers are in chatbot_tools.go, in the same order.
type chatbotHandler struct {
	reg toolRegistrar
}

func ChatbotHandler(reg toolRegistrar) *chatbotHandler {
	h := &chatbotHandler{reg: reg}
	h.register()
	return h
}

// chatbotItem is the list projection.
//
// The widget key is deliberately absent. It is not a secret — it ends up in the
// HTML of whatever page hosts the widget, and allowedOrigins is what actually
// restricts use — but a listing has no reason to hand out every embed key on the
// instance at once. A single-chatbot lookup returns it, matching REST.
type chatbotItem struct {
	ID      uint64 `json:"chatbotID,string"`
	Handle  string `json:"handle,omitempty"`
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled"`
}

// findChatbot resolves a chatbot by ID or handle. There is no FindByAny, and
// ChatbotFilter carries no name field, so the name is not resolvable.
func findChatbot(ctx context.Context, ref string) (*sysTypes.Chatbot, error) {
	if ref == "" {
		return nil, toolkit.Requiredf("chatbot is required")
	}

	if id, err := strconv.ParseUint(ref, 10, 64); err == nil && id > 0 {
		c, err := sysService.DefaultChatbot.FindByID(ctx, id)
		if err != nil {
			return nil, toolkit.Errf("chatbot lookup", err)
		}
		return c, nil
	}

	set, _, err := sysService.DefaultChatbot.Search(ctx, sysTypes.ChatbotFilter{
		Handle: ref,
		Paging: filter.Paging{Limit: 2},
	})
	if err != nil {
		return nil, toolkit.Errf("chatbot lookup", err)
	}

	switch len(set) {
	case 0:
		return nil, fmt.Errorf("no chatbot with handle %q — call system_chatbot_lookup with no arguments to list them", ref)
	case 1:
		return set[0], nil
	default:
		return nil, fmt.Errorf("handle %q matches more than one chatbot; use the chatbotID instead", ref)
	}
}

func (h *chatbotHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	if ref := toolkit.Str(args, "chatbot"); ref != "" {
		c, err := findChatbot(ctx, ref)
		if err != nil {
			return nil, err
		}
		return toolkit.JSONResultWith(c, chatbotLinks(c))
	}

	f := sysTypes.ChatbotFilter{
		Query:  toolkit.Str(args, "query"),
		Handle: toolkit.Str(args, "handle"),
	}
	if toolkit.Bool(args, "includeDeleted") {
		f.Deleted = filter.StateInclusive
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := sysService.DefaultChatbot.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("chatbot search", err)
	}

	items := make([]chatbotItem, 0, len(set))
	for _, c := range set {
		items = append(items, chatbotItem{ID: c.ID, Handle: c.Handle, Name: c.Name, Enabled: c.Enabled})
	}

	return toolkit.JSONResult(map[string]any{
		"chatbots":       items,
		"nextPageCursor": out.NextPage,
	})
}

// chatbotReadinessNote names what will stop the widget answering, for the
// conditions a person can legitimately fix afterwards.
//
// A conversation scenario is driven by the system rather than by the visitor:
// the session runs the agent as the scenario's runAs user. A scenario without
// one is therefore unreachable, and the visitor gets "agent not configured" on
// their first message while the saved chatbot looks entirely healthy. Building
// the chatbot before choosing the user is a reasonable order to work in, so
// this is said rather than refused.
func chatbotReadinessNote(c *sysTypes.Chatbot) string {
	if c == nil {
		return ""
	}

	var blocked []string
	for _, sc := range c.Scenarios {
		if sc.Type != "conversation" || sc.AgentID == 0 {
			continue
		}

		if sc.RunAs == 0 {
			blocked = append(blocked, fmt.Sprintf("%q", sc.ID))
		}
	}

	if len(blocked) == 0 {
		return ""
	}

	return "Stored, but this chatbot cannot answer yet: scenario " +
		strings.Join(blocked, ", ") +
		" names no runAs user. A chatbot runs its agent as that user, so set \"runAs\" (a user ID) on the " +
		"scenario or the visitor's first message fails with \"agent not configured\"."
}

func (h *chatbotHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	name, err := toolkit.ReqStr(args, "name")
	if err != nil {
		return nil, err
	}

	c := &sysTypes.Chatbot{
		Name:       name,
		Handle:     toolkit.Str(args, "handle"),
		Enabled:    toolkit.Bool(args, "enabled"),
		SessionTTL: toolkit.Str(args, "sessionTTL"),
	}

	if err = applyChatbotSections(c, args); err != nil {
		return nil, err
	}

	c, err = sysService.DefaultChatbot.Create(ctx, c)
	if err != nil {
		return nil, toolkit.Errf("chatbot creation", err)
	}

	return toolkit.JSONResultWith(c, withNote(chatbotLinks(c), chatbotReadinessNote(c)))
}

func (h *chatbotHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "chatbot")
	if err != nil {
		return nil, err
	}

	// Loaded first because the service writes the whole record.
	c, err := findChatbot(ctx, ref)
	if err != nil {
		return nil, err
	}

	if v, ok := toolkit.OptStr(args, "name"); ok && v != "" {
		c.Name = v
	}
	if v, ok := toolkit.OptStr(args, "handle"); ok {
		c.Handle = v
	}
	if v, ok := toolkit.OptBool(args, "enabled"); ok {
		c.Enabled = v
	}
	if v, ok := toolkit.OptStr(args, "sessionTTL"); ok {
		c.SessionTTL = v
	}

	if err = applyChatbotSections(c, args); err != nil {
		return nil, err
	}

	c, err = sysService.DefaultChatbot.Update(ctx, c)
	if err != nil {
		return nil, toolkit.Errf("chatbot update", err)
	}

	return toolkit.JSONResultWith(c, withNote(chatbotLinks(c), chatbotReadinessNote(c)))
}

func (h *chatbotHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "chatbot")
	if err != nil {
		return nil, err
	}

	c, err := findChatbot(ctx, ref)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultChatbot.DeleteByID(ctx, c.ID); err != nil {
		return nil, toolkit.Errf("chatbot delete", err)
	}

	return toolkit.TextResult("chatbot %d deleted — restore it with system_chatbot_undelete chatbotID=%d", c.ID, c.ID), nil
}

func (h *chatbotHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	id, err := toolkit.ReqID(args, "chatbotID")
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultChatbot.UndeleteByID(ctx, id); err != nil {
		return nil, toolkit.Errf("chatbot undelete", err)
	}

	return toolkit.TextResult("chatbot %d restored", id), nil
}

func (h *chatbotHandler) regenerateKey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "chatbot")
	if err != nil {
		return nil, err
	}

	c, err := findChatbot(ctx, ref)
	if err != nil {
		return nil, err
	}

	c, err = sysService.DefaultChatbot.RegenerateWidgetKey(ctx, c.ID)
	if err != nil {
		return nil, toolkit.Errf("chatbot widget key regeneration", err)
	}

	return toolkit.JSONResultWith(c, chatbotLinks(c))
}

// applyChatbotSections replaces each configuration section the caller sent and
// leaves the rest alone (§8.2). Same shape as applyAgentSections, and for the
// same two reasons: decoding over the existing value would merge, and zeroing
// before the presence check would wipe a section on an empty argument.
func applyChatbotSections(c *sysTypes.Chatbot, args map[string]any) error {
	var (
		handoff   sysTypes.ChatbotHandoff
		styling   sysTypes.ChatbotStyling
		scenarios sysTypes.ChatbotScenarios
		origins   sysTypes.ChatbotAllowedOrigins
	)

	sections := []struct {
		key     string
		subject string
		into    any
		assign  func()
	}{
		{"handoff", "chatbot handoff", &handoff, func() { c.Handoff = handoff }},
		{"styling", "chatbot styling", &styling, func() { c.Styling = styling }},
		{"scenarios", "chatbot scenarios", &scenarios, func() { c.Scenarios = scenarios }},
		{"allowedOrigins", "chatbot allowedOrigins", &origins, func() { c.AllowedOrigins = origins }},
	}

	for _, s := range sections {
		present, err := jsonSection(args, s.key, s.subject, s.into)
		if err != nil {
			return err
		}
		if present {
			s.assign()
		}
	}

	return nil
}

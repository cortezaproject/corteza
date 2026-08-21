package agentic

import (
	"context"
	"fmt"
	"strconv"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for these handlers are in agent_tools.go, in the same order.
type agentHandler struct {
	reg toolRegistrar
}

func AgentHandler(reg toolRegistrar) *agentHandler {
	h := &agentHandler{reg: reg}
	h.register()
	return h
}

// agentItem is the list projection. The five configuration sections are the
// heavy part of an agent and none of them helps pick one out of a list, so a
// listing carries only what identifies it. A single-agent lookup returns the
// service type whole.
type agentItem struct {
	ID     uint64 `json:"agentID,string"`
	Handle string `json:"handle,omitempty"`
	Short  string `json:"short,omitempty"`
	Status string `json:"status,omitempty"`
}

// findAgent resolves an agent by ID or handle.
//
// There is no FindByAny for agents, so this is the strategy the brief fixes:
// all-digits is an ID, anything else is a handle. Deliberately not falling back
// to Meta.Short — it is a label, is not unique, and nothing indexes it, so
// resolving by it would pick an arbitrary agent.
func findAgent(ctx context.Context, ref string) (*sysTypes.Agent, error) {
	if ref == "" {
		return nil, fmt.Errorf("agent is required")
	}

	if id, err := strconv.ParseUint(ref, 10, 64); err == nil && id > 0 {
		a, err := sysService.DefaultAgent.FindByID(ctx, id)
		if err != nil {
			return nil, toolkit.Errf("agent lookup", err)
		}
		return a, nil
	}

	set, _, err := sysService.DefaultAgent.Search(ctx, sysTypes.AgentFilter{
		Handle: ref,
		Paging: filter.Paging{Limit: 2},
	})
	if err != nil {
		return nil, toolkit.Errf("agent lookup", err)
	}

	switch len(set) {
	case 0:
		return nil, fmt.Errorf("no agent with handle %q — call system_agent_lookup with no arguments to list them", ref)
	case 1:
		return set[0], nil
	default:
		return nil, fmt.Errorf("handle %q matches more than one agent; use the agentID instead", ref)
	}
}

func (h *agentHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// The reference param is never Required (§8.1), so an empty ref lists.
	if ref := toolkit.Str(args, "agent"); ref != "" {
		a, err := findAgent(ctx, ref)
		if err != nil {
			return nil, err
		}
		return toolkit.JSONResultWith(a, agentLinks(a))
	}

	f := sysTypes.AgentFilter{
		Query:  toolkit.Str(args, "query"),
		Handle: toolkit.Str(args, "handle"),
		Status: toolkit.Str(args, "status"),
	}
	if toolkit.Bool(args, "includeDeleted") {
		f.Deleted = filter.StateInclusive
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := sysService.DefaultAgent.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("agent search", err)
	}

	items := make([]agentItem, 0, len(set))
	for _, a := range set {
		items = append(items, agentItem{ID: a.ID, Handle: a.Handle, Short: a.Meta.Short, Status: a.Status})
	}

	return toolkit.JSONResult(map[string]any{
		"agents":         items,
		"nextPageCursor": out.NextPage,
	})
}

func (h *agentHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	a := &sysTypes.Agent{
		Handle: toolkit.Str(args, "handle"),
		Status: toolkit.Str(args, "status"),
	}

	if err = applyAgentSections(a, args); err != nil {
		return nil, err
	}

	a, err = sysService.DefaultAgent.Create(ctx, a)
	if err != nil {
		return nil, toolkit.Errf("agent creation", err)
	}

	return toolkit.JSONResultWith(a, agentLinks(a))
}

func (h *agentHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "agent")
	if err != nil {
		return nil, err
	}

	// Loaded first because the service writes the whole record: it has no
	// field-level merge, so anything not carried over here is erased.
	a, err := findAgent(ctx, ref)
	if err != nil {
		return nil, err
	}

	if v, ok := toolkit.OptStr(args, "handle"); ok {
		a.Handle = v
	}
	if v, ok := toolkit.OptStr(args, "status"); ok && v != "" {
		a.Status = v
	}

	if err = applyAgentSections(a, args); err != nil {
		return nil, err
	}

	a, err = sysService.DefaultAgent.Update(ctx, a)
	if err != nil {
		return nil, toolkit.Errf("agent update", err)
	}

	return toolkit.JSONResultWith(a, agentLinks(a))
}

func (h *agentHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "agent")
	if err != nil {
		return nil, err
	}

	a, err := findAgent(ctx, ref)
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultAgent.DeleteByID(ctx, a.ID); err != nil {
		return nil, toolkit.Errf("agent delete", err)
	}

	// The ID is in the acknowledgement because it is the only way back: a
	// deleted agent is excluded from the listing, so its handle stops resolving.
	return toolkit.TextResult("agent %d deleted — restore it with system_agent_undelete agentID=%d", a.ID, a.ID), nil
}

func (h *agentHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	id, err := toolkit.ReqID(args, "agentID")
	if err != nil {
		return nil, err
	}

	if err = sysService.DefaultAgent.UndeleteByID(ctx, id); err != nil {
		return nil, toolkit.Errf("agent undelete", err)
	}

	return toolkit.TextResult("agent %d restored", id), nil
}

// applyAgentSections replaces each configuration section the caller sent and
// leaves the rest alone (§8.2).
//
// Each section decodes into a FRESH zero value which is then assigned, for two
// reasons. Decoding over the existing struct would merge field by field, so
// {"short":"x"} would keep the old description — a merge nothing documents.
// And zeroing the section up front would wipe it whenever the argument turns
// out not to be present after all, which an empty string is.
func applyAgentSections(a *sysTypes.Agent, args map[string]any) error {
	var (
		meta       sysTypes.AgentMeta
		behavior   sysTypes.AgentBehavior
		execution  sysTypes.AgentExecution
		access     sysTypes.AgentAccess
		invocation sysTypes.AgentInvocation
	)

	sections := []struct {
		key     string
		subject string
		into    any
		assign  func()
	}{
		{"meta", "agent meta", &meta, func() { a.Meta = meta }},
		{"behavior", "agent behavior", &behavior, func() { a.Behavior = behavior }},
		{"execution", "agent execution", &execution, func() { a.Execution = execution }},
		{"access", "agent access", &access, func() { a.Access = access }},
		{"invocation", "agent invocation", &invocation, func() { a.Invocation = invocation }},
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

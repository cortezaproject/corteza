package agentic

import (
	"context"
	"fmt"
	"strconv"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/agentic/toolkit"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

type (
	agentService interface {
		Search(ctx context.Context, f sysTypes.AgentFilter) (sysTypes.AgentSet, sysTypes.AgentFilter, error)
		Update(ctx context.Context, upd *sysTypes.Agent) (*sysTypes.Agent, error)
	}

	namespaceHandler struct {
		reg    toolRegistrar
		agents agentService
	}

	// nsItem is the slim list projection. A namespace's meta, labels and
	// translations are incidental to picking one out of a list, so a list drops
	// them; a single-namespace lookup returns the full service type.
	nsItem struct {
		ID   uint64 `json:"namespaceID,string"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
)

// Declarations for these handlers are in namespace_tools.go, in the same order.
func NamespaceHandler(reg toolRegistrar, agents agentService) *namespaceHandler {
	h := &namespaceHandler{reg: reg, agents: agents}
	h.register()
	return h
}

// searchNamespaces returns a page of the slim projection plus the cursor for the
// next page, or nil when there is none.
func (h *namespaceHandler) searchNamespaces(ctx context.Context, page toolkit.Paging) ([]nsItem, *filter.PagingCursor, error) {
	f := cmpTypes.NamespaceFilter{}

	var err error
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := cmpService.DefaultNamespace.Search(ctx, f)
	if err != nil {
		return nil, nil, toolkit.Errf("namespace list", err)
	}

	items := make([]nsItem, len(set))
	for i, ns := range set {
		items[i] = nsItem{ID: ns.ID, Name: ns.Name, Slug: ns.Slug}
	}

	return items, out.NextPage, nil
}

func (h *namespaceHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	page := toolkit.Page(args)

	// The reference is optional: without one this lists, which is how a caller
	// discovers what exists in the first place.
	nsRef := toolkit.Str(args, "namespace")
	if nsRef == "" {
		items, next, err := h.searchNamespaces(ctx, page)
		if err != nil {
			return nil, err
		}

		// The cursor goes in as the value, not as a string: marshalling a
		// PagingCursor emits the encoded form that pageCursor accepts back,
		// while its String method is a debug rendering that does not round-trip.
		return toolkit.JSONResult(map[string]any{
			"namespaces":     items,
			"nextPageCursor": next,
		})
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, nsRef)
	if err != nil {
		// Namespace not found — return the list alongside the error so the caller
		// can pick the correct one without a second round trip.
		items, next, listErr := h.searchNamespaces(ctx, page)
		if listErr != nil {
			return nil, toolkit.Errf("namespace lookup", err)
		}

		return toolkit.JSONResult(map[string]any{
			"error":          fmt.Sprintf("namespace %q not found", nsRef),
			"namespaces":     items,
			"nextPageCursor": next,
		})
	}

	return toolkit.JSONResult(ns)
}

func (h *namespaceHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	name, err := toolkit.ReqStr(args, "name")
	if err != nil {
		return nil, err
	}

	slug, err := toolkit.ReqStr(args, "slug")
	if err != nil {
		return nil, err
	}

	ns, err := cmpService.DefaultNamespace.Create(ctx, &cmpTypes.Namespace{
		Name:           name,
		Slug:           slug,
		Enabled:        parseBoolArg(args["enabled"], true),
		CreatedByAgent: a.GetAgentIDFromContext(ctx),
	})
	if err != nil {
		return nil, toolkit.Errf("namespace creation", err)
	}

	return toolkit.JSONResult(nsResult(ns))
}

func (h *namespaceHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsRef, err := toolkit.ReqStr(args, "namespace")
	if err != nil {
		return nil, err
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, nsRef)
	if err != nil {
		return nil, toolkit.Errf("namespace lookup", err)
	}

	// An absent argument leaves the field unchanged; one sent as an empty string
	// clears it. The assertion fails for a missing key and for a JSON null, both
	// of which count as absent.
	if v, ok := args["name"].(string); ok {
		ns.Name = v
	}
	if v, ok := args["slug"].(string); ok {
		ns.Slug = v
	}
	if v, ok := args["enabled"]; ok {
		ns.Enabled = parseBoolArg(v, ns.Enabled)
	}

	ns, err = cmpService.DefaultNamespace.Update(ctx, ns)
	if err != nil {
		return nil, toolkit.Errf("namespace update", err)
	}

	return toolkit.JSONResult(nsResult(ns))
}

func (h *namespaceHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsRef, err := toolkit.ReqStr(args, "namespace")
	if err != nil {
		return nil, err
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, nsRef)
	if err != nil {
		return nil, toolkit.Errf("namespace lookup", err)
	}

	if err = cmpService.DefaultNamespace.DeleteByID(ctx, ns.ID); err != nil {
		return nil, toolkit.Errf("namespace delete", err)
	}

	if h.agents != nil {
		h.pruneNamespaceFromAgents(ctx, ns.ID)
	}

	return toolkit.TextResult("namespace %d deleted", ns.ID), nil
}

// pruneNamespaceFromAgents removes all allow-list entries referencing the deleted
// namespace from every agent's tool configuration.
func (h *namespaceHandler) pruneNamespaceFromAgents(ctx context.Context, namespaceID uint64) {
	// Use service-user context so RBAC checks inside Search/Update pass regardless
	// of who invoked the agent that triggered the delete.
	svcCtx := a.SetIdentityToContext(ctx, a.ServiceUser())

	agents, _, err := h.agents.Search(svcCtx, sysTypes.AgentFilter{})
	if err != nil {
		return
	}

	for _, ag := range agents {
		changed := false
		for i := range ag.Access.Tools {
			filtered := ag.Access.Tools[i].Allow[:0]
			for _, allow := range ag.Access.Tools[i].Allow {
				if allow.NamespaceID != namespaceID {
					filtered = append(filtered, allow)
				} else {
					changed = true
				}
			}
			ag.Access.Tools[i].Allow = filtered
		}
		if changed {
			_, _ = h.agents.Update(svcCtx, ag)
		}
	}
}

// nsResult is the result shape for the write ops: the list projection plus the
// enabled flag, with the ID as a string to survive a JavaScript client.
func nsResult(ns *cmpTypes.Namespace) map[string]any {
	return map[string]any{
		"namespaceID": strconv.FormatUint(ns.ID, 10),
		"name":        ns.Name,
		"slug":        ns.Slug,
		"enabled":     ns.Enabled,
	}
}

// parseBoolArg returns the bool value of v, or fallback if v is absent or unrecognised.
func parseBoolArg(v any, fallback bool) bool {
	switch b := v.(type) {
	case bool:
		return b
	case string:
		parsed, err := strconv.ParseBool(b)
		if err != nil {
			return fallback
		}
		return parsed
	default:
		return fallback
	}
}

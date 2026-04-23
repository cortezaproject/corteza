package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	a "github.com/crusttech/human/server/pkg/auth"
	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

type agentService interface {
	Search(ctx context.Context, filter sysTypes.AgentFilter) (sysTypes.AgentSet, sysTypes.AgentFilter, error)
	Update(ctx context.Context, upd *sysTypes.Agent) (*sysTypes.Agent, error)
}

type namespaceHandler struct {
	reg     toolRegistrar
	agents  agentService
}

func NamespaceHandler(reg toolRegistrar, agents agentService) *namespaceHandler {
	h := &namespaceHandler{reg: reg, agents: agents}
	h.register()
	return h
}

func (h *namespaceHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_lookup",
			mcp.WithDescription("Look up a specific namespace by name, handle, or slug. Only call this if you do not already know the namespace."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
		),
		"Lookup namespace",
		h.lookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_create",
			mcp.WithDescription("Create a new namespace. A namespace is a top-level container for modules and records in Corteza Compose."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Display name for the namespace")),
			mcp.WithString("slug", mcp.Required(), mcp.Description("URL-friendly identifier (lowercase letters, digits, and hyphens only)")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the namespace is enabled (default: true)")),
		),
		"Create namespace",
		h.create,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_update",
			mcp.WithDescription("Update an existing namespace's name, slug, or enabled state."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("name", mcp.Description("New display name")),
			mcp.WithString("slug", mcp.Description("New URL-friendly identifier")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the namespace should be enabled")),
		),
		"Update namespace",
		h.update,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_delete",
			mcp.WithDescription("Delete a namespace by name, handle, slug, or ID. This permanently removes the namespace and all its contents."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
		),
		"Delete namespace",
		h.del,
	)
}

func (h *namespaceHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]interface{})

	nsRef, _ := args["namespace"].(string)
	if nsRef == "" {
		set, _, err := cmpService.DefaultNamespace.Find(ctx, cmpTypes.NamespaceFilter{})
		if err != nil {
			return nil, fmt.Errorf("namespace list failed: %w", err)
		}
		type nsItem struct {
			ID   uint64 `json:"namespaceID,string"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		}
		items := make([]nsItem, len(set))
		for i, ns := range set {
			items[i] = nsItem{ID: ns.ID, Name: ns.Name, Slug: ns.Slug}
		}
		out, err := json.Marshal(items)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal namespaces: %w", err)
		}
		return mcp.NewToolResultText(string(out)), nil
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, nsRef)
	if err != nil {
		// Namespace not found — return full list so the LLM can pick the correct one
		set, _, listErr := cmpService.DefaultNamespace.Find(ctx, cmpTypes.NamespaceFilter{})
		if listErr != nil {
			return nil, fmt.Errorf("namespace lookup failed: %w", err)
		}
		type nsItem struct {
			ID   uint64 `json:"namespaceID,string"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		}
		items := make([]nsItem, len(set))
		for i, n := range set {
			items[i] = nsItem{ID: n.ID, Name: n.Name, Slug: n.Slug}
		}
		out, _ := json.Marshal(map[string]any{
			"error":      fmt.Sprintf("namespace %q not found", nsRef),
			"namespaces": items,
		})
		return mcp.NewToolResultText(string(out)), nil
	}
	out, err := json.Marshal(ns)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal namespace: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
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

func (h *namespaceHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	name, _ := args["name"].(string)
	slug, _ := args["slug"].(string)
	if name == "" || slug == "" {
		return nil, fmt.Errorf("name and slug are required")
	}

	enabled := parseBoolArg(args["enabled"], true)

	ns, err := cmpService.DefaultNamespace.Create(ctx, &cmpTypes.Namespace{
		Name:    name,
		Slug:    slug,
		Enabled: enabled,
	})
	if err != nil {
		return nil, fmt.Errorf("namespace creation failed: %w", err)
	}

	out, err := json.Marshal(map[string]any{
		"namespaceID": strconv.FormatUint(ns.ID, 10),
		"name":        ns.Name,
		"slug":        ns.Slug,
		"enabled":     ns.Enabled,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal namespace: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *namespaceHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	if v, ok := args["name"].(string); ok && v != "" {
		ns.Name = v
	}
	if v, ok := args["slug"].(string); ok && v != "" {
		ns.Slug = v
	}
	if v, ok := args["enabled"]; ok {
		ns.Enabled = parseBoolArg(v, ns.Enabled)
	}

	ns, err = cmpService.DefaultNamespace.Update(ctx, ns)
	if err != nil {
		return nil, fmt.Errorf("namespace update failed: %w", err)
	}

	out, err := json.Marshal(map[string]any{
		"namespaceID": strconv.FormatUint(ns.ID, 10),
		"name":        ns.Name,
		"slug":        ns.Slug,
		"enabled":     ns.Enabled,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal namespace: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *namespaceHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	if err = cmpService.DefaultNamespace.DeleteByID(ctx, ns.ID); err != nil {
		return nil, fmt.Errorf("namespace delete failed: %w", err)
	}

	if h.agents != nil {
		h.pruneNamespaceFromAgents(ctx, ns.ID)
	}

	return mcp.NewToolResultText(fmt.Sprintf("namespace %d deleted", ns.ID)), nil
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

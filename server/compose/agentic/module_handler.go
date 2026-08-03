package agentic

import (
	"context"
	"encoding/json"
	"fmt"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/agentic/toolkit"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

type moduleHandler struct {
	reg    toolRegistrar
	agents agentService
}

// Declarations for these handlers are in module_tools.go, in the same order.
func ModuleHandler(reg toolRegistrar, agents agentService) *moduleHandler {
	h := &moduleHandler{reg: reg, agents: agents}
	h.register()
	return h
}

// resolveModule resolves the namespace and module refs the single-module tools
// all take. The module ref is required here — the lookup tool, which may omit
// it, resolves the namespace on its own.
func (h *moduleHandler) resolveModule(ctx context.Context, args map[string]any) (*cmpTypes.Module, error) {
	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, toolkit.Errf("namespace lookup", err)
	}

	modRef, err := toolkit.ReqStr(args, "module")
	if err != nil {
		return nil, err
	}

	mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, modRef)
	if err != nil {
		return nil, toolkit.Errf("module lookup", err)
	}

	return mod, nil
}

func (h *moduleHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, toolkit.Errf("namespace lookup", err)
	}

	// Single-module mode returns the raw service type, fields and config included.
	if modRef := toolkit.Str(args, "module"); modRef != "" {
		mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, modRef)
		if err != nil {
			return nil, toolkit.Errf("module lookup", err)
		}
		return toolkit.JSONResult(mod)
	}

	f := cmpTypes.ModuleFilter{NamespaceID: ns.ID}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := cmpService.DefaultModule.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("module list", err)
	}

	// List mode is a slim projection and carries no fields at all: a module's
	// field list is the heavy part, and it is incidental when the caller is
	// picking one module out of many. Fetch a single module for the definition.
	type modItem struct {
		ID          uint64 `json:"moduleID,string"`
		NamespaceID uint64 `json:"namespaceID,string"`
		Name        string `json:"name"`
		Handle      string `json:"handle"`
	}

	items := make([]modItem, len(set))
	for i, mod := range set {
		items[i] = modItem{ID: mod.ID, NamespaceID: mod.NamespaceID, Name: mod.Name, Handle: mod.Handle}
	}

	return toolkit.JSONResult(map[string]any{
		"modules":        items,
		"nextPageCursor": out.NextPage,
	})
}

func (h *moduleHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	name, err := toolkit.ReqStr(args, "name")
	if err != nil {
		return nil, err
	}

	handle, err := toolkit.ReqStr(args, "handle")
	if err != nil {
		return nil, err
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, toolkit.Errf("namespace lookup", err)
	}

	mod := &cmpTypes.Module{
		NamespaceID:    ns.ID,
		Name:           name,
		Handle:         handle,
		CreatedByAgent: a.GetAgentIDFromContext(ctx),
	}

	if rawFields, ok := args["fields"]; ok && rawFields != nil {
		mod.Fields, err = parseModuleFields(rawFields)
		if err != nil {
			return nil, fmt.Errorf("invalid fields: %w", err)
		}
	}

	if rawCfg, ok := args["config"]; ok && rawCfg != nil {
		mod.Config, err = parseModuleConfig(rawCfg)
		if err != nil {
			return nil, fmt.Errorf("invalid config: %w", err)
		}
	}

	mod, err = cmpService.DefaultModule.Create(ctx, mod)
	if err != nil {
		return nil, toolkit.Errf("module creation", err)
	}
	return toolkit.JSONResult(mod)
}

func (h *moduleHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	mod, err := h.resolveModule(ctx, args)
	if err != nil {
		return nil, err
	}

	// Absent leaves the value alone. Name and handle are not clearable — a
	// module without either is invalid — so an empty string is ignored rather
	// than treated as a clear.
	if v := toolkit.Str(args, "name"); v != "" {
		mod.Name = v
	}
	if v := toolkit.Str(args, "handle"); v != "" {
		mod.Handle = v
	}

	// Fields are the documented merge exception: incoming fields are matched by
	// name against the existing set instead of replacing it, and removeFields is
	// the escape hatch that makes deletion expressible.
	if rawFields, ok := args["fields"]; ok && rawFields != nil {
		incoming, err := parseModuleFields(rawFields)
		if err != nil {
			return nil, fmt.Errorf("invalid fields: %w", err)
		}
		mod.Fields = mergeModuleFields(mod.Fields, incoming)
	}

	if rawRemove, ok := args["removeFields"]; ok && rawRemove != nil {
		toRemove, err := parseStringArray(rawRemove)
		if err != nil {
			return nil, fmt.Errorf("invalid removeFields: %w", err)
		}
		mod.Fields = removeModuleFields(mod.Fields, toRemove)
	}

	// Config replaces wholesale: present-and-empty clears it.
	if rawCfg, ok := args["config"]; ok && rawCfg != nil {
		mod.Config, err = parseModuleConfig(rawCfg)
		if err != nil {
			return nil, fmt.Errorf("invalid config: %w", err)
		}
	}

	// Re-assign Place after any merges/removals
	for i, f := range mod.Fields {
		f.Place = i
	}

	mod, err = cmpService.DefaultModule.Update(ctx, mod)
	if err != nil {
		return nil, toolkit.Errf("module update", err)
	}
	return toolkit.JSONResult(mod)
}

func (h *moduleHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	mod, err := h.resolveModule(ctx, args)
	if err != nil {
		return nil, err
	}

	if err = cmpService.DefaultModule.DeleteByID(ctx, mod.NamespaceID, mod.ID); err != nil {
		return nil, toolkit.Errf("module delete", err)
	}

	if h.agents != nil {
		h.pruneModuleFromAgents(ctx, mod.ID)
	}

	return toolkit.TextResult("module %d deleted", mod.ID), nil
}

func (h *moduleHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsRef, err := toolkit.ReqRef(args, "namespace")
	if err != nil {
		return nil, err
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, nsRef)
	if err != nil {
		return nil, toolkit.Errf("namespace lookup", err)
	}

	// The module ref is an ID here rather than the name-or-handle every other
	// module tool takes: a deleted module is excluded from every search path, so
	// FindByAny cannot resolve one and there is nothing to resolve against.
	modID, err := toolkit.ReqID(args, "moduleID")
	if err != nil {
		return nil, err
	}

	if err = cmpService.DefaultModule.UndeleteByID(ctx, ns.ID, modID); err != nil {
		return nil, toolkit.Errf("module undelete", err)
	}

	return toolkit.TextResult("module %d restored", modID), nil
}

func (h *moduleHandler) pruneModuleFromAgents(ctx context.Context, moduleID uint64) {
	svcCtx := a.SetIdentityToContext(ctx, a.ServiceUser())

	agents, _, err := h.agents.Search(svcCtx, sysTypes.AgentFilter{})
	if err != nil {
		return
	}

	for _, ag := range agents {
		changed := false
		for i := range ag.Access.Tools {
			for j, allow := range ag.Access.Tools[i].Allow {
				filtered := ag.Access.Tools[i].Allow[j].ModuleIDs[:0]
				for _, mid := range allow.ModuleIDs {
					if mid != moduleID {
						filtered = append(filtered, mid)
					} else {
						changed = true
					}
				}
				ag.Access.Tools[i].Allow[j].ModuleIDs = filtered
			}
		}
		if changed {
			_, _ = h.agents.Update(svcCtx, ag)
		}
	}
}

// mergeModuleFields merges incoming fields into existing ones by name.
// Existing fields not present in incoming are kept. New fields are appended.
func mergeModuleFields(existing, incoming cmpTypes.ModuleFieldSet) cmpTypes.ModuleFieldSet {
	result := make(cmpTypes.ModuleFieldSet, len(existing))
	copy(result, existing)

	for _, in := range incoming {
		found := false
		for _, ex := range result {
			if ex.Name == in.Name {
				// Preserve IDs and timestamps — update everything else
				in.ID = ex.ID
				in.ModuleID = ex.ModuleID
				in.NamespaceID = ex.NamespaceID
				in.CreatedAt = ex.CreatedAt
				*ex = *in
				found = true
				break
			}
		}
		if !found {
			result = append(result, in)
		}
	}

	return result
}

// removeModuleFields drops fields whose names are in the remove list.
func removeModuleFields(fields cmpTypes.ModuleFieldSet, names []string) cmpTypes.ModuleFieldSet {
	remove := make(map[string]bool, len(names))
	for _, n := range names {
		remove[n] = true
	}
	out := fields[:0]
	for _, f := range fields {
		if !remove[f.Name] {
			out = append(out, f)
		}
	}
	return out
}

// parseStringArray parses a JSON string or native array into []string.
func parseStringArray(raw interface{}) ([]string, error) {
	var data []byte
	switch v := raw.(type) {
	case string:
		data = []byte(v)
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return nil, err
		}
	}
	var out []string
	return out, json.Unmarshal(data, &out)
}

// parseModuleConfig unmarshals a JSON string or object into ModuleConfig.
func parseModuleConfig(raw interface{}) (cmpTypes.ModuleConfig, error) {
	var data []byte
	switch v := raw.(type) {
	case string:
		data = []byte(v)
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return cmpTypes.ModuleConfig{}, fmt.Errorf("cannot encode config: %w", err)
		}
	}
	var cfg cmpTypes.ModuleConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cmpTypes.ModuleConfig{}, fmt.Errorf("config must be a JSON object: %w", err)
	}
	return cfg, nil
}

// parseModuleFields unmarshals a JSON string or array directly into ModuleFieldSet.
// Place is set from the array index since it has json:"-".
func parseModuleFields(raw interface{}) (cmpTypes.ModuleFieldSet, error) {
	var data []byte
	switch v := raw.(type) {
	case string:
		data = []byte(v)
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return nil, fmt.Errorf("cannot encode fields: %w", err)
		}
	}

	var fields cmpTypes.ModuleFieldSet
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, fmt.Errorf("fields must be a JSON array: %w", err)
	}

	for i, f := range fields {
		if f.Name == "" {
			return nil, fmt.Errorf("field[%d]: name is required", i)
		}
		if f.Kind == "" {
			return nil, fmt.Errorf("field[%d] %q: kind is required", i, f.Name)
		}
		f.Place = i
		if f.Options == nil {
			f.Options = cmpTypes.ModuleFieldOptions{}
		}
	}

	return fields, nil
}

package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	cmpService "github.com/cortezaproject/corteza/server/compose/service"
	cmpTypes "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type (
	toolRegistrar interface {
		RegisterTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc)
	}

	recordHandler struct {
		reg toolRegistrar
	}
)

func RecordHandler(reg toolRegistrar) *recordHandler {
	h := &recordHandler{reg: reg}
	h.register()
	return h
}

func (h *recordHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_lookup",
			mcp.WithDescription("Look up a specific namespace by name, handle, or slug. Only call this if you do not already know the namespace."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
		),
		"Lookup namespace",
		h.namespaceLookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_module_lookup",
			mcp.WithDescription("Look up a specific module by name or handle. Only call this if you do not already know the module's field names."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
		),
		"Lookup module",
		h.moduleLookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_lookup",
			mcp.WithDescription("Look up a record by ID, or list/filter records in a module. Provide 'recordID' to fetch one record. Omit 'recordID' and use 'filter' to search by field values (e.g. \"name = 'John'\"). Do NOT call this before creating a record — only use it when the user explicitly asks to search or check for existing records."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("recordID", mcp.Description("Record ID (as string to prevent precision loss). Omit to list/filter instead.")),
			mcp.WithString("filter", mcp.Description("Filter expression when no recordID given, e.g. \"name = 'John'\" or \"status = 'open'\".")),
		),
		"Lookup record",
		h.lookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_create",
			mcp.WithDescription("Create a new record. If you do not know the field names, call compose_module_lookup first to get them."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("values", mcp.Required(), mcp.Description("JSON object of field name-value pairs")),
		),
		"Create record",
		h.create,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_update",
			mcp.WithDescription("Update an existing record. Requires a record ID — use discovery_search to find it if unknown."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("recordID", mcp.Required(), mcp.Description("Record ID (as string to prevent precision loss)")),
			mcp.WithString("values", mcp.Required(), mcp.Description("JSON object of field name-value pairs to update")),
		),
		"Update record",
		h.update,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_delete",
			mcp.WithDescription("Delete a record by ID. Requires a record ID — use discovery_search to find it if unknown."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("recordID", mcp.Required(), mcp.Description("Record ID (as string to prevent precision loss)")),
		),
		"Delete record",
		h.del,
	)
}

func (h *recordHandler) namespaceLookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (h *recordHandler) moduleLookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	modRef, _ := args["module"].(string)
	if modRef == "" {
		set, _, err := cmpService.DefaultModule.Find(ctx, cmpTypes.ModuleFilter{NamespaceID: ns.ID})
		if err != nil {
			return nil, fmt.Errorf("module list failed: %w", err)
		}
		type fieldItem struct {
			Name string `json:"name"`
			Kind string `json:"kind"`
		}
		type modItem struct {
			ID          uint64      `json:"moduleID,string"`
			NamespaceID uint64      `json:"namespaceID,string"`
			Name        string      `json:"name"`
			Handle      string      `json:"handle"`
			Fields      []fieldItem `json:"fields"`
		}
		items := make([]modItem, len(set))
		for i, mod := range set {
			fields := make([]fieldItem, len(mod.Fields))
			for j, f := range mod.Fields {
				fields[j] = fieldItem{Name: f.Name, Kind: f.Kind}
			}
			items[i] = modItem{ID: mod.ID, NamespaceID: mod.NamespaceID, Name: mod.Name, Handle: mod.Handle, Fields: fields}
		}
		out, err := json.Marshal(items)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal modules: %w", err)
		}
		return mcp.NewToolResultText(string(out)), nil
	}

	mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, modRef)
	if err != nil {
		return nil, fmt.Errorf("module lookup failed: %w", err)
	}
	out, err := json.Marshal(mod)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal module: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func parseValues(raw interface{}) (map[string]string, error) {
	switch v := raw.(type) {
	case string:
		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(v), &raw); err != nil {
			return nil, fmt.Errorf("invalid values JSON: %w", err)
		}
		m := make(map[string]string, len(raw))
		for k, val := range raw {
			m[k] = fmt.Sprintf("%v", val)
		}
		return m, nil
	case map[string]interface{}:
		m := make(map[string]string, len(v))
		for key, val := range v {
			m[key] = fmt.Sprintf("%v", val)
		}
		return m, nil
	default:
		return nil, fmt.Errorf("invalid values: expected JSON string or object")
	}
}

func (h *recordHandler) resolveNsMod(ctx context.Context, args map[string]interface{}) (nsID, modID uint64, err error) {
	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return 0, 0, fmt.Errorf("namespace lookup failed: %w", err)
	}

	mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, args["module"])
	if err != nil {
		return 0, 0, fmt.Errorf("module lookup failed: %w", err)
	}

	return ns.ID, mod.ID, nil
}

func (h *recordHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	if recIDStr, _ := args["recordID"].(string); recIDStr != "" {
		recID, err := strconv.ParseUint(recIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid recordID: %w", err)
		}
		rec, _, err := cmpService.DefaultRecord.FindByID(ctx, nsID, modID, recID)
		if err != nil {
			return nil, fmt.Errorf("record lookup failed: %w", err)
		}
		out, err := json.Marshal(rec)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal record: %w", err)
		}
		return mcp.NewToolResultText(string(out)), nil
	}

	filter, _ := args["filter"].(string)
	set, _, err := cmpService.DefaultRecord.Find(ctx, cmpTypes.RecordFilter{
		NamespaceID: nsID,
		ModuleID:    modID,
		Query:       filter,
	})
	if err != nil {
		return nil, fmt.Errorf("record list failed: %w", err)
	}
	out, err := json.Marshal(set)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal records: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *recordHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	valuesMap, err := parseValues(args["values"])
	if err != nil {
		return nil, err
	}

	rec := &cmpTypes.Record{
		NamespaceID: nsID,
		ModuleID:    modID,
	}
	for name, value := range valuesMap {
		rec.Values = append(rec.Values, &cmpTypes.RecordValue{Name: name, Value: value})
	}

	rec, _, err = cmpService.DefaultRecord.Create(ctx, rec)
	if err != nil {
		return nil, fmt.Errorf("record creation failed: %w", err)
	}

	out, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *recordHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	recID, err := strconv.ParseUint(args["recordID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid recordID: %w", err)
	}

	valuesMap, err := parseValues(args["values"])
	if err != nil {
		return nil, err
	}

	rec := &cmpTypes.Record{
		ID:          recID,
		NamespaceID: nsID,
		ModuleID:    modID,
	}
	for name, value := range valuesMap {
		rec.Values = append(rec.Values, &cmpTypes.RecordValue{Name: name, Value: value})
	}

	rec, _, err = cmpService.DefaultRecord.Update(ctx, rec)
	if err != nil {
		return nil, fmt.Errorf("record update failed: %w", err)
	}

	out, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *recordHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	recID, err := strconv.ParseUint(args["recordID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid recordID: %w", err)
	}

	if err = cmpService.DefaultRecord.DeleteByID(ctx, nsID, modID, recID); err != nil {
		return nil, fmt.Errorf("record delete failed: %w", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("record %d deleted", recID)), nil
}

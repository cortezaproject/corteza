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

type toolRegistrar interface {
	RegisterTool(tool mcp.Tool, handler server.ToolHandlerFunc)
}

type recordHandler struct {
	reg toolRegistrar
}

func RecordHandler(reg toolRegistrar) *recordHandler {
	h := &recordHandler{reg: reg}
	h.register()
	return h
}

func (h *recordHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_list",
			mcp.WithDescription("List all available Corteza namespaces"),
		),
		h.namespaceList,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_module_list",
			mcp.WithDescription("List all modules in a Corteza namespace"),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace ID, handle, or slug")),
		),
		h.moduleList,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_lookup",
			mcp.WithDescription("Look up a Corteza namespace by ID, handle, or slug"),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace ID, handle, or slug")),
		),
		h.namespaceLookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_module_lookup",
			mcp.WithDescription("Look up a Corteza module by ID, handle, or name within a namespace"),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace ID, handle, or slug")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module ID, handle, or name")),
		),
		h.moduleLookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_lookup",
			mcp.WithDescription("Look up a record by ID"),
			mcp.WithString("namespaceID", mcp.Required(), mcp.Description("Namespace ID")),
			mcp.WithString("moduleID", mcp.Required(), mcp.Description("Module ID")),
			mcp.WithString("recordID", mcp.Required(), mcp.Description("Record ID")),
		),
		h.lookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_create",
			mcp.WithDescription("Create a new record"),
			mcp.WithString("namespaceID", mcp.Required(), mcp.Description("Namespace ID")),
			mcp.WithString("moduleID", mcp.Required(), mcp.Description("Module ID")),
			mcp.WithString("values", mcp.Required(), mcp.Description("JSON object of field name-value pairs")),
		),
		h.create,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_update",
			mcp.WithDescription("Update an existing record"),
			mcp.WithString("namespaceID", mcp.Required(), mcp.Description("Namespace ID")),
			mcp.WithString("moduleID", mcp.Required(), mcp.Description("Module ID")),
			mcp.WithString("recordID", mcp.Required(), mcp.Description("Record ID")),
			mcp.WithString("values", mcp.Required(), mcp.Description("JSON object of field name-value pairs")),
		),
		h.update,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_delete",
			mcp.WithDescription("Delete a record by ID"),
			mcp.WithString("namespaceID", mcp.Required(), mcp.Description("Namespace ID")),
			mcp.WithString("moduleID", mcp.Required(), mcp.Description("Module ID")),
			mcp.WithString("recordID", mcp.Required(), mcp.Description("Record ID")),
		),
		h.del,
	)
}

func (h *recordHandler) namespaceList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (h *recordHandler) moduleList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

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

func (h *recordHandler) namespaceLookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
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

	mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, args["module"])
	if err != nil {
		return nil, fmt.Errorf("module lookup failed: %w", err)
	}

	out, err := json.Marshal(mod)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal module: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *recordHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	nsID, err := strconv.ParseUint(args["namespaceID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid namespaceID: %w", err)
	}

	modID, err := strconv.ParseUint(args["moduleID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid moduleID: %w", err)
	}

	recID, err := strconv.ParseUint(args["recordID"].(string), 10, 64)
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

func (h *recordHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	nsID, err := strconv.ParseUint(args["namespaceID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid namespaceID: %w", err)
	}

	modID, err := strconv.ParseUint(args["moduleID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid moduleID: %w", err)
	}

	var valuesMap map[string]string
	if err := json.Unmarshal([]byte(args["values"].(string)), &valuesMap); err != nil {
		return nil, fmt.Errorf("invalid values JSON: %w", err)
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

	nsID, err := strconv.ParseUint(args["namespaceID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid namespaceID: %w", err)
	}

	modID, err := strconv.ParseUint(args["moduleID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid moduleID: %w", err)
	}

	recID, err := strconv.ParseUint(args["recordID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid recordID: %w", err)
	}

	var valuesMap map[string]string
	if err := json.Unmarshal([]byte(args["values"].(string)), &valuesMap); err != nil {
		return nil, fmt.Errorf("invalid values JSON: %w", err)
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

	nsID, err := strconv.ParseUint(args["namespaceID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid namespaceID: %w", err)
	}

	modID, err := strconv.ParseUint(args["moduleID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid moduleID: %w", err)
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

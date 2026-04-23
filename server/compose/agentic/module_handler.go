package agentic

import (
	"context"
	"encoding/json"
	"fmt"

	a "github.com/crusttech/human/server/pkg/auth"
	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

type moduleHandler struct {
	reg    toolRegistrar
	agents agentService
}

func ModuleHandler(reg toolRegistrar, agents agentService) *moduleHandler {
	h := &moduleHandler{reg: reg, agents: agents}
	h.register()
	return h
}

func (h *moduleHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_module_lookup",
			mcp.WithDescription("Look up modules and their fields in a namespace. To list ALL modules in a namespace, omit the module argument. To get details on a specific module, provide its name or handle. Never guess module names — always list first if you are unsure what exists."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Description("Module name, handle, or ID. Omit to list all modules in the namespace.")),
		),
		"Lookup module",
		h.lookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_module_delete",
			mcp.WithDescription("Delete a module by name, handle, or ID. This permanently removes the module and all its records."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID")),
		),
		"Delete module",
		h.del,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_module_update",
			mcp.WithDescription("Update an existing module's name, handle, fields, or configuration. Fields are merged: provided fields are added or updated by name, existing fields not mentioned are kept. To remove specific fields use removeFields."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID")),
			mcp.WithString("name", mcp.Description("New display name for the module")),
			mcp.WithString("handle", mcp.Description("New handle for the module")),
			mcp.WithString("fields", mcp.Description("JSON array of fields to add or update. Same format as compose_module_create. Existing fields not listed are preserved.")),
			mcp.WithString("removeFields", mcp.Description("JSON array of field names to remove, e.g. [\"fieldA\",\"fieldB\"]")),
			mcp.WithString("config", mcp.Description(`JSON object for module-level configuration. Supports: recordDeDup (duplicate detection), recordRevisions (audit trail), privacy (data sensitivity). Example: {"recordDeDup":{"rules":[{"name":"unique-email","strict":true,"constraints":[{"attribute":"email","modifier":"ignore-case"}]}]},"recordRevisions":{"enabled":true},"privacy":{"usageDisclosure":"Used for customer contact only"}}`)),
		),
		"Update module",
		h.update,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_module_create",
			mcp.WithDescription("Create a new module in a namespace. A module defines a data structure (like a table) with typed fields. As a developer acting on behalf of the user, proactively add config where appropriate: enable duplicate detection for modules storing contacts/leads/customers (match on email or phone), enable recordRevisions for important transactional data, and set privacy disclosure for modules holding personal information."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Display name for the module")),
			mcp.WithString("handle", mcp.Required(), mcp.Description("URL-friendly identifier (lowercase letters, digits, and underscores only)")),
			mcp.WithString("fields", mcp.Description(`JSON array of field definitions. Each field: {"name":"fieldName","kind":"String","label":"Field Label","isRequired":false,"isMulti":false,"options":{},"defaultValue":[{"name":"fieldName","value":"default"}],"expressions":{"value":"","sanitizers":[],"validators":[],"formatters":[]}}. Supported kinds: String, Number, DateTime, Bool, Record, User, File, Select, Email, Url, Currency, Duration. Kind-specific options: Select→{"options":[{"value":"a","text":"A"}]}, Record→{"moduleID":"123","recordLabelField":"name"}, Number→{"precision":2}, DateTime→{"onlyDate":false,"onlyTime":false}, Bool→{"trueLabel":"Yes","falseLabel":"No"}.`)),
			mcp.WithString("config", mcp.Description(`JSON object for module-level configuration. recordDeDup rules: {"name":"rule-name","strict":true,"constraints":[{"attribute":"fieldName","modifier":"ignore-case|case-sensitive|fuzzy-match|sounds-like","multiValue":"one-of|equal"}]}. recordRevisions: {"enabled":true}. privacy: {"usageDisclosure":"text","sensitivityLevelID":"123"}.`)),
		),
		"Create module",
		h.create,
	)
}

func (h *moduleHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func (h *moduleHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	if err = cmpService.DefaultModule.DeleteByID(ctx, ns.ID, mod.ID); err != nil {
		return nil, fmt.Errorf("module delete failed: %w", err)
	}

	if h.agents != nil {
		h.pruneModuleFromAgents(ctx, mod.ID)
	}

	return mcp.NewToolResultText(fmt.Sprintf("module %d deleted", mod.ID)), nil
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

func (h *moduleHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	if v, ok := args["name"].(string); ok && v != "" {
		mod.Name = v
	}
	if v, ok := args["handle"].(string); ok && v != "" {
		mod.Handle = v
	}

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
		return nil, fmt.Errorf("module update failed: %w", err)
	}

	out, err := json.Marshal(mod)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal module: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *moduleHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	name, _ := args["name"].(string)
	handle, _ := args["handle"].(string)
	if name == "" || handle == "" {
		return nil, fmt.Errorf("name and handle are required")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	mod := &cmpTypes.Module{
		NamespaceID: ns.ID,
		Name:        name,
		Handle:      handle,
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
		return nil, fmt.Errorf("module creation failed: %w", err)
	}

	out, err := json.Marshal(mod)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal module: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
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

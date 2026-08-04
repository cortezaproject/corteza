package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations only. Implementations are in module_handler.go, in this order.

func (h *moduleHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_module_lookup",
			mcp.WithDescription(
				"Look up modules in a namespace, or get one module's full definition. A module is a data "+
					"structure (like a table) with typed fields. "+
					"To list ALL modules in a namespace, omit the 'module' argument; the list is a slim view — "+
					"module identity plus each field's name and kind. To get everything else about one module "+
					"— field labels, options, default values, expressions and module config — provide its name, "+
					"handle or ID. "+
					"Never guess module or field names: list first if you are unsure what exists, then fetch the "+
					"one you need. "+
					"Use 'limit' and 'pageCursor' when a namespace holds more modules than one response can carry.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Description("Module name, handle, or ID (as string to prevent precision loss). Omit to list all modules in the namespace.")),
			mcp.WithString("limit", mcp.Description("Maximum modules to return when listing, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup module",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_module_create",
			mcp.WithDescription("Create a new module in a namespace. A module defines a data structure (like a table) with typed fields. As a developer acting on behalf of the user, proactively add config where appropriate: enable duplicate detection for modules storing contacts/leads/customers (match on email or phone), enable recordRevisions for important transactional data, and set privacy disclosure for modules holding personal information."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Display name for the module")),
			mcp.WithString("handle", mcp.Required(), mcp.Description("URL-friendly identifier (lowercase letters, digits, and underscores only)")),
			mcp.WithString("fields", mcp.Description(`JSON array of field definitions. Each field: {"name":"fieldName","kind":"String","label":"Field Label","isRequired":false,"isMulti":false,"options":{},"defaultValue":[{"name":"fieldName","value":"default"}],"expressions":{"value":"","sanitizers":[],"validators":[],"formatters":[]}}. Supported kinds: String, Number, DateTime, Bool, Record, User, File, Select, Email, Url, Currency, Duration. Kind-specific options: Select→{"options":[{"value":"a","text":"A"}]}, Record→{"moduleID":"123","recordLabelField":"name"}, Number→{"precision":2}, DateTime→{"onlyDate":false,"onlyTime":false}, Bool→{"trueLabel":"Yes","falseLabel":"No"}.`)),
			mcp.WithString("config", mcp.Description(`JSON object for module-level configuration. recordDeDup rules: {"name":"rule-name","strict":true,"constraints":[{"attribute":"fieldName","modifier":"ignore-case|case-sensitive|fuzzy-match|sounds-like","multiValue":"one-of|equal"}]}. recordRevisions: {"enabled":true}. privacy: {"usageDisclosure":"text","sensitivityLevelID":"123"}.`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create module",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_module_update",
			mcp.WithDescription(
				"Update an existing module's name, handle, fields, or configuration. Arguments you omit are "+
					"left unchanged. "+
					"Fields are the exception to the usual replace-wholesale rule: they are MERGED by field "+
					"name — fields you pass are added or updated, existing fields you do not mention are kept, "+
					"and passing an empty array removes nothing. To drop fields, name them in 'removeFields'. "+
					"'config', by contrast, replaces the whole module config when provided.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("name", mcp.Description("New display name for the module")),
			mcp.WithString("handle", mcp.Description("New handle for the module")),
			mcp.WithString("fields", mcp.Description("JSON array of fields to add or update. Same format as compose_module_create. Existing fields not listed are preserved.")),
			mcp.WithString("removeFields", mcp.Description("JSON array of field names to remove, e.g. [\"fieldA\",\"fieldB\"]")),
			mcp.WithString("config", mcp.Description(`JSON object for module-level configuration. Replaces the existing config. Supports: recordDeDup (duplicate detection), recordRevisions (audit trail), privacy (data sensitivity). Example: {"recordDeDup":{"rules":[{"name":"unique-email","strict":true,"constraints":[{"attribute":"email","modifier":"ignore-case"}]}]},"recordRevisions":{"enabled":true},"privacy":{"usageDisclosure":"Used for customer contact only"}}`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update module",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_module_delete",
			mcp.WithDescription(
				"Delete a module by name, handle, or ID. The delete is soft: nothing is erased — the module "+
					"stops appearing in lookups and its records become unreachable through it, but both are "+
					"retained. "+
					"Pages, charts and record fields that point at this module will no longer resolve it, so "+
					"check what references it before deleting. compose_module_undelete reverses this, but note "+
					"the module ID first — a deleted module can no longer be found by name or handle.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete module",
		h.del,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_module_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted module, reversing compose_module_delete. A delete only sets a marker: "+
					"the module definition, its fields and every record in it were retained, so the module "+
					"comes back exactly as it was and the pages, charts and Record fields that reference it "+
					"resolve again. Records that were deleted individually stay deleted — use "+
					"compose_record_undelete for those. "+
					"Requires the numeric moduleID, and nothing else will do: a deleted module is excluded from "+
					"every lookup path, so compose_module_lookup can no longer resolve it by name or handle, and "+
					"that tool exposes no includeDeleted-style filter. Take the ID from what "+
					"compose_module_delete reported, or from a compose_module_lookup made before the delete. "+
					"Calling this on a module that is not deleted is accepted and changes nothing.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("moduleID", mcp.Required(), mcp.Description("ID of the deleted module (as string to prevent precision loss). A name or handle will not work — deleted modules are not resolvable by either.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete module",
		h.undelete,
	)
}

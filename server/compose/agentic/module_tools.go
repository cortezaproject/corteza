package agentic

import (
	"fmt"
	"strings"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// The kind list is read from types.ModuleFieldKinds rather than written out
// again here. Spelling it twice is what let this documentation advertise a
// Currency and a Duration the platform never had, and omit Geometry, which it
// did — with nothing to catch the drift, since an unknown kind is stored as
// written and rendered as text.
var fieldsParamDoc = fmt.Sprintf(
	`JSON array of field definitions. Each field: {"name":"fieldName","kind":"String",`+
		`"label":"Field Label","isRequired":false,"isMulti":false,"options":{},`+
		`"defaultValue":[{"name":"fieldName","value":"default"}],`+
		`"expressions":{"value":"","sanitizers":[],"validators":[],"formatters":[]}}. `+
		`Supported kinds: %s — these are the whole set, and a kind outside it is not a field `+
		`type the webapp can render. `+
		`Kind-specific options: Select→{"options":[{"value":"a","text":"A"}]}, `+
		`Record→{"moduleID":"123","labelField":"name","queryFields":["name"],"selectType":"default"} `+
		`(labelField is what the picker and every viewer show; without it they fall back to the `+
		`module's first field. recordLabelField is only the second-level label for when labelField `+
		`itself points at a Record field), `+
		`Number→{"precision":2,"format":"0,0.00","prefix":"","suffix":""} `+
		`(precision rounds what is STORED; display comes from format, so precision alone drops the `+
		`decimals it kept), `+
		`DateTime→{"onlyDate":false,"onlyTime":false}, Bool→{"trueLabel":"Yes","falseLabel":"No"}, `+
		`Geometry→{"center":[46.05,14.51],"zoom":7}. `+
		`"defaultValue" fills the field when a record is created without it; the "name" key is `+
		`optional and is stored empty, because the field already says which field it is for. `+
		expressionsDoc,
	strings.Join(cmpTypes.ModuleFieldKinds, ", "),
)

// writeDetailDoc is the same knob compose_module_lookup carries, on the two
// tools that write. A module echoed whole is a few thousand tokens of storage
// configuration per call — every field's DAL encoding strategy, privacy block
// and revision flag — and a caller that has just written the module wants
// confirmation of what it asked for, not the storage it did not.
const writeDetailDoc = `How much of the stored module to echo back. "summary" (the default) ` +
	`confirms what was written — handle, name, and each field's name, kind, label, ` +
	`required/multi, select options and value expression. "full" adds every field's ID, ` +
	`timestamps and DAL storage config, which is thousands of tokens and is only useful when ` +
	`you need the field IDs. Issues with a field expression are reported either way.`

// The four expression slots do not share a scope and two of them are not
// predicates at all, so writing one from the shape of another produces a module
// that stores clean and misbehaves at record save. Everything here is stated
// because a caller with only these tools has nowhere else to read it.
const expressionsDoc = `EXPRESSIONS. Each slot gets a different scope. ` +
	`String literals need DOUBLE quotes — a single-quoted string is a syntax error. ` +
	`"value" (the field value expression) sees the record's own fields by their BARE names ` +
	`(stage != "hired" && stage != "rejected"), plus "new" and "old" as whole records ` +
	`(new.values.stage, new.recordID, old.values.stage; on a create every "old" field is null). ` +
	`There is no "record" and no bare "values" — record.values.stage fails every save with ` +
	`'unknown parameter record.values'. A value expression OVERWRITES whatever the caller sent ` +
	`for that field, so do not send it. "isRequired" on such a field means the expression must ` +
	`produce a value, not that a caller must supply one. ` +
	`"validators":[{"test":"…","error":"…"}] — test sees "value" (the value being saved, as a ` +
	`string), "oldValue" and "values.<field>". READ THIS ONE TWICE: test names the condition ` +
	`under which the value is REJECTED. test "value >= 0 && value <= 5" REFUSES 3 and stores 7 — ` +
	`the exact opposite of what it reads like, while showing an error message asserting the ` +
	`range. Write the rule you want as its rejection: "value < 0 || value > 5". ` +
	`"sanitizers" and "formatters" are transforms, not tests: each sees only "value" and its ` +
	`RESULT REPLACES the value — trim(value), toUpper(value). A sanitizer runs before ` +
	`validation and is stored; a formatter runs on the way out. ` +
	`The write result reports what cannot work under "issues"; no "issues" key is the clean ` +
	`result.`

// The words a caller brings for what a module tool configures. Field
// expressions in particular had no route in: their whole documentation lives on
// the "fields" parameter, and search reads names, keywords and the tool
// description only — so "what is in scope for a value expression" matched every
// tool that mentions expressions in passing and none that answers it.
var moduleKeywords = hmcp.WithKeywords(
	"table", "schema", "field", "column",
	"expression", "formula", "calculated", "validator", "validation", "sanitizer", "formatter",
	"default value",
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
			// Reading a definition is what a data agent does before it can read the
			// data: a "usage" grant that cannot see a module has no way to name one
			// or to interpret the values it gets back. Writing one stays configuring.
			mcp.WithString("detail", mcp.Description(
				"How much of each field to return. \"summary\" (the default) gives name, kind, label, required/multi, "+
					"select options and any value expression — what you need to read or filter data. \"full\" adds "+
					"field IDs, timestamps and storage config, which you only need when editing the module itself. "+
					"Summary is a fraction of the size; a namespace of modules at full detail is tens of kilobytes.")),
			moduleKeywords,
			hmcp.InGroup(hmcp.GroupConfiguring, hmcp.GroupUsage),
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
			mcp.WithString("fields", mcp.Description(fieldsParamDoc)),
			mcp.WithString("config", mcp.Description(`JSON object for module-level configuration. Supports: recordDeDup (duplicate detection), recordRevisions (audit trail), privacy (data sensitivity). Rules live under a "rules" array — a bare rule object is accepted and silently stored as {}. Example: {"recordDeDup":{"rules":[{"name":"unique-email","strict":true,"constraints":[{"attribute":"email","modifier":"ignore-case|case-sensitive|fuzzy-match|sounds-like","multiValue":"one-of|equal"}]}]},"recordRevisions":{"enabled":true},"privacy":{"usageDisclosure":"text","sensitivityLevelID":"123"}}`)),
			mcp.WithString("detail", mcp.Description(writeDetailDoc)),
			moduleKeywords,
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
			hmcp.NeedsFullDocs(),
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
			mcp.WithString("handle", mcp.Description("New handle for the module. Use snake_case (lowercase letters, digits, underscores); a hyphen is the subtraction operator wherever an identifier is parsed.")),
			mcp.WithString("fields", mcp.Description("JSON array of fields to add or update. Same format as compose_module_create. Existing fields not listed are preserved. "+expressionsDoc)),
			mcp.WithString("removeFields", mcp.Description("JSON array of field names to remove, e.g. [\"fieldA\",\"fieldB\"]")),
			mcp.WithString("config", mcp.Description(`JSON object for module-level configuration. Replaces the existing config. Supports: recordDeDup (duplicate detection), recordRevisions (audit trail), privacy (data sensitivity). Example: {"recordDeDup":{"rules":[{"name":"unique-email","strict":true,"constraints":[{"attribute":"email","modifier":"ignore-case"}]}]},"recordRevisions":{"enabled":true},"privacy":{"usageDisclosure":"Used for customer contact only"}}`)),
			mcp.WithString("detail", mcp.Description(writeDetailDoc)),
			moduleKeywords,
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

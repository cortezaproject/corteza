package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations only. Implementations are in record_handler.go, in this order.

// The generic words for the thing a module holds.
var recordKeywords = hmcp.WithKeywords("data", "row", "entry", "item")

func (h *recordHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_lookup",
			mcp.WithDescription(
				"Look up a record by ID, or list and filter records in a module. Provide 'recordID' to fetch "+
					"one; omit it and use 'filter' to search by field values. "+
					"Do NOT call this before creating a record — only use it when the user explicitly asks to "+
					"search or check for existing records. "+
					"Records are returned in full, including their values, so use 'limit' and narrow with "+
					"'filter' rather than listing a whole module: a large module will exceed the result size "+
					"limit and the call will fail. "+
					"A Record or User value is stored as the bare ID of what it points at. The response carries "+
					"'refs', a dictionary of every such ID to the name a person would see — read the name from "+
					"there rather than looking each reference up one by one.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("recordID", mcp.Description("Record ID (as string to prevent precision loss). Omit to list or filter instead.")),
			mcp.WithString("recordIDs", mcp.Description("Comma-separated record IDs to fetch in one call, e.g. \"101,102,103\" — use this instead of one call per ID, and instead of an IN expression in 'filter', which the query language does not support. Takes precedence over 'filter'.")),
			mcp.WithString("filter", mcp.Description("Filter expression when no recordID is given, e.g. \"name = 'John'\" or \"status = 'open'\".")),
			mcp.WithString("limit", mcp.Description("Maximum records to return, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupUsage),
			recordKeywords,
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup record",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_record_report",
			mcp.WithDescription(
				"Aggregate records server-side: sums, averages, extremes and counts, optionally grouped by a field. "+
					"Use this for ANY question that is answered by a number over more than a handful of records — "+
					"a total, an average, a count, a breakdown. Reading the records with compose_record_lookup and "+
					"adding them up yourself is a wrong answer waiting to happen, and it silently drops everything "+
					"past the page you were given. "+
					"Returns {\"rows\": [...], \"units\": {...}}: one row per group, each with the metrics you asked "+
					"for plus 'count', the number of records in that group. With no 'dimension' you get a single "+
					"row for the whole set, where dimension_0 is \"*\". "+
					"'units' gives the prefix or suffix the module puts on each aggregated field — use it when "+
					"stating the figure, and do not supply a currency or unit it does not name.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("metrics", mcp.Description(
				"Comma-separated aggregate expressions over numeric fields, e.g. \"SUM(line_value) AS total, AVG(price) AS avg_price\". "+
					"Functions: SUM, AVG, MIN, MAX, COUNT. The alias after AS is the key in the result; without one the "+
					"expression itself is the key. 'count' is always returned and needs no metric.")),
			mcp.WithString("dimension", mcp.Description(
				"A single field name ON THIS MODULE to group by, e.g. \"rarity\" — not a dotted path through a "+
					"Record reference: \"card.rarity\" is not a field and the call fails. To group by something held "+
					"on a referenced record, aggregate that module instead, or denormalise the field. "+
					"Omit for one row covering every record. "+
					"Date fields can be bucketed with DATE(field), and a chart's modifiers (QUARTER, YEAR) are not "+
					"available here — group by the raw field and combine the rows yourself if you need coarser buckets.")),
			mcp.WithString("filter", mcp.Description("Filter expression narrowing which records are aggregated, e.g. \"rarity = 'Mythic'\". Same syntax as compose_record_lookup's filter.")),
			hmcp.InGroup(hmcp.GroupUsage),
			recordKeywords,
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Aggregate records",
		h.report,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_record_create",
			mcp.WithDescription("Create a new record. If you do not know the field names, call compose_module_lookup first to get them."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("values", mcp.Required(), mcp.Description(recordValuesDoc)),
			hmcp.InGroup(hmcp.GroupUsage),
			recordKeywords,
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create record",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_record_update",
			mcp.WithDescription(
				"Update an existing record. Requires a record ID — use compose_record_lookup with a filter to "+
					"find it if unknown. 'values' REPLACES the record's whole value set: any field you leave out "+
					"is cleared, not kept. Read the record first and send every field you intend to survive.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("recordID", mcp.Required(), mcp.Description("Record ID (as string to prevent precision loss)")),
			mcp.WithString("values", mcp.Required(), mcp.Description(recordValuesDoc+" On update these REPLACE the record's values: send the whole set a field should end up with, not just what changed. Adding one value to a multi-value field means reading the record with compose_record_lookup and writing the full list back.")),
			hmcp.InGroup(hmcp.GroupUsage),
			recordKeywords,
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update record",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_record_delete",
			mcp.WithDescription(
				"Delete a record by ID. Requires a record ID — use compose_record_lookup with a filter to find "+
					"it if unknown. The delete is soft: the record stops appearing in lookups but is retained, "+
					"and compose_record_undelete brings it back.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("recordID", mcp.Required(), mcp.Description("Record ID (as string to prevent precision loss)")),
			hmcp.InGroup(hmcp.GroupUsage),
			recordKeywords,
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete record",
		h.del,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_record_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted record, reversing compose_record_delete. Deleting a record in Human "+
					"only marks it deleted, so nothing was lost and the record comes back with its values "+
					"intact. "+
					"Requires the record ID. A deleted record is left out of every filtered listing and "+
					"compose_record_lookup offers no includeDeleted-style filter — but fetching it directly by "+
					"'recordID' with compose_record_lookup does still return it, with 'deletedAt' set, which is "+
					"how you confirm you have the right record before restoring it. Failing that, use the ID "+
					"compose_record_delete reported. "+
					"This restores one record per call. It writes to the record — the revision counter moves and "+
					"undelete automation runs — so do not call it speculatively on a record that is not deleted.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss). The module itself must not be deleted — restore it first with compose_module_undelete if it is.")),
			mcp.WithString("recordID", mcp.Required(), mcp.Description("Record ID of the deleted record (as string to prevent precision loss)")),
			hmcp.InGroup(hmcp.GroupUsage),
			recordKeywords,
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete record",
		h.undelete,
	)
}

const recordValuesDoc = `JSON object of field name to value: {"title":"Kickoff","attendees":42,"done":true}. ` +
	`A MULTI-VALUE field takes an array, and each element becomes one of the record's values in the ` +
	`order given: {"tags":["red","blue"]}. A field whose value is itself structured takes an object, ` +
	`stored as its JSON — a Geometry point is {"geo":{"coordinates":[46.05,14.51]}} (latitude first). ` +
	`Anything else is refused rather than guessed at. ` +
	`A field you leave out is stored as no value at all, which is a different state from a false or ` +
	`an empty one and does not match a query for it: omitting a Bool rather than sending false means ` +
	`"done = false" finds none of those records, and the same goes for a prefilter or a Metric block ` +
	`filter built on that field. Send every field a filter or chart will group on, including the ` +
	`false ones.`

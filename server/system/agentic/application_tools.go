package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Applications are the entries of this instance's app selector. They are a flat
// global list — there is no parent or namespace dimension — which is why
// system_application_reorder legitimately re-weights every application and why
// the lookup filter needs no scoping argument.
//
// This file holds declarations only. Implementations are in
// application_handler.go, in the same order.

func (h *applicationHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_application_lookup",
			mcp.WithDescription(
				"Look up one application by name or ID, or list and filter applications. An application is an "+
					"entry in this instance's app selector: an administrative name, an on/off switch, an "+
					"ordering weight, and an 'unify' block holding the label, URL, icon, logo and configuration "+
					"the selector renders. "+
					"Provide 'application' to fetch one, which returns the full record including the unify "+
					"block; omit it to list, which returns a compact form (applicationID, name, enabled, "+
					"weight, flags, deletedAt) with the unify block left out because it is bulky and useless "+
					"for picking one out of a list — fetch by name or ID when you need it. "+
					"'application' matches an exact name or an ID, never a partial name; use 'query' to search "+
					"by fragment. "+
					"Deleted applications are hidden unless you set 'includeDeleted', which is how you find the "+
					"ID that system_application_undelete requires. "+
					"Weights are listed in ascending order of appearance; call this before "+
					"system_application_reorder to get the current order and the complete set of IDs.",
			),
			mcp.WithString("application", mcp.Description("Exact application name, or an application ID as a string to prevent precision loss. Omit to list instead of fetching one.")),
			mcp.WithString("query", mcp.Description("List mode only. Case-insensitive substring match against the application name.")),
			mcp.WithString("flags", mcp.Description(`List mode only. Return only applications carrying any of these flags. JSON array of strings, e.g. ["pinned"]. Matches your own flags plus global ones, exactly as the app selector sees them.`)),
			mcp.WithBoolean("includeDeleted", mcp.Description("Include soft-deleted applications alongside live ones. Deletes are reversible with system_application_undelete.")),
			mcp.WithString("limit", mcp.Description("Maximum results, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup application",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_application_create",
			mcp.WithDescription(
				"Create an application, adding an entry to this instance's app selector. 'name' is the "+
					"administrative name used in listings; the 'unify' block is what users actually see, and an "+
					"application with no unify URL has nothing to open. "+
					"A new application is disabled and unlisted unless you say otherwise: pass enabled true and "+
					"unify.listed true for it to reach users. It is also placed last in the ordering — this "+
					"tool sets no weight, use system_application_reorder to position it. "+
					"With unify.kind \"custom\" the application is its own HTML document rather than a link to "+
					"a section: leave unify.url out and it is filled in with app/<applicationID> once the ID "+
					"exists, and system_application_source_set is what gives it a body. "+
					"Call system_application_lookup on an existing application first if you are unsure what a "+
					"working unify block looks like on this instance.",
			),
			mcp.WithString("name", mcp.Required(), mcp.Description("Administrative name of the application, shown in admin listings.")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the application is active. Defaults to false, which hides it from every user regardless of the unify block.")),
			mcp.WithString("unify", mcp.Description(
				`App selector configuration, as a JSON object. Keys: "name" (label in the selector, falls back `+
					`to the application name), "listed" (boolean, whether users see it), "url" (where it opens), `+
					`"config" (free-form configuration string), "icon" and "logo" (URLs), "iconID" and "logoID" `+
					`(attachment IDs as strings), "kind" ("custom" for an application whose UI is one HTML `+
					`document, shown by the app view at /app/<applicationID>; with no "url" of your own the `+
					`server sets it to app/<applicationID>). Example: `+
					`{"name":"Reports","listed":true,"url":"/compose/ns/reports"}`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create application",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_application_update",
			mcp.WithDescription(
				"Update an application. Omit a field to leave it unchanged. "+
					"'unify' is merged key by key: the keys you send are applied, the ones you leave out keep "+
					"their current value, and sending a key with an empty value clears it — so you can flip "+
					"listed without resending the URL. "+
					"This tool never changes the ordering weight; use system_application_reorder for that. "+
					"Setting unify.kind to \"custom\" on an application with no unify.url points the URL at "+
					"app/<applicationID>, which is where the app view serves the stored HTML. "+
					"To take an application away from users without deleting it, set enabled false — that is "+
					"reversible and keeps every setting.",
			),
			mcp.WithString("application", mcp.Required(), mcp.Description("Exact application name, or an application ID as a string to prevent precision loss.")),
			mcp.WithString("name", mcp.Description("New administrative name. Omit to leave unchanged; it cannot be set to an empty string.")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the application is active. Omit to leave unchanged.")),
			mcp.WithString("unify", mcp.Description(
				`App selector configuration to merge in, as a JSON object. Same keys as `+
					`system_application_create, "kind" included. Only the keys present are touched, so `+
					`{"listed":false} hides the application from the selector and leaves the URL, icon and `+
					`logo alone.`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update application",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_application_source_get",
			mcp.WithDescription(
				"Read a custom application's HTML back, exactly as stored, with the meta the app view "+
					"reads instead of the document: hash, byte size, and the namespace and modules the app "+
					"is allowed to query. "+
					"Neither system_application_lookup nor the listing carries the source — this is the only "+
					"way to see it. Call it before patching with system_application_source_set, so "+
					"'old_string' is copied from what is actually stored rather than from what you believe "+
					"you wrote. "+
					"Needs 'access' on the application, the same grant that lets a user open it.",
			),
			mcp.WithString("application", mcp.Required(), mcp.Description("Exact application name, or an application ID as a string to prevent precision loss.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Read custom application source",
		h.sourceGet,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_application_source_set",
			mcp.WithDescription(
				"Store the HTML of a custom application — the whole document the app view renders inside "+
					"its sandbox. The application must already carry unify.kind \"custom\"; any other "+
					"application is refused. "+
					"Send it one of two ways: 'source' replaces the document outright, or 'old_string' with "+
					"'new_string' patches the stored one, which is how you edit an app without resending it. "+
					"A patch is applied only when 'old_string' matches exactly once — zero or several "+
					"matches are refused with the count, so include enough surrounding text to be "+
					"unambiguous. "+
					"'namespace' and 'modules' declare the data the app may read: the bridge refuses every "+
					"module not named here. A patch that sends neither keeps the declaration already "+
					"stored. "+
					"The document is plain HTML with inline script. JSX, ES modules and a React import are "+
					"refused, and so are fetch, XMLHttpRequest and WebSocket — the sandbox's CSP blocks all "+
					"of them, so an app using them fails in front of a user instead of here. Load the "+
					"custom_app skill with system_skill_lookup before writing one.",
			),
			mcp.WithString("application", mcp.Required(), mcp.Description("Exact application name, or an application ID as a string to prevent precision loss.")),
			mcp.WithString("source", mcp.Description("The whole HTML document, replacing whatever is stored. Capped at 256 KB; aim well under that. Cannot be combined with old_string.")),
			mcp.WithString("old_string", mcp.Description("Patch mode: the exact text to replace, whitespace included. It must occur exactly once in the stored source — read it with system_application_source_get first.")),
			mcp.WithString("new_string", mcp.Description("Patch mode: what old_string becomes. Pass an empty string to delete the matched text. Required whenever old_string is given.")),
			mcp.WithString("namespace", mcp.Description("Handle of the compose namespace the app reads from. Omit in patch mode to keep the stored one.")),
			mcp.WithString("modules", mcp.Description(`JSON array of module handles the app may query, as strings, e.g. ["Lead","Deal"]. This is the allowlist the bridge enforces; a module left out is refused at runtime. Omit in patch mode to keep the stored list.`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
			hmcp.NeedsFullDocs(),
		),
		"Set custom application source",
		h.sourceSet,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_application_delete",
			mcp.WithDescription(
				"Delete an application, removing it from the app selector for every user. "+
					"The delete is soft: the record is kept, is visible again by passing includeDeleted to "+
					"system_application_lookup, and is restored by system_application_undelete. "+
					"Prefer system_application_update with enabled false when you only want to take the "+
					"application out of circulation — that leaves it in the admin listing where someone can "+
					"find it again.",
			),
			mcp.WithString("application", mcp.Required(), mcp.Description("Exact application name, or an application ID as a string to prevent precision loss.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete application",
		h.delete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_application_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted application, putting it back in the app selector with all of its "+
					"settings intact. "+
					"This one takes an ID rather than a name, because name lookup only ever sees live "+
					"applications: call system_application_lookup with includeDeleted set to find the "+
					"applicationID of the deleted entry first.",
			),
			mcp.WithString("applicationID", mcp.Required(), mcp.Description("Application ID as a string, to prevent precision loss. Get it from system_application_lookup with includeDeleted set.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete application",
		h.undelete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_application_reorder",
			mcp.WithDescription(
				"Set the order applications appear in in the app selector. The applications you list are "+
					"weighted in the order given. "+
					"Applications are a single flat list for the whole instance, so pass the COMPLETE ordered "+
					"list: anything you leave out is pushed behind what you listed, keeping its current "+
					"relative order. "+
					"Call system_application_lookup first for the full set of applicationIDs and the current "+
					"weights. This changes ordering only — nothing else about an application — and it is the "+
					"only way to set weight, which the create and update tools deliberately do not expose.",
			),
			mcp.WithString("order", mcp.Required(), mcp.Description(`JSON array of application IDs in the desired order, as strings to prevent precision loss, e.g. ["123","456","789"]`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Reorder applications",
		h.reorder,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_application_flag",
			mcp.WithDescription(
				"Put a flag on an application. Flags are free-form labels the app selector reads, "+
					"\"pinned\" being the usual one; check what an existing application carries with "+
					"system_application_lookup before inventing a name. "+
					"'mode' is required and decides who the flag is for, so state it deliberately: \"own\" sets "+
					"a flag only you see, while \"global\" writes shared state that every user sees and needs a "+
					"separate, higher permission. You cannot flag on another user's behalf in either mode. "+
					"Flagging something already flagged in that mode fails rather than doing nothing.",
			),
			mcp.WithString("application", mcp.Required(), mcp.Description("Exact application name, or an application ID as a string to prevent precision loss.")),
			mcp.WithString("flag", mcp.Required(), mcp.Description("Flag name to set, e.g. \"pinned\".")),
			mcp.WithString("mode", mcp.Required(), mcp.Enum("own", "global"), mcp.Description("\"own\" flags it for you alone; \"global\" flags it for every user and requires the global-flag permission. There is no default — pick one.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Flag application",
		h.flag,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_application_unflag",
			mcp.WithDescription(
				"Take a flag off an application, reversing system_application_flag. "+
					"'mode' is required and must match the flag you are removing: \"own\" removes only the flag "+
					"you set yourself and never touches a global one, while \"global\" removes the flag for "+
					"every user and needs the global-flag permission. "+
					"One quirk worth knowing: if an application carries a global flag and you have also set the "+
					"same flag in \"own\" mode, unflagging in \"own\" mode suppresses the global flag for you "+
					"alone and leaves everyone else's view untouched. "+
					"Removing a flag that is not set fails rather than doing nothing.",
			),
			mcp.WithString("application", mcp.Required(), mcp.Description("Exact application name, or an application ID as a string to prevent precision loss.")),
			mcp.WithString("flag", mcp.Required(), mcp.Description("Flag name to remove, e.g. \"pinned\".")),
			mcp.WithString("mode", mcp.Required(), mcp.Enum("own", "global"), mcp.Description("\"own\" removes your own flag; \"global\" removes the shared flag for every user and requires the global-flag permission. There is no default — pick one.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Unflag application",
		h.unflag,
	)
}

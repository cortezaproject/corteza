package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Permissions are not a resource — RESOURCES.md keeps them out of the resource
// grammar and says the shape of any tool over them is decided case by case and
// recorded there first. It is recorded there now.
//
// Four tools, and the one thing to understand about all of them:
//
// **An agent may grant only what it already holds.** Reading is free; every
// write checks the operation against the calling user's own access on the same
// resource and refuses what they do not have. That is not a policy layered on
// top, it is the whole mechanism: the check runs the same RBAC evaluation the
// server would run if the caller performed the operation themselves, so a
// caller who cannot do X cannot give X away. Nothing here can raise anybody
// above the person driving it.
//
// The service layer's own gate is coarser than that — Grant checks one
// component-wide "grant" permission and then writes whatever rule it is handed,
// which is how the admin UI works and why a tool cannot simply pass a rule
// through. Both apply: the caller needs the component's grant permission AND
// the operation being granted.
//
// This file holds declarations only. The implementation is in
// permission_handler.go.

const permissionRulesDoc = `JSON array of rules. One rule: {"resource":"corteza::compose:module/511/512","operation":"record.create","access":"allow"}
resource is the full RBAC resource string, exactly as system_permission_schema prints it — component, resource type and every path segment. A short form like "module/*" is refused by the server, and so is a path with the wrong number of segments: a module is "corteza::compose:module/<namespaceID>/<moduleID>", a record is "corteza::compose:record/<namespaceID>/<moduleID>/<recordID>". A component-wide rule has an empty path and a trailing slash: "corteza::compose/".
A path segment may be "*" to cover every resource at that level; a rule written that way is only accepted if you hold the operation at that same breadth.
operation is one of the operations system_permission_schema lists for that resource type. Operations are per resource type, not global: "record.create" lives on the module, not on the record.
access is "allow" or "deny". Omit it and it is "allow". To remove a rule, use system_permission_revoke — writing "inherit" here is refused, so that taking access away is never something this tool does by accident.
Rules for different components may be mixed in one call; they are grouped and applied per component.`

func (h *permissionHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_permission_schema",
			mcp.WithDescription(
				"List every permission that can be granted: each resource type, the shape of its resource "+
					"string, and the operations valid on it. This is a fixed reference list built into "+
					"Human — read-only, the same on every instance, and not something a caller adds to. "+
					"Read it before system_permission_grant or system_permission_revoke. Both take a full "+
					"RBAC resource string and an operation, neither is guessable, and a wrong one is "+
					"refused by the server rather than stored — so this is the only place to get them "+
					"right. Two mistakes it prevents: a shortened resource ('module/*' instead of "+
					"'corteza::compose:module/*/*'), which the server rejects outright; and an operation "+
					"looked for on the wrong resource type — 'record.create' is an operation on the "+
					"module, not on the record, and letting a role into an application needs 'access' and "+
					"'read' together, on the application. "+
					"The 'resource' field of each entry is the wildcard form covering every resource of "+
					"that type. Replace the '*' segments with real IDs to scope a rule to one thing: "+
					"'corteza::compose:module/511/512' is one module, 'corteza::compose:module/511/*' is "+
					"every module in namespace 511. "+
					"Around 40 resource types across three components, returned whole by default; "+
					"'component' and 'resourceType' narrow it.",
			),
			mcp.WithString("component", mcp.Description("Narrow to one component: \"system\", \"compose\" or \"automation\". Accepts the full form (\"corteza::compose\") too. Omit for all three.")),
			mcp.WithString("resourceType", mcp.Description("Narrow to one resource type. Either the full form (\"corteza::compose:module\") or the bare tail (\"module\"). Case-insensitive. A value matching nothing answers with every known type rather than an empty result.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup grantable permissions",
		h.schema,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_permission_lookup",
			mcp.WithDescription(
				"Read the permission rules a role holds, and what you yourself may grant. "+
					"Pass 'role' — an ID, handle or name — for that role's rules; omit it to sweep every "+
					"role you can see, which is how you find out who already has access to something. "+
					"'resource' narrows to one resource string, and also accepts a prefix: "+
					"'corteza::compose:module/511' returns the rules on every module in that namespace. "+
					"A rule's 'access' is 'allow' or 'deny'. A permission with no rule at all is not "+
					"listed: it inherits, which for most things means denied. So an empty result means "+
					"the role holds nothing here, not that something failed. "+
					"When you pass 'resource', the answer also carries 'yourAccess' — the operations you "+
					"personally hold on it. That list is exactly the ceiling on what system_permission_grant "+
					"will let you give away, so reading it first turns a refusal into a decision. "+
					"Reading rules needs the 'grant' permission on the component the rules belong to; a "+
					"component you cannot manage is left out of the answer and named under "+
					"'componentsSkipped' rather than silently omitted.",
			),
			mcp.WithString("role", mcp.Description("Role ID as a string (to prevent precision loss), handle or name. Omit to cover every role you can see.")),
			mcp.WithString("resource", mcp.Description("Full RBAC resource string, or a prefix of one. Use system_permission_schema for the shapes. Omit for every resource.")),
			mcp.WithString("limit", mcp.Description("Maximum rules returned, default 50, capped at 200.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup permission rules",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_permission_grant",
			mcp.WithDescription(
				"Give a role permissions. This is the second half of access control: "+
					"system_role_create makes an empty container and system_role_member_add puts people "+
					"in it, but until a rule is granted the role allows nothing. "+
					"YOU CAN ONLY GRANT WHAT YOU HOLD. Every rule is checked against your own access on "+
					"that exact resource, evaluated the same way the server evaluates it when you act "+
					"yourself. An operation you cannot perform is refused, by name, and the whole call is "+
					"rejected rather than partly applied. This is structural, not a policy: there is no "+
					"argument, no role and no resource that lets this tool hand out more than the person "+
					"driving it already has. You also need the 'grant' permission on the component — "+
					"holding an operation is not the same as being allowed to delegate it. "+
					"Get the resource strings and operation names from system_permission_schema; they are "+
					"not guessable and a wrong one is refused rather than stored. Check what you can give "+
					"with system_permission_lookup, which returns 'yourAccess' for a resource. "+
					"A rule replaces any existing rule for the same role, resource and operation. Rules "+
					"for resources you did not name are untouched — unlike system_role_clone_rules, which "+
					"replaces a role's whole rule set. "+
					"Effect is immediate for new sessions, but an already signed-in user keeps the access "+
					"their session was built with until they sign in again. "+
					"'deny' is a rule like any other and outranks an allow inherited from anywhere else, "+
					"so use it deliberately. To take a rule away rather than override it, use "+
					"system_permission_revoke.",
			),
			mcp.WithString("role", mcp.Required(), mcp.Description("Role to grant to: ID as a string (to prevent precision loss), handle or name. Create one with system_role_create.")),
			mcp.WithString("rules", mcp.Required(), mcp.Description(permissionRulesDoc)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Grant permissions to a role",
		h.grant,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_permission_revoke",
			mcp.WithDescription(
				"Take permission rules off a role. The rule is deleted rather than set to 'deny': the "+
					"permission goes back to inheriting, which for most things means the role no longer "+
					"allows it, but a rule inherited from elsewhere — a wildcard rule, or the "+
					"Authenticated role — takes over again. If you need the operation blocked outright, "+
					"grant 'deny' with system_permission_grant instead. "+
					"THE SAME CEILING APPLIES AS ON GRANT: you can only revoke an operation you yourself "+
					"hold on that resource, and you need the component's 'grant' permission. Revoking is "+
					"a privilege change like any other and is not a way around the boundary. "+
					"Read the role's rules with system_permission_lookup first — revoking a rule that was "+
					"never there succeeds and changes nothing, which reads as if it worked. "+
					"Nothing here deletes a role or removes a member; that is system_role_delete and "+
					"system_role_member_remove.",
			),
			mcp.WithString("role", mcp.Required(), mcp.Description("Role to revoke from: ID as a string (to prevent precision loss), handle or name.")),
			mcp.WithString("rules", mcp.Required(), mcp.Description("JSON array of rules to remove: [{\"resource\":\"corteza::compose:module/511/512\",\"operation\":\"record.create\"}]. An 'access' key is ignored here — revoking always means back to inherit. Resource and operation must match the granted rule exactly; system_permission_lookup prints them.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Revoke permissions from a role",
		h.revoke,
	)
}

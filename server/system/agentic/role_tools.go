package agentic

import (
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// Role is the privilege-escalation surface of the tool registry. A role is a
// permission container: whatever is granted to it is held by every user and
// user group that is a member of it, so create, update, member_add and
// clone_rules are the four ways an agent could widen access — its own included.
// Full CRUD is deliberately in scope, which makes the descriptions load-bearing:
// each write says plainly what it grants, and none of them is phrased as
// routine bookkeeping.
//
// Two service facts every description leans on. System roles (the bypass,
// authenticated and anonymous roles the instance runs on) refuse delete,
// archive, unarchive and undelete outright, and refuse any change to handle or
// name; closed and context roles refuse membership changes. Both are visible as
// isSystem / isClosed in system_role_lookup, so a caller can tell before it
// tries.
//
// This file holds declarations only. Implementations are in role_handler.go, in
// the same order.

func (h *roleHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_role_lookup",
			mcp.WithDescription(
				"Look up a role by ID, handle or name, or list and filter roles. Provide 'role' to fetch "+
					"one, which returns the whole role plus isSystem and isClosed; omit it to list, which "+
					"returns a compact form (roleID, name, handle, isSystem, isClosed, archivedAt, deletedAt). "+
					"Read this before any role write. A role is a permission container: whatever it is granted "+
					"is held by everyone in it, so knowing which role you are about to touch matters. isSystem "+
					"marks the roles the instance itself runs on — the bypass, authenticated and anonymous "+
					"roles — which refuse delete, archive and undelete and cannot be renamed; isClosed marks "+
					"roles whose membership cannot be changed. "+
					"Use system_role_member_list to see who is in a role, and the 'member' filter here to see "+
					"which roles a user holds. "+
					"Deleted and archived roles are excluded unless you pass includeDeleted or "+
					"includeArchived, and a deleted role can then only be fetched by ID, because handle and "+
					"name lookups skip deleted roles.",
			),
			mcp.WithString("role", mcp.Description("Role ID as a string, handle, or name. Omit to list instead.")),
			mcp.WithString("query", mcp.Description("Free-text match against handle and name.")),
			mcp.WithString("name", mcp.Description("Exact name match.")),
			mcp.WithString("handle", mcp.Description("Exact handle match.")),
			mcp.WithString("member", mcp.Description("Only roles this user is a member of. User ID as a string, handle or email. Cannot be combined with userGroup.")),
			mcp.WithString("userGroup", mcp.Description("Only roles this user group is a member of. Group ID as a string or handle. Cannot be combined with member.")),
			mcp.WithBoolean("includeDeleted", mcp.Description("Include soft-deleted roles. This is how you find the roleID that system_role_undelete needs.")),
			mcp.WithBoolean("includeArchived", mcp.Description("Include archived roles, which are otherwise hidden.")),
			mcp.WithString("limit", mcp.Description("Maximum results, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup role",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_role_create",
			mcp.WithDescription(
				"Create a role. This is the first half of granting access: a role is a permission container, "+
					"and every user or user group later added to it receives everything the role is permitted "+
					"to do. Creating one requires the create-role permission and is recorded in the action log. "+
					"A new role starts with no permissions and no members, so on its own it grants nothing — it "+
					"becomes real when permissions are assigned to it (system_role_clone_rules copies another "+
					"role's entire permission set onto it) and members are added with system_role_member_add. "+
					"'handle' is the stable identifier other tools accept in place of an ID; a role without one "+
					"can only be referenced by ID or name. Handle and name must each be unique.",
			),
			mcp.WithString("name", mcp.Required(), mcp.Description("Human-readable name, e.g. \"Support agents\". Must be unique.")),
			mcp.WithString("handle", mcp.Description("Stable identifier, e.g. \"support-agents\". Letters, digits, dash and underscore, at least 2 characters. Must be unique.")),
			mcp.WithString("description", mcp.Description("What the role is for. Free text; say what it grants, since nothing else records that.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create role",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_role_update",
			mcp.WithDescription(
				"Rename a role or change its description. Handle and name are what other tools and stored "+
					"configuration use to refer to a role, so changing them can break anything that names the "+
					"old value. "+
					"System roles are refused a rename: the bypass, authenticated and anonymous roles reject "+
					"any change to handle or name, and only a description-only edit goes through. Check "+
					"isSystem with system_role_lookup first rather than calling this and retrying. "+
					"Omit a field to leave it unchanged; pass it empty to clear it. "+
					"This tool changes neither the role's permissions nor its members — use "+
					"system_role_clone_rules for permissions and system_role_member_add or "+
					"system_role_member_remove for members.",
			),
			mcp.WithString("role", mcp.Required(), mcp.Description("Role ID as a string, handle, or name.")),
			mcp.WithString("name", mcp.Description("New name. Empty clears it. Must stay unique.")),
			mcp.WithString("handle", mcp.Description("New handle. Empty clears it, which leaves the role referenceable only by ID or name. Must stay unique.")),
			mcp.WithString("description", mcp.Description("New description. Empty clears it. Other role metadata, including any context expression, is left untouched.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update role",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_role_delete",
			mcp.WithDescription(
				"Delete a role. Everyone who held access through this role loses it: a deleted role stops "+
					"being applied, so every permission it granted stops taking effect for every one of its "+
					"members at once. "+
					"The delete is soft — the role and its permission rules are retained, system_role_lookup "+
					"with includeDeleted lists it, and system_role_undelete restores it with its rules intact. "+
					"System roles (bypass, authenticated, anonymous) are refused. "+
					"If you want the role out of use but expect to bring it back, prefer system_role_archive. "+
					"If you want one user to stop holding the role, use system_role_member_remove rather than "+
					"deleting the role for everybody.",
			),
			mcp.WithString("role", mcp.Required(), mcp.Description("Role ID as a string, handle, or name.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete role",
		h.delete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_role_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted role, and with it every permission that role granted to its members. "+
					"This re-grants access — it does not merely make the role visible again — so check what "+
					"the role is before restoring it. "+
					"Takes 'roleID' rather than a handle or name because a deleted role is only reachable by "+
					"ID: handle and name lookups skip deleted roles. Get the ID from system_role_lookup with "+
					"includeDeleted. System roles are refused.",
			),
			mcp.WithString("roleID", mcp.Required(), mcp.Description("Role ID as a string. Handles and names do not resolve to deleted roles.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete role",
		h.undelete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_role_archive",
			mcp.WithDescription(
				"Archive a role: it stops being applied, so its members stop receiving what it grants, and it "+
					"drops out of default listings — but nothing is removed and system_role_unarchive puts it "+
					"straight back. Use this to retire a role you may want later. "+
					"Archiving is not deleting: system_role_delete marks the role deleted and hides it from "+
					"everything, archiving parks it. It is also not suspending a user — it affects every "+
					"member of the role at once, so to stop one person holding the role use "+
					"system_role_member_remove. System roles are refused.",
			),
			mcp.WithString("role", mcp.Required(), mcp.Description("Role ID as a string, handle, or name.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Archive role",
		h.archive,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_role_unarchive",
			mcp.WithDescription(
				"Return an archived role to use. Its permissions start applying to its members again, so this "+
					"widens access by exactly as much as the role grants — read the role with "+
					"system_role_lookup before calling it. "+
					"Has no effect on a role that was never archived, and cannot recover a deleted one: use "+
					"system_role_undelete for that. System roles are refused.",
			),
			mcp.WithString("role", mcp.Required(), mcp.Description("Role ID as a string, handle, or name. Archived roles need includeArchived to show up in system_role_lookup.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Unarchive role",
		h.unarchive,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_role_member_list",
			mcp.WithDescription(
				"List the members of a role — the users and user groups that hold whatever the role grants. "+
					"Each entry is a kind (user or userGroup) and an ID; resolve those with the user or user "+
					"group tools if you need names. "+
					"Returns every member in one response with no paging, so call it for one role at a time. "+
					"Requires read access to the role, and is refused for closed roles and for context roles, "+
					"whose membership is computed from an expression rather than stored. "+
					"To go the other way — which roles one user holds — list roles with system_role_lookup and "+
					"its 'member' filter.",
			),
			mcp.WithString("role", mcp.Required(), mcp.Description("Role ID as a string, handle, or name.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"List role members",
		h.memberList,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_role_member_add",
			mcp.WithDescription(
				"Add a user or a user group to a role, granting that user — or every member of that group — "+
					"everything the role is permitted to do. This is a privilege change, not bookkeeping: read "+
					"the role first and do not add anyone to a role whose permissions you have not inspected. "+
					"Pass exactly one of 'user' or 'userGroup'. Requires permission to manage members on the "+
					"role. Closed roles and context roles refuse membership changes; system_role_lookup shows "+
					"isClosed. "+
					"Reversible with system_role_member_remove.",
			),
			mcp.WithString("role", mcp.Required(), mcp.Description("Role ID as a string, handle, or name. This is the role whose permissions the member will hold.")),
			mcp.WithString("user", mcp.Description("User ID as a string, handle, or email address. Pass this or userGroup, not both.")),
			mcp.WithString("userGroup", mcp.Description("User group ID as a string, or handle. Adds the whole group, so every current and future member of it holds the role.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Add role member",
		h.memberAdd,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_role_member_remove",
			mcp.WithDescription(
				"Remove a user or a user group from a role, withdrawing everything that role granted them. "+
					"Any other role they hold is unaffected, so this does not necessarily leave them without "+
					"access — check system_role_lookup with the 'member' filter if that is the goal. "+
					"Pass exactly one of 'user' or 'userGroup'. Removing a user does not remove them from a "+
					"user group that is itself a member of the role: if their access comes through the group, "+
					"remove the group or remove them from it. "+
					"Requires permission to manage members on the role; closed and context roles refuse "+
					"membership changes. Reversible with system_role_member_add.",
			),
			mcp.WithString("role", mcp.Required(), mcp.Description("Role ID as a string, handle, or name.")),
			mcp.WithString("user", mcp.Description("User ID as a string, handle, or email address. Pass this or userGroup, not both.")),
			mcp.WithString("userGroup", mcp.Description("User group ID as a string, or handle. Removes the whole group's membership of the role.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Remove role member",
		h.memberRemove,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_role_clone_rules",
			mcp.WithDescription(
				"Copy every permission rule from one role onto another, replacing the target's rules. This is "+
					"the widest-reaching write in this tool surface: the target ends up with exactly the "+
					"source's access, so every member of the target immediately holds everything members of "+
					"the source hold. Never use it to 'start from' a powerful role, and confirm with a human "+
					"before pointing it at any role you do not fully understand. "+
					"It is destructive for the target as well: the target's existing rules are discarded, not "+
					"merged with the source's, and there is no undo. The source role is left untouched. "+
					"Requires permission to manage permissions across the instance, which is a higher bar than "+
					"editing a role. To review either side first, read both roles with system_role_lookup.",
			),
			mcp.WithString("role", mcp.Required(), mcp.Description("Source role to copy permission rules FROM. ID as a string, handle, or name. Unchanged by this call.")),
			mcp.WithString("targetRole", mcp.Required(), mcp.Description("Target role to copy permission rules ONTO. ID as a string, handle, or name. Its own rules are discarded first.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			// Destructive, not write: rbac/service.go sets every existing rule
			// on the target to Inherit before copying, so the target's prior
			// permission set is discarded with no undo. §4's test is "removes,
			// or loses data" and this loses data.
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Clone role permission rules",
		h.cloneRules,
	)
}

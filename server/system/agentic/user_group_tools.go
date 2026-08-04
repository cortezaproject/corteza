package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// User group tools. See mcp/CONVENTIONS.md and mcp/briefs/system_user_group.md.
//
// Two things shape every description here and are worth stating once:
//
//   - Membership feeds RBAC. The group graph is mirrored into the permission
//     graph (rbacUserGroupService in system/service/user_group.go), so adding a
//     member or re-parenting a group changes what somebody can do, immediately.
//   - A user belongs to exactly one group. system/types/User carries a single
//     UserGroupID, which is why there is a member_add and no member_remove: the
//     service exposes no removal operation at all.
//
// This file holds declarations only. Implementations are in
// user_group_handler.go, in the same order.

func (h *userGroupHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_user_group_lookup",
			mcp.WithDescription(
				"Look up a user group by ID or handle, or list and filter user groups. User groups are the "+
					"hierarchy that permissions are granted through: a group inherits from its parents, and "+
					"every member of a group gets what the group and its ancestors are allowed to do. "+
					"Provide 'userGroup' to fetch one group in full, including its parent links, description "+
					"and labels; omit it to list groups as {userGroupID, handle, short, isRoot, parentIDs, "+
					"archivedAt, deletedAt}. "+
					"A user group has no name field — the handle is its identifier and 'short' is the only "+
					"human label it carries. "+
					"Members are never returned here, by either mode; call system_user_group_member_list for "+
					"those. To find which group a given user is in, read the 'userGroupID' field on the user "+
					"instead — this tool cannot filter by member. "+
					"Archived and deleted are two distinct states and both are excluded by default, so a "+
					"group you archived or deleted will not appear, and will not resolve by handle either, "+
					"until you set 'includeArchived' or 'includeDeleted'.",
			),
			mcp.WithString("userGroup", mcp.Description("User group handle, or ID as a string to prevent precision loss. Omit to list groups instead.")),
			mcp.WithString("query", mcp.Description("Substring match against the handle. Only the handle is searched — descriptions are not indexed.")),
			mcp.WithBoolean("includeDeleted", mcp.Description("Include soft-deleted groups. Needed to find a group before calling system_user_group_undelete.")),
			mcp.WithBoolean("includeArchived", mcp.Description("Include archived groups. Archiving is separate from deleting; a group can be one, both or neither.")),
			mcp.WithString("limit", mcp.Description("Maximum groups to return when listing, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup user group",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_group_create",
			mcp.WithDescription(
				"Create a user group. Groups form a hierarchy that permissions flow down, so where you "+
					"attach a new group decides what its members will be able to do. "+
					"Every group except the instance's single root group must name at least one parent in "+
					"'parents'; creating one without a parent is rejected as an invalid parent reference. "+
					"Call system_user_group_lookup first to find the parent you want. "+
					"This creates an empty group — members are added afterwards, one at a time, with "+
					"system_user_group_member_add.",
			),
			mcp.WithString("handle", mcp.Required(), mcp.Description("Unique handle, at least 2 characters, starting with a letter and ending alphanumeric; letters, digits, underscore, hyphen and dot in between. This is how the group is referenced everywhere.")),
			mcp.WithString("parents", mcp.Description("JSON array of parent groups, e.g. [\"engineering\"] or [{\"parent\":\"engineering\",\"name\":\"reporting-line\"}]. Each 'parent' is a group handle or ID as a string. Give every entry a distinct 'name' when there is more than one parent — duplicate or missing names are rejected. Omit only for the root group.")),
			mcp.WithString("short", mcp.Description("Short human label for the group, shown where a handle would read badly.")),
			mcp.WithString("description", mcp.Description("Longer description of what the group is for.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create user group",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_group_update",
			mcp.WithDescription(
				"Update a user group's handle, labels or parent links. Omit a field to leave it unchanged; "+
					"pass it as an empty string to clear it. 'parents' replaces the existing parent links "+
					"wholesale rather than merging, so send the full set you want. "+
					"Re-parenting is not cosmetic: the group moves in the permission graph and its members "+
					"immediately gain or lose whatever the old and new ancestors granted. "+
					"This cannot update the root group — the service requires at least one parent on every "+
					"update and the root group has none, so an update of it fails with a missing-parent "+
					"error whatever you send. "+
					"Members are not touched here; use system_user_group_member_add.",
			),
			mcp.WithString("userGroup", mcp.Required(), mcp.Description("User group handle, or ID as a string to prevent precision loss.")),
			mcp.WithString("handle", mcp.Description("New handle. Must stay unique across all groups.")),
			mcp.WithString("parents", mcp.Description("Replacement JSON array of parent groups, e.g. [\"engineering\"] or [{\"parent\":\"engineering\",\"name\":\"reporting-line\"}]. Replaces every existing parent link. Omit to leave the current parents alone.")),
			mcp.WithString("short", mcp.Description("New short label. Empty clears it.")),
			mcp.WithString("description", mcp.Description("New description. Empty clears it.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update user group",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_group_delete",
			mcp.WithDescription(
				"Delete a user group. The delete is soft and reversible with system_user_group_undelete, "+
					"but it takes the group out of the permission graph straight away, so everyone in it "+
					"loses whatever that group and its ancestors granted them. "+
					"Users are not deleted and still point at the group, so undeleting restores their "+
					"permissions. Child groups are not deleted either and are left pointing at a group that "+
					"is no longer in the graph — re-parent them first with system_user_group_update. "+
					"A group used as the security user group of an auth client cannot be deleted; that is "+
					"reported as a permission error, so check the auth clients before assuming you lack "+
					"rights.",
			),
			mcp.WithString("userGroup", mcp.Required(), mcp.Description("User group handle, or ID as a string to prevent precision loss.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete user group",
		h.delete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_group_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted user group and put it back into the permission graph, which "+
					"restores what its members could do. Find the group first with "+
					"system_user_group_lookup and 'includeDeleted' set — a deleted group does not resolve "+
					"by handle or ID without it. "+
					"This is the way back from system_user_group_delete; it has no effect on a group that "+
					"was never deleted.",
			),
			mcp.WithString("userGroup", mcp.Required(), mcp.Description("User group handle, or ID as a string to prevent precision loss.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete user group",
		h.undelete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_group_member_list",
			mcp.WithDescription(
				"List the users who are members of one user group, as {userID, handle, name}. "+
					"Membership is deliberately kept out of system_user_group_lookup because it is the "+
					"heavy part of a group, so this is the only way to read it. "+
					"Direct members only — users in child groups are not included, even though they inherit "+
					"the group's permissions. "+
					"This returns every member in one response and cannot be paged, so a very large group "+
					"may exceed the result size limit. Email addresses and other user detail are not "+
					"returned; look the user up by ID if you need them.",
			),
			mcp.WithString("userGroup", mcp.Required(), mcp.Description("User group handle, or ID as a string to prevent precision loss.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"List user group members",
		h.memberList,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_group_member_add",
			mcp.WithDescription(
				"Add a user to a user group. This changes what that user is allowed to do: group membership "+
					"is mirrored into the permission graph, so the user immediately gains everything this "+
					"group and its ancestors grant. Treat it as a privilege change and confirm the group is "+
					"the one intended before calling. "+
					"A user belongs to exactly one group at a time, so this moves them: whatever group they "+
					"were in before, they are no longer in, and they lose the permissions that came with it. "+
					"There is no companion removal tool because the service exposes no removal operation — "+
					"to take a user out of a group, add them to a different one.",
			),
			mcp.WithString("userGroup", mcp.Required(), mcp.Description("Target user group handle, or ID as a string to prevent precision loss.")),
			mcp.WithString("user", mcp.Required(), mcp.Description("User handle, email, or ID as a string to prevent precision loss.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Add user group member",
		h.memberAdd,
	)
}

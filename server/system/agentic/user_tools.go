package agentic

import (
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCP tools for system/user. See mcp/CONVENTIONS.md for the rules and
// mcp/briefs/system_user.md for the per-resource decisions.
//
// Two operations the user service exposes are deliberately absent:
//
//   - SetPassword. It sets a credential, and CONVENTIONS.md §8.6b rules that a
//     tool which sets or returns a credential is not written. The service
//     authorizes it correctly; the objection is that no agent workflow should
//     be setting a person's password without a human in the loop, and there is
//     no confirmation step on this surface.
//   - CreateWithAvatar / UpdateWithAvatar / UploadAvatar. They take an
//     io.Reader or a multipart file header, and a tool call carries neither, so
//     there is nothing honest to declare.
//
// This file holds declarations only. Implementations are in user_handler.go,
// in the same order.

func (h *userHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_user_lookup",
			mcp.WithDescription(
				"Look up one user by ID, handle or email address, or list and filter the user accounts on "+
					"this Human instance. Provide 'user' to fetch one, which returns the full record including "+
					"meta and labels; omit it to list, which returns a compact form (userID, email, handle, "+
					"name, kind, suspendedAt, deletedAt) — fetch by ref when you need more than that. "+
					"Listing hides deleted and suspended accounts unless you ask for them, and returns only "+
					"ordinary users unless you set 'kind' or 'allKinds'. This is also the only way to obtain "+
					"the numeric userID that system_user_undelete requires.",
			),
			mcp.WithString("user", mcp.Description("User reference: numeric ID as a string, a handle, or an email address — anything containing '@' is read as an email. Usernames are not resolved. Omit to list instead.")),
			mcp.WithString("query", mcp.Description("Free-text search, matched as a substring against email, username, handle and name.")),
			mcp.WithString("email", mcp.Description("Exact email address. Unlike 'query' this must match in full.")),
			mcp.WithString("handle", mcp.Description("Exact handle. Unlike 'query' this must match in full.")),
			mcp.WithString("username", mcp.Description("Exact username. A legacy field that is empty on most accounts; prefer 'handle' or 'email'.")),
			mcp.WithString("role", mcp.Description("Only users who are members of this role. Accepts a role ID as a string, a handle, or a role name.")),
			mcp.WithString("userGroup", mcp.Description("Only users belonging to this user group. Accepts a group ID as a string, a handle, or a group name.")),
			mcp.WithString("kind", mcp.Description("User kind. Empty (the default) means ordinary user accounts; \"sys\" means built-in system accounts. Those are the only two values.")),
			mcp.WithBoolean("allKinds", mcp.Description("Return ordinary and system accounts together, ignoring 'kind'.")),
			mcp.WithBoolean("includeDeleted", mcp.Description("Also return soft-deleted users. Off by default. Set this to find the numeric userID that system_user_undelete needs.")),
			mcp.WithBoolean("includeSuspended", mcp.Description("Also return suspended users. Off by default, so a suspended account is invisible to a plain list.")),
			mcp.WithString("limit", mcp.Description("Maximum results, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup user",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_create",
			mcp.WithDescription(
				"Create a user account. Only the email address is required; when you omit 'handle' one is "+
					"derived from the name and email. The account is created with its email marked confirmed "+
					"and with no password, so the person cannot sign in until a credential exists — setting "+
					"one is deliberately not available through this tool surface, so send them through the "+
					"normal sign-in or password-reset flow. "+
					"To put the new user in a user group afterwards call system_user_group_member_add rather "+
					"than system_user_update: only that call also refreshes the permission graph.",
			),
			mcp.WithString("email", mcp.Required(), mcp.Description("Email address. Must be a valid address and unique among users that are not deleted. This is what the person signs in with.")),
			mcp.WithString("name", mcp.Description("Display name, e.g. \"Ada Lovelace\". Optional but used to derive the handle and the avatar initials.")),
			mcp.WithString("handle", mcp.Description("URL-safe handle, letters/digits/._- starting with a letter, at least 2 characters. Omit to have one generated from the name and email.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create user",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_update",
			mcp.WithDescription(
				"Update a user's email address, display name or handle. Omit a field to leave it unchanged; "+
					"pass 'name' or 'handle' as an empty string to clear them. "+
					"'email' cannot be cleared and must stay a valid, unique address: it is the identity the "+
					"person signs in with and the address notifications go to, so changing it changes both. "+
					"It does not send or require a new confirmation mail — use "+
					"system_user_set_email_confirmed if you need the confirmation flag changed. "+
					"Kind, avatar and user group are not editable here; to move someone between groups call "+
					"system_user_group_member_add, which also refreshes the permission graph.",
			),
			mcp.WithString("user", mcp.Required(), mcp.Description("User reference: numeric ID as a string, a handle, or an email address.")),
			mcp.WithString("email", mcp.Description("New email address. Must be valid and unique. Cannot be cleared.")),
			mcp.WithString("name", mcp.Description("New display name. Empty clears it.")),
			mcp.WithString("handle", mcp.Description("New handle. Empty clears it. Must be unique among users that are not deleted.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update user",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_delete",
			mcp.WithDescription(
				"Delete a user account. The delete is soft — the record is kept and system_user_undelete "+
					"restores it — but the account disappears from every listing and its access tokens are "+
					"revoked immediately, so treat this as removing the person. "+
					"If you only want to stop someone signing in while keeping the account visible and easily "+
					"reversible, use system_user_suspend instead. "+
					"Before deleting, note that undelete needs the numeric userID, which afterwards is only "+
					"obtainable from system_user_lookup with includeDeleted set.",
			),
			mcp.WithString("user", mcp.Required(), mcp.Description("User reference: numeric ID as a string, a handle, or an email address.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete user",
		h.delete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted user, reversing system_user_delete. "+
					"This takes the numeric user ID and nothing else: handle and email lookups do not resolve "+
					"deleted accounts, so a reference that worked before the delete will fail here. Get the ID "+
					"from system_user_lookup with includeDeleted set. "+
					"The restore is refused if the email, handle or username of the deleted account now "+
					"collides with a live user, or if it would take the instance over its user limit.",
			),
			mcp.WithString("userID", mcp.Required(), mcp.Description("Numeric user ID as a string, to prevent precision loss. Obtain it from system_user_lookup with includeDeleted set.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete user",
		h.undelete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_suspend",
			mcp.WithDescription(
				"Suspend a user account: the person can no longer sign in and their existing access tokens "+
					"are revoked, but the account, its data, and its group and role membership are all kept. "+
					"Reversible with system_user_unsuspend. "+
					"Prefer this over system_user_delete whenever the person may come back or the account "+
					"merely needs freezing — suspend is about access, delete is about removal. "+
					"A suspended account drops out of system_user_lookup listings unless you pass "+
					"includeSuspended, though it stays resolvable by ID, handle or email.",
			),
			mcp.WithString("user", mcp.Required(), mcp.Description("User reference: numeric ID as a string, a handle, or an email address.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Suspend user",
		h.suspend,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_unsuspend",
			mcp.WithDescription(
				"Lift a suspension so the user can sign in again, reversing system_user_suspend. "+
					"Accepts the same references as suspend, because a suspended account is still resolvable "+
					"by ID, handle or email. "+
					"This cannot recover a deleted account — use system_user_undelete for that — and it fails "+
					"if reactivating the user would take the instance over its configured user limit.",
			),
			mcp.WithString("user", mcp.Required(), mcp.Description("User reference: numeric ID as a string, a handle, or an email address.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Unsuspend user",
		h.unsuspend,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_user_set_email_confirmed",
			mcp.WithDescription(
				"Mark a user's email address as confirmed, or withdraw that confirmation. This flips the "+
					"stored flag only: it neither sends a confirmation mail nor validates one, and it does not "+
					"change the address itself — use system_user_update for that. "+
					"Use it to unblock someone stuck in the sign-up confirmation flow, or to force a "+
					"re-confirmation after an address is corrected. Accounts made with system_user_create are "+
					"already confirmed, so a fresh user rarely needs this.",
			),
			mcp.WithString("user", mcp.Required(), mcp.Description("User reference: numeric ID as a string, a handle, or an email address.")),
			mcp.WithBoolean("confirmed", mcp.Required(), mcp.Description("True marks the address confirmed, false marks it unconfirmed. Must be given explicitly.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Set user email confirmation",
		h.setEmailConfirmed,
	)
}

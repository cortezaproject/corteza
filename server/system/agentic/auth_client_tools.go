package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Auth clients are the OAuth2 credentials third-party applications use against
// this instance, so this family draws a line the other families do not have to:
// no tool here reads back or replaces the secret of a client that already
// exists. ExposeSecret and RegenerateSecret are both properly access-checked,
// but read permission is a lower bar than an existing credential deserves and
// the disclosure is irreversible — a secret that reaches a model's context is in
// the transcript and in the logs for good.
//
// Create is on the other side of that line, and deliberately: the credential
// does not exist until the caller asks for it, the caller is the only party it
// is ever shown to, and the ceiling is the same RBAC check every other write
// tool answers to. It is the one place a secret leaves the server, it says so,
// and it is a ruling rather than a slip. See CONVENTIONS.md §8.6b and
// RESOURCES.md.
//
// This file holds declarations only. Implementations are in
// auth_client_handler.go, in the same order.

func (h *authClientHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_auth_client_lookup",
			mcp.WithDescription(
				"Look up an OAuth2 auth client by ID or handle, or list and filter auth clients. An auth "+
					"client is the identity an external application, integration or script uses to obtain a "+
					"token from this instance. Provide 'authClient' to fetch one, which returns its full "+
					"configuration; omit it to list, which returns a compact form (ID, handle, name, enabled, "+
					"isDefault, validFrom, expiresAt, deletedAt). "+
					"No response from this tool ever includes a client secret. A client's secret is shown "+
					"once, by system_auth_client_create, and there is no tool that reads it back or "+
					"regenerates it — a client whose secret has been lost is replaced rather than recovered, "+
					"or its secret is read by a person in the admin UI. "+
					"Soft-deleted clients are hidden unless you set 'includeDeleted'.",
			),
			mcp.WithString("authClient", mcp.Description("Auth client ID as a string, to prevent precision loss, or the client's handle. Omit to list instead.")),
			mcp.WithString("handle", mcp.Description("List only the client whose handle is exactly this. Ignored when 'authClient' is given; unlike a lookup by handle, an unknown value returns an empty list rather than an error.")),
			mcp.WithBoolean("includeDeleted", mcp.Description("Include soft-deleted auth clients. They are hidden by default, so this is how you find one to restore with system_auth_client_undelete.")),
			mcp.WithString("limit", mcp.Description("Maximum results, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup auth client",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_auth_client_create",
			mcp.WithDescription(
				"Create an OAuth2 auth client — the identity an external application, integration or "+
					"script uses to obtain tokens from this instance. "+
					"THIS IS THE ONE TOOL THAT HANDS BACK A CREDENTIAL. The server generates the client "+
					"secret during creation and this returns it, once. Nothing reads it back afterwards: "+
					"there is no tool that exposes or regenerates a secret, so a secret that is lost means "+
					"creating a new client. Put it where it is going before you go on, and do not repeat it "+
					"anywhere it does not need to be. "+
					"Which grant to use decides what else you must supply. \"authorization_code\" is the "+
					"flow where a person signs in and is sent back to 'redirectURI', which is then required. "+
					"\"client_credentials\" has no person in it: the client authenticates as itself and acts "+
					"as one nominated user, so 'impersonateUser' (from system_user_lookup) is required and everything the client does "+
					"is done with that user's permissions. Choosing a user with more access than you have is "+
					"how a client ends up more powerful than the person who made it — pick the narrowest "+
					"account that can do the job. "+
					"'scope' is what the token may be used for and is almost always \"profile api\": the API "+
					"middleware requires the 'api' scope, so a client without it authenticates and is then "+
					"refused by every endpoint. "+
					"The client is enabled on creation unless you say otherwise. 'validFrom' and 'expiresAt' "+
					"bound when it works at all and are the clean way to issue a credential that stops "+
					"working on its own. "+
					"Nothing here sets the client's permitted, prohibited or forced roles; those are a "+
					"person's job in the admin UI, and system_auth_client_update does not touch them either.",
			),
			mcp.WithString("name", mcp.Required(), mcp.Description("Display name, shown in listings. An auth client must have one.")),
			mcp.WithString("handle", mcp.Description("Short URL-safe identifier used in OAuth2 requests, e.g. \"reporting_service\". Optional, but a client without one can only ever be referenced by its numeric ID.")),
			mcp.WithString("description", mcp.Description("What this client is for and who runs it. Worth writing: a credential nobody can attribute is a credential nobody dares revoke.")),
			mcp.WithString("validGrant", mcp.Description("OAuth2 grant type: \"authorization_code\" (a person signs in; needs redirectURI) or \"client_credentials\" (no person; needs impersonateUser). Omit to leave it unset, which means the client cannot obtain a token until system_auth_client_update sets one.")),
			mcp.WithString("impersonateUser", mcp.Description("The user a client_credentials client acts as: user ID as a string (to prevent precision loss), handle or email. REQUIRED when validGrant is client_credentials, and rejected by the server otherwise. Everything the client does carries this user's permissions.")),
			mcp.WithString("redirectURI", mcp.Description("Absolute URL the person is sent back to after authorizing. Required in practice for authorization_code; meaningless for client_credentials.")),
			mcp.WithString("scope", mcp.Description("Space-separated OAuth2 scopes the client may request. Use \"profile api\" unless you have a reason not to — without the 'api' scope every API call is refused after a successful login.")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the client may obtain tokens. Defaults to true. Pass false to create it dormant and enable it later with system_auth_client_update.")),
			mcp.WithBoolean("trusted", mcp.Description("A trusted client skips the consent screen a person is otherwise shown when authorizing it. Set it only for clients this instance itself owns.")),
			mcp.WithString("validFrom", mcp.Description("Client cannot be used before this time. RFC3339, e.g. 2026-08-03T09:00:00Z. Omit for no lower bound.")),
			mcp.WithString("expiresAt", mcp.Description("Client cannot be used after this time. RFC3339. Omit for a credential that never expires on its own.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create auth client",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_auth_client_update",
			mcp.WithDescription(
				"Update an auth client's configuration: its handle, name and description, whether it is "+
					"enabled and trusted, its OAuth2 grant type, redirect URI, scope and validity window. "+
					"Omit a field to leave it unchanged; pass it empty to clear it. "+
					"This cannot read, set or regenerate the client secret, so an update never disturbs the "+
					"credential the client is already using. A secret is shown once, by "+
					"system_auth_client_create, and never again. "+
					"It also leaves the client's security settings alone (impersonated user, permitted, "+
					"prohibited and forced roles): those grant privileges and are set when the client is "+
					"created or by a person in the admin UI, which is why setting 'validGrant' to "+
					"client_credentials is rejected unless an impersonation user is already configured. "+
					"The instance's default auth client — the one 'isDefault' marks in a lookup — cannot have "+
					"its handle changed and cannot be disabled.",
			),
			mcp.WithString("authClient", mcp.Required(), mcp.Description("Auth client ID as a string, to prevent precision loss, or the client's handle.")),
			mcp.WithString("handle", mcp.Description("New handle: a short URL-safe identifier used in OAuth2 requests. Empty clears it. Cannot be changed on the default client.")),
			mcp.WithString("name", mcp.Description("New display name. Cannot be cleared: an auth client must have a name, so an empty value is rejected.")),
			mcp.WithString("description", mcp.Description("New description of what the client is for. Empty clears it.")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the client may be used to obtain tokens at all. Disabling is the reversible alternative to deleting; the default client cannot be disabled.")),
			mcp.WithBoolean("trusted", mcp.Description("Trusted clients skip the consent step users are otherwise shown when authorizing. Grant this only to clients this instance owns.")),
			mcp.WithString("validGrant", mcp.Description("OAuth2 grant type, typically \"authorization_code\" or \"client_credentials\". Empty clears it. client_credentials needs security settings this tool does not set.")),
			mcp.WithString("redirectURI", mcp.Description("Absolute URL the user is sent back to after authorizing. Empty clears it.")),
			mcp.WithString("scope", mcp.Description("Space-separated OAuth2 scopes the client may request, e.g. \"profile api\". Empty clears it.")),
			mcp.WithString("validFrom", mcp.Description("Client is unusable before this time. RFC3339, e.g. 2026-08-03T09:00:00Z. Empty removes the restriction.")),
			mcp.WithString("expiresAt", mcp.Description("Client is unusable after this time. RFC3339. Empty removes the expiry.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update auth client",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_auth_client_delete",
			mcp.WithDescription(
				"Delete an auth client. Every application, integration or script authenticating with it stops "+
					"being able to obtain tokens, so establish what uses it before calling. "+
					"The delete is soft: the client is retained with its secret intact, is listed again by "+
					"passing includeDeleted to system_auth_client_lookup, and is restored by "+
					"system_auth_client_undelete. "+
					"If you only want to stop the client working while keeping it obviously present, set "+
					"enabled=false with system_auth_client_update instead. "+
					"The instance's default auth client cannot be deleted.",
			),
			mcp.WithString("authClient", mcp.Required(), mcp.Description("Auth client ID as a string, to prevent precision loss, or the client's handle.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete auth client",
		h.delete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_auth_client_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted auth client. It becomes usable again with the same ID, handle and "+
					"secret as before, so applications configured against it start working without being "+
					"reconfigured. "+
					"Find the client first with system_auth_client_lookup and includeDeleted set — a deleted "+
					"client is invisible to a plain lookup. A handle works here even though a handle search "+
					"normally skips deleted clients: this tool searches them deliberately.",
			),
			mcp.WithString("authClient", mcp.Required(), mcp.Description("Auth client ID as a string, to prevent precision loss, or the client's handle.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete auth client",
		h.undelete,
	)
}

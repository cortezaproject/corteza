package agentic

import (
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// Auth clients are the OAuth2 credentials third-party applications use against
// this instance, so this family is shaped by one rule above the usual ones: no
// tool here returns or sets a client secret. The service can do both —
// ExposeSecret hands back a working credential and RegenerateSecret mints a new
// one — and both are properly access-checked, but read permission is a lower
// bar than a credential deserves and the disclosure is irreversible: a secret
// that reaches a model's context is in the transcript and in the logs for good.
// Creation is excluded for the same reason, because the service generates the
// secret during create and returns it on the new client. See CONVENTIONS.md
// §8.6b; the gaps are recorded in TOOLS.md.
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
					"No response from this tool ever includes a client secret, and there is deliberately no "+
					"tool that reveals or regenerates one — a secret is readable only by a person through the "+
					"admin UI, so do not offer to fetch it. "+
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
		mcp.NewTool("system_auth_client_update",
			mcp.WithDescription(
				"Update an auth client's configuration: its handle, name and description, whether it is "+
					"enabled and trusted, its OAuth2 grant type, redirect URI, scope and validity window. "+
					"Omit a field to leave it unchanged; pass it empty to clear it. "+
					"This cannot read, set or regenerate the client secret — no tool can — so an update never "+
					"disturbs the credential the client is already using. "+
					"It also leaves the client's security settings alone (impersonated user, permitted, "+
					"prohibited and forced roles): those grant privileges and are set by a person in the "+
					"admin UI, which is why setting 'validGrant' to client_credentials is rejected unless an "+
					"impersonation user has already been configured there. "+
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

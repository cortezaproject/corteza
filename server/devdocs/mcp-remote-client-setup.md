# Connecting Claude to Human's MCP surface

How someone with no checkout of this repository points Claude at a Human
instance and gets the tool surface, authenticated as themselves.

There are two client surfaces and they differ in one way that decides
everything else: **who makes the HTTP request to Human**.

| Surface                                                         | Request originates from                              | Works against an on-site Human                      |
| --------------------------------------------------------------- | ---------------------------------------------------- | --------------------------------------------------- |
| claude.ai chat, desktop and mobile apps, Claude Code on the web | Anthropic's infrastructure, egress `160.79.104.0/21` | Only if the instance is reachable over public HTTPS |
| Claude Code CLI on the user's own machine                       | That machine                                         | Yes — LAN, VPN or localhost all work                |

The server-side setup below is the same for both. Only the last section differs.

## What the server does

An MCP client that gets a `401` has to find its way to a login screen without
being told anything in advance. Human answers that in two steps:

1. `POST /api/mcp` without a token returns `401` with

   ```
   WWW-Authenticate: Bearer resource_metadata="https://<host>/.well-known/oauth-protected-resource/api/mcp", scope="api profile"
   ```

2. That document names the authorization server:

   ```json
   {
     "resource": "https://<host>/api/mcp",
     "authorization_servers": ["https://<host>/auth"],
     "scopes_supported": ["api", "profile"],
     "bearer_methods_supported": ["header"]
   }
   ```

The client then reads `https://<host>/auth/.well-known/openid-configuration`,
runs an ordinary authorization-code flow with PKCE, and calls `/api/mcp` with
the resulting access token. RBAC applies to the user who logged in, exactly as
it would in the webapp.

The `resource` value is built from the incoming request, so it is correct on
whatever hostname the instance is served from, with no extra configuration. It
must match the URL the user typed into their client character for character —
that is the client's own check, not ours.

## Prerequisites

- **`AUTH_BASE_URL` must be the externally correct URL.** It is published as the
  issuer and the client will fetch it. If it says `localhost`, a remote client
  follows it to its own machine and the flow dies there.
- **TLS.** Assume any host other than `localhost` needs real HTTPS. A private
  certificate authority has not been tested with either client.
- **For the claude.ai surfaces only:** the instance must be reachable from
  `160.79.104.0/21`, and so must the authorization server, since discovery
  requests come from the same range. A WAF in front of `/auth` breaks the flow
  even when `/api/mcp` is reachable.

## Rollout order

Do these in order. Each step is checkable on its own, so a failure has one
cause rather than four candidates.

1. Deploy a build containing this handshake to the host testers will use.
2. Verify the handshake from **outside** the network, before touching any
   client — the three curl commands under
   [When it does not work](#when-it-does-not-work).
3. Create the auth client (next section) and note its ID.
4. An administrator adds the connector once, for the whole organization.
5. Each tester clicks Connect and signs into Human with their own account.

Then tell them what to expect on first use: **every tool is listed, but each is
described in one line.** That is deliberate — sending the full descriptions
costs roughly 40,000 tokens per request, so Claude fetches the detail for a
tool when it needs it. Testers should describe what they want done rather than
go looking for a tool by name.

This used to work the other way round: five tools were listed and the rest were
pulled in by searching. It broke on the claude.ai connector, which does not
re-list when told the tool set changed, so writes were advertised as loaded and
then refused — and because the five always-on tools were all reads, the whole
connector read as read-only. If a tester on an older build reports that Claude
"can only read", that is the cause, and the fix is a build with the slim
listing rather than anything in their configuration.

## One-time setup: the auth client

Human has no dynamic client registration endpoint, so Claude cannot register
itself. An administrator creates one OAuth client and its ID is entered once per
organization on the client side.

In the admin area, **System → Auth Clients → New**
(`/admin/system/auth-clients/new`):

| Field        | Value                                                                       | Why                                                                                                                                                                                                       |
| ------------ | --------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Valid grant  | `authorization_code`                                                        | Claude requires a user consent step. A pure `client_credentials` machine token is refused by the hosted surfaces regardless of what the server would issue.                                               |
| Scope        | `api profile`                                                               | `api` is what `HttpTokenValidator` checks on `/api/mcp`; without it every call is rejected as out of scope.                                                                                               |
| Redirect URI | `https://claude.ai/api/mcp/auth_callback http://localhost http://127.0.0.1` | First is the hosted surfaces. The other two are Claude Code, which is a native client and redirects to a loopback port that changes every run — the values are matched by prefix, so the port is covered. |
| Enabled      | yes                                                                         |                                                                                                                                                                                                           |
| Trusted      | yes, for a demo instance                                                    | Skips the per-user consent screen. Leave it off if you want testers to see the consent step.                                                                                                              |

Note the client ID, and the secret if one is generated. Both are entered on the
client side; the secret is optional there and only needed if this client
requires confidential authentication.

## What testers point at

Hand out the bare endpoint:

```
https://<host>/api/mcp
```

That surface is the configuring and usage layers — modules, pages, roles, TAQs,
workflows, records, reports, automation runs. Development-layer tools are not on
it at all; they live in the repo-side MCP under `dev/mcp` and are not reachable
over HTTP.

Two narrowing options exist and neither is wanted for general testing, but they
are worth knowing about:

- `/api/mcp/configuring` or `/api/mcp/usage` list only that layer. Narrower than
  what a services person needs, since their work spans both.
- `?maxRisk=read` or `?maxRisk=write` caps what can be called, refused at
  dispatch rather than merely hidden. Useful for pointing someone at production.

Metadata is served for these variants too, so a client can use them without
anything extra.

## Surface A: claude.ai connector

Nicest onboarding, needs the public HTTPS instance. On Team and Enterprise
plans only an administrator can add it, and once added everyone in the
organization has it.

1. Go to [claude.ai/customize/connectors](https://claude.ai/customize/connectors)
   (**Customize → Connectors**) and choose **Add custom connector**.
2. URL: `https://<host>/api/mcp`.
3. Under advanced settings, enter the OAuth client ID, and the secret if the
   client has one.
4. **Connect** — a Human login opens in the browser.

The connector then appears in claude.ai chat, in the desktop and mobile apps,
and in Claude Code — both web sessions and the CLI under `/mcp` — for anyone
signed in with that account.

**For non-developers, point them at claude.ai chat, not Claude Code.**
Configuring a system and building reports is not a coding task, and Claude Code
on the web additionally requires a connected GitHub account and a repository to
clone, which someone who never touches code will not have. Chat needs neither.

## Surface B: Claude Code CLI

Works against an instance that is only reachable from the tester's own network,
because the CLI and its OAuth browser redirect both run on their machine. No
checkout and no `.mcp.json` are involved.

```bash
claude mcp add --transport http --scope user human https://<host>/api/mcp
```

Then inside Claude Code:

```
/mcp
```

and follow the browser login. `claude mcp login human` does the same thing from
the shell. `--scope user` writes to `~/.claude.json`, so the server is available
in every directory rather than tied to a project.

If the CLI reports that the server does not support dynamic client
registration, pass the client from the previous section explicitly:

```bash
claude mcp add --transport http --scope user \
  --client-id <client-id> --client-secret \
  human https://<host>/api/mcp
```

`--client-secret` prompts with masked input, and a paste into that prompt can
arrive mangled with no complaint at all — the symptom is a browser login that
succeeds followed by `Client authentication failed` in the terminal, because the
stored secret is not the one Human holds. Set it in the environment instead,
which is also the only way to script the command:

```bash
MCP_CLIENT_SECRET='<secret>' claude mcp add --transport http --scope user \
  --client-id <client-id> --client-secret \
  human https://<host>/api/mcp
```

## Testing it without a public host

Most of this can be exercised against the local dev server, because Claude Code
is a native client: it runs on the tester's machine, redirects to a loopback
port, and never needs to be reachable from outside. Point it at
`http://localhost:1043/api/mcp` and the whole handshake — discovery, the
metadata document, authorization code with PKCE, the redirect prefix match,
scope checking, tool listing and tool calls — runs for real.

It needs an auth client the dev server does not have after `bootstrap.sh`: that
one uses `client_credentials`, which the hosted surfaces refuse and which skips
the browser flow entirely. Create a second one per the table above, with the
loopback redirect prefixes.

What a local run still cannot tell you:

- whether Anthropic's egress reaches the instance, since nothing external is
  involved
- whether the claude.ai connector dialog accepts the client, which only the
  hosted flow exercises
- whether the deployment's proxy sets `X-Forwarded-Proto` and
  `X-Forwarded-Host` correctly. Approximate it by sending those headers by
  hand and checking what `resource` comes back as.

To close the first two without deploying, put a tunnel in front of the dev
server:

```bash
cloudflared tunnel --url http://localhost:1043
```

Then set `AUTH_BASE_URL=https://<tunnel-host>/auth` and restart, or the metadata
sends Claude to `localhost` and the login dies on the tester's machine. The
hostname changes on every tunnel restart, and the connector has to be re-added
each time.

## When it does not work

Run these against the instance in order. The first one that misbehaves is the
problem.

```bash
# 1. Does the 401 say where to authenticate?
curl -si -X POST https://<host>/api/mcp | grep -i www-authenticate

# 2. Does the metadata document exist, and does `resource` match the URL
#    the tester typed, exactly?
curl -s https://<host>/.well-known/oauth-protected-resource/api/mcp

# 3. Is the authorization server it names actually there?
curl -s https://<host>/auth/.well-known/openid-configuration
```

Known causes, in the order they bite:

- **Metadata returns an empty `200` instead of JSON.** The MCP surface is
  disabled, so the route is not registered and the request fell through to the
  catch-all, which answers every unrouted path with an empty `200`. That is also
  why the `WWW-Authenticate` pointer matters: a client probing well-known paths
  blind finds nothing here that fails cleanly.
- **`resource` does not match what the tester entered.** Usually a proxy that
  rewrites the host without setting `X-Forwarded-Host`, or a tester using a
  different hostname than the one the connector was added with.
- **`authorization_servers` points at `localhost`.** `AUTH_BASE_URL` is wrong
  for external use.
- **Token accepted but every tool call is rejected as out of scope.** The auth
  client's scope is missing `api`.
- **Claude Code's callback is refused.** The auth client's redirect URI list is
  missing the loopback prefixes; the port varies per run and cannot be
  registered individually.
- **The browser reports success and the terminal reports
  `Client authentication failed`.** The code was issued and the token exchange
  was refused: the secret stored on the tester's machine is not the one Human
  holds. Almost always a mangled paste — re-add with `MCP_CLIENT_SECRET` set.
  Human's own error distinguishes the two cases, so this is checkable: a wrong
  secret answers `invalid_client`, a wrong code answers `invalid_grant`.

## What is deliberately not here

- **No dynamic client registration.** Adding an RFC 7591 `registration_endpoint`
  would remove the client-ID handout entirely, at the cost of an open
  registration surface that needs gating and an `auth_clients` table that grows
  a row per fresh connection. Advertising a Client ID Metadata Document is the
  cheaper alternative and needs the token endpoint to accept a public client.
- **No resource indicators (RFC 8707).** Claude may send a `resource` parameter
  on the authorization request; the OAuth server ignores unknown parameters, so
  it neither helps nor breaks anything today.

# Agent dev toolkit

Deterministic scripts that let an agent (or you) talk to the **local** dev
server API without re-deriving auth every session. Local-only by design:
every script refuses non-localhost `HUMAN_API` hosts, and nothing here is
ever provisioned automatically — `server/provision/` must never reference
`seed/dev-agent.yaml`.

## Setup (once, idempotent)

```sh
dev/agent/bootstrap.sh
```

Requires the dev server running (`cd server && make watch`, API on
`http://localhost:1043/api` — override with `HUMAN_API`). Creates:

- user `agent-dev` / `agent@local.dev`, member of `super-admin` (no password,
  token-only; keeps agent actions distinguishable in action logs)
- auth client `dev-agent` (`client_credentials` grant, impersonates
  `agent-dev`), managed via REST; its server-generated secret is cached in
  `.state/secret` (gitignored)

Self-healing: re-running bootstrap repairs a misconfigured client (PUT with
the full desired config) and re-caches the secret from the expose endpoint.

Requires: `curl`, `python3` (no jq dependency), a built server binary in
`server/build/` (gin's `make watch` provides one).

## Daily use

| script | purpose |
| --- | --- |
| `smoke.sh` | server up? version? token OK? who am I? authed API works? |
| `token.sh` | print a valid bearer token (cached ~2h; oauth2 flow, CLI-jwt fallback) |
| `api.sh METHOD PATH [curl args…]` | authenticated request, jq-pretty output |
| `seed.sh [--force] [fixture…]` | import `dev/fixtures/` (skips seeded fixtures unless `--force`) |
| `cleanup.sh` | delete all `agent-*` compose namespaces |
| `pagebuild.py SLUG SPEC.json` | build/refresh charts + pages via REST (spec format in its header) |
| `mcp.py tools\|schema\|call` | call Human's own MCP server (`/api/mcp`, 24 tools: compose CRUD, TAQ/workflow exec, discovery) |

For interactive Claude Code sessions, `.mcp.json` registers the `human` MCP
server; it needs `HUMAN_MCP_TOKEN` exported before starting Claude Code:
`export HUMAN_MCP_TOKEN=$(dev/agent/token.sh)` (2h lifetime — re-export when
stale, or use `mcp.py` which self-authenticates per call).

```sh
dev/agent/smoke.sh
dev/agent/api.sh GET '/system/users/?limit=5'
dev/agent/api.sh POST /compose/namespace/ -d '{"name":"Agent Test","slug":"agent-test"}'
```

## API gotchas (learned the hard way)

- **Trailing slash matters** on collection endpoints: `/system/users/` works,
  `/system/users` is a 404 (chi routing).
- **Errors come back as HTTP 200** with `{"error":{"message":…}}` — `api.sh`
  detects this and exits non-zero.
- Auth clients must leave `validFrom`/`expiresAt` unset —
  `AuthClient.Verify()` has inverted comparisons and treats set values as
  expired/not-yet-valid.
- Tokens need `scope=profile api` — the API middleware requires the `api`
  scope.
- Envoy YAML import cannot configure auth clients fully (`validGrant` is
  silently dropped, `security.impersonateUser` refs don't resolve) — manage
  auth clients via REST instead, as bootstrap does.
- The server CLI (`auth jwt`, cobra `cmd.Println`) prints to **stderr**;
  capture with `2>&1`.
- **Compose updates are POST, not PUT** (`POST …/module/{id}` updates it);
  system auth clients use PUT. Check the `handlers/*.go` route registrations.
- CSV record sources import only when the fixture **directory** is passed to
  `import` (single-file decode never registers the CSV providers — nil-panic).
- CLI `import` bypasses the running server's DAL registry — record queries
  then fail with "model does not exist" until models re-register (seed.sh
  does a no-op module POST per module; a server restart also works).
- **Envoy YAML cannot build working pages**: block option refs
  (module/chart handles) are never resolved to IDs — the encoder's ref map
  is keyed by resource kind while the SetValue write-back expects
  `Blocks.N.Options.…` paths (`server/compose/envoy/store_decode.go`
  `toEnvoyRefs`). Pages/charts go through `pagebuild.py` instead.

## Conventions for agent-created data

- Every resource an agent creates for testing uses an `agent-` handle/slug
  prefix (namespaces, projects, modules, TAQs, users, …). Anything so
  prefixed is disposable; never touch unprefixed data.
- Good agent-built systems get harvested into fixtures:
  `server_cli export compose-namespace <handle>` (see `common.sh` for
  `server_cli`).

## State

`.state/` (gitignored): client secret, cached token. Delete it any time;
`bootstrap.sh` regenerates.

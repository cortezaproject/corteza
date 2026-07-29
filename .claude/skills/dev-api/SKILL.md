---
name: dev-api
description: Make authenticated API calls against the local dev server using the dev/agent toolkit — token minting, curl wrapper, smoke check. Use whenever a task needs to read/write real server data, verify an implementation end-to-end, or inspect API behavior.
---

# /dev-api

Deterministic scripts in `dev/agent/` handle all auth mechanics. Never
hand-roll the oauth dance or guess ports — use the toolkit.

## Workflow

1. **Always start with** `dev/agent/smoke.sh` — answers: server up (which
   version), token mintable, identity, authed API access. If it fails, run
   `dev/agent/bootstrap.sh` (idempotent) and retry; if the server is down,
   tell the user to run `cd server && make watch`.
2. Make calls with `dev/agent/api.sh METHOD PATH [extra curl args…]`:

   ```sh
   dev/agent/api.sh GET '/system/users/?limit=5'
   dev/agent/api.sh POST /compose/namespace/ -d '{"name":"X","slug":"agent-x"}'
   ```

   Paths are relative to `http://localhost:1043/api`. Output is pretty JSON;
   API errors exit non-zero with the message on stderr.
3. Raw token when needed (websockets, custom curl): `dev/agent/token.sh`.

## Rules

- **`agent-` prefix**: every handle/slug you create for testing must start
  with `agent-`. Prefixed data is disposable (`dev/agent/cleanup.sh` wipes
  it); unprefixed data belongs to the human — never modify it.
- Local-only: the toolkit refuses non-localhost hosts. Never point it (or
  copy its patterns) at shared/production instances.

## Gotchas (each cost a debug cycle — don't re-learn them)

- Collection endpoints need a **trailing slash**: `/system/users/` — without
  it you get a 404.
- Errors often return **HTTP 200** with `{"error":{"message":…}}`; `api.sh`
  already detects this.
- **Compose updates are POST, not PUT** (e.g. `POST …/module/{id}` updates);
  some system endpoints (auth clients) use PUT. Check
  `server/*/rest/handlers/*.go` route registrations when unsure.
- Tokens need scope `profile api`.
- Server CLI (`server_cli` in `dev/agent/common.sh`) prints JWTs and command
  output to **stderr**; capture with `2>&1`.
- No `jq` on this machine — pipe JSON to `python3` (see `json_get` in
  `dev/agent/common.sh`).

Full reference: `dev/agent/README.md`.

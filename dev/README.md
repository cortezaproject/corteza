# Claude Code on @Human

The path from a fresh clone to an agent that can build, test, drive and commit
this repo. Follow it top to bottom; every step is idempotent, so re-running one
is always safe.

Reference material is linked rather than restated: `dev/agent/README.md`
(toolkit + gotchas), `dev/mcp/SPEC.md` (why the developer MCP exists),
`CLAUDE.md` (conventions and the commit rule), `server/.env.min.example` (every
server option that matters, with what breaks if you omit it).

## 0. Prerequisites

| need          | version | note                                      |
| ------------- | ------- | ----------------------------------------- |
| Go            | 1.24+   | deps are vendored — no `go mod download`  |
| Node          | 22+     | `engines` in the root `package.json`      |
| pnpm          | 10+     | npm and yarn are refused                  |
| PostgreSQL    | 14+     | or SQLite, see step 2                     |
| python3, curl | any     | the agent toolkit uses them instead of jq |
| jq            | any     | `make setup` installs it if it can        |

Linux, macOS or WSL2. No submodules, nothing to `git clone --recursive`.

## 1. Install

```sh
make setup
```

pnpm install, the playwright chromium browser, `client/web/unify/.env.e2e` from
its example, and jq if it is missing.

## 2. Create the database

Postgres is what the team runs. Create the database; the server creates its own
tables on first boot.

```sh
sudo -u postgres createuser -P human      # password: human
sudo -u postgres createdb -O human human
```

SQLite works too if you want zero infrastructure — see the commented `DB_DSN` in
`server/.env.min.example`. With **no** `DB_DSN` at all the server still boots,
on an in-memory database that vanishes on restart: fine for a smoke test, no
good for development.

## 3. `server/.env`

```sh
cd server && cp .env.min.example .env
```

Read it — it is short and it explains each value. Three lines decide whether the
rest of this guide works:

- **`ENVIRONMENT=dev`** — anything else is treated as production, and the server
  then refuses to provision super users. Without it `dev/agent/bootstrap.sh`
  cannot create its user and the agent toolkit can never authenticate.
- **`DOMAIN` must agree with `HTTP_ADDR`** — it goes into generated links,
  cookies and OAuth redirects, so a mismatch gives you logins that appear to
  work and then bounce.
- **`HTTP_API_BASE_URL=/api`** — `dev/agent/*` and the webapp both assume it.

Add two more that the minimal file leaves out:

```sh
# your own login: the password is the email address itself
AUTH_PROVISION_SUPER_USER=you@example.tld

# every fresh browser context does an oauth roundtrip; the default 60/min
# per-IP limit fails a playwright run in a way that looks like a crashed server
AUTH_REQUEST_RATE_LIMIT=0
```

The super user is created on boot, gets every bypass role, and is skipped if the
email or handle already exists.

## 4. `client/web/unify/public/config.js`

**`make setup` does not create this one, and the webapp cannot reach the API
without it.**

```sh
cd client/web/unify && cp public/config.example.js public/config.js
```

Then point it at the local server:

```js
window.HumanAPI = 'http://localhost:1043/api'
```

`public/config.js` is gitignored, as is `.env.e2e`. `client/web/unify/.env` is
tracked and needs no edit.

## 5. Start the stack

Two long-running processes:

```sh
cd server && make watch            # API on :1043, gin live-reload
cd client/web/unify && pnpm dev    # webapp on :5173
```

`make dev-all` from the repo root starts both under one trap if you prefer a
single terminal.

First boot takes a while — a Go build plus schema creation. Then open
http://localhost:5173 and log in as the super user from step 3 (password = the
email address).

Two things about `make watch` worth knowing before they confuse you: gin builds
on the **first request** to its proxy, so `curl -s localhost:1043/api/` before
believing any check; and it never respawns on its own once the webapp bypasses
its proxy, so a Go change needs a touched `.go` file and ~15s.

## 6. Provision the agent identities

With the server up:

```sh
dev/agent/bootstrap.sh   # idempotent, self-healing
dev/agent/smoke.sh       # server up? token? who am I? authed call?
```

`bootstrap.sh` creates the `dev_agent` client-credentials auth client and three
identities. They are not interchangeable:

| identity             | credential                           | reached by                                                                                                                  |
| -------------------- | ------------------------------------ | --------------------------------------------------------------------------------------------------------------------------- |
| `agent@local.dev`    | client secret + `.state/ui-password` | everything — `api.sh`, `token.sh`, `mcp.py`, `seed.sh`, the `human` MCP, `verify-ui.mjs`, `drive.mjs`, `dev_ui_verify`, e2e |
| `agent-ro@local.dev` | password in `.state/ro-password`     | permission-gate checks, and nothing else                                                                                    |

One identity covers API and browser alike: the auth client impersonates
`agent@local.dev` for token calls, and the same user logs into the webapp for
browser checks. The `client_credentials` flow loads the impersonated user by ID
and does not care that it also holds a password.

`agent-ro` is the one that differs in kind. Its role allows a deliberate subset
(users and roles, read only) and it belongs to no user group, so everything not
explicitly allowed resolves to deny. That is the only way to see what a user
_without_ permission gets: the admin sidebar drops the entries it cannot reach,
and an editor opened by deep link renders with every field disabled.

Both password files live under `dev/agent/.state/` and are gitignored.

**`smoke.sh` is the first thing to run in any session that touches the server.**
It distinguishes "server down" from "token stale" from "your call is wrong",
which otherwise look identical.

## 7. Finish `.env.e2e`

`make setup` copied the example; fill it from the password step 6 just set:

```sh
E2E_BASE_URL=http://localhost:5173
E2E_USER=agent@local.dev
E2E_PASS=<contents of dev/agent/.state/ui-password>
```

## 8. Launch Claude Code

The `human` MCP server needs a token exported **before** Claude Code starts:

```sh
export HUMAN_MCP_TOKEN=$(dev/agent/token.sh)   # ~2h lifetime
claude
```

Approve both MCP servers once in the session (`/mcp`); they are already listed
in `.claude/settings.local.json`.

---

# Working with it

## The two MCP servers

| server      | transport | needs a running Human? | what it is                             |
| ----------- | --------- | ---------------------- | -------------------------------------- |
| `human-dev` | stdio     | no                     | the **developer** layer — this repo    |
| `human`     | http      | yes                    | the **configurator** — a running Human |

### `human-dev` — repo tools

Built on demand by `dev/mcp/run.sh`: it rebuilds when a source file is newer
than the binary, so restarting the MCP client is all a tool edit needs.

| tool                                            | use it for                                          |
| ----------------------------------------------- | --------------------------------------------------- |
| `dev_test_run`                                  | the touched package's suite — returns failures only |
| `dev_format_run`                                | gofmt/prettier on named files                       |
| `dev_commit_create`                             | commit with the convention enforced                 |
| `dev_branch_status`                             | branch, base, ahead/behind, working tree            |
| `dev_intent_governing` / `_check` / `_affected` | intent docs, drift, affected e2e specs              |
| `dev_server_status` / `dev_server_logs`         | is the dev server up, how stale, filtered log tail  |
| `dev_ui_verify`                                 | render-check a webapp path in a real browser        |
| `dev_fixture_cleanup`                           | remove this session's seeded data                   |

`dev/mcp/SPEC.md` also describes `dev_diff_survey`, `dev_lint_run`,
`dev_fixture_seed`, `dev_scratch_build` and `dev_e2e_run` — those are **planned,
not registered**. Use the scripts for those jobs today.

Every repo-acting tool takes an optional `worktree`. It is stateless by design:
a session sitting in the primary checkout while its work lives in a worktree
must pass it on **every** call, or the tests pass on the wrong tree and the
commit lands on the wrong branch, both reporting success.

### `human` — configure a running Human

Compose/system/automation CRUD against the local server: namespaces, modules,
pages, records, TAQs, workflows. A stale `HUMAN_MCP_TOKEN` means the server
simply never connects — re-export and restart Claude Code, or use
`dev/agent/mcp.py tools|schema|call`, which re-authenticates per call and is the
better bet in a long session.

The tool surface is behind progressive disclosure: `human_tool_search` finds a
tool, `human_tool_load` fetches its schema.

## Skills — what to invoke when

| situation                                    | skill           |
| -------------------------------------------- | --------------- |
| any piece of work bigger than a one-liner    | `/dev-task`     |
| one change, plan already agreed              | `/dev-change`   |
| read/write real server data                  | `/dev-api`      |
| need known test data, or a clean slate       | `/dev-seed`     |
| build a system _inside_ Human (not its code) | `/sys-design`   |
| several independent issues at once           | `/orchestrate`  |
| code under a locked contract or WIP zone     | `/intent-task`  |
| drift audit of an intent-covered area        | `/intent-audit` |

Two rules about this list: the **intent system is opt-in** — never touch a
`*.intent.md` or run intent tooling unless `/intent-task` or `/intent-audit` was
invoked; and `/orchestrate` is the _only_ thing in this repo that authorises
spawning subagents.

`/dev-task` and `/dev-change` reply in the standard block — sections, confidence
grade, open questions asked as an interview. The format is
`.claude/reporting.md`.

## Testing

| layer                                    | command                                         |
| ---------------------------------------- | ----------------------------------------------- |
| everything                               | `make test` (lib, client, test-compile, server) |
| unify webapp                             | `cd client/web/unify && npx vitest run`         |
| lib/js                                   | mocha, not vitest — `cd lib/js && pnpm test`    |
| Go, one package                          | `cd server && go test ./compose/service/...`    |
| every package still _compiles_ its tests | `make test-compile`                             |
| e2e (playwright)                         | `make e2e` — needs the stack up and `.env.e2e`  |

Prefer `dev_test_run` for the inner loop: it dispatches by target and returns
only the failures, with `file:line`.

Browser checks come in three sizes — `dev_ui_verify` / `dev/agent/verify-ui.mjs`
for "did this path render clean"; `dev/agent/drive.mjs` for anything needing a
second step (click, assert where you landed, count rows); playwright specs under
`client/web/unify/e2e` for the real suite. `drive.mjs` records native dialogs
rather than dismissing them — playwright's default dismiss turns an
unsaved-changes confirm into a phantom "the button does nothing" — and fails a
check on console errors, page errors and HTTP 500s.

## Working in parallel

One task, one checkout, one server, one webapp, one database:

```sh
dev/agent/worktree.sh new alpha    # checkout + DB clone + ports + env files
dev/agent/worktree.sh up alpha     # first up is minutes; after that, seconds
dev/agent/worktree.sh land alpha   # rebase onto main, ff-merge, remove
dev/agent/worktree.sh gc           # find leftovers; --reap removes them
```

`new` writes that slot's `server/.env`, `public/config.js` and `.env.e2e` for
you — steps 3, 4 and 7 are a one-time cost for the primary checkout only. Ports
derive from the slot (API `1043+slot*100`, vite `5173+slot`), the primary is
slot 0, and nothing cleans itself up when a terminal closes: only `rm` and
`land` remove anything.

## The gotchas that cost the most time

- **Trailing slash matters**: `/system/users/` works, `/system/users` is a 404.
- **Errors arrive as HTTP 200** with `{"error":{"message":…}}`; `api.sh` detects
  this and exits non-zero.
- **Compose updates are POST, not PUT** — a PUT is a bare 405.
- **Vite never picks up edits under `lib/`** — it serves its startup transform
  forever, so a UI check reports old behaviour for a fix that is really there.
  Confirm with `curl -s 'localhost:5173/@fs<abs-path>' | head`.
- **An unknown webapp path renders the home page**, not an error — so a wrong
  path can come back as a pass with a screenshot of a healthy page nobody asked
  about.
- **`stale: false` can still mean a stale binary** when a rebuild started before
  your edit and finished after it. `touch` the file and re-check.
- **Never start, stop or restart the dev servers yourself** — report the need
  and let the human relaunch.
- **Clean up what you created and only that.** `cleanup.sh` deletes what this
  session's ledger (`dev/agent/.state/created.jsonl`) holds. Data the session
  did not create is off-limits whatever it is called.
- **Names are snake_case** everywhere an identifier is parsed — a hyphen lexes
  as subtraction.

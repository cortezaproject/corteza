# Claude Code on @Human

The path from a fresh clone to an agent that can build, test, drive and commit
this repo. Follow it top to bottom; every step is idempotent, so re-running one
is always safe.

Reference material is linked rather than restated: `dev/agent/README.md`
(toolkit + gotchas), `dev/mcp/SPEC.md` (why the developer MCP exists),
`AGENTS.md` (conventions and the commit rule), `server/.env.min.example` (every
server option that matters, with what breaks if you omit it).

## 0. Prerequisites

| need          | version | note                                      |
| ------------- | ------- | ----------------------------------------- |
| Go            | 1.24+   | deps are vendored — no `go mod download`  |
| Node          | 22+     | `engines` in the root `package.json`      |
| pnpm          | 10+     | npm and yarn are refused                  |
| PostgreSQL    | 14+     | or SQLite, see step 1                     |
| python3, curl | any     | the agent toolkit uses them instead of jq |
| jq            | any     | `make setup` installs it if it can        |

Linux, macOS or WSL2. No submodules, nothing to `git clone --recursive`.
`make setup` checks every row and says which are missing.

## 1. Set it up

```sh
make setup
```

Idempotent, and it never edits a file you already have — anything it finds, it
checks and leaves alone. What it does:

| does                                | detail                                                                                |
| ----------------------------------- | ------------------------------------------------------------------------------------- |
| dependencies                        | `pnpm install`, the playwright chromium browser, jq                                   |
| `server/.env`                       | from `.env.min.example`, with a pinned `AUTH_JWT_SECRET`                              |
| `client/web/unify/public/config.js` | `window.HumanAPI` pointed at this checkout's own API port                             |
| `client/web/unify/.env.e2e`         | base URL, and `agent@local.dev` as the test login                                     |
| the database                        | creates the role and database in `DB_DSN`, or prints the two commands to run yourself |

All three files are gitignored. `server/.env.min.example` is short and explains
every value in it; the ones `make setup` will not let you get wrong are
`ENVIRONMENT=dev` (anything else is production, and the server then refuses to
provision super users), `DOMAIN` agreeing with `HTTP_ADDR` (a mismatch gives you
logins that appear to work and then bounce), `HTTP_API_BASE_URL=/api` (both
`dev/agent/*` and the webapp assume it), and `AUTH_REQUEST_RATE_LIMIT=0`
(a playwright run trips the default per-IP limit and fails on blank 429 pages).

With **no** `DB_DSN` at all the server still boots, on an in-memory database
that vanishes on restart: fine for a smoke test, no good for development.
SQLite works too — see the commented `DB_DSN` in `server/.env.min.example`.

### If 1043 or 5173 is taken

Two knobs, and everything else follows them. Set both **before** `make setup`:

```sh
cp server/.env.min.example server/.env
#   HTTP_ADDR=:8080  and  DOMAIN=localhost:8080   — they must agree
echo 'VITE_PORT=3000' > client/web/unify/.env.local
make setup
```

`config.js`, `.env.e2e`, the `human-local` MCP URL, playwright, `dev_ui_verify`
and every worktree slot are derived from those two — nothing else names a port.
`.env.local` is gitignored and is what vite reads; `strictPort` is on, so a
clash fails at startup rather than drifting onto the next port while the rest of
the toolkit keeps naming this one.

Order matters, because `make setup` writes `config.js` and `.env.e2e` from the
resolved values and never rewrites a file you already have. Change the ports
afterwards and `make doctor` flags both: delete those two files and re-run it.

`DB_DSN` is yours entirely — user, password, host, port and name are carried
into every worktree's own database. The role in it usually cannot `CREATE
DATABASE`, though, and cloning a worktree needs one that can:

```sh
export PGSUPERUSER=postgres PGSUPERPASS=postgres
```

```sh
make doctor
```

is every check `make setup` makes, writing nothing, and non-zero when something
is off. Run it when the stack behaves as though it is talking to the wrong
thing.

## 2. Start the stack, and make yourself a login

Two long-running processes:

```sh
cd server && make watch            # API, rebuilt and restarted on every edit
cd client/web/unify && pnpm dev    # webapp
```

`make dev-all` from the repo root starts both under one trap if you prefer a
single terminal. First boot takes a while — a Go build plus schema creation.

Then open the webapp and **sign up**. The first non-system user in the database
is auto-promoted to super-administrator (`system/service/auth.go`, `autoPromote`),
and step 3 is about to create `agent@local.dev`, which would take that slot. For
a login made later:

```sh
cd server && ./build/dev-bin --env-file .env roles useradd super-admin you@example.tld
```

`make watch` runs `server/cmd/devwatch`, which rebuilds and restarts the server
on any `.go` write. There is no proxy and no second port: nothing has to make a
request to provoke a build, and an edit that lands _during_ one is built by the
cycle after it. A build takes ~15s, and the previous server keeps serving until
the new binary is ready — so the check that matters is whether the running
process started after your edit, which `dev_server_status` reports on.

## 3. Provision the agent identities

With the server up:

```sh
make setup-agent
```

`dev/agent/bootstrap.sh` (idempotent, self-healing), then `E2E_PASS` written
into `.env.e2e`, then `dev/agent/smoke.sh`. It creates the `dev_agent`
client-credentials auth client and two identities, which are not
interchangeable:

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

## 4. Launch Claude Code

```sh
make claude
make claude MCP=1                               # also attach human-local (see below)
make claude-yolo                                # --dangerously-skip-permissions
make claude -- --model opus --verbose           # any other flags, after `--`
```

`.mcp.json` registers `human-dev` only. `human-local` — Human's own `/api/mcp`,
some 150 configurator tools — is opt-in with `MCP=1`: the launcher mints the
token it needs, which a config file cannot do for itself, and passes
`dev/human-local.mcp.json` with `--mcp-config`. Without it a session reaches
the dev server through `dev/agent/api.sh` and `mcp.py`, which authenticate per
call and record what they create in the cleanup ledger. A flag containing `=`
has to go through `ARGS="…"` instead, since make reads it as a variable
assignment.

---

# Working with it

## The MCP servers, and which Human they reach

| tools                  | from                       | reaches                                           |
| ---------------------- | -------------------------- | ------------------------------------------------- |
| `mcp__human-local__*`  | `make claude MCP=1`, http  | **this checkout's** `/api/mcp` — the configurator |
| `mcp__human-dev__*`    | `.mcp.json`, stdio         | **this repo** — the developer layer               |
| any other `*_Human__*` | an account-level connector | **a remote instance**, NOT this checkout          |
| `dev/agent/mcp.py`     | the shell                  | this checkout's `/api/mcp`, re-auth'd each call   |

**A connector is not the dev server.** A Human instance attached to your own
account shows up with tool names that read as "the Human MCP" and carries the
same configurator surface — compose CRUD, TAQ and workflow exec, users and
roles — against somebody's live data. Nothing it writes is in
`dev/agent/.state/created.jsonl`, so `cleanup.sh` cannot undo it. Local work
goes through `human-local`, `mcp.py` or `api.sh`; reach for a connector only
when the remote instance is the point, and say so.

`human-local` is why `make claude MCP=1` exists. `dev/human-local.mcp.json`
can interpolate an environment variable into its Authorization header but
cannot run a command to produce one, and Human offers no dynamic client
registration (`registration_endpoint` is absent from both discovery documents),
so a token has to be in the environment before Claude Code starts. The launcher
mints it; `AUTH_OAUTH2_ACCESS_TOKEN_LIFETIME=720h` keeps it alive past the end
of a session, because nothing can refresh it in flight.

### `human-dev` — repo tools

Built on demand by `dev/mcp/run.sh`: it rebuilds when a source file is newer
than the binary, so restarting the MCP client is all a tool edit needs.

| tool                                            | use it for                                                                                                                       |
| ----------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `dev_test_run`                                  | the touched package's suite — returns failures only; `related` runs just the tests covering named source files, for the baseline |
| `dev_format_run`                                | gofmt/prettier on named files                                                                                                    |
| `dev_commit_create`                             | commit with the convention enforced                                                                                              |
| `dev_branch_status`                             | branch, base, ahead/behind, working tree                                                                                         |
| `dev_intent_governing` / `_check` / `_affected` | intent docs, drift, affected e2e specs                                                                                           |
| `dev_server_status` / `dev_server_logs`         | up, how stale (`wait` blocks), filtered log tail                                                                                 |
| `dev_ui_verify`                                 | render-check a webapp path in a real browser                                                                                     |
| `dev_fixture_cleanup`                           | remove this session's seeded data                                                                                                |

From the shell, `related` is `npx vitest related --run <files>` inside the
workspace (`cd lib/vue && npx vitest related --run src/components/input/CInputModule.vue`),
and the whole suite is `npx vitest run` there; `lib/js` runs mocha.

`dev/mcp/SPEC.md` also describes `dev_diff_survey`, `dev_lint_run`,
`dev_fixture_seed`, `dev_scratch_build` and `dev_e2e_run` — those are **planned,
not registered**. Use the scripts for those jobs today.

Every repo-acting tool takes an optional `worktree`. It is stateless by design:
a session sitting in the primary checkout while its work lives in a worktree
must pass it on **every** call, or the tests pass on the wrong tree and the
commit lands on the wrong branch, both reporting success.

### `human-local` — configure a running Human

Compose/system/automation CRUD against the local server: namespaces, modules,
pages, records, TAQs, workflows. Attached only by `make claude MCP=1`; a stale
token means the server simply never connects — relaunch, or use
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

`new` writes that slot's `server/.env`, `public/config.js`, `.env.local` and
`.env.e2e` for you, so a worktree needs no `make setup`. Ports derive from the
slot, counted off whatever the primary serves (API `+slot*100`, vite `+slot`);
the primary is slot 0, and nothing cleans itself up when a terminal closes:
only `rm` and `land` remove anything.

Every script reads those files rather than a literal, so a call made inside a
worktree reaches that worktree — `dev/agent/stack.sh` is the one resolver, and
printing it is the fastest way to see which stack a checkout is pointed at:

```sh
dev/agent/stack.sh    # HUMAN_API · HUMAN_BASE · HUMAN_AUTH · HUMAN_WEBAPP
```

## The gotchas that cost the most time

- **Trailing slash matters**: `/system/users/` works, `/system/users` is a 404.
- **Errors arrive as HTTP 200** with `{"error":{"message":…}}`; `api.sh` detects
  this and exits non-zero.
- **Compose updates are POST, not PUT** — a PUT is a bare 405.
- **Vite never picks up edits under `lib/`** — it serves its startup transform
  forever, so a UI check reports old behaviour for a fix that is really there.
  Confirm with `curl -s "$(dev/agent/stack.sh | sed -n 's/^HUMAN_WEBAPP=//p')/@fs<abs-path>" | head`.
- **An unknown webapp path renders the home page**, not an error — so a wrong
  path can come back as a pass with a screenshot of a healthy page nobody asked
  about.
- **`stale: false` can still mean a stale binary** when a rebuild started before
  your edit and finished after it. `touch` the file and re-check.
- **Never start, stop or restart the dev servers yourself** — report the need
  and let the human relaunch.
- **A checkout answers for its own stack.** `dev/agent/stack.sh` resolves the
  API and webapp from this checkout's `server/.env` and `.env.e2e`; exporting
  `HUMAN_API` still overrides it. Nothing needs the port typed in any more.
- **`.state/` is one directory for every checkout.** `bootstrap.sh` caches the
  secret of the auth client it made in whatever database it reached, so running
  it against a stack on a database that was not cloned from the primary leaves
  the primary minting tokens through the CLI fallback and nothing says so.
- **Clean up what you created and only that.** `cleanup.sh` deletes what this
  session's ledger (`dev/agent/.state/created.jsonl`) holds. Data the session
  did not create is off-limits whatever it is called.
- **Names are snake_case** everywhere an identifier is parsed — a hyphen lexes
  as subtraction.

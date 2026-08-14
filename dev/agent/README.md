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

| script                                       | purpose                                                                                                                                                                                      |
| -------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `smoke.sh`                                   | server up? version? token OK? who am I? authed API works?                                                                                                                                    |
| `token.sh`                                   | print a valid bearer token (cached ~2h; oauth2 flow, CLI-jwt fallback)                                                                                                                       |
| `api.sh METHOD PATH [curl args…]`            | authenticated request, jq-pretty output                                                                                                                                                      |
| `seed.sh [--force] [fixture…]`               | import `dev/fixtures/` (skips seeded fixtures unless `--force`)                                                                                                                              |
| `cleanup.sh [--all\|--session ID] [--purge]` | delete namespaces recorded in `.state/created.jsonl` for this session (`--all` for every session); `--purge` also hard-deletes soft-deleted corpses (dev-only server CLI, NOT ledger-scoped) |
| `logs.sh [-n N] [-f] [PATTERN]`              | read the dev server log (`make watch` tees to `server/build/dev.log`)                                                                                                                        |
| `pagebuild.py SLUG SPEC.json`                | build/refresh charts + pages via REST (spec format in its header)                                                                                                                            |
| `mcp.py tools\|schema\|call`                 | call Human's own MCP server (`/api/mcp`: compose CRUD incl. charts, TAQ/workflow exec)                                                                                                       |
| `verify-ui.mjs PATH…`                        | render-verify webapp paths in headless Chromium as agent-ui; screenshots + console/page errors                                                                                               |
| `drive.mjs SUITE.mjs [--only N]`             | multi-step browser checks: navigate, click, assert where you landed. Logged in already, dialogs recorded, console/network collected per check                                                |
| `ids.sh [SLUG…]`                             | handle → ID map for a namespace (modules/pages/charts/records) cached in `.state/ids.json`; `seed.sh` refreshes it                                                                           |

For interactive Claude Code sessions, `.mcp.json` registers the `human` MCP
server; it needs `HUMAN_MCP_TOKEN` exported before starting Claude Code:
`export HUMAN_MCP_TOKEN=$(dev/agent/token.sh)` (2h lifetime — re-export when
stale, or use `mcp.py` which self-authenticates per call).

```sh
dev/agent/smoke.sh
dev/agent/api.sh GET '/system/users/?limit=5'
dev/agent/api.sh POST /compose/namespace/ -d '{"name":"Sandbox","slug":"sandbox_demo"}'
```

## Checking the webapp

`verify-ui.mjs` answers "did this path render clean". Anything needing a second
step — click Back and assert where you land, search and count rows, leave an
editor and see whether it warns — goes through `drive.mjs`:

```js
// dev/agent/checks/mine.mjs
import { drive, check, expectPath, ids } from '../drive.mjs'
const id = ids('catalogue') // handles, not pasted IDs

drive('back from a deep link reaches the list', async page => {
  await page.open(`/compose/namespace/catalogue/admin/modules/${id.module.catalogue_field}/edit`)
  await page.back()
  expectPath(page, '/compose/namespace/catalogue/admin/modules')
})
```

```sh
dev/agent/ids.sh catalogue                       # once, or via seed.sh
node dev/agent/drive.mjs dev/agent/checks/mine.mjs
node dev/agent/drive.mjs dev/agent/checks/mine.mjs --only "back"
```

Each check gets its own browser context, so a deep link genuinely has no
history to go back to. Native dialogs are **recorded, not dismissed** —
playwright's default dismiss turns an unsaved-changes confirm into a phantom
"the button does nothing". Console errors, page errors and HTTP 500s fail the
check unless it passes `{ allowProblems: true }`. Screenshots are written only
on failure.

In a dev build the app also exposes `window.__human` — `routes()` (name → path,
so route paths need not be grepped), `stores()` (every Pinia store, state
readable) and `globals()` (`$Settings`, `$Auth`, …). Reach for it instead of
editing a component to park state on `window`:

```js
const ns = await page.evaluate(() => window.__human.stores().namespace.set.length)
```

Shell landmarks carry test ids: `app-sidebar`, `app-topbar`, `editor-actions`,
`editor-back`.

Read a container's text with `textContent`, not `innerText`. `innerText` is
layout-dependent and returns only part of the subtree for panels inside a
dialog — a check asking "does this form have a Page Layout field" answers no
while the screenshot of the same moment shows the field. `textContent` sees the
DOM as it is; use a visibility assertion when visibility is the question.

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
- **Re-running `pagebuild.py` over an existing page does not extend its
  layout.** The tool writes `blocks` and nothing else; the layout the server
  gave the page at creation still lists the original blockIDs, so blocks a
  later run adds are stored but never rendered — the page quietly shows a
  subset while the tool prints "page … updated". Delete the page and let it be
  created again when the block set changes.
- **Do not let the formatter touch a generated file.** `dev_format_run` with no
  arguments formats everything git reports as changed, and prettier realigns
  the markdown tables in `server/system/agentic/mcp/TOOLS.md` — which is
  generated, and compared byte-for-byte by `TestToolsMatrix`. The test then
  fails on whitespace with a diff thousands of lines long. Regenerate with
  `cd server && go test ./tests/mcp/ -run TestToolsMatrix -update` and pass
  explicit paths to the formatter afterwards. `dev_commit_create` formats
  before staging too, so it mangles the file **at commit time**, after the
  suite you just ran went green — commit a regenerated `TOOLS.md` with plain
  `git`, and check `git show --stat` afterwards.
- **A webapp path no route matches renders the home page, not an error.** The
  router's catch-all (`client/web/unify/src/router/index.js`) redirects an
  unknown path to `/`, which loads cleanly — so a wrong path used to come back
  from `verify-ui.mjs` / `dev_ui_verify` as a pass, with a screenshot of a
  healthy page nobody asked about. Both now compare where the app settled
  against what was requested and say so. Every compose route is under
  `/compose` (`sections/compose/index.js` prefixes them): a page is
  `/compose/namespace/<slug>/pages/<pageID>`, addressed by namespace **slug**,
  not `/compose/ns/…`.
- **`stale: false` can still mean a stale binary.** The check compares the
  binary's build time against source mtimes, so a rebuild that the watcher
  started _before_ you edited, and finished _after_, reports as fresh while
  running the older code. It bites hardest when editing a file repeatedly in
  quick succession — mutation-testing a file, say, where a live check then
  reports behaviour matching a mutation you already reverted. Cost an
  investigation once. When a live result contradicts a passing unit test,
  `touch` the file, wait for the next build, and re-check before believing
  either.

## Conventions for agent-created data

- **Clean up what you created, and only that.** `api.sh` and `mcp.py` record
  every namespace they create in `.state/created.jsonl`, tagged with the
  session, and `cleanup.sh` deletes that session's entries and nothing else.
  Anything the session did not create is off-limits, whoever made it and
  whatever it is called. No name prefix is needed or wanted.
- **Names are snake_case** — namespace slugs, module and page and chart and
  TAQ handles, and every module field name. A hyphen is the subtraction
  operator wherever an identifier is parsed, so `close-date` lexes as `close`
  minus `date`.
- Good agent-built systems get harvested into fixtures:
  `server_cli export compose-namespace <handle>` (see `common.sh` for
  `server_cli`).

## State

`.state/` (gitignored): client secret, cached token. Delete it any time;
`bootstrap.sh` regenerates.

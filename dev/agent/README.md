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

`make setup-agent` runs this, fills `E2E_PASS` and smokes it; call the script
directly to re-run just this part. Requires the dev server running
(`cd server && make watch`). Creates:

- user `agent@local.dev`, member of `super-admin`, with a password in
  `.state/ui-password` — the one identity for both API and browser
- auth client `dev_agent` (`client_credentials` grant, impersonates
  `agent@local.dev`), managed via REST; its server-generated secret is cached
  in `.state/secret` (gitignored)

Self-healing: re-running bootstrap repairs a misconfigured client (PUT with
the full desired config) and re-caches the secret from the expose endpoint.

Requires: `curl`, `python3` (no jq dependency), a built server binary in
`server/build/` (`make watch` provides one).

Which server it reaches is not typed anywhere. `stack.sh` resolves it from the
checkout's own `server/.env` and `.env.e2e`, so a worktree's scripts answer for
the worktree; every script here sources it through `common.sh`, and the node
ones through `stack.mjs`. Print it to see where a checkout points:

```sh
dev/agent/stack.sh   # HUMAN_API · HUMAN_BASE · HUMAN_AUTH · HUMAN_WEBAPP
```

An exported `HUMAN_API` or `HUMAN_WEBAPP` still overrides it, but nothing in
normal use needs one.

## Daily use

| script                                       | purpose                                                                                                                                                                                      |
| -------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `smoke.sh`                                   | server up? version? token OK? who am I? authed API works?                                                                                                                                    |
| `token.sh`                                   | print a valid bearer token (cached ~2h; oauth2 flow, CLI-jwt fallback)                                                                                                                       |
| `api.sh [--json] METHOD PATH [curl args…]`   | authenticated request, jq-pretty output; `--json` guarantees JSON on stdout for every outcome, so a parser never dies on a diagnostic                                                        |
| `seed.sh [--force] [fixture…]`               | import `dev/fixtures/` (skips seeded fixtures unless `--force`)                                                                                                                              |
| `cleanup.sh [--all\|--session ID] [--purge]` | delete namespaces recorded in `.state/created.jsonl` for this session (`--all` for every session); `--purge` also hard-deletes soft-deleted corpses (dev-only server CLI, NOT ledger-scoped) |
| `logs.sh [-n N] [-f] [PATTERN]`              | read the dev server log (`make watch` tees to `server/build/dev.log`)                                                                                                                        |
| `pagebuild.py SLUG SPEC.json`                | build/refresh charts + pages via REST (spec format in its header)                                                                                                                            |
| `mcp.py tools\|schema\|call`                 | call Human's own MCP server; failures print `{"error":{…}}` on stdout and 5xx retries once (`/api/mcp`: compose CRUD incl. charts, TAQ/workflow exec)                                        |
| `verify-ui.mjs PATH…`                        | render-verify webapp paths in headless Chromium as the agent user; per-path text report (clipped blocks, cut-off columns, raw IDs) + screenshots + console/page errors                       |
| `drive.mjs SUITE.mjs [--only N]`             | multi-step browser checks: navigate, click, assert where you landed. Logged in already, dialogs recorded, console/network collected per check                                                |
| `ids.sh [SLUG…]`                             | handle → ID map for a namespace (modules/pages/charts/records) cached in `.state/ids.json`; `seed.sh` refreshes it                                                                           |

For interactive Claude Code sessions, `make claude MCP=1` attaches
`human-local`, this checkout's `/api/mcp`, through `dev/human-local.mcp.json`
and mints the token it needs — a config file can interpolate an environment
variable but cannot run a command to produce one. `mcp.py` needs none of that;
it authenticates per call and is the better bet from a script.

```sh
dev/agent/smoke.sh
dev/agent/api.sh GET '/system/users/?limit=5'
dev/agent/api.sh POST /compose/namespace/ -d '{"name":"Sandbox","slug":"sandbox_demo"}'
```

## Two browser logins

`bootstrap.sh` provisions both, idempotently; the passwords live in
`.state/` and are gitignored.

| user                 | password file        | sees                                                         |
| -------------------- | -------------------- | ------------------------------------------------------------ |
| `agent@local.dev`    | `.state/ui-password` | everything — super-admin; the default for UI checks          |
| `agent-ro@local.dev` | `.state/ro-password` | users and roles, read only — role `agent_readonly`, no group |

The first is the same user the auth client impersonates, so a browser check and
an API call are one actor. Nothing in the `client_credentials` flow objects to
it also holding a password — it loads the impersonated user by ID.

The read-only user exists to check what someone **without** permission gets,
which a super-admin account cannot show: the admin sidebar drops the entries it
cannot reach (5 of 30 remain), and an editor opened by deep link renders its
banner with every field disabled. A list row does not open the editor at all
for such a user, so deep-link to it.

Its role allows a deliberate SUBSET, so a gate that has stopped working shows up
as entries that should have gone. Two things about the rules are easy to get
wrong: the component resource is `corteza::system/` (`corteza::system:component`
is refused), and without `access` on `corteza::system:application/*` the shell
bounces every admin route to `/?denied=admin`.

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

- **The dev watcher rebuilds on `.go` writes only.** A skill in
  `server/system/agentic/skills/library/*.md` is `go:embed`ded, so editing it
  changes nothing on the running server until some `.go` file is written; a
  `touch` is not a write.
- **A fixture directory must be named after its namespace slug.** `seed.sh`
  imports the directory, then looks the namespace up by the directory name; a
  mismatch reports "namespace not found" after the import has already created
  modules with no records behind them. `seed.sh --force <slug>` cleans it up.
- **`TOOLS.md` is a golden file.** `dev_commit_create` runs prettier on
  Markdown, which pads the table and makes `TestToolsMatrix` fail; commit it
  with `skipFormat`.

- **Trailing slash matters** on collection endpoints: `/system/users/` works,
  `/system/users` is a 404 (chi routing).
- **`DOMAIN_WEBAPP` is where MCP tool results say to go.** `server/pkg/weburl`
  builds every `url`/`editUrl` a tool result carries from it. Here the webapp is
  vite's port, not the API's, so `server/.env` needs `DOMAIN_WEBAPP=localhost:5173`
  on the primary (`worktree.sh new` writes the slot's own vite port). Without
  it the links come back pointing at the API, which does not serve the webapp.
- **Errors come back as HTTP 200** with `{"error":{"message":…}}` — `api.sh`
  detects this and exits non-zero.
- **`api.sh` cannot upload a file.** It always sends
  `Content-Type: application/json`, so a `-F upload=@file` body is parsed as
  JSON and fails with `invalid character '-' in numeric literal`. Multipart
  endpoints (settings logos, attachments) take a raw curl: the token from
  `token.sh`, the `/api` prefix (`http://localhost:<api>/api/system/settings/ui.main-logo`;
  without `/api` the auth server answers with its HTML), and `-F upload=@file`.
- **A user update needs a fresh read first.** `PUT /system/users/{id}` is
  optimistically locked: a body saved from an earlier GET fails with "Someone
  else, or a workflow, changed this user after you opened it". GET, patch, PUT
  in one go.
- **`api.sh` shows an error's English message, the webapp its translation.**
  The webapp's JSON requests make the server translate an error from
  `locale/en/human-server/`, and a missing string comes back as the raw key
  (`user-group.errors.notFound`). Check wording the way the webapp gets it:
  `curl -H 'Accept: application/json' -H 'Accept-Language: en'` with the token.
- **Automation `input` takes only typed envelopes.** `ngAutomationExec`,
  `workflowExec` and session resume decode `input` into `expr.Vars`, which reads
  `{"@type":…,"@value":…}` and nothing else, so one bare string or array rejects
  the whole request before the automation is looked up — as an HTTP 200 carrying
  `json: cannot unmarshal … into expr.typedValueWrap`. A list goes as
  `{"@type":"Array","@value":[…]}` with its items raw; wrapping an item in its
  own envelope turns it into a `Vars` with the keys `@type` and `@value`.
- **`/automation/sessions/` lists only unfinished sessions by default.** A
  workflow that ran and finished looks like it never started; `completed=1`
  includes finished (completed or failed) sessions, `completed=2` lists only them.
- **A step or trigger created without `meta.visual` does not render in the
  editor.** The codec keys nodes by `meta.visual.id` and places them by
  `meta.visual.xywh`; give both when a browser check needs to click the node.
- **`make -C server codegen-legacy` writes `*_actions.gen.go` and then panics**
  on the docs output (`../docs/src/modules` missing). The actions it wrote are
  complete; diff them and discard anything else it touched.
- **A workflow trigger is its own resource, and its step field is
  `workflowStepID`.** `workflowCreate`/`workflowUpdate` ignore a nested
  `triggers` array, so POST `/automation/triggers/` separately — and `stepID`
  there is accepted, ignored, and echoed back as `0`, which reads exactly like
  the server refusing the value.
- **An RBAC check on a wildcard resource always answers "no".**
  `rbac.service.checkValidity` refuses any resource containing `*`, so
  `Can(ses, op, "corteza::compose:record/1/2/*")` is false however the rules
  read, and `/permissions/effective?resource=…/*` reports everything denied.
  `Trace` is the one entry point that evaluates a wildcard honestly — rules are
  matched with `path.Match` against the requested pattern, so a rule at the same
  or a broader scope matches and a narrower one does not. It skips the
  user-group branch, so it can only under-report.
- Tokens need `scope=profile api` — the API middleware requires the `api`
  scope.
- Envoy YAML import cannot configure auth clients fully (`validGrant` is
  silently dropped, `security.impersonateUser` refs don't resolve) — manage
  auth clients via REST, as bootstrap does, or via
  `system_auth_client_create`, which sets both.
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
  `git`, and check `git show --stat` afterwards. The same goes for
  `server/pkg/codegen/resource_schema.gen.json` and `server/system/rest.yaml`:
  prettier reflows the whole file, thousands of lines, for a three-line edit —
  pass `skipFormat` to `dev_commit_create` when either is in the commit. A
  `server/*/service/*_actions.yaml` gets every double quote rewritten to single
  the same way.
- **The dev tools' `files` argument is space-separated.** A comma-separated
  list is read as one path, and `dev_intent_governing` answers it with "no
  intent doc governs this path" — a clean negative, not an error, for files
  that are covered. `dev_format_run` likewise reports "nothing to format".
- **`make codegen` exits 0 when the REST codegen failed.** `codegen-legacy`
  runs `$GOPATH/bin/human-codegen`, which is built only when missing, so a
  binary older than `server/pkg/codegen/assets/*.tpl` dies on the template
  (`can't evaluate field … in type *codegen.eventProps`) before it reaches
  `rest.yaml` — `rest/handlers` and `rest/request` stay stale while the make
  target reports success. Grep the output for `failed to process`; a fresh
  `go build -o <tmp> ./cmd/codegen/main.go` run from `server/` with `-v` fixes it.
- **A server integration suite passing on the primary can still fail in a
  worktree at `could not find en in loaded languages`.** The primary may hold a
  gitignored `server/pkg/locale/src/en/` from an old `make -C pkg/locale`, and
  that copy gets built into the binary. It hides a `LOCALE_PATH` that points
  nowhere, and a worktree has no such copy. `tests/helpers` resolves a relative
  `LOCALE_PATH` from `server/.env` against `server/`, so look for
  `language overloaded … imported` in the suite's log. `embedded: true` alone
  means the path missed. In `server/tests/compose`, error assertions need
  `Header("Accept", "application/json")`, since without it the body is plain
  text and `AssertError` fails to decode it. The suite runs with no locale
  loaded, so `AssertError` compares the translation key (`module.errors.…`),
  not the message.
- **`.p-select.p-disabled` in a `drive.mjs` check matches selects on hidden
  tabs.** The module editor renders the unique-values tab's disabled Column
  select up front. Scope the locator to the row, e.g.
  `.p-inputgroup:has(.p-select.p-disabled)`.
- **A webapp path no route matches renders the home page, not an error.** The
  router's catch-all (`client/web/unify/src/router/index.js`) redirects an
  unknown path to `/`, which loads cleanly — so a wrong path used to come back
  from `verify-ui.mjs` / `dev_ui_verify` as a pass, with a screenshot of a
  healthy page nobody asked about. Both now compare where the app settled
  against what was requested and say so. Every compose route is under
  `/compose` (`sections/compose/index.js` prefixes them): a page is
  `/compose/namespace/<slug>/pages/<pageID>`, addressed by namespace **slug**,
  not `/compose/ns/…`.
- **A build is ~15s, and the old server serves through it.** `make watch`
  rebuilds and restarts on any `.go` write, and no edit is dropped however fast
  they land — but for those seconds the API answers with the previous build.
  That bites hardest when editing a file repeatedly in quick succession:
  mutation-testing, say, where a live check reports behaviour matching a
  mutation you already reverted. Cost an investigation once. `dev_server_status`
  with `wait` blocks until the running process's start time beats your edit; a
  timeout means the build failed, and `dev_server_logs` has the compiler
  output.
- **Vite does not pick up edits under `lib/`.** The webapp resolves
  `@planetcrust/human-js` / `human-vue` to their TypeScript sources, but the
  running dev server keeps serving the transform it made at startup — a
  `touch` does not invalidate it, and the file on disk and the module the
  browser gets disagree indefinitely. A UI check then reports the old
  behaviour against a fix that is genuinely there, which reads as the fix not
  working. Confirm what is actually served before believing a UI result about
  a lib change:
  `curl -s "$(dev/agent/stack.sh | sed -n 's/^HUMAN_WEBAPP=//p')/@fs<abs-path-to-file>" | head`.
  Only a vite restart clears it — ask the human, never restart it yourself.
- **A same-origin `/api/…` fetch from inside a drive check hits the SPA, not
  the API.** Vite proxies exactly two paths (`/custom.css`,
  `/code-snippets.js`, `client/web/unify/vite.config.js`); everything else
  falls through to the history fallback, so `fetch('/api/…')` comes back
  **HTTP 200 with `index.html`**. `r.json()` then throws on the doctype, and a
  check reading `r.ok` or `r.status` instead calls a nonexistent endpoint
  healthy. The webapp itself talks to an absolute origin (`window.HumanAPI`,
  `public/config.js`), which a page probe can use — but the API state a check
  wants to assert is cheaper read from the shell with `api.sh` after the run.

- **`performance.getEntriesByType('resource')` overflows before your XHR.**
  The buffer holds 250 entries by default and a cold vite page fills it with
  module scripts, so a probe asking "was this request made" comes back empty
  for every page and every request — an assertion that nothing was fetched
  then passes without testing anything. Call
  `performance.setResourceTimingBufferSize()` and `clearResourceTimings()`
  first, re-trigger the request, and prove the probe by asserting a call
  **is** seen somewhere in the same run.

- **MCPJam 3.x cannot render an MCP App headless.** It opens on an account
  sign-in before any server is shown. The ext-apps repo's
  `examples/basic-host` is the host to drive instead: build it against the
  published `@modelcontextprotocol/ext-apps` `dist`, run `serve.ts` with `tsx`
  (it assumes bun), and put a local proxy in front of `/api/mcp` that adds the
  bearer and CORS, exposing `mcp-session-id`. Its page never reaches
  `networkidle` (the SSE stream stays open), so wait for `load`; the view sits
  two iframes down.

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

## Your own space: `worktree.sh`

One task, one checkout, one server, one webapp, one database. Nothing a
worktree does reaches another session, and nothing another session does can
change what it is testing.

```
worktree.sh new  NAME [--base REF]   checkout + DB clone + ports + env files
worktree.sh up   [NAME]              start its server and webapp
worktree.sh down [NAME] [--force]    stop them
worktree.sh list                     every slot and what is running
worktree.sh land NAME [--keep]       rebase onto main, merge, remove
worktree.sh gc   [--reap]            find abandoned worktrees and databases
worktree.sh rm   NAME                stop, drop the DB, remove the checkout
```

Ports come from the slot, so two worktrees cannot collide — API `+slot*100`
and vite `+slot`, counted off whatever the primary itself serves rather than
off a fixed number. The primary is slot 0 and is never reassigned. `new` writes
`server/.env`, `public/config.js`, `.env.local` (the vite port) and `.env.e2e`
pointed at those ports; all four are gitignored, so nothing shows up in the
worktree's `git status`, and they are what `stack.sh` reads back — so a
worktree needs no `make setup` and no `HUMAN_API` in its environment.

The whole flow works in there — unit tests, API, browser, and e2e. Playwright
reads `E2E_BASE_URL` from the worktree's `.env.e2e`, so `npx playwright test`
run from the worktree hits the worktree.

### Nothing cleans itself up

Closing a tab reaps none of this. The checkout, branch, database and slot all
survive, and the servers are their own process group, so they outlive the
terminal too. Only `rm` and `land` remove anything, and a session that ends
mid-task calls neither.

`worktree.sh gc` is what finds the leftovers — an entry whose checkout is gone,
a database no entry claims, a worktree that is clean and fully merged and
serving nothing. It reports; `--reap` removes. A worktree holding uncommitted
work or unmerged commits is reported in full and never touched: that work
exists nowhere else.

### Getting the work onto main

A worktree commits to its **own branch**, never to main. `land` is the whole
sequence: rebase onto current main, merge fast-forward, remove the worktree.

It is mostly refusals, and deliberately so — main lives in the primary's
working tree, which usually holds another session's uncommitted work. `land`
stops before touching anything if the worktree is dirty, if there is nothing to
land, if the branch does not rebase cleanly, or **if the primary has
uncommitted changes in a file the branch also touches**. `--keep` merges but
leaves the worktree running, for landing an increment and carrying on.

An agent lands a High-confidence change itself and says so; at Medium or Low
it asks first. A refusal is always reported, never worked around.

If you `rm` before landing, the branch survives with your commits and `rm`
says so, with the command to merge it.

What catches people out:

- **A worktree checks out HEAD.** Uncommitted work in the primary does not come
  with it. `new` warns and names the count.
- **First `up` is minutes, not seconds** — `pnpm install` and a Go build. After
  that it is seconds.
- **Do not pipe `up` into `tail` or `head`.** The server it starts keeps the
  pipe open, so the pipeline never ends and anything chained after it never
  runs. Redirect to a file instead.
- **The server rebuilds and restarts itself**, here as on the primary. Judge by
  the process start time, as on the primary. `touch` on a `.go` file does not
  restart it; to check what survives a restart, `worktree.sh down` then `up`.
- **The shared token works against every worktree** — same JWT secret, and the
  cloned DB has the same user IDs. No re-bootstrap for the API.
- **A browser login to a fresh worktree can be refused** ("invalid username and
  password combination") while the token works. `bootstrap.sh` run from the
  worktree resets `agent@local.dev`'s password in its database to `.state/`'s.
  The primary can drift the same way (`page.waitForURL` timeout from
  `dev_ui_verify`); `bootstrap.sh` from the primary is the same fix. It prints
  `read-only login ready` for `agent-ro` even when the server refuses that
  password as not secure enough, so an `agent-ro` login that still fails needs
  its password set by hand, on that checkout's database only.
- **`dev_ui_verify` takes no worktree argument.** It drives the webapp named in
  the `.env.e2e` of the checkout its MCP server was started from — the primary,
  for a session launched there. Against a worktree, run
  `node dev/mcp/uiverify.mjs` with `baseURL` in its JSON argument;
  `npx playwright test` run from the worktree already reads the worktree's own
  `.env.e2e`.
- **`up` also starts Corredor** when the checkout has `corredor/`: on the
  slot's gRPC port (API port + 50000), before the server, with
  `CORREDOR_ENABLED=true` and `CORREDOR_ADDR` written into the slot's
  `server/.env`. It mounts `dev/fixtures/corredor` (scripts against the
  `agent-sandbox` namespace — seed it) plus `corredor/usr/*`; log at
  `.run/corredor.log`. `down` stops it last.
- **`rm` refuses while the checkout is dirty.** It will not discard your work.
- **A worktree needs no `cleanup.sh`** — `rm` drops its whole database, so
  nothing it created can outlive it.
- **`rm` hands the checkout's open todo to the global pool**, so the queue is
  the one thing it does not take with it.

## Does a skill still work? `skill_eval.py`

A skill is text a model reads, so the only proof it still works is a model
following it. `skill_eval.py` runs each prompt in `evals/<skill>.json` in a
fresh `claude -p` session that sees only this checkout's MCP server — no
CLAUDE.md, nothing pasted — and scores what the session produced.

    dev/agent/seed.sh contacts_crm agent-sandbox   # the fixtures the prompts use
    dev/agent/skill_eval.py                        # every prompt
    dev/agent/skill_eval.py --only tempting-extras # one

For `custom_app` a prompt passes when the session read the skill before
writing, and the page is plain HTML with the bridge snippet unchanged, carries
`window.SAMPLE` and a sample/live badge, and is accepted by the real deploy
guard (on a scratch app the script deletes again). Transcripts and pages land in
`.state/evals/<timestamp>/`. Sessions run on Sonnet unless `--model` says
otherwise, never on Fable.

Run it after any change to a skill, to the MCP instructions, or to a tool a
skill names. The server embeds skills at build time and the watcher rebuilds on
`.go` writes only, so rebuild before an eval that follows a skill edit.

What the custom app rounds taught, worth keeping when the probe changes:

- **A check must be able to fail.** Write briefs once passed without writing
  anything, because the probe only pressed buttons. The `saves` check now asks
  the store whether anything changed, and a brief the probe cannot drive fails
  as unproven rather than passing quietly.
- **Drive a page the way a person does**: press what opens something, fill
  what is now on screen, then press what commits. Never press `close`,
  `cancel`, `dismiss` or `×`, and click a switch's label, not its hidden input.
- **Where the bridge names something, the name is the API's, or it takes
  both.** Every product bug the rounds found was a bridge name differing from
  the tool the skill sends authors to (`label`/`text`, `dimensions`/`dimension`,
  a six-key theme against the tool's twelve). Check a new operation against its
  tool's parameters before it ships.
- **The score is what must be true, not craft.** Two models scoring the same
  can still differ in quality, so judging a page still means opening its
  screenshot.
- **Reseed between models.** They share a database, and a second session that
  finds the first one's module is right to refuse to duplicate it.

## Deferred work: `backlog.sh`

The queue of things a task decided not to do, so they survive the turn that
found them. It has two tiers, and which one an item lands in is the whole
design.

```
backlog.sh add TEXT [--why W] [--files F,F] [--task T] [--global]
backlog.sh promote ID                       local → global
backlog.sh list [--global|--both|--orphaned] [--closed] [--files F]
backlog.sh show ID
backlog.sh done ID [--note N]   /   backlog.sh drop ID [--note N]
```

**Local is this checkout's todo.** An item belongs there when doing it would
change a file the current task is already changing, or when it follows
directly from that change. `list` with no flags shows exactly that, and it is
what triage and the end of a task read.

**Everything else is `--global`.** A real defect found while looking at
something else is filed, named in the report, and then out of the way. The
pool is read when somebody asks for it — `list --global` — and at
`/orchestrate` intake, which is asking for it.

The two tiers are one file, so nothing is lost by filing globally; it is only
not put in front of the next turn.

**A worktree owns its local items**, so the queue survives the session that
opened the checkout, and `worktree.sh rm` promotes what is still open rather
than dropping it with the directory. On the primary there is no such boundary,
so the owner is the session.

Every scope the current script did not write is global: the old shared pool,
and the per-session items from before the tiers. `--orphaned` is the recovery
path for local items whose worktree is gone or whose session ended.

It lives in the shared `.state`, which every worktree symlinks, so both tiers
are one file for every session on this machine. Append-only JSONL: two sessions
writing at once interleave lines rather than corrupting each other, and closing
an item is another append, not a rewrite.

Each item records where it came from — the task, the session, the commit HEAD
was on. `--files` is what makes it findable later, and `list --global --files
<path>` is the deliberate check for whether anyone deferred something where you
are about to work.

## State

`.state/` (gitignored): client secret, cached token, the backlog, and the
worktree registry. Delete the secret and token any time; `bootstrap.sh`
regenerates them.

Scratch inside it is **per session** — `.state/sessions/<session>/` holds
screenshots, browser storage and drive output. A fixed name there means a peer
session's screenshot arrives under your filename, which has happened.

A worktree's `dev/agent/.state` is a **symlink** to the primary's, so token,
ledger and backlog are shared while scratch stays split. That is why
`dev/agent/.gitignore` says `.state` and not `.state/` — a trailing slash
matches directories only, and git sees a symlink as a file.

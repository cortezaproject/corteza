# Dev toolkit gotchas

Facts about the `dev/agent` scripts, the dev MCP tools, worktrees, running tests and e2e, git habits in a shared checkout, fixtures and the cleanup ledger — each one a trap that reports success while checking nothing.

## Scripts and servers

### The dev watcher is devwatch

`make watch` runs `server/cmd/devwatch`: no proxy, no port 3001, binary `server/build/dev-bin`. No write during a build is lost, a second watcher per checkout is refused by a flock on `dev-bin.lock`, and a failed build keeps the old server up with compiler output in `server/build/dev.log`. The server has its own process group, so stopping it means signalling the listener's group; `worktree.sh down` does.

**Why:** a stale server can mean a failed build, and `processStartedAt` can disagree with the log.

**How to apply:** `dev_server_status` with `wait` blocks until the process start beats your newest Go source; judge by `processStartedAt`, never `binaryBuiltAt`, and let the log win a disagreement (`dev_server_logs`, or `dev/agent/logs.sh`, errors first, `-a` for the raw tail).

### `stack.sh` resolves which server a checkout talks to

`dev/agent/stack.sh` resolves `HUMAN_API`, `HUMAN_BASE`, `HUMAN_AUTH` and `HUMAN_WEBAPP` from the checkout's own files, so a call made in a worktree reaches that worktree. `common.sh` sources it and `stack.mjs` execs it for the node scripts; `dev/mcp/tools` reads the same files in Go (`HTTP_ADDR` from `server/.env`, `E2E_BASE_URL` from `.env.e2e`). Precedence: environment, then `server/.env` (`HTTP_ADDR`, `HTTP_API_BASE_URL`) for the API; for the webapp `VITE_PORT` (`.env.local` over `.env`), then `.env.e2e` (`E2E_BASE_URL`), then the worktree registry; then the defaults 1043/5173.

**Why:** a hardcoded port makes a worktree's checks pass against the primary's database.

**How to apply:** never type a dev port or export `HUMAN_API` for a worktree; run `stack.sh` bare to see where a checkout points. Reading an env file by hand, take the last uncommented assignment, and skip `//` lines in `public/config.js`.

### `bootstrap.sh` against a separate database overwrites the shared secret

`dev/agent/.state/` is one directory for every checkout (a worktree symlinks the primary's). `bootstrap.sh` creates the `dev_agent` client in whatever database the server it hits uses and caches that secret in `.state/secret`. Pointed at a non-cloned database, it replaces the primary's secret; `token.sh` then fails `mint_oauth` and silently falls back to `mint_cli`. Worktrees from `worktree.sh new` clone the primary's DB and are not affected.

**Why:** nothing errors; the oauth path just stops being exercised.

**How to apply:** after bootstrapping a stack on a non-cloned database, run `bootstrap.sh` once from the primary and delete `.state/token` and `.state/token-exp`.

### The cleanup ledger holds namespaces only

`.state/created.jsonl` records namespace creates and nothing else: `api.sh` calls `ledger_record` only for a namespace, and `mcp.py`'s only hook is `ledger_record_namespace`. A workflow, role, user, agent or auth client created through either is invisible to `cleanup.sh`. `seed.sh` writes nothing to the ledger; fixtures are meant to persist and re-seed with `--force`.

**How to apply:** namespaces go through `cleanup.sh`; track and delete every other resource you create yourself.

### `cleanup.sh` can scope to the wrong session

`cleanup.sh` deletes the current session's ledger entries, keyed by `CLAUDE_CODE_SESSION_ID`. In a background job the id it computes can differ from the one `api.sh` stamped, and it prints `cleanup done — 0 recorded namespace(s) in scope`, which reads as a successful no-op. `dev_fixture_cleanup` has the same defect and words it as an error about the server restarting.

**Why:** the namespace stays on the server while the output says done.

**How to apply:** check the count against `.state/created.jsonl`; if it says 0 but the ledger holds your entries, re-run with `--session <id-from-ledger>`, never `--all`. Confirm each namespace is gone with a `GET`.

### Logging in as another user from the shell

`dev/agent` has no "act as user X" helper; auth is a session-plus-CSRF HTML flow. With one cookie jar against the API origin: `GET /auth/oauth2/authorize?response_type=code&client_id=…&redirect_uri=…&scope=profile+api` (or `GET -L /auth/oauth2/default-client?…`), scrape `same-site-authenticity-token` from `GET /auth/login`, `POST /auth/login` with email, password and token, follow the redirects to `…?code=`, then `POST /auth/oauth2/default-client` with `code`, `redirect_uri` and a fresh token for JSON with `access_token`. Landing on `/auth` means the flow was rejected; the alert is in the HTML.

**How to apply:** set a password with `POST /system/users/{id}/password`; toggle confirmation with a JSON-patch `PATCH /system/users/{id}`. Users you create are not in the ledger; delete them by hand.

### Vite and edits under `lib/`

The unify vite server resolves `@planetcrust/human-js` and `@planetcrust/human-vue` to their TypeScript sources through symlinks outside its root. `client/web/unify/vite.config.js` carries a dev-only `human:watch-lib-sources` plugin that directory-watches `lib/js/src` and `lib/vue/src`; without it a per-file watch dies on the first git write and that module stays frozen until restart. A frozen `lib/vue/src/index.ts` barrel missing a new export blanks the whole app with `does not provide an export named …`.

**How to apply:** a `/@fs` curl cannot tell fresh from frozen, and the served transform strips comments; probe by behaviour through `window.__human`. Only a vite restart clears a frozen module, and that is the human's call; verify lib behaviour with the package's own suite meanwhile.

### A nested Claude Code session is the toolkit's end-to-end test

`claude mcp list` from a checkout reports whether `human-dev` actually connects
(`✔ Connected` vs `CONNECTION_CLOSED`), which a stdio probe of `run.sh` cannot
show. `claude -p "<task>" --output-format json --max-turns 25 --model sonnet
--mcp-config .mcp.json --strict-mcp-config` runs one task and returns
`num_turns`, `duration_ms`, `usage` and `total_cost_usd`, so the same prompt in
two checkouts is an A/B of the toolkit. Inside a Claude Code session the CLI
refuses to nest until `CLAUDECODE` is unset (`env -u CLAUDECODE claude …`).
`--strict-mcp-config` matters: without it the nested session also gets every
account-level connector, remote Human instances included.

**How to apply:** export `HUMAN_API` (and `HUMAN_DEV_LOG`) so a worktree's
scripts reach a running server; compare turns, input tokens and cost, and read
the answers, since a cheaper wrong answer is not an improvement.

### docs/ is the product site and renders Vue

`docs/` is VitePress: every `.md` under it becomes a public page and the llms
plugin folds it into `docs/llms-full.txt`, and `{{ … }}` in prose is a Vue
expression (a literal `{{ $t('k') }}` broke the build). Developer notes live in
`dev/gotchas/`, not under `docs/`. Files other than `.md` are served only from
`docs/public/` (the schemas live there, at `/schemas/*.json`); `make docs-llms`
refreshes the committed `llms*.txt` and `dev/agent/docs-llms.test.mjs` fails
when they are stale.

## Dev MCP tools

### `files` arguments are whitespace-separated

`dev_format_run`, `dev_commit_create` and the intent tools split `files` with `strings.Fields`. A comma list is read as one path that does not exist: `dev_format_run` answers `"nothing to format"` and `dev_intent_governing` answers "no intent doc governs this path".

**How to apply:** pass `"a.go b.go"`. When a format run skips everything, check the separator before believing the files were clean.

### The format step reflows generated and locale files

`dev_format_run`, and the format step inside `dev_commit_create`, send every `.ts .js .vue .json .yaml .yml .md .css .scss` file to `prettier --write`; the skip list is only `.gen.go`, `/vendor/` and `/node_modules/`, and `.prettierignore` covers only `docs/**/*.gen.md`. Naming `pnpm-lock.yaml` reflows the whole lockfile; hand-written locale yaml gets unrelated quotes flipped. Same for `*.gen.ts`, `*.gen.json`, `.intent/intent.lock.json` and `server/system/agentic/mcp/TOOLS.md`.

**Why:** the format runs after you inspect `git diff`, and the tool reports success either way.

**How to apply:** pass `skipFormat: true` for generated files and `locale/**/*.yaml`, or stage them with git; read `git show --stat` after committing. Repair a reflowed lockfile by restoring it from the parent and regenerating it with its own tool.

### Tool edits need an MCP reconnect

`dev/mcp/run.sh` rebuilds `build/dev-mcp` when a source (or `server/pkg/mcpkit`) is newer, but only at server start; `/mcp` reconnect is enough. A failed rebuild keeps serving the previous binary, with the compiler output on stderr and in `dev/mcp/build/build.log`. Each Claude session runs its own `dev-mcp` child, so only the reconnecting session gets the new code. `dev/mcp/uiverify.mjs` runs as `node uiverify.mjs` per call and takes effect at once, though a new field it returns is dropped until the Go struct in `tools/ui.go` ships.

**Why:** stale tool output reads as "my fix did not work", and a broken build still connects.

**How to apply:** verify a `dev/mcp/tools` change in the same session from a temporary Go test in `package tools`, and tell the human a reconnect is needed.

### Which checkout a tool acts on

The checkout tools (`dev_format_run`, `dev_test_run`, `dev_commit_create`, `dev_branch_status`, the intent tools) take a `worktree` argument (`dev/mcp/tools/tools.go`, `checkoutOption`). Omitted, they act on the checkout the MCP server was launched from and report success against files you never edited. `dev_ui_verify`, `dev_server_status`, `dev_server_logs` and `dev_fixture_cleanup` take none and always answer for the launch checkout.

**How to apply:** name `worktree` whenever the work lives in one. For the UI against a worktree, run `node dev/mcp/uiverify.mjs` with `baseURL`, or playwright from the worktree. A `.vue` that will not compile shows as HTTP 500 from `/@fs<abs-path>` on vite.

### `dev_test_run`'s `run` filters test names

`run` is vitest's `-t` name filter, not a file selector. A path fragment matches no test, zero tests execute, and the tool returns `passed: true`; the only tell is `skipped` dwarfing `packages`.

**How to apply:** use `run` only for a real `describe`/`it` string. For a file, put it in `target` after the workspace (`client/web/unify/src/….test.js`), or name the sources under test in `related`.

`target` takes one path after the workspace. Two or more (`lib/vue src/a src/b`) fail with `produced no report` and run nothing; call it once per path, or run `npx vitest run <paths…>` in the package.

### `dev_ui_verify` is for one-shot checks

Each call writes a uniquely named screenshot under `dev/mcp/.state/shots/` (`dev/mcp/shots.mjs`, pid plus timestamp) and reuses one storage state per origin, `dev/mcp/.state/ui-session-<origin>.json`. Its `steps` grammar is `click` and `fill` only, and it returns no input values or DOM state.

**How to apply:** for anything longer, or a mouse drag, drive playwright directly: resolve `@playwright/test` through `client/web/unify/package.json`, reuse the session file as `storageState`, write screenshots to scratch. An SVG `<pattern>`/`<defs>` node is never visible; wait with `{ state: 'attached' }`.

## Worktrees and git

### Recovering a half-finished `land` or `rm`

`worktree.sh` resolves the primary once at start and refuses to run `dropdb` without a database user. If a `land` or `rm` is interrupted after the merge, the database, slot file, registry JSON or branch can be left behind. Never reach for `gc --reap`: it drops any database it thinks unclaimed, and a peer's worktree mid-`new` can look that way.

**How to apply:** run `land` and `rm` with the primary's `dev/agent/worktree.sh`. To clean up by hand, drop exactly that worktree's DB with the primary DSN's user, remove its `.state/worktrees/<name>.json` and `.slot-N`, then `git worktree prune` and `git branch -d <name>`.

### A peer session can commit your in-flight edits

The primary checkout is shared. A peer that names a file in its commit takes whatever is in it, your uncommitted edits included; the sign is a file dropping off `git status` without being reverted. The mirror case: a `git status`-derived file list sweeps the peer's work into yours.

**How to apply:** re-check your file list immediately before committing (`git diff --quiet -- "$f" && echo "SWEPT: $f"`), name files explicitly, baseline failures against current HEAD (`git worktree add --detach <tmp> HEAD`), and report a swept change rather than rewriting the peer's commit. A worktree avoids the whole class.

### A peer's amend or reset can swallow your commits

On the shared primary HEAD moves under you. A peer's `git commit --amend` can amend your commit, a `reset HEAD~2` can drop it, and the index is shared, so `git add <your files>` does not make the commit only yours.

**How to apply:** commit with `git commit -- <paths>`. After committing, check `git merge-base --is-ancestor <sha> HEAD` for each recent commit and read `git reflog` if anything looks off. Guard any rewrite with `[ "$(git rev-parse HEAD)" = "$EXPECT" ] || exit 1` in the same invocation, and prefer a worktree for multi-commit work.

### `git revert` sweeps staged peer files

`git revert --no-edit <hash>` and a bare `git commit` commit whatever a peer has staged, even seconds after a clean `git status`. The hash being reverted may also no longer be an ancestor once the peer has rewritten main.

**How to apply:** on the primary, `git revert -n <hash>` then `git commit -m … -- <paths>`; check `git log -3` just before. Repair with `git reset --soft HEAD~1`, `git restore --staged` the peer's files, recommit with paths, and say so in the report.

### `git rebase --autosquash` needs `-i`

Without `-i`, `git rebase --autosquash <base>` reports success and does nothing; the `fixup!` commit stays in history.

**How to apply:** `GIT_SEQUENCE_EDITOR=true GIT_EDITOR=true git rebase -i --autosquash --autostash <sha>~1` runs non-interactively. Verify with `git log --oneline`, and check HEAD first, since a peer may have committed.

## Tests and e2e

### `lib/js` runs on mocha

`lib/js` tests run on mocha with chai `expect`: `lib/js/.mocharc.cjs` sets `spec: ['src/**/*.test.ts']`, `tsx/cjs` and `bail: true`. Run `cd lib/js && pnpm run test:unit`, or `npx mocha <file>`. `dev_test_run` picks the runner from the workspace's `.mocharc.*`, so it handles `lib/js` too. `client/web/unify` and `lib/vue` are the vitest packages.

**How to apply:** with `bail: true` the run stops at the first failure, so later tests are absent, not passing. If the tool says "nothing matched" for a path you trust, suspect the runner before the path.

### `lib/vue` collects `*.test.ts` only, and has no `matchMedia`

`lib/vue/vitest.config` includes `src/**/*.test.ts`; a `*.test.js` there is silently not collected. PrimeVue overlays (`Menu`, `TieredMenu`, anything using `bindMatchMediaListener`) throw `matchMedia is not a function` on mount, and `src/test/setup.ts` does not stub it.

**How to apply:** name tests `*.test.ts`, and stub `window.matchMedia` (returning `matches: false` plus no-op listeners) before mounting an overlay.

### Vitest can pass every test and still exit 1

An unhandled promise rejection is reported under **Unhandled Errors**, not against a test. The summary reads `Tests N passed` and the process exits 1.

**Why:** a mutation check that greps the summary calls a covered line uncovered.

**How to apply:** judge a teeth check by the exit code (`npx vitest run <file> >/dev/null 2>&1; echo $?`), and grep for `Unhandled` to learn why. `dev_test_run` reports these as `passed: false`.

### `node --test` hangs on an open handle

`node --test` waits for the event loop to drain, so a test leaving any socket, timer or stream open hangs silently after every assertion passed. The classic source is `createServer().listen(0)` to pick a free port. Detached `setsid` children with `stdio: 'ignore'` do not cause it.

**How to apply:** suspect the harness before the subject; run the fixture steps in plain bash to bisect. Check a port with `ss -ltn "sport = :$p"` instead, as `dev/agent/worktree.test.mjs` does.

### PrimeVue does not resolve in compose view tests

The compose view suites (`RecordView.*.test.js` and siblings) mount `shallow: true` with no PrimeVue plugin, so `Button`, `Message`, `Form` and friends render as literal lowercase tags with props as attributes (`<button label="…">`, `<message>`). Only components from the mocked `@planetcrust/human-vue` get `-stub` tags, and shallow re-stubs them whatever template the mock gives. The "Failed to resolve component" warnings are noise. `positionedBlocks` is `layout.blocks ∩ page.blocks` by `blockID`, so a fixture needs the block on both sides.

**How to apply:** query `findAll('button')` plus `attributes('label')`, `find('message')`, and lib components by kebab stub tag.

### View tests need `CViewContainer` in the `human-vue` mock

View tests mock `@planetcrust/human-vue` by hand. A view wrapped in `CViewContainer` renders nothing unless the mock includes `CViewContainer: { template: '<div><slot /></div>' }`.

### Array-form stub props never receive `true`

Vue casts a valueless attribute to `true` only for a prop declared `type: Boolean`. Stubs here usually declare props as an array (for example `FieldPicker` in `client/web/unify/src/sections/compose/components/PageBlocks/Configurators/fieldPickers.test.js`), so the stub gets `''`.

**Why:** the failure names the assertion, sending you to fix correct code.

**How to apply:** give the stub the object form for that prop (`{ type: Boolean, default: false }`), or assert `toBeTruthy()`.

### Running the unify e2e suite

`client/web/unify` runs playwright (`npx playwright test`) against the live dev stack and needs `client/web/unify/.env.e2e` (gitignored). `make setup-agent` writes it with `E2E_USER=agent@local.dev` and the password from `dev/agent/.state/ui-password`; `worktree.sh new` writes a worktree's own. `auth.setup.ts` timing out on `waitForURL`, with every spec "did not run", is either the auth rate limit (`AUTH_REQUEST_RATE_LIMIT`, 60/min per IP) or a crashed or restarting server; page text "Human server initializing" means the latter.

**How to apply:** check `dev_server_status` before assuming the rate limit. In specs, wait on `getByTestId('app-topbar')`, and keep that wait under the per-test timeout so a slow cold route names its element.

### Two ways an e2e spec goes green without testing

`getByRole('textbox').first()` is the sidebar's resource search, which renders before the route content, so a form filled that way never saves and fails steps later. `locator.count()` does not auto-wait, and `page.goto()` resolves while the SPA still boots, so a `test.skip` on a count of 0 skips the spec and reports green.

**How to apply:** address form controls by id (`input#name`) or a scoped locator, wait for a landmark before counting, and after writing a spec, mutate what it catches and confirm it fails. In `describe.serial`, mutate one fix at a time.

### Playwright empties `e2e/.results` on every run

The unify playwright config sets `outputDir: './e2e/.results'`, wiped at the start of each run. An ad-hoc script importing `@playwright/test` must live in the unify package to resolve it, but not there.

**How to apply:** put such scripts at the package root under a dotted temp name and delete them afterwards.

## Browser checks

### Writing `drive.mjs` checks against fixtures

Reuse a fixture rather than creating one per run: soft-deleting a user does not free its email, and users have no hard delete. Look it up first (`userList({ email, deleted: 1 })`). Address rows by something unique to the probe, or playwright's strict mode meets leftovers from earlier runs. Admin routes live under `/admin` (`/admin/system/users`); an unprefixed path lands on `/` with the topbar present, so `page.open` looks successful.

**How to apply:** check `page.path()` after opening; use `window.__human.globals()` for the app's API clients; clean up yourself, since only namespaces are ledgered.

### The `drive.mjs` page wrapper, `hasText` and scrolling menus

`drive()` hands a wrapper: `page.open()` is on it, but `locator`, `getByRole`, `setViewportSize`, `waitForTimeout` and `screenshot` are on `page.raw`. A `hasText` regex sees concatenated `textContent`, so `<span>Contact</span><Tag>Record page</Tag>` reads `ContactRecord page`; filter with `{ has: p.getByText(title, { exact: true }) }`. In a long `CResourceList`, playwright scrolling to a row's menu item closes the popup menu, so the click times out as "not stable" or "detached".

**How to apply:** narrow a long list with its search box until the row is near the top before opening the row menu.

### A PrimeVue button icon is a `<span>`

`<Button icon="pi pi-cog" />` renders `<span class="p-button-icon pi pi-cog">`, so `button:has(i.pi-cog)` matches nothing and raises no error. Bare `<i class="pi …">` appears only in hand-written markup.

**How to apply:** use `button:has(.pi-cog)`. To see what a page renders: `page.locator('button .pi').evaluateAll(els => [...new Set(els.map(e => e.className))])`.

### The unify sidebar and launcher in headless checks

The left sidebar is a PrimeVue `Drawer` with `data-testid="app-sidebar"` (`.p-drawer-left`). Its open state is per section in `localStorage` `ui.sidebar.expanded`, defaulting to the section's `sidebarExpandedByDefault`, so a fresh context can show zero rows; an open drawer also puts its search box first in the DOM. The app launcher is `.right-sidebar`, not a PrimeVue overlay. `offsetParent` is null for these fixed-position panels, so it is no visibility test.

**How to apply:** set `ui.sidebar.expanded` before navigating (`drive.mjs` has `page.expandSidebar(section)`), scope form selectors to the form, use `getClientRects().length > 0` for visibility, and assert the launcher is shut first and on something only it carries.

### Capturing a record list's own requests

When a `drive.mjs` check records `/record/` requests, the first column-header click is followed by an unrelated `limit=1` request with no `sort`, so "last request wins" reads an empty sort. The prev/next ID list is `limit=50`; a page-2 fetch carries `pageCursor`.

**How to apply:** filter to requests whose `limit` equals the block's perPage.

### A probe that reuses one resource trips domain rules

A loop exercising a control several times against the same resource stops testing the control after the first pass. In the workflow editor a trigger holds one outbound edge (`getMaxOutbound` in `connectionRules.js`), so a three-target probe reports targets 2 and 3 as broken handles.

**Why:** a rule rejection looks exactly like a broken control.

**How to apply:** give each attempt fresh state, and run the probe against the unmodified code before calling anything a defect, ideally before you edit. To undo a teeth-check mutation, restore from a `.bak` copy, never `git checkout` or `git stash`.

### e2e right after `worktree.sh up` fails in auth.setup

Straight after `up` the worktree's API can still be coming up, and `auth.setup.ts` then fails with a bare 30s timeout on `waitForURL(/\/auth\//)`, which reads like a broken login.

**How to apply:** wait for `curl -s http://localhost:<api>/healthcheck` to return 200 and for the webapp to answer before the first run, and re-run once before believing an auth failure.

### `$SystemAPI.userCreate` refuses a user without a group

The generated JS client throws `field userGroupID is empty` before any request is sent, even though the server accepts a user with no group. The user editor sends `'0'` for that reason.

**How to apply:** an e2e spec that creates a user through the API passes `userGroupID: '0'`.

### lib/js runs mocha, not vitest

`npx vitest run` inside `lib/js` fails every suite with `describe is not defined`, which reads like broken tests. The workspace's runner is mocha (`pnpm test`), with `bail` set.

**How to apply:** run it through `dev_test_run lib/js`, which picks the runner and says when bail cut the run short.

### `backlog.sh add --help` files an item called `--help`

`add` takes its first argument as the item text, flags included. The usage is in the script's header comment.

**How to apply:** read `sed -n 1,40p dev/agent/backlog.sh`; drop a stray item with `backlog.sh drop ID`.

### verify-ui calls a lone iframe-backed block empty

On a page whose only block draws its content inside an iframe (Custom, IFrame), `verify-ui.mjs` reports `empty block [single] — rendered nothing` while the screenshot shows it full. The report cannot see into frames.

**How to apply:** for frame-backed blocks, read the screenshot rather than the report line.

### A stale `human-codegen` binary breaks codegen-legacy

`make codegen-legacy` runs `$(GOPATH)/bin/human-codegen`, and `Makefile.inc` builds it only when it is missing, so an old binary is never replaced. An old one fails with `stat …/docs/src/modules: no such file or directory`, or silently skips what was added since — the OpenAPI documents in `server/docs/*.yaml` among them, which `TestOpenAPIDocsAreCurrent` then catches.

**How to apply:** rebuild it before codegen — `cd server && go build -o $(go env GOPATH)/bin/human-codegen ./cmd/codegen/main.go` — then `make codegen` from the root.

### A fixture-dependent e2e spec on an unseeded stack fails as a hook timeout

`corredor-record-scripts.spec.ts` needs the `agent-sandbox` fixture, which the dev database does not carry by default. Its `beforeAll` waits 30s for contact rows, as long as the hook's own timeout, so the run ends in `"beforeAll" hook timeout` before the spec can say what is missing. Corredor's `ENOENT … unify.client-scripts.js` in the same log is unrelated: no fixture defines a `unify` client bundle.

**How to apply:** before attributing a `beforeAll` timeout, check the fixture the spec names (`dev/agent/ids.sh <slug>`), seed it with `dev/agent/seed.sh <slug>`, and re-run that spec alone.

### Editing an embedded file does not rebuild the dev server

devwatch rebuilds on writes to `.go` files only (`server/cmd/devwatch/main.go`), and `touch` is not a write. A skill in `server/system/agentic/skills/library/*.md`, or anything else `go:embed`s, keeps serving its old text until some `.go` file in the checkout is written.

**How to apply:** after editing an embedded file, write a `.go` file in its package with the same content (`cat f.go > f.go` through a copy), wait for `serving` in the server log, and re-read through the tool rather than trusting the file on disk.

### `make codegen` reorders `resource_schema.gen.json`

`server/pkg/codegen/resource_schema.gen.json` comes out in a different order on each run, so one new field shows as a diff of thousands of lines. The file also holds duplicate `"-"` keys, so a JSON load and dump silently merges them and changes an unrelated entry.

**How to apply:** keep the committed file and add the new line by text, next to its sibling field; compare the old and new content as sets of entries, not as text, to see what really changed.

### An empty array kills a `set -u` script on macOS

macOS ships bash 3.2 as `/bin/bash`, where `"${arr[@]}"` on an empty array is an unbound-variable error under `set -u` — the script exits before the command runs. Bash 4.4+ (every Linux box here) expands it to nothing, so the breakage never shows locally.

**How to apply:** expand an array that can be empty as `${arr[@]+"${arr[@]}"}`. Check a script with `docker run --rm -v "$PWD":/w -w /w bash:3.2 bash <script>`.

### `backlog.sh add --help` files an item called `--help`

`add` takes its first argument as the item's text, flags included. The usage is the comment at the top of `dev/agent/backlog.sh`.

**How to apply:** read the header for usage; drop a stray item with `backlog.sh drop ID`.

# Working in this repo

The rules for any coding agent in this checkout. Tool-specific notes for Claude
Code are in `CLAUDE.md`, which imports this file.

## Intent system (opt-in)

The intent system (`.intent/SPEC.md`, `*.intent.md` contract docs) is
**opt-in**: it applies only when the human asks for an intent-governed change or
a drift audit. When not opted in, do not edit `*.intent.md` files or run intent
tooling. `node .intent/intent.mjs check` fails on the baseline and is not run
in CI; `check --changed` is the scoped form.

## Dev server & agent toolkit

- To call the local dev API (verify implementations, inspect data): use
  `dev/agent/` — `smoke.sh` first, `shape.sh compose/record update` for a
  call's method, path, params and traps, then `api.sh METHOD PATH`. Setup:
  `bootstrap.sh` (idempotent). Reference: `dev/agent/README.md`.
- **Never start, stop or restart a dev server** (vite, `make watch`) yourself,
  even one that is stale or misbehaving: it may be another session's. Say what
  is wrong and let the human relaunch it.
- **Non-obvious behaviour by area lives in `dev/gotchas/`** (server, compose,
  TAQ, agentic, project, frontend, dev-toolkit, intent). Read the area file
  before working in it; add an entry when something cost a turn. `docs/` is the
  product documentation for people who use, deploy or configure Human; nothing
  developer-facing goes there.
- Test data: versioned fixtures in `dev/fixtures/`; `dev/agent/seed.sh` /
  `cleanup.sh`. Pages/charts are built via `dev/agent/pagebuild.py` (never
  via envoy YAML — block refs don't resolve).
- **Only `dev/agent/api.sh`, `dev/agent/mcp.py` and this checkout's own MCP
  server reach the local dev server.** Any other MCP server offering Human
  tools is a remote instance: it writes to somebody's live data, and nothing it
  creates is in the cleanup ledger. Use one only when the remote instance is
  the point.
- `api.sh` and `mcp.py` record the **namespaces** a session creates in
  `dev/agent/.state/created.jsonl`, and `cleanup.sh` deletes only what that
  ledger holds for the current session. Anything else you create (a user, a
  role, a workflow) is yours to delete by hand. No name prefix is needed or
  wanted; data the session did not create is off-limits, whatever it is called.
- **Testing is scripts and code.** Never build or propose a check that runs a
  model; a skill or a tool description is held to the code by a deterministic
  test and reviewed by hand.

## Conventions

- **Formatting**: prettier (FE) / gofmt (Go) on changed files only, after edits.
- **Comments** (code and templates): short, and about what the thing **is** —
  never why it changed. No "used to", no "now that", no dated rulings, no
  pointer to the decision that produced the line. History is git's job and
  rationale is the intent doc's; a comment earns its place only by saying
  something the code next to it cannot.
- **Task isolation**: a task that runs the server or spans files gets its own
  checkout, server, webapp and database — `dev/agent/worktree.sh new <name>`,
  then `up`/`down`/`rm`. The primary is slot 0 and is never reassigned.
- **Deferred work**: `dev/agent/backlog.sh add` queues anything left undone
  that should outlive the turn. An item is **local** — this checkout's todo —
  only when doing it would change a file the task is already changing.
  Everything else is `--global`: filed, named in the report, and read only when
  the human asks for the pool.
- **Questions**: a hand-back costs the human's turnaround, not yours — ~8 min
  median, more than a full e2e run. So never ask what you could have priced:
  every option carries its cost, its effort, and your confidence in that
  number. Batch questions into rounds rather than asking as things occur to you.
- **Commits**: see **[Commit convention](#commit-convention)** below — the one
  statement of it.
- **Teeth checks**: undo a mutation from a `.bak` copy you took first, never
  `git checkout <file>` (reverts your own edits to HEAD) or `git stash` (sweeps a
  peer session's in-flight files along with yours).
- **Editors open clean**: an editor for an existing resource must not start
  dirty. Resolve defaults for display (a computed with a fallback getter), never
  write them into the draft on load; only create forms may start dirty.
- **Never truncate text, wrap it**: `whitespace-normal break-words` inside a
  bounded width, not Tailwind's `truncate`. An ellipsis hides exactly what
  distinguishes two similar options; dense ID or filename chips are the one
  arguable exception.
- **Tests**: run the touched package's suite (e.g. `cd client/web/unify && npx
vitest run`) before committing behavior changes. Baseline every suite you will
  verify with, e2e included, **before** editing — a failure first seen
  afterwards costs a revert and a second run to attribute, and an e2e suite is
  minutes. When one does surface late, revert and re-run **that spec alone**;
  never the whole set twice.
- **i18n**: single merged `human-webapp` bundle — edit `locale/en/human-webapp/`;
  prefer verbose, explicit wording.
- **Codegen**: `make codegen` (server first, then lib). Never hand-edit
  generated files (`*.gen.*`, `lib/js/src/api-clients/`, `server/docs/*.yaml`,
  `docs/public/schemas/`, `docs/llms*.txt`); regenerate and commit.
- **Vocabulary**: TAQ = Trigger Action Query (never "Task Queue"); the glossary
  is `docs/get-started/glossary.md`.
- **Compose routes are frozen** — renames require redirects (see
  `sections/compose/compose.intent.md`).

## Commit convention

The whole rule lives here.

- **Author**: the human user. No `Co-Authored-By`, no generated-with line, no
  robot emoji, anywhere in the message.
- **Message**: one imperative subject, capitalised, under 72 characters, no
  trailing period — usually the whole message. A body is at most **one short
  sentence naming what changed**, for when the subject cannot hold it. Never
  prose, never the reasoning: why belongs in the intent doc, the diff shows
  what.
- **Atomic**: docs, bugfix and cleanup are separate commits. A product fix and
  a tooling fix are two commits even when made in one sitting. Name the files
  rather than sweeping the tree — it often holds someone else's in-flight edits.
- **Except intent**: `*.intent.md` and `.intent/intent.lock.json` go in the
  **same** commit as the code they govern, the way tests do. A doc states what
  its change made true, so splitting them leaves history self-contradicting.
  Reconcile the doc, `node .intent/intent.mjs sync <files>`, commit the lot.
- **Except translations**: `locale/` strings go in the **same** commit as the
  code that reads them. A commit adding a key nothing uses, or using a key it
  has not added, is broken at that commit.
- **Except generated references**: a change that moves `server/docs/*.yaml`,
  `docs/public/schemas/`, `docs/llms*.txt`, `TOOLS.md` or `mcp-tools.gen.md`
  carries the regenerated file in the same commit; a test fails otherwise.
- **When**: once the work is verified _and_ nothing about it is still open — no
  question waiting on the human, no doc flagged and left, no check named and
  skipped. Verified is not the same as finished, and a commit made over a live
  thread is one the next turn has to amend. Judge the scope: something
  unrelated that surfaced is the next task, not an open thread. Commit then,
  without being asked.
- **Never push, never touch the tracker.** A local commit is reversible;
  publishing is not. Pushes, PRs, and closing or commenting on an issue are the
  human's call, every time — a fix landing locally is not the issue's verdict,
  so don't ask about closing it either. Naming the issue a commit addresses is
  enough.

## Layout

- `client/web/unify` — the webapp (sections: home, agentic, workflow, taq,
  admin, compose, chatbot, project); `lib/vue` + `lib/js` — shared libraries
  (one-way dependency: apps import libs, never the reverse); `server` — Go
  backend (documented by intent docs, not enforced); `def/` — cue codegen
  sources; `dev/` — developer tooling and notes; `docs/` — the product site.

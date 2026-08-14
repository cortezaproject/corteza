# Working in this repo

## Intent system (opt-in)

The intent system (`.intent/SPEC.md`, `*.intent.md` contract docs) is
**opt-in**: it applies only when the human invokes `/intent-task` (changes) or
`/intent-audit` (drift audits) — those skills carry the full procedure. When
not opted in, do not edit `*.intent.md` files or run intent tooling.

## Dev server & agent toolkit

- To call the local dev API (verify implementations, inspect data): use
  `dev/agent/` — `smoke.sh` first, then `api.sh METHOD PATH`. Setup:
  `bootstrap.sh` (idempotent). Skills: `/dev-api`, `/dev-seed`. Reference +
  gotchas: `dev/agent/README.md`.
- Test data: versioned fixtures in `dev/fixtures/`; `dev/agent/seed.sh` /
  `cleanup.sh`. Pages/charts are built via `dev/agent/pagebuild.py` (never
  via envoy YAML — block refs don't resolve).
- Building whole systems (datamodel + pages) in Human: `/sys-design` skill.
- Everything an agent creates on the dev server is recorded in
  `dev/agent/.state/created.jsonl`, and `cleanup.sh` deletes only what that
  ledger holds for the current session. No name prefix is needed or wanted;
  data the session did not create is off-limits, whatever it is called.

## Conventions

- **Formatting**: prettier (FE) / gofmt (Go) on changed files only, after edits.
- **Comments** (code and templates): short, and about what the thing **is** —
  never why it changed. No "used to", no "now that", no dated rulings, no
  pointer to the decision that produced the line. History is git's job and
  rationale is the intent doc's; a comment earns its place only by saying
  something the code next to it cannot.
- **Commits**: see **[Commit convention](#commit-convention)** below — the one
  statement of it. Skills and `dev_commit_create` point here rather than
  restating it.
- **Tests**: run the touched package's suite (e.g. `cd client/web/unify && npx
vitest run`) before committing behavior changes.
- **i18n**: single merged `human-webapp` bundle — edit `locale/en/human-webapp/`;
  prefer verbose, explicit wording.
- **Codegen**: `make codegen` (server first, then lib). Never hand-edit
  generated files (`*.gen.*`, `lib/js/src/api-clients/`).
- **Vocabulary**: TAQ = Trigger Action Query (never "Task Queue").
- **Compose routes are frozen** — renames require redirects (see
  `sections/compose/compose.intent.md`).

## Commit convention

The whole rule lives here. `dev_commit_create` enforces the mechanical half and
warns on the rest; `/dev-change`, `/dev-task` and `/intent-task` point here.

- **Author**: the human user. No `Co-Authored-By`, no generated-with line, no
  robot emoji, anywhere in the message.
- **Message**: one imperative subject, capitalised, under 72 characters, no
  trailing period. Body only when the why isn't obvious from the diff — 2–3
  lines, never more.
- **Atomic**: docs, bugfix and cleanup are separate commits. A product fix and
  a tooling fix are two commits even when made in one sitting. Name the files
  in `files` rather than sweeping the tree — it often holds someone else's
  in-flight edits.
- **Except intent**: `*.intent.md` and `.intent/intent.lock.json` go in the
  **same** commit as the code they govern, the way tests do. A doc states what
  its change made true, so splitting them leaves history self-contradicting.
  Reconcile the doc, `node .intent/intent.mjs sync <files>`, commit the lot.
- **When**: once the work is verified _and_ nothing about it is still open — no
  question waiting on the human, no doc flagged and left, no check named and
  skipped. Verified is not the same as finished, and a commit made over a live
  thread is one the next turn has to amend. Judge the scope: something
  unrelated that surfaced is the next task, not an open thread. Commit then,
  without being asked.
- **Never push.** A local commit is reversible; publishing is not. Pushes and
  PRs are the human's call, every time.

## Layout

- `client/web/unify` — the webapp (sections: home, agentic, workflow, taq,
  admin, compose, chatbot, project); `lib/vue` + `lib/js` — shared libraries
  (one-way dependency: apps import libs, never the reverse); `server` — Go
  backend (not yet intent-covered); `def/` — cue codegen sources.

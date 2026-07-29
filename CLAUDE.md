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
- Anything an agent creates on the dev server uses an `agent-` handle/slug
  prefix — prefixed data is disposable, unprefixed data is off-limits.

## Conventions

- **Formatting**: prettier (FE) / gofmt (Go) on changed files only, after edits.
- **Commits**: as the human user — no AI co-author trailers. Atomic logical
  commits (docs, bugfix, cleanup separately). Imperative subject.
- **Tests**: run the touched package's suite (e.g. `cd client/web/unify && npx
  vitest run`) before committing behavior changes.
- **i18n**: single merged `human-webapp` bundle — edit `locale/en/human-webapp/`;
  prefer verbose, explicit wording.
- **Codegen**: `make codegen` (server first, then lib). Never hand-edit
  generated files (`*.gen.*`, `lib/js/src/api-clients/`).
- **Vocabulary**: TAQ = Trigger Action Query (never "Task Queue").
- **Compose routes are frozen** — renames require redirects (see
  `sections/compose/compose.intent.md`).

## Layout

- `client/web/unify` — the webapp (sections: home, agentic, workflow, taq,
  admin, compose, chatbot, project); `lib/vue` + `lib/js` — shared libraries
  (one-way dependency: apps import libs, never the reverse); `server` — Go
  backend (not yet intent-covered); `def/` — cue codegen sources.

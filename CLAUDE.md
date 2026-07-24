# Working in this repo

This repo runs on the **intent system** (`.intent/SPEC.md` is the constitution,
`.intent/TODO.md` the outstanding work). Every covered folder has a
`<foldername>.intent.md`, load-bearing files have `<name>.intent.md` sidecars.
These docs are verifiable contracts, not notes.

## The law (every session, every task)

1. **Before editing a covered file**: read its governing intent doc — the
   sibling `<name>.intent.md` if present, else the nearest folder doc up the
   tree. Section-level constitutions (e.g. `sections/project/project.intent.md`)
   record human-ruled **locked contracts** and **WIP zones**.
2. **Locked contracts change only by human ruling** — ask via AskUserQuestion
   (interview style: targeted questions, concrete options) before touching one.
   `> **WIP:**` zones: do not rely on them; interview before building on them.
   `> **DRIFT:**` notes: recorded intent the code doesn't meet yet.
3. **New features are intent-first**: update the intent docs, get confirmation,
   then implement.
4. **After changing covered code**: reconcile the governing doc (update it, or
   confirm intent is unchanged), then `node .intent/intent.mjs sync <file>`.
   `node .intent/intent.mjs check` must pass before committing — pre-commit and
   CI enforce it. Fresh clone: `make intent-hooks` once.
5. Use `/intent-task` for changes and `/intent-audit` for drift audits — they
   encode the standard procedures; don't improvise your own flow.

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

# AGENTS.md — bridge to the intent system

This repo uses the **intent system** (see `.intent/SPEC.md`): every covered folder
carries a `<foldername>.intent.md` and load-bearing files carry `<name>.intent.md` sidecars.
These are verifiable contracts — intention, data touched, tests — not just notes.

Rules for any AI agent working here:

1. Before changing a file, read its governing intent doc: the sibling
   `<name>.intent.md` if present, else the nearest folder doc (`<foldername>.intent.md`) up the tree.
2. After changing covered code, reconcile the doc (update it, or confirm intent is
   unchanged), then run `node .intent/intent.mjs sync <file>`.
3. `node .intent/intent.mjs check` must pass before committing (CI enforces it).
   Fresh clones: run `make intent-hooks` once to activate the pre-commit check.
4. New features are intent-first: update the intent docs, then implement.

Build/test basics: pnpm workspace (`pnpm install` at root), web apps under
`client/web/*` (Vite), Go server under `server/`, `make codegen` after changing
cue defs or rest.yaml. Format changed files with prettier (FE) / gofmt (Go).

---
kind: folder
covers: recursive
owner: fe
depends-on: []
touched-by:
  - client/web/unify/playwright.config.ts
tests: []
---

# E2E suite (Playwright)

## Intention

Browser-level verification of the intent docs' UX contracts. Specs mirror the
section tree (`e2e/sections/<section>/…`) so every intent doc's `tests:` field
can point at the spec covering its screen; `intent affected` maps changed files
to the specs to run.

## Contract

- Specs attach to a running dev stack (`E2E_BASE_URL`, default localhost:5173)
  — they never boot or seed anything, and must tolerate existing dev data
  (skip, don't fail, when required data is absent).
- Auth: `auth.setup.ts` logs in once via the real auth form (credentials from
  gitignored `.env.e2e`) and shares the session via storageState.
- Serial workers by design: the stack's data is shared, runs stay deterministic.
- A spec asserts the CONTRACT of a screen (what its intent doc locks), not
  implementation details — prefer roles/labels over CSS internals.

## When changing this

- New specs must be referenced from the `tests:` field of the intent doc(s)
  they verify — an unreferenced spec is invisible to `intent affected`.
- CI cannot run these yet (no boot/seed story) — see .intent/TODO.md.

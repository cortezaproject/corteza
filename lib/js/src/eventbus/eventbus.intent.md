---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - lib/js/src/guards.ts
touched-by:
  - client/web/unify
  - lib/vue
tests:
  - lib/js/src/eventbus/constraints.test.ts
  - lib/js/src/eventbus/eventbus.test.ts
  - lib/js/src/eventbus/handlers.test.ts
---

# EventBus (browser-side)

## Intention

Client-side bus for Human events — runs in the browser, never on the corredor
server. Two flows: (1) corredor-built client-script bundles register handlers
against it; (2) explicit server scripts are registered as wrapper handlers
that forward execution to the API (and on to corredor).

## Contract

- `Dispatch(ev, script?)` is the single entry point: without `script` it fires implicit event handlers; with one it runs that manually invoked (`onManual`) script by name.
- Handlers match on resource + event + constraints; constraint values support minimatch globs (`constraints.ts`), mirroring server-side trigger constraint semantics.
- Matching handlers run ordered by ascending `weight`; a manual dispatch additionally matches on `scriptName`.

## When changing this

- Constraint matching must stay in sync with server trigger semantics — a client-only divergence silently changes which automation runs.
- Apps register handlers at bootstrap; changing registration or dispatch signatures touches every web app's automation wiring.

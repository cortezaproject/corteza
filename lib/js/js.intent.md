---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - server
touched-by:
  - lib/vue
  - client/web/unify
  - client/web/chatbot-widget
tests:
  - lib/js/src/cast.test.ts
  - lib/js/src/guards.test.ts
---

# @planetcrust/human-js — shared JS/TS library

## Intention

Framework-agnostic core library for all Human web clients: resource classes
(compose, system, automation), generated API clients, validation, formatting,
the browser eventbus, and corredor scripting support. `lib/vue` and the web
apps consume it (`workspace:*`); dependency is strictly one-way — nothing in
here may import from `lib/vue` or any app.

## Map

- `src/index.ts` — public surface: namespaced exports (`compose`, `system`, `automation`, `apiClients`, `eventbus`, `corredor`, `validator`, `shared`, `fmt`) plus `NoID` and `renderMarkdown`.
- `src/cast.ts` — casting helpers (`Apply`, `HumanID`, `ISO8601Date`, `NoID`); every resource class builds on these.
- `src/guards.ts` — structural type guards (`IsOf`, `AreObjectsOf`, …).
- `src/markdown.ts` — `renderMarkdown`: markdown-it + DOMPurify pipeline shared by all chat surfaces (agent chat, chatbot inbox, embeddable widget).
- `src/api-clients/` — CODEGEN-OWNED (see below); each subsystem's axios client class.
- Subfolders each carry their own intent doc (validator, corredor, compose, system, automation, eventbus, formatting, shared).

## api-clients regeneration rule

`src/api-clients/` is generated from the server's `rest.yaml` definitions via
`tools/codegen/human-api-client.js` (`pnpm codegen` here, driven by root
`make codegen`). Never hand-edit those files and never document them in the
intent system. Chain order matters: server codegen runs first (finalizing
`rest.yaml`), then lib codegen, so clients always derive from the final spec.

## When changing this

- `NoID = '0'` is a backend contract: IDs are uint64 serialized as strings; never treat IDs as numbers.
- Anything exported from `src/index.ts` is public API for three apps plus lib/vue — removing/renaming exports needs a consumer sweep.
- Endpoint changes start in server `rest.yaml` + codegen, never in `src/api-clients/`.

---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/plugins/auth.ts
touched-by:
  - client/web/unify/src/App.vue
tests: []
---

# Libs

## Intention

Non-Vue infrastructure primitives: URL construction shared by auth/API
plugins, and the realtime websocket client the shell uses to receive server
push (notifications, workflow prompts, reminders).

## Map

- url.ts — `Make({ url, query, hash, ref })`: builds absolute URLs from relative/schemaless inputs with qs-style query serialization.
- websocket.ts — `endpoint()` derives the ws(s) URL from `window.HumanWebsocket`/`HumanAPI`; `RealtimeClient` authenticates with the access token, auto-reconnects (bounded attempts), re-auths on token refresh.

## When changing this

- `RealtimeClient` message payloads are `{'@type', '@value'}` envelopes;
  consumers (notifications store, shell event bus) dispatch on `@type` — keep
  the envelope stable.
- endpoint derivation must keep working with only `window.HumanAPI` set
  (no explicit `HumanWebsocket`); protocol upgrades follow the page
  (`https` → `wss`).
- `Make` is what auth.ts uses to build callback/flow URLs — query
  serialization changes ripple into the OAuth2 flow.

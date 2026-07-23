---
kind: folder
covers: recursive
owner: fe
depends-on:
  - client/web/unify/src/sections/index.js
touched-by:
  - client/web/unify/src/plugins/index.js
tests: []
---

# Router

## Intention

One `vue-router` instance for the whole app, assembled purely from section route
declarations — the router itself owns no feature routes.

## Contract

- Routes come only from `sections/index.js` (`routes` export); sections prefix their own paths and namespace their route names.
- Unknown paths redirect to `/` (home section) — never a 404 page.
- History mode with `BASE_URL`; deep links must survive the async boot in `main.js`.

## When changing this

- Adding feature routes here (instead of in a section) violates the section contract.
- Route-name collisions across sections surface here first — names are namespaced by convention only, nothing enforces them.

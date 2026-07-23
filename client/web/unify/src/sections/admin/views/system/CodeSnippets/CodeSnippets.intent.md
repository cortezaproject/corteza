---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
touched-by: []
tests: []
---

# CodeSnippets

## Intention

Single settings page where an admin maintains code snippets injected into
the platform (e.g. custom scripts/markup). Not a CRUD resource — snippets
are stored as system settings.

## Data touched

- `$SystemAPI.settingsList({ prefix: 'code-snippets' })` /
  `settingsUpdate` — flat key/value settings, no resource class; all
  snippets live in the single `code-snippets` setting value.

## Map

- `Index.vue` — the one snippets screen (route target, own sidecar).

## When changing this

- The `code-snippets` settings key is the contract with the server/consumer
  side; renaming it orphans existing stored snippets.

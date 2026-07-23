---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify
  - client/web/chatbot-widget
tests: []
---

# lib/vue — @planetcrust/human-vue

## Intention

The shared Vue 3 library for all Human web apps: components, composables,
Pinia stores, and app-bootstrap plugins that must behave identically across
apps. Consumed as source (`main: src/index.ts`), no build step.

## Core invariant

Dependency is strictly one-way: apps import from lib/vue; lib/vue NEVER
imports from any app. Anything app-specific reaches components via props,
provide/inject (`$SystemAPI`, `$ComposeAPI`, `$Auth`, `$Settings`, `$toast`,
...), or translation objects — never via direct app imports.

## Map of src/

- `components/` — shared components (see its own intent docs)
- `composables/` — reusable composition fns (useResourceList, usePermissions, useTheme, ...)
- `stores/` — shared Pinia stores (user, record, module, notifications, agent chat, ...)
- `plugins/` — app bootstrap: auth, human-api, i18n, settings, toast, event-bus, primevue-components
- `filters/` — date formatting helpers
- `libs/` — url + websocket helpers
- `utils/` — app icons, internal navigation helpers
- `assets/`, `test/` — static assets, test setup
- `index.ts` — package entry re-exporting the public surface; `vue.d.ts` — ambient types

Subfolders carrying their own `*.intent.md` override this doc for their subtree.

## When changing this

Exports from `src/index.ts` are cross-app API — removing or renaming one is a
breaking change for every consuming app. New externals must go into
package.json dependencies (apps do not hoist for the lib).

---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/compose/views/Admin/Modules/List.vue
tests: []
---

# Module components

## Intention

Non-route helpers for module administration. Currently only the module import
dialog used from the modules list.

## Data touched

Creates modules in the current namespace via the Compose API from an uploaded
export file.

## Map

- ModuleImporter.vue — upload/parse a module export, preview, and import into the namespace; emits so the list view can refresh

## When changing this

Import must stay compatible with the export format produced elsewhere in the
product (namespace/module export). Failure states should be visible per module
rather than aborting the whole import silently.

> **DRIFT:** they are not. `ModuleImporter` creates modules in one loop inside
> a single try/catch, so the first failure stops the remaining modules, shows
> one generic toast and closes the dialog — leaving a half-imported namespace
> with no indication of which module failed.

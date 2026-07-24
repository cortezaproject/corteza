---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/components/permissions
touched-by:
  - client/web/unify/src/sections/project/views/Wizard.vue
tests: []
---

# Page dialogs

## Intention

The page resource-kind's implementation of the project resource dialog
standard: CreateDialog + DetailDialog, rendered once by the Wizard and opened
via `createResource('page')` / `inspectResource('page', id)`. The dialogs own
only page metadata; real page building happens in the full admin page builder,
reached by deep link.

## Data touched

- `useProjectsStore`: `pagesFor`, `addPage`, `updatePage`.
- Router: resolves `admin.pages.builder` (`slug` = project `namespaceID`) for
  the "open builder" links.

## Map

- `PageCreateDialog.vue` — name-only create with two commit buttons: plain
  Create, and Create-and-open-builder. The builder tab is opened synchronously
  inside the click (before the await) so popup blockers don't eat it; it is
  pointed at the builder once the page exists, or closed on failure. Never
  auto-opens the detail dialog.
- `PageDetailDialog.vue` — staged name/description/visible (committed on
  Save); embeds `ResourcePermissionsSection` (`kind="page"`); footer
  deep-links to the builder (new tab) when the project has a namespace.

## When changing this

- Keep the synchronous `window.open` pattern when touching the create flow —
  it exists solely to survive popup blocking.
- The builder links are gated on `project.hasNamespace`; a project without a
  provisioned namespace must degrade to no link, not a broken one.
- Dialogs open only through the Wizard's provides (dialog standard).

---
kind: file
covers: Edit.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/usePageStore.js
  - lib/vue/src/stores/usePageLayoutStore.js
  - client/web/unify/src/sections/compose/components/Admin/Page/PageTranslator.vue
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests:
  - client/web/unify/src/sections/compose/views/Admin/Pages/Edit.layout-title.test.js
---

# Page Edit view

## Intention

Edit everything about a page except its blocks: metadata (title, handle,
description), navigation presentation (icon, visibility, sub-page expansion,
record notifications), and the page's set of layouts with their conditions and
record-toolbar configuration. Blocks are the builder's job.

## UX capabilities

- Title required, handle validated; save gated by `canUpdatePage`; unsaved guard covers both page and layout edits.
- Icon dialog: upload to the shared icon library, pick from it, or use an external link; the choice persists via `pageUpdateIcon` on save.
- Layouts (edit mode only): add/remove/drag-reorder rows, jump to the builder per layout; a config dialog edits title (`useTitle` turns that title into the page's own, interpolated — hence the `ƒ` affordance; `RecordView` renders it against the open record, `View` against the signed-in user), visibility expression + roles, record-toolbar button toggles, and custom actions (`toLayout` / `toURL` with placement/variant/open-in).
- Deleting a page with children forks into rebase (children move up) or cascade; leaf pages get a plain confirm.
- Save-as-copy (non-record pages) clones the page and seeds a `primary` layout from its blocks; create auto-creates an empty `primary` layout then redirects to edit.
- Layout persistence on save: removed layouts deleted first, new ones created, `_updated` ones updated, then reordered via `pageLayoutReorder` and re-fetched.

## Routes

`admin.pages.create` at `admin/pages/create`, `admin.pages.edit` at `admin/pages/:pageID/edit`. Topbar (edit): module edit (record pages), builder, view page, translator. Load failure bounces to `admin.pages`.

## When changing this

- Seeded layouts must use the `compose.PageLayout` type (default toolbar buttons enabled) and — when cloning — carry the page's blocks; an empty layout hides every block in the public view.
- The layout config dialog edits a deep clone; nothing applies until its Save splices the clone back and marks `_updated`.
- Layout builder links pass `?layoutID`; deleting layouts is deferred to the main save (they accumulate in `removedLayouts`).
- `meta.title` does double duty: the layout's name in this editor (required) and, when `useTitle` is on, the page's interpolated title. Changing one changes the other — so both places that edit it, the list row and the config dialog, swap to the expression input on that flag, with one hint under the list serving the row last focused.
- `CInputExpression` carries the required-title `invalid` state itself; this list checks its own rule against `validationTriggered` rather than through the Form resolver, so the flag has to reach whichever input is rendered.
- `useTitle` is offered on every page kind; what varies is the scope `useExpressionScope` hands the input and its hint — record variables only where there is a record.

- The record-page condition hint (`record.values.fieldName`) is accurate only because `RecordView.vue`/`View.vue` resolve the record before picking a layout (see their sidecars) — keep hint copy and actual evaluation timing in sync.

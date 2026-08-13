---
kind: file
covers: Edit.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/PageBlocks/Grid.vue
  - lib/vue/src/stores/useModuleStore.js
  - lib/vue/src/stores/useRecordStore.js
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Admin Record View/Edit view

## Intention

Inspect or edit a single record from the admin area through a synthetic
all-fields Record block. One component serves both the view and edit routes;
the mode is derived from the matched route name.

## UX capabilities

- View mode: back, delete (`canDeleteRecord`), edit (`canUpdateRecord`); edit mode: cancel (back to view), delete, save.
- Edit works on a clone of the pristine record; cancel and view↔edit transitions restore/re-clone without refetching. Navigating away while in edit mode always prompts an unsaved-changes confirm — there is no dirty comparison, so it prompts even when nothing was touched.
- Required-field validation, server field errors mapped onto the form, pending file uploads flushed before update.
- Delete replaces to the admin record list; save replaces to the view route.
- Topbar links to the module editor and admin record list.

## Routes

`admin.modules.record.view` at `admin/modules/:moduleID/records/:recordID` and `admin.modules.record.edit` at `.../edit` — same component; `isEditMode` keys off `route.name === 'admin.modules.record.edit'`. Reloads only when module or recordID changes (not on mode switch).

## When changing this

- Records are loaded with `force: true` so stale store entries never render.
- Renaming either route name breaks the mode detection and the post-save/cancel `router.replace` targets.
- Keep `recordViewContext` / `$fileUploadContext` shapes in sync with `Pages/RecordView.vue`.

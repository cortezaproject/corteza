---
kind: file
covers: Builder.vue
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/PageBlocks/Grid.vue
  - client/web/unify/src/sections/compose/components/PageBlocks/Configurators
  - lib/vue/src/stores/usePageStore.js
  - lib/vue/src/stores/usePageLayoutStore.js
  - client/web/unify/src/sections/compose/lib/resource-translations.ts
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Page Builder view

## Intention

Visually compose one page layout at a time: place, size, configure, clone, and
remove blocks on the same grid the public view renders, so the builder shows
exactly what users will see. Blocks belong to the page; the layout only
positions a subset of them.

## UX capabilities

- Editable `Grid` with per-block overlay: drag handle, edit, clone (offset below, blank ID), translate (saved blocks only), remove-from-layout (confirm; block stays on the page as an orphan).
- Layout switcher in the topbar when a page has several layouts; layout-level actions: save, save-as-copy, delete (falls back to the next layout or the "add layout" empty state), create when the page has none. Switching layout while there are staged edits confirms first — it rebuilds the working set from the saved page, discarding exactly what leaving would.
- Add-block dialog: block-type tiles (record-page-only kinds shown contextually) plus the page's orphan blocks, which can be re-placed or permanently deleted from the page (removes them from every layout too).
- New blocks are staged: they join the page/layout only on the first save of their configurator; cancelling discards them entirely.
- Block configurator dialog with three tabs: General (title, description, custom ID/CSS class, header style, magnify, wrap/border, refresh rate, visibility expression + roles), a per-kind configurator (registry of ~18 kinds; draft shared via `blockDraft` provide), and Automation (RecordList selection buttons only).
- Tabs blocks manage child blocks through the provided `$pageBuilder` API (edit/remove/clone tab entries); children are kept `meta.hidden` in sync so they render only inside their tab.
- Saving a record page is blocked (error toast, no request) while required module fields are not covered by any Record block.
- Leaving with staged edits — a block added, configured, cloned, removed, dragged or resized, or an orphan deletion staged — confirms first; an untouched builder leaves silently.

## Routes

`admin.pages.builder` at `admin/pages/:pageID/builder` under `namespace.view`. Topbar: module edit (record pages), view page (record pages open create-record), page edit, translator.

## When changing this

- Save contract: page blocks = preserved (not in this layout) first + working set tail; new blockIDs are mapped back by index from the response, Tabs children remapped in a second save if needed, then the layout is updated with `{blockID, xywh}` only. Breaking the preserved/tail order corrupts the ID mapping and can orphan or delete other layouts' blocks.
- Unsaved blockIDs are NoID (`'0'`); identity falls back to `meta.tempID` — all lookups must go through that fallback.
- `?layoutID` deep link (passed by the page editor's layout-builder buttons) selects that layout on load; falls back to the previously selected, then first, layout when absent or unknown.
- Deleting an orphan block is staged like every other edit: applied on Save (dropped from the page payload and every referencing layout), discarded when leaving without saving.
- The dirty baseline (`useDraftGuard`) is captured on load, on layout switch and after a save reloads the page; staged orphan deletions ride in `extra`. Block positions are part of the compared draft — gridstack writes `xywh` back into the working set, so a drag is caught by the same comparison as a config edit.
- The layout `Select` is remounted (`layoutSelectKey`) when its confirm is refused: it keeps its own copy of the selection, so a switch backed out of would otherwise leave it displaying a layout the builder is not editing.

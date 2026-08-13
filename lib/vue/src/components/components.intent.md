---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify
tests: []
---

# components/ root

## Intention

Barrel + the few components too small to warrant a family folder. Each family
subfolder (agent, field, navigation, ...) carries its own intent doc; this doc
governs only `index.ts` and the loose root components.

## Map

- `index.ts` — the public component export surface. Everything re-exported here
  (family barrels + loose components + `emojiData`) is cross-app API; removing
  an export is a breaking change. New shared components must be added here to
  become importable by apps.
- `CEmojiPicker.vue` — virtual-scrolled emoji grid with search, quick reactions,
  and a frequently-used row persisted in localStorage (`human:emoji:frequently-used`).
- `CEmptyState.vue` — slot-only bordered placeholder for empty lists/tables;
  callers own the v-if deciding when it shows.
- `CResizeHandle.vue` — vertical drag handle emitting `mousedown`; pairs with
  useRightSidebarResize's drag-start handler.
- `CViewContainer.vue` — standard page wrapper for section views. Two shapes:
  default = fixed-height, content scrolls internally; `scroll` = column that
  grows and scrolls as a whole (edit forms), with `gap` (4 or 5).

## When changing this

This doc `covers: '.'`, so a family subfolder without its own doc (`tag/`,
`chip/`) falls through to `lib/vue/vue.intent.md`, not to here. Keep `index.ts`
additive; deprecate before deleting exports.

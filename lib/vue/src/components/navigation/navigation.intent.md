---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/plugins/settings.ts
  - lib/vue/src/composables/useInternalLink.ts
touched-by: []
tests: []
---

# navigation/

## Intention

The app chrome every Human app shares: topbar, sidebar shell + nav tree, app
switcher, and toolbar. Behavior is driven by admin-configured `ui.*` settings
so branding/visibility changes need no code.

## Map

- `CTopbar.vue` — header: title/tools teleport targets, page buttons, home /
  app-selector / agent / notifications / help / profile menus.
- `CSidebar.vue` — collapsible sidebar shell (v-model expanded, `$Settings` logo).
- `CSidebarNav.vue` / `CSidebarNavItem.vue` — generic nav tree; item id/parent/
  label/icon/weight fields are configurable via `*Key` props (defaults are the
  compose-page shape: `pageID`/`selfID`/`title`/`weight`); emits `select`.
- `CAppList.vue` / `CAppListSidebar.vue` — application switcher (grid / sidebar).
- `CToolbar.vue` — slot-only bottom bar (`start`/`center`/`end`).
- `CRouterLinkButton.vue` — Button-styled router-link.

## Contracts consumers rely on

- CTopbar consumes `ui.topbar` settings via `$Settings.get('ui.topbar')`,
  shallow-merged under the `settings` prop (prop wins): `hide*` flags,
  `profileLinks`, and `pageButtons`. Page buttons need `label`+`url` and are
  urlMatch-filtered: `urlMatch` shows the button when the current href contains
  it; no `urlMatch` → shown only at root path. Kept reactive across SPA
  navigation by patching `history.pushState/replaceState` (restored on unmount).
- `#tools` slot and the `#topbar-title` / `#topbar-tools` div ids are teleport
  targets apps mount page titles/actions into — ids are API.
- All string labels come in via the required `labels` prop (no i18n inside).

## When changing this

Chrome is shared by every app: test topbar/sidebar in unify at minimum, and
against settings-driven hide flags (each `hide*` toggle is admin-reachable).

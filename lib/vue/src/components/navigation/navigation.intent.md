---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/plugins/settings.ts
  - lib/vue/src/composables/useInternalLink.ts
touched-by: []
tests:
  - lib/vue/src/components/navigation/CSidebarNav.expand.test.ts
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
  A trailing marker comes from `badgeKey` (`{ value, severity?, title? }`),
  overridable through the `#badge` slot — the tree never learns what the
  marker means (project's status chip is the live example).
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
- Topbar rows wrap content-driven (ruled 2026-07-24): the title keeps its
  natural width and the tools + right-icon cluster drops to its own
  right-aligned row when the bar runs out of room — the title truncates only
  past a full row's width. There is no forced layout switch at any width: the
  ≤550px container query (which measures the topbar's own width — window
  minus sidebar — not the viewport) applies mobile trims only: `#topbar-tools`
  hidden and a smaller title font. The bar is `min-height` sized, so wrapping
  grows it; fixed-position elements offset by `--topbar-height` (right
  sidebar, toasts, search overlay) accept the overlap on wrapped bars.
- All string labels come in via the required `labels` prop (no i18n inside).
- The chrome carries `data-testid` landmarks — `app-sidebar` on the sidebar
  drawer, `app-topbar` on the header — because everything else about them is a
  styling detail: the sidebar is a PrimeVue Drawer, so any structural selector
  for it is really a guess about PrimeVue's markup.
- `expandAll` is a standing instruction, not a starting state. Navs are fed by
  stores that are still loading at setup, so a group opens when its items
  arrive, and again whenever the prop turns back on (search toggles it); a
  group the user collapsed stays collapsed until then.

## When changing this

Chrome is shared by every app: test topbar/sidebar in unify at minimum, and
against settings-driven hide flags (each `hide*` toggle is admin-reachable).

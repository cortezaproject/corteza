---
kind: file
covers: Index.vue
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue/src/components/navigation/CTopbar.vue
tests: []
---

# Navigation Index view

## Intention

Configure the app shell's topbar chrome for the whole instance: which built-in
buttons appear, plus custom profile-menu links and per-page buttons — all
stored in the single `ui.topbar` setting object.

## UX capabilities

- Hide toggles: app selector, drafts, notifications, search (only offered when `discovery.enabled`), profile menu.
- Profile-menu options: hide profile / change-password / theme-selector links; editable list of custom profile links (handle, url, newTab).
- Page-buttons list: label, url, urlMatch (route substring that makes the button appear), newTab, description — rendered live by `CTopbar` (lib/vue), which matches urlMatch against the current URL.
- Single save writes the whole `ui.topbar` object back.

## Routes

`ui.navigation` at `/ui/navigation`; no params.

## When changing this

- Drafts and search are stored positively (`showDrafts` / `showSearch`) but presented as hide-checkboxes — keep the inversion when adding toggles.
- Section topbar overrides merge over `ui.topbar` at runtime — renaming keys inside the object is a cross-app contract change.

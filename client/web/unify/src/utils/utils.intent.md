---
kind: folder
covers: recursive
owner: fe
depends-on: []
touched-by: []
tests: []
---

# Utils

## Intention

Small app-wide helpers with no state and no section affiliation. Anything reusable
beyond this app belongs in `lib/vue` / `lib/js` instead.

## Map

- `appIcons.js` — maps application/section ids to their icon assets for the app list & topbar.
- `appReachable.js` — `appReachable` / `useAppReachable`: whether the app menu
  should offer a registry application. A `unify.url` this shell serves is
  judged by the section it lands in (the same rule the router's gate applies,
  so what is offered and what is admitted cannot drift); a url no section
  serves — an address elsewhere, or an app served outside the shell — has no
  section to speak for it and is judged by its own `canAccessApplication`. So
  is a custom application at `app/<id>`, whose section (`perApp`) is keyed to no
  application of its own. `enabled` is CAppList's own check and is applied on
  top of all three.
- `documentTitle.js` — the browser tab title.

## When changing this

- Keep this folder tiny; if a helper grows state or feature logic it belongs in a section or a store.

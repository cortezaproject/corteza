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

## When changing this

- Keep this folder tiny; if a helper grows state or feature logic it belongs in a section or a store.

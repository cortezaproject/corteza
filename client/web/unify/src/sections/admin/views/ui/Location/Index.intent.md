---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Location Index view

## Intention

Choose the geosearch provider (and its API key when needed) that map
components across the webapp use for address lookup — stored as the single
`ui.location` setting object.

## UX capabilities

- Provider dropdown over a fixed catalog (OpenStreetMap, OpenCage, Esri, Geoapify, Geocode Earth, Google Maps, LocationIQ, Mapbox, Pelias).
- API-key input shown only for providers that require one; changing provider clears the stored key.
- Save writes `{ geoSearchProvider, geoSearchApiKey }` back via `$SystemAPI.settingsUpdate`; missing setting falls back to OpenStreetMap.

## Routes

`ui.location` at `/ui/location`; no params.

## When changing this

- Map components consume `ui.location` at runtime — provider values are a cross-app contract; adding a provider means knowing whether it needs a key.
- Keep the clear-key-on-provider-change behavior; a stale key for the wrong provider is worse than none.

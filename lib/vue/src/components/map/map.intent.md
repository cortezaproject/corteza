---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by: []
tests: []
---

# map/

## Intention

The one shared Leaflet map (OSM tiles) so geometry fields, location inputs,
and map blocks all get identical geosearch, themed controls, and marker
rendering. Exported from the components barrel like every other shared
component (historically deep-imported; fixed as an oversight).

## Map

- `CMap.vue` — LMap wrapper + geosearch box (leaflet-geosearch/OSM), PrimeVue
  zoom/locate controls, imperative marker/polygon rendering into a layerGroup.

## Contracts consumers rely on

- Declarative in, evented out: `markers` (`{value:[lat,lng], title?, color?}`)
  and `polygons` (`{latLngs, color?, fillOpacity?}`) props are re-rendered on
  change; invalid entries silently dropped. Emits `ready` (the leaflet object),
  `map-click`, `marker-click`, `location-found`, `update:center/zoom/bounds`.
- Geosearch picks emit a synthetic `map-click` — pick-a-location flows get it
  for free; value commit stays the caller's job.
- Exposes `invalidateSize()` / `fitBounds()`; size is auto-invalidated via
  ResizeObserver.

## When changing this

- Gotcha (verified in code): vue-leaflet 0.10's LMap does not reliably forward
  leaflet events as Vue events — CMap binds `click`/`moveend` directly on the
  leaflet object inside the `@ready` handler. Never rely on `@click` on LMap.
- `use-global-leaflet` must stay `true`: vue-leaflet and the direct `L.*`
  calls must share one leaflet module, else bounds math crosses two prototype
  chains and throws.
- Marker clicks stop propagation so they don't double-fire `map-click`.
